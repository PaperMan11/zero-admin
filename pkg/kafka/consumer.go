package kafkaclient

import (
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stat"
	"io"
	"sync"
	"time"
)

const (
	defaultCommitInterval = time.Second
	defaultMaxWait        = time.Second
)

type (
	ConsumeErrorHandler func(ctx context.Context, msg kafka.Message, err error)
	ConsumeHandler      interface {
		Consume(ctx context.Context, key, value string) error
	}

	KafkaConsumer struct {
		ctx          context.Context
		cancel       context.CancelFunc
		reader       *kafka.Reader
		c            KqConf
		handler      ConsumeHandler
		errorHandler ConsumeErrorHandler
		metrics      *stat.Metrics
		forceCommit  bool
		channel      chan *kafka.Message
		wg           sync.WaitGroup
	}

	consumerOptions struct {
		commitInterval time.Duration
		maxWait        time.Duration
		metrics        *stat.Metrics
		errorHandler   ConsumeErrorHandler
	}
	ConsumerOptions func(*consumerOptions)
)

func WithCommitInterval(interval time.Duration) ConsumerOptions {
	return func(o *consumerOptions) {
		o.commitInterval = interval
	}
}

func WithMaxWait(maxWait time.Duration) ConsumerOptions {
	return func(o *consumerOptions) {
		o.maxWait = maxWait
	}
}

func NewKafkaConsumer(c KqConf, handler ConsumeHandler, opts ...ConsumerOptions) *KafkaConsumer {
	var options consumerOptions
	for _, opt := range opts {
		opt(&options)
	}
	ensureConsumerOptions(c, &options)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        c.Brokers,
		Topic:          c.Topic,
		GroupID:        c.Group,
		MinBytes:       c.MinBytes,
		MaxBytes:       c.MaxBytes,
		MaxWait:        options.maxWait,
		CommitInterval: options.commitInterval,
	})

	ctx, cancel := context.WithCancel(context.Background())
	return &KafkaConsumer{
		ctx:          ctx,
		cancel:       cancel,
		reader:       r,
		c:            c,
		handler:      handler,
		errorHandler: options.errorHandler,
		metrics:      options.metrics,
		forceCommit:  c.ForceCommit,
		channel:      make(chan *kafka.Message, c.Processors*2),
	}
}

func ensureConsumerOptions(c KqConf, options *consumerOptions) {
	if options.commitInterval == 0 {
		options.commitInterval = defaultCommitInterval
	}
	if options.maxWait == 0 {
		options.maxWait = defaultMaxWait
	}
	if options.metrics == nil {
		options.metrics = stat.NewMetrics(c.ServerName)
	}
	if options.errorHandler == nil {
		options.errorHandler = func(ctx context.Context, msg kafka.Message, err error) {
			logc.Errorf(ctx, "consume: %s, error: %v", string(msg.Value), err)
		}
	}
}

func (c *KafkaConsumer) Stop() error {
	c.cancel()
	close(c.channel)
	c.reader.Close()
	c.wg.Wait()
	return nil
}

func (c *KafkaConsumer) Start() error {
	c.startProducers()
	c.startConsumers()
	return nil
}

func (c *KafkaConsumer) startConsumers() {
	for i := 0; i < c.c.Processors; i++ {
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			for {
				select {
				case <-c.ctx.Done():
					return
				default:
				}
				message := <-c.channel
				if err := c.handler.Consume(c.ctx, string(message.Key), string(message.Value)); err != nil {
					if c.errorHandler != nil {
						c.errorHandler(c.ctx, *message, err)
					}
					if !c.forceCommit {
						continue
					}
				}
				if err := c.reader.CommitMessages(c.ctx, *message); err != nil {
					logc.Errorf(c.ctx, "commit message: %s, error: %v", string(message.Value), err)
				}
			}
		}()
	}
}

func (c *KafkaConsumer) startProducers() {
	for i := 0; i < c.c.Consumers; i++ {
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			for {
				select {
				case <-c.ctx.Done():
					return
				default:
				}
				message, err := c.reader.FetchMessage(c.ctx)
				// io.EOF means consumer closed
				// io.ErrClosedPipe means committing messages on the consumer,
				// kafka will refire the messages on uncommitted messages, ignore
				if err == io.EOF || errors.Is(err, io.ErrClosedPipe) {
					logx.Errorf("Error on reading message, %q", err.Error())
					return
				}
				if err != nil {
					logx.Errorf("Error on reading message, %q", err.Error())
					continue
				}
				c.channel <- &message
			}
		}()
	}
}
