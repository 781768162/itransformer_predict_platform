package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"gateway/applications/router"
	"gateway/config"
	"gateway/internal/database"
	"gateway/internal/mq"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 初始化数据库连接，增加重试以等待容器内 MySQL 就绪
	var dbErr error
	for i := 1; i <= 10; i++ {
		dbErr = database.NewDBHandler()
		if dbErr == nil {
			break
		}
		log.Printf("init db failed (attempt %d/10): %v", i, dbErr)
		time.Sleep(2 * time.Second)
	}
	if dbErr != nil {
		panic(dbErr)
	}

	cfg := mq.Config{
		Brokers:       config.Settings.Kafka.Brokers,
		ConsumerTopic: config.Settings.Kafka.ConsumerTopic,
		GroupID:       config.Settings.Kafka.GroupID,
	}

	mq.StartConsumerWithRetry(ctx, cfg, mq.HandleTaskResult) // 启动消费者协程

	r := router.MustNewRouter()

	r.Run(config.Settings.Server.Addr)
}
