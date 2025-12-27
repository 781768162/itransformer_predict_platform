package logic

import (
	"context"
	"encoding/json"

	"gateway/internal/code"
	"gateway/internal/database/crud"
	"gateway/internal/database/model"
	"gateway/internal/mq"

	"github.com/segmentio/kafka-go"
)

var (
	cfg mq.Config
	prod *kafka.Writer = mq.NewProducer(cfg)
)

type taskPayload struct {
    TaskID     int64              `json:"task_id"`
    PassData   [13][72]float64  `json:"pass_data"`
    FutureData [12][24]float64  `json:"future_data"`
}

func CreateTaskLogic(ctx context.Context, userId int64, passData [13][72]float64, futureData [12][24]float64) (int, error) {
	t := &model.Task{
		UserID: userId,
		Status: "pending",
	}
	err := crud.CreateTask(ctx, t) // 插入记录
	if err != nil {
		return 0, code.ErrDatabase
	}

	body, err := json.Marshal(taskPayload{
        TaskID:     t.TaskID,
        PassData:   passData,
        FutureData: futureData,
    })
    if err != nil {
        return 0, code.ErrJsonMarshal
    }

	err = prod.WriteMessages(ctx, kafka.Message{ // kafka投递消息
		Value: body,
	})
	if err != nil {
		return 0, code.ErrMessageQueue
	}

	return int(t.TaskID), nil
}

func GetTaskLogic(ctx context.Context, taskId int) (string, [24]float64, error) {
	status, err := crud.GetTaskStatus(ctx, int64(taskId))
	if err != nil {
		return "", [24]float64{}, code.ErrDatabase
	}

	result, err := crud.GetTaskOutputs(ctx, int64(taskId))
	if err != nil {
		return "", [24]float64{}, code.ErrDatabase
	}

	return status, [24]float64(result), nil
}