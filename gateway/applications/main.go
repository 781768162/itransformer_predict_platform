package main

import (
	"context"
	"os/signal"
	"syscall"

	"gateway/applications/router"
	"gateway/config"
	"gateway/internal/mq"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := mq.Config{
		Brokers:       config.Settings.Kafka.Brokers,
		ConsumerTopic: config.Settings.Kafka.ConsumerTopic,
		GroupID:       config.Settings.Kafka.GroupID,
	}

	mq.StartConsumerWithRetry(ctx, cfg, mq.HandleTaskResult) // 启动消费者协程

	r := router.MustNewRouter()

	r.Run(config.Settings.Server.Addr)
}
