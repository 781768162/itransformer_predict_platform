package crud

import (
	"context"

	"gateway/internal/database"
	"gateway/internal/database/model"
)

// CreateTask 插入一条任务记录，TaskID 为自增主键。
func CreateTask(ctx context.Context, t *model.Task) error {
	return database.DB.WithContext(ctx).Create(t).Error
}

// UpdateTaskStatus 根据 task_id 更新任务状态。
func UpdateTaskStatus(ctx context.Context, taskID int64, status string) error {
	return database.DB.WithContext(ctx).
		Model(&model.Task{}).
		Where("task_id = ?", taskID).
		Update("status", status).Error
}

// ClearTaskOutputs 删除指定任务的全部输出。
func ClearTaskOutputs(ctx context.Context, taskID int64) error {
	return database.DB.WithContext(ctx).
		Where("task_id = ?", taskID).
		Delete(&model.TaskOutput{}).Error
}

// CreateTaskOutputs 批量插入 task_output 记录。
func CreateTaskOutputs(ctx context.Context, outputs []model.TaskOutput) error {
	if len(outputs) == 0 {
		return nil
	}
	return database.DB.WithContext(ctx).Create(&outputs).Error
}

// GetTaskStatus 查询任务状态。
func GetTaskStatus(ctx context.Context, taskID int64) (string, error) {
	var task model.Task
	if err := database.DB.WithContext(ctx).
		Select("status").
		Where("task_id = ?", taskID).
		First(&task).Error; err != nil {
		return "", err
	}
	return task.Status, nil
}

// GetTaskOutputs 查询任务对应的 task_output value 列表（按 time_index 升序）。
func GetTaskOutputs(ctx context.Context, taskID int64) ([]float64, error) {
	var outputs []model.TaskOutput
	if err := database.DB.WithContext(ctx).
		Select("value").
		Where("task_id = ?", taskID).
		Order("time_index asc").
		Find(&outputs).Error; err != nil {
		return nil, err
	}

	vals := make([]float64, 0, len(outputs))
	for _, o := range outputs {
		vals = append(vals, o.Value)
	}

	return vals, nil
}
