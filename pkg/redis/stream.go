package redisclient

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

const (
	BlockTime        = time.Second * 5
	MaxMsgSize       = 100000
	BatchSize        = 100
	CleanupInterval  = time.Minute * 10
	CleanupBatchSize = 1000
	RetentionPeriod  = time.Hour * 24 * 7
)

func StreamProducer(client *redis.Client, stream string, message map[string]interface{}) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	return client.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		MaxLen: MaxMsgSize,
		Approx: true,
		Values: message,
	}).Err()
}

// --------------------------------------------------------------------------------------------------------------------

// StreamConsumerConfig 配置参数
type StreamConsumerConfig struct {
	Stream           string        // stream 名称
	Group            string        // 消费者组名
	Consumer         string        // 消费者名
	MaxMsgSize       int64         // stream 最大长度
	BlockTime        time.Duration // 拉取消息阻塞时间
	BatchSize        int64         // 每次拉取的消息数量
	CleanupInterval  time.Duration // 清理已确认消息的间隔时间
	RetentionPeriod  time.Duration // 消息保留时间，超过此时长确认的消息将被删除
	CleanupBatchSize int64         // 每次清理的消息批次大小
}

// StreamConsumer 结构体
type StreamConsumer struct {
	cfg    *StreamConsumerConfig
	client *redis.Client
	ctx    context.Context
	cancel context.CancelFunc
}

// StreamMessageHandler 消息处理函数
type StreamMessageHandler func(message map[string]interface{}) error

// NewStreamConsumer 创建一个新的消费者
func NewStreamConsumer(client *redis.Client, cfg *StreamConsumerConfig) *StreamConsumer {
	// 设置默认清理间隔为10分钟
	if cfg.CleanupInterval <= 0 {
		cfg.CleanupInterval = CleanupInterval
	}

	// 设置默认消息保留时间为7天
	if cfg.RetentionPeriod <= 0 {
		cfg.RetentionPeriod = RetentionPeriod
	}

	// 设置默认清理批次大小
	if cfg.CleanupBatchSize <= 0 {
		cfg.CleanupBatchSize = CleanupBatchSize
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &StreamConsumer{
		cfg:    cfg,
		client: client,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start 启动消费者
func (c *StreamConsumer) Start(handler StreamMessageHandler) {
	// 确保消费者组存在
	if err := c.EnsureStreamAndGroupExist(); err != nil {
		zap.L().Fatal("初始化消费者组失败", zap.Error(err))
		return
	}

	// 处理未确认消息
	c.processPendingMessages(handler)

	// 正常消费新消息
	go c.pollNewMessages(handler)

	// 启动定期清理已确认消息的goroutine
	go c.startCleanupTask()

	zap.L().Debug("Redis Stream 消费者已启动",
		zap.String("stream", c.cfg.Stream),
		zap.String("group", c.cfg.Group),
		zap.Duration("cleanup_interval", c.cfg.CleanupInterval))
}

// 确保消费者组存在
func (c *StreamConsumer) EnsureStreamAndGroupExist() (err error) {
	if err = c.ensureStreamExist(); err != nil {
		return
	}
	return
}

func (c *StreamConsumer) ensureStreamExist() error {
	err := c.client.XGroupCreateMkStream(c.ctx, c.cfg.Stream, c.cfg.Group, "0-0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return fmt.Errorf("初始化Stream失败: %v", err)
	}
	return nil
}

func (c *StreamConsumer) ensureGroupExist() error {
	err := c.client.XGroupCreate(c.ctx, c.cfg.Stream, c.cfg.Group, "0-0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		zap.L().Error("创建消费者组失败", zap.Error(err))
		return err
	}
	return nil
}

// 处理未确认消息
//
//	1.同一消费者重启（组内消费者名不变）: 主动读取 PEL 中的消息：使用 XPENDING 查看未确认消息，再通过 XREADGROUP 以具体消息 ID 拉取（如指定 START 0 回溯）
//		# 查看组内未确认消息
//		XPENDING mystream mygroup
//		# 拉取未确认的消息（包括之前未处理的）
//		XREADGROUP GROUP mygroup consumer1 STREAMS mystream 0  # 0 表示从最早未确认的消息开始
//
//	2.消费者组内其他消费者接管: 若原消费者崩溃且未重启，组内其他消费者（如 consumer2）可通过 XCLAIM 命令将 PEL 中的未确认消息 “认领” 到自己名下，进行处理
//		# 将超时未处理的消息转移给 consumer2（例如 30 秒未确认）
//		XCLAIM mystream mygroup consumer2 30000 1620000000000-0
//
//	- 未确认的消息不会因消费者重启而丢失，而是保存在消费者组的 PEL 中。
//	- 重启后的消费者（无论是否同名）可以重新消费这些未确认的消息，前提是主动读取 PEL 或通过 XCLAIM 认领。
//	- 若消费者不主动处理 PEL 消息，这些消息会一直留在列表中，直到被确认（XACK）或被修剪（如 Stream 触发 MAXLEN 清理）。
//	与 Kafka 的对比
//	- Redis Stream 中，未确认消息通过 PEL 显式追踪，重启后需主动处理（或配置自动认领）。
//	- Kafka 中，未提交偏移量的消息会在消费者重启 / 重平衡后，从上次提交的偏移量开始重新消费，无需显式 “认领”。
//	- 两者均保证未确认 / 未提交的消息可被重新处理，但实现机制不同。
func (c *StreamConsumer) processPendingMessages(handler StreamMessageHandler) {
	pending, _ := c.client.XPending(c.ctx, c.cfg.Stream, c.cfg.Group).Result()

	if pending.Count > 0 {
		zap.L().Debug("处理未确认消息", zap.Int64("count", pending.Count))

		// 获取未确认消息详情
		pendingMsgs, err := c.client.XPendingExt(c.ctx, &redis.XPendingExtArgs{
			Stream: c.cfg.Stream,
			Group:  c.cfg.Group,
			Start:  "-",
			End:    "+",
			Count:  pending.Count,
		}).Result()

		if err != nil {
			zap.L().Error("获取未确认消息详情失败", zap.Error(err))
			return
		}

		// 处理未确认消息
		for _, msg := range pendingMsgs {
			select {
			case <-c.ctx.Done():
				return
			default:
			}
			zap.L().Debug("处理未确认消息", zap.String("message_id", msg.ID))

			// 读取具体消息内容
			message, err := c.client.XRange(c.ctx, c.cfg.Stream, msg.ID, msg.ID).Result()
			if err != nil {
				zap.L().Error("读取消息失败", zap.String("message_id", msg.ID), zap.Error(err))
				continue
			}

			// 消息如果被xdel或自动裁剪则找不到该消息
			if len(message) > 0 {
				// 处理消息
				if err = handler(message[0].Values); err != nil {
					zap.L().Error("解析消息失败", zap.String("message_id", msg.ID), zap.Any("raw_values", message[0].Values), zap.Error(err))
				}
			} else {
				zap.L().Error("消息已被删除", zap.String("message_id", msg.ID))
			}

			// 确认消息
			if err := c.client.XAck(c.ctx, c.cfg.Stream, c.cfg.Group, msg.ID).Err(); err != nil {
				zap.L().Error("确认消息失败", zap.String("message_id", msg.ID), zap.Error(err))
			}
		}
	}
}

// 启动定期清理任务
func (c *StreamConsumer) startCleanupTask() {
	ticker := time.NewTicker(c.cfg.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			zap.L().Debug("清理已确认消息的任务已停止", zap.String("stream", c.cfg.Stream))
			return
		case <-ticker.C:
			// 执行清理
			deletedCount, err := c.cleanupConfirmedMessages()
			if err != nil {
				zap.L().Error("清理已确认消息失败", zap.String("stream", c.cfg.Stream), zap.Error(err))
			} else if deletedCount > 0 {
				zap.L().Debug("清理已确认消息完成",
					zap.String("stream", c.cfg.Stream),
					zap.Int64("deleted_count", deletedCount))
			}
		}
	}
}

// 清理已确认消息
func (c *StreamConsumer) cleanupConfirmedMessages() (int64, error) {
	var totalDeleted int64 = 0

	// 计算截止时间：当前时间减去保留时间
	cutoffTime := time.Now().Add(-c.cfg.RetentionPeriod).UnixMilli()
	cutoffID := fmt.Sprintf("%d-0", cutoffTime)

	// 循环删除直到没有更多符合条件的消息
	for {
		// 获取截止时间之前的消息
		messages, err := c.client.XRange(c.ctx, c.cfg.Stream, "-", cutoffID).Result()
		if err != nil {
			return totalDeleted, fmt.Errorf("获取消息范围失败: %v", err)
		}

		if len(messages) == 0 {
			break // 没有更多消息需要删除
		}

		// 收集消息ID
		messageIDs := make([]string, len(messages))
		for i, msg := range messages {
			messageIDs[i] = msg.ID
		}

		// 检查这些消息是否已被确认
		pendingMsgs, err := c.client.XPendingExt(c.ctx, &redis.XPendingExtArgs{
			Stream: c.cfg.Stream,
			Group:  c.cfg.Group,
			Start:  messageIDs[0],
			End:    messageIDs[len(messageIDs)-1],
			Count:  c.cfg.CleanupBatchSize,
		}).Result()
		if err != nil && err != redis.Nil {
			return totalDeleted, fmt.Errorf("获取未确认消息失败: %v", err)
		}

		// 找出已确认的消息（不在未确认列表中的消息）
		pendingIDs := make(map[string]bool)
		for _, msg := range pendingMsgs {
			pendingIDs[msg.ID] = true
		}

		var confirmedIDs []string
		for _, id := range messageIDs {
			if !pendingIDs[id] {
				confirmedIDs = append(confirmedIDs, id)
			}
		}

		if len(confirmedIDs) > 0 {
			// 删除已确认的消息
			deleted, err := c.client.XDel(c.ctx, c.cfg.Stream, confirmedIDs...).Result()
			if err != nil {
				return totalDeleted, fmt.Errorf("删除消息失败: %v", err)
			}

			totalDeleted += deleted
		}

		// 如果获取的消息数量小于批次大小，说明没有更多消息了
		if int64(len(messages)) < c.cfg.CleanupBatchSize {
			break
		}
	}

	//// 同时修剪流的长度，确保不会无限增长
	//if c.cfg.MaxMsgSize > 0 {
	//	if err := c.client.XTrim(c.ctx, c.cfg.Stream, c.cfg.MaxMsgSize).Err(); err != nil {
	//		zap.L().Warn("修剪流长度失败", zap.String("stream", c.cfg.Stream), zap.Error(err))
	//	}
	//}

	return totalDeleted, nil
}

func (c *StreamConsumer) pollNewMessages(handler StreamMessageHandler) {
	retryInterval := time.Second
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			results, err := c.client.XReadGroup(c.ctx, &redis.XReadGroupArgs{
				Group:    c.cfg.Group,
				Consumer: c.cfg.Consumer,
				Streams:  []string{c.cfg.Stream, ">"}, // '>'表示最新未分配消息, '$'表示最新消息
				Block:    c.cfg.BlockTime,
				Count:    c.cfg.BatchSize,
				NoAck:    false,
			}).Result()

			if err != nil {
				if err == redis.Nil {
					time.Sleep(retryInterval)
					continue
				}
				time.Sleep(retryInterval)
				if retryInterval < 10*time.Second {
					retryInterval *= 2
				}
				continue
			}
			retryInterval = time.Second

			for _, stream := range results {
				for _, msg := range stream.Messages {
					zap.L().Debug("收到消息", zap.String("message_id", msg.ID))

					// 处理消息
					if err = handler(msg.Values); err != nil {
						zap.L().Error("解析消息失败", zap.String("message_id", msg.ID), zap.Any("raw_values", msg.Values), zap.Error(err))
						continue
					}

					ackErr := c.client.XAck(c.ctx, c.cfg.Stream, c.cfg.Group, msg.ID).Err()
					if ackErr != nil {
						zap.L().Error("确认消息失败", zap.String("message_id", msg.ID), zap.Error(ackErr))
					}
				}
			}
		}
	}
}

// Stop 停止消费者
func (c *StreamConsumer) Stop() {
	c.cancel()
	zap.L().Debug("Redis Stream 消费者已停止", zap.String("stream", c.cfg.Stream), zap.String("group", c.cfg.Group))
}
