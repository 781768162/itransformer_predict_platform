package model

import "time"

type User struct {
	UserID    int64     `gorm:"column:user_id;primaryKey;autoIncrement" json:"user_id"`
	UserName  string    `gorm:"column:user_name;size:64;unique;not null" json:"user_name"`
	Password  string    `gorm:"column:password;size:255;not null" json:"-"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string { return "user" }
