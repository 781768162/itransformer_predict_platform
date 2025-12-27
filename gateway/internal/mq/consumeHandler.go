package mq

import (
	"context"
	"encoding/json"

	"gateway/internal/code"
	"gateway/internal/database/crud"
	"gateway/internal/database/model"

	"github.com/segmentio/kafka-go"
)

// 假设 payload 形如：{"task_id":123,"status":"success","result":[...24 floats...]}
type TaskResultPayload struct {
	TaskID int64    `json:"task_id"`
	Status string    `json:"status"`
	Result []float64 `json:"result"`
}

// HandleTaskResult 解析消息并写入数据库：更新 task 状态，重写 task_output。
func HandleTaskResult(ctx context.Context, msg kafka.Message) error {
	var payload TaskResultPayload
	if err := json.Unmarshal(msg.Value, &payload); err != nil {
		return code.ErrInvalidParam
	}

	if payload.TaskID == 0 {
		return code.ErrInvalidParam
	}
	if len(payload.Result) != 24 {
		return code.ErrInvalidParam
	}

	err := crud.UpdateTaskStatus(ctx, payload.TaskID, payload.Status)
	if err != nil {
		return code.ErrDatabase
	}

	err = crud.ClearTaskOutputs(ctx, payload.TaskID)
	if err != nil {
		return code.ErrDatabase
	}

	outputs := make([]model.TaskOutput, 0, len(payload.Result))
	for idx, v := range payload.Result {
		outputs = append(outputs, model.TaskOutput{
			TaskID:    payload.TaskID,
			TimeIndex: uint16(idx),
			Value:     v,
		})
	}

	err = crud.CreateTaskOutputs(ctx, outputs)
	if err != nil {
		return code.ErrDatabase
	}

	return nil
}
