package database

import (
	"fmt"
	"time"

	"gateway/config"
	"gateway/internal/database/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type Config struct {
	User string
	Pass string
	Host string
	Port string
	Name string
}

// 拼接数据库连接串
func NewDSN(c Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Pass, c.Host, c.Port, c.Name)
}

// 初始化数据库句柄
func InitDB(dsn string) error {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db

	// 自动建表（如果不存在）
	if err := DB.AutoMigrate(&model.Task{}, &model.TaskOutput{}, &model.User{}); err != nil {
		return err
	}
	return nil
}

func NewDBHandler() error {
	dsn := NewDSN(Config{
		User: config.Settings.DB.User,
		Pass: config.Settings.DB.Pass,
		Host: config.Settings.DB.Host,
		Port: config.Settings.DB.Port,
		Name: config.Settings.DB.Name,
	})
	return InitDB(dsn)
}
