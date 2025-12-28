package mq

import (
	"context"
	"encoding/json"
	"time"

	"gateway/internal/code"
	"gateway/internal/database/crud"
	"gateway/internal/database/model"
	"gateway/pkg/logger"

	"github.com/segmentio/kafka-go"
)

// payload 格式：{"task_id":123,"status":"success","result":[...24 floats...]}
type TaskResultPayload struct {
	TaskID int64    `json:"task_id"`
	Status string    `json:"status"`
	Result []float64 `json:"result"`
}

// HandleTaskResult 解析消息并写入数据库：更新 task 状态，重写 task_output。
func HandleTaskResult(ctx context.Context, msg kafka.Message) error {
	var payload TaskResultPayload
	if err := json.Unmarshal(msg.Value, &payload); err != nil {
		logger.Errorf("Unmarshal error: %v", err)
		return code.ErrInvalidParam
	}

	if payload.TaskID == 0 {
		logger.Errorf("TaskID is zero")
		return code.ErrInvalidParam
	}
	if len(payload.Result) != 24 {
		logger.Errorf("len of Result is Invalid")
		return code.ErrInvalidParam
	}

	err := crud.UpdateTaskStatus(ctx, payload.TaskID, payload.Status) // 更新Task状态
	if err != nil {
		logger.Errorf("UpdateTaskStatus TaskID: %d Status: %s error: %v", payload.TaskID, payload.Status, err)
		return code.ErrDatabase
	}

	err = crud.ClearTaskOutputs(ctx, payload.TaskID) // 清空对应TaskOutputs
	if err != nil {
		logger.Errorf("ClearTaskOutputs TaskID: %d error: %v", payload.TaskID, err)
		return code.ErrDatabase
	}

	outputs := make([]model.TaskOutput, 0, len(payload.Result))
	for idx, v := range payload.Result {
		outputs = append(outputs, model.TaskOutput{
			TaskID:    payload.TaskID,
			TimeIndex: uint16(idx),
			TS: time.Now(),	
			Value:     v,
		})
	}

	err = crud.CreateTaskOutputs(ctx, outputs) // 写入对应TaskOutputs
	if err != nil {
		logger.Errorf("ClearTaskOutputs TaskID: %d error: %v", payload.TaskID, err)
		return code.ErrDatabase
	}

	return nil
}
