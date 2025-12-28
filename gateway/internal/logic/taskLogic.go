package logic

import (
	"context"
	"encoding/json"
	"errors"

	"gateway/config"
	"gateway/internal/code"
	"gateway/internal/database/crud"
	"gateway/internal/database/model"
	"gateway/internal/mq"
	"gateway/pkg/logger"

	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

var (
	prod *kafka.Writer
)

func init() {
	cfg := mq.Config{
		Brokers:       config.Settings.Kafka.Brokers,
		ProducerTopic: config.Settings.Kafka.ProducerTopic,
	}
	prod = mq.NewProducer(cfg)
}

type taskPayload struct {
    TaskID     int64              `json:"task_id"`
    PassData   [13][72]float64  `json:"pass_data"`
    FutureData [12][24]float64  `json:"future_data"`
}

func CreateTaskLogic(ctx context.Context, userId int64, date string, passData [13][72]float64, futureData [12][24]float64) (int, error) {
	t := &model.Task{
		UserID: userId,
		Date:   date,
		Status: "pending",
	}
	err := crud.CreateTask(ctx, t) // 插入记录
	if err != nil {
		logger.Errorf("CreateTask error: %v", err)
		return 0, code.ErrDatabase
	}

	body, err := json.Marshal(taskPayload{
        TaskID:     t.TaskID,
        PassData:   passData,
        FutureData: futureData,
    })
    if err != nil {
		logger.Errorf("Marshal error: %v", err)
        return 0, code.ErrJsonMarshal
    }

	err = prod.WriteMessages(ctx, kafka.Message{ // kafka投递消息
		Value: body,
	})
	if err != nil {
		logger.Errorf("WriteMessages message: %s error: %v", string(body), err)
		return 0, code.ErrMessageQueue
	}

	return int(t.TaskID), nil
}

func GetTaskLogic(ctx context.Context, taskId int) (string, string, [24]float64, error) {
	status, date, err := crud.GetTaskStatusAndDate(ctx, int64(taskId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Errorf("taskId %d not found", taskId)
			return "", "", [24]float64{}, code.ErrNotFound
		}else {
			logger.Errorf("GetTaskStatusAndDate error: %v", err)
			return "", "", [24]float64{}, code.ErrDatabase
		}
	}

	result, err := crud.GetTaskOutputs(ctx, int64(taskId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Errorf("taskId %d outputs not found", taskId)
			return "", "", [24]float64{}, code.ErrNotFound
		}else {
			logger.Errorf("GetTaskOutputs error: %v", err)
			return "", "", [24]float64{}, code.ErrDatabase
		}
	}

	var arr [24]float64
	for i := 0; i < len(result) && i < 24; i++ {
		arr[i] = result[i]
	}

	return status, date, arr, nil
}