package redisclient

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"testing"
	"time"
)

// 测试用的Redis客户端
func getTestRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "192.168.241.128:6379", // 默认Redis地址
		Password: "123456",               // 无密码
		DB:       0,                      // 使用测试专用DB，避免影响正式数据
	})
}

// 清理测试数据
func cleanupTestData(t *testing.T, client *redis.Client, stream string) {
	ctx := context.Background()
	// 删除测试流
	if err := client.Del(ctx, stream).Err(); err != nil {
		t.Logf("清理测试数据失败: %v", err)
	}
}

// 测试生产者发送消息
func TestStreamProducer(t *testing.T) {
	client := getTestRedisClient()
	defer client.Close()

	testStream := "test_stream_producer"
	cleanupTestData(t, client, testStream)
	defer cleanupTestData(t, client, testStream)

	ctx := context.Background()
	testData := map[string]interface{}{
		"id":   "123",
		"name": "test message",
		"time": time.Now().Unix(),
	}

	// 测试发送消息
	err := StreamProducer(client, testStream, testData)
	if err != nil {
		t.Fatalf("发送消息失败: %v", err)
	}

	// 验证消息是否存在
	result, err := client.XRange(ctx, testStream, "-", "+").Result()
	if err != nil {
		t.Fatalf("查询消息失败: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("消息数量不正确，期望1条，实际%v条", len(result))
	}

	msg := result[0]
	if msg.Values["id"] != testData["id"] {
		t.Errorf("消息内容不匹配，期望id=%v，实际id=%v", testData["id"], msg.Values["id"])
	}
}

// 测试消费者基本功能
func TestStreamConsumer_Basic(t *testing.T) {
	InitTestLog(t)
	client := getTestRedisClient()
	defer client.Close()

	testStream := "test_stream_consumer_basic"
	testGroup := "test_group"
	testConsumer := "test_consumer"
	cleanupTestData(t, client, testStream)
	defer cleanupTestData(t, client, testStream)

	// 配置消费者
	cfg := &StreamConsumerConfig{
		Stream:           testStream,
		Group:            testGroup,
		Consumer:         testConsumer,
		MaxMsgSize:       100,
		BlockTime:        2 * time.Second,
		BatchSize:        10,
		CleanupInterval:  time.Minute,
		RetentionPeriod:  time.Hour,
		CleanupBatchSize: 100,
	}

	// 创建消费者
	consumer := NewStreamConsumer(client, cfg)
	defer consumer.Stop()

	// 用于验证消息处理的通道
	messageChan := make(chan map[string]interface{}, 10)

	// 启动消费者
	go consumer.Start(func(message map[string]interface{}) error {
		messageChan <- message
		return nil
	})

	// 等待消费者启动
	time.Sleep(100 * time.Millisecond)

	// 发送测试消息
	testData := map[string]interface{}{"test_key": "test_value", "num": 123}
	if err := StreamProducer(client, testStream, testData); err != nil {
		t.Fatalf("发送测试消息失败: %v", err)
	}

	// 等待消息处理
	select {
	case receivedMsg := <-messageChan:
		if receivedMsg["test_key"] != testData["test_key"] {
			t.Errorf("消息内容不匹配，期望%v，实际%v", testData["test_key"], receivedMsg["test_key"])
		}
	case <-time.After(5 * time.Second):
		t.Fatal("超时未收到消息")
	}
}

// 测试未确认消息处理
func TestStreamConsumer_PendingMessages(t *testing.T) {
	InitTestLog(t)
	client := getTestRedisClient()
	defer client.Close()

	testStream := "test_stream_pending"
	testGroup := "test_group_pending"
	testConsumer := "test_consumer_pending"
	cleanupTestData(t, client, testStream)
	defer cleanupTestData(t, client, testStream)

	ctx := context.Background()

	// 手动创建消费者组
	if err := client.XGroupCreateMkStream(ctx, testStream, testGroup, "0-0").Err(); err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		t.Fatalf("创建消费者组失败: %v", err)
	}

	// 发送测试消息
	_, err := client.XAdd(ctx, &redis.XAddArgs{
		Stream: testStream,
		Values: map[string]interface{}{"pending": "true", "data": "test"},
	}).Result()
	if err != nil {
		t.Fatalf("发送消息失败: %v", err)
	}

	// 手动将消息标记为已读取但未确认（模拟消费者崩溃场景）
	_, err = client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    testGroup,
		Consumer: testConsumer,
		Streams:  []string{testStream, ">"},
		Count:    1,
	}).Result()
	if err != nil {
		t.Fatalf("读取消息失败: %v", err)
	}

	// 配置消费者
	cfg := &StreamConsumerConfig{
		Stream:    testStream,
		Group:     testGroup,
		Consumer:  testConsumer,
		BlockTime: time.Second,
		BatchSize: 10,
	}

	// 创建并启动消费者
	consumer := NewStreamConsumer(client, cfg)
	defer consumer.Stop()

	messageChan := make(chan string, 1)

	go consumer.Start(func(message map[string]interface{}) error {
		if val, ok := message["data"]; ok && val == "test" {
			messageChan <- "received"
		}
		return nil
	})

	// 等待未确认消息被处理
	select {
	case <-messageChan:
		// 验证消息已被确认
		pending, err := client.XPending(ctx, testStream, testGroup).Result()
		if err != nil {
			t.Fatalf("查询未确认消息失败: %v", err)
		}
		if pending.Count != 0 {
			t.Errorf("未确认消息数量不为0，实际为%v", pending.Count)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("超时未处理未确认消息")
	}
}

// 测试消息清理功能
func TestStreamConsumer_Cleanup(t *testing.T) {
	InitTestLog(t)
	client := getTestRedisClient()
	defer client.Close()

	testStream := "test_stream_cleanup"
	testGroup := "test_group_cleanup"
	testConsumer := "test_consumer_cleanup"
	cleanupTestData(t, client, testStream)
	//defer cleanupTestData(t, client, testStream)

	ctx := context.Background()

	// 配置消费者，使用短清理间隔和保留时间便于测试
	cfg := &StreamConsumerConfig{
		Stream:           testStream,
		Group:            testGroup,
		Consumer:         testConsumer,
		BlockTime:        time.Second,
		BatchSize:        10,
		CleanupInterval:  1 * time.Second, // 1秒清理一次
		RetentionPeriod:  2 * time.Second, // 消息保留2秒
		CleanupBatchSize: 100,
	}

	// 创建消费者
	consumer := NewStreamConsumer(client, cfg)
	defer consumer.Stop()

	// 启动消费者
	go consumer.Start(func(message map[string]interface{}) error {
		// 立即确认消息
		return nil
	})

	// 等待消费者启动
	time.Sleep(100 * time.Millisecond)

	// 发送测试消息
	for i := 0; i < 5; i++ {
		if err := StreamProducer(client, testStream, map[string]interface{}{"id": i}); err != nil {
			t.Fatalf("发送消息失败: %v", err)
		}
	}

	// 等待消息被处理和确认
	time.Sleep(1 * time.Second)

	// 等待清理任务执行
	time.Sleep(3 * time.Second) // 超过保留时间

	// 检查消息是否被清理
	result, err := client.XRange(ctx, testStream, "-", "+").Result()
	if err != nil {
		t.Fatalf("查询消息失败: %v", err)
	}

	if len(result) > 0 {
		t.Errorf("消息未被正确清理，剩余%v条", len(result))
	}
}

// 测试消息清理功能
func TestStreamConsumer_Cleanup2(t *testing.T) {
	InitTestLog(t)
	client := getTestRedisClient()

	defer client.Close()

	testStream := "test_stream_cleanup"
	testGroup := "test_group_cleanup"
	testConsumer := "test_consumer_cleanup"
	cleanupTestData(t, client, testStream)
	//defer cleanupTestData(t, client, testStream)

	ctx := context.Background()

	// 配置消费者，使用短清理间隔和保留时间便于测试
	cfg := &StreamConsumerConfig{
		Stream:           testStream,
		Group:            testGroup,
		Consumer:         testConsumer,
		BlockTime:        time.Second,
		BatchSize:        10,
		CleanupInterval:  1 * time.Second, // 1秒清理一次
		RetentionPeriod:  2 * time.Second, // 消息保留2秒
		CleanupBatchSize: 100,
	}

	// 创建消费者
	consumer := NewStreamConsumer(client, cfg)
	defer consumer.Stop()

	// 启动消费者
	go consumer.Start(func(message map[string]interface{}) error {
		// 立即确认消息
		id := message["id"].(string)
		if id == "1" {
			return errors.New("处理消息失败") // 模拟消息消息处理失败
		}
		return nil
	})

	// 等待消费者启动
	time.Sleep(100 * time.Millisecond)

	// 发送测试消息
	for i := 0; i < 5; i++ {
		if err := StreamProducer(client, testStream, map[string]interface{}{"id": i}); err != nil {
			t.Fatalf("发送消息失败: %v", err)
		}
	}

	// 等待消息被处理和确认
	time.Sleep(1 * time.Second)

	// 等待清理任务执行
	time.Sleep(3 * time.Second) // 超过保留时间

	// 检查消息是否被清理
	result, err := client.XRange(ctx, testStream, "-", "+").Result()
	if err != nil {
		t.Fatalf("查询消息失败: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("消息未被正确清理，剩余%v条", len(result))
	}
}

// 测试并发场景
func TestStreamConsumer_Concurrent(t *testing.T) {
	InitTestLog(t)
	client := getTestRedisClient()
	defer client.Close()

	testStream := "test_stream_concurrent"
	testGroup := "test_group_concurrent"
	cleanupTestData(t, client, testStream)
	defer cleanupTestData(t, client, testStream)

	msgCount := 100
	received := make(chan struct{}, msgCount)

	// 启动多个消费者
	consumerCount := 3
	for i := 0; i < consumerCount; i++ {
		cfg := &StreamConsumerConfig{
			Stream:    testStream,
			Group:     testGroup,
			Consumer:  "consumer_" + string(rune(i)),
			BlockTime: time.Second,
			BatchSize: 10,
		}

		consumer := NewStreamConsumer(client, cfg)
		defer consumer.Stop()

		go consumer.Start(func(message map[string]interface{}) error {
			received <- struct{}{}
			return nil
		})
	}

	// 等待消费者启动
	time.Sleep(100 * time.Millisecond)

	// 并发发送消息
	for i := 0; i < msgCount; i++ {
		go func(id int) {
			StreamProducer(client, testStream, map[string]interface{}{"id": id})
		}(i)
	}

	// 等待所有消息被处理
	timeout := time.After(10 * time.Second)
	for i := 0; i < msgCount; i++ {
		select {
		case <-received:
			continue
		case <-timeout:
			t.Fatalf("超时，只收到%v条消息，期望%v条", i, msgCount)
		}
	}

	t.Logf("成功处理所有%d条消息", msgCount)
}

// 初始化测试日志
func InitTestLog(t *testing.T) {
	// 使用zaptest简化测试日志
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)
}
