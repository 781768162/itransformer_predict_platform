package mq

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// Config 保存 Kafka 连接配置。
type Config struct {
	Brokers       []string
	ProducerTopic string // 生产端使用的 topic（如 task_input）
	ConsumerTopic string // 消费端使用的 topic（如 task_result）
	GroupID       string
}

func NewProducer(cfg Config) *kafka.Writer {
	// 生产者
	prod := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Topic:        cfg.ProducerTopic,
		Balancer:     &kafka.Hash{},
		RequiredAcks: kafka.RequireAll,
		Async:        false,
	}

	return prod
}

func NewConsumer(cfg Config) *kafka.Reader {
	// 异步消费者
	cons := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Brokers,
		GroupID:        cfg.GroupID,
		Topic:          cfg.ConsumerTopic,
		MinBytes:       1 << 10,  // 1KB
		MaxBytes:       10 << 20, // 10MB
		CommitInterval: time.Second,
	})

	return cons
}

// ConsumerLoop 启动异步消费循环。
// handler 处理消息：返回 nil 表示提交 offset，返回非 nil 不提交（消息将被重投）。
func ConsumerLoop(ctx context.Context, c *kafka.Reader, handler func(context.Context, kafka.Message) error) error {
	for {
		msg, err := c.FetchMessage(ctx)
		if err != nil {
			return err
		}
		if err := handler(ctx, msg); err != nil {
			continue
		}
		if err := c.CommitMessages(ctx, msg); err != nil {
			return err
		}
	}
}

// StartConsumerWithRetry 启动消费者循环，异常退出后在 ctx 未取消时自动重试。
func StartConsumerWithRetry(ctx context.Context, cfg Config, handler func(context.Context, kafka.Message) error) {
	go func() {
		for {
			consumer := NewConsumer(cfg)
			if err := ConsumerLoop(ctx, consumer, handler); err != nil {
				log.Printf("consumer stopped: %v", err)
			}
			consumer.Close()

			if ctx.Err() != nil {
				return
			}

			time.Sleep(time.Second)
		}
	}()
}
