package model

import "time"

type Task struct {
	TaskID        int64    `gorm:"column:task_id;primaryKey;autoIncrement" json:"task_id"`
	UserID        int64    `gorm:"column:user_id;not null" json:"user_id"`
	Status        string    `gorm:"column:status;type:enum('pending','success','failed');not null;default:'pending'" json:"status"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (Task) TableName() string { return "task" }

type TaskOutput struct {
	ID        int64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TaskID    int64    `gorm:"column:task_id;not null" json:"task_id"`
	TimeIndex uint16    `gorm:"column:time_index;not null" json:"time_index"`
	TS        time.Time `gorm:"column:ts;not null" json:"ts"`
	Value     float64   `gorm:"column:value;not null" json:"value"`
}

func (TaskOutput) TableName() string { return "task_output" }
