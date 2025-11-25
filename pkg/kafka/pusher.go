package kafkaclient

import (
	"context"
	"github.com/segmentio/kafka-go"
	"strconv"
	"time"
)

type (
	KafkaPusher struct {
		topic  string
		writer *kafka.Writer
	}

	pushOptions struct {
		// kafka.Writer options
		allowAutoTopicCreation bool
		balancer               kafka.Balancer

		// syncPush is used to enable sync push
		syncPush bool
	}

	PushOption func(options *pushOptions)
)

func NewKafkaPusher(addrs []string, topic string, opts ...PushOption) *KafkaPusher {
	producer := &kafka.Writer{
		Addr:        kafka.TCP(addrs...),
		Topic:       topic,
		Balancer:    &kafka.LeastBytes{},
		Compression: kafka.Snappy,
	}

	var options pushOptions
	for _, opt := range opts {
		opt(&options)
	}

	// apply kafka.Writer options
	producer.AllowAutoTopicCreation = options.allowAutoTopicCreation
	if options.balancer != nil {
		producer.Balancer = options.balancer
	}

	pusher := &KafkaPusher{
		writer: producer,
		topic:  topic,
	}

	// if syncPush is true, return the pusher directly
	if options.syncPush {
		producer.BatchSize = 1
		return pusher
	}

	return pusher
}

func (p *KafkaPusher) Push(ctx context.Context, v string) error {
	return p.PushWithKey(ctx, strconv.FormatInt(time.Now().UnixNano(), 10), v)
}

// PushWithKey sends a message with the given key to the Kafka topic.
func (p *KafkaPusher) PushWithKey(ctx context.Context, key, v string) error {
	msg := kafka.Message{
		Key:   []byte(key),
		Value: []byte(v),
	}

	// todo: inject trace context into message

	return p.writer.WriteMessages(ctx, msg)
}

func WithAllowAutoTopicCreation(allowAutoTopicCreation bool) PushOption {
	return func(options *pushOptions) {
		options.allowAutoTopicCreation = allowAutoTopicCreation
	}
}

func WithBalancer(balancer kafka.Balancer) PushOption {
	return func(options *pushOptions) {
		options.balancer = balancer
	}
}
