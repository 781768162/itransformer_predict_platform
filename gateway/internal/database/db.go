package database

import (
	"fmt"
	"os"
	"time"

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
	return nil
}

func NewDBHandler() error {
	cfg := Config{
		User: getenv("DB_USER", "root"),
		Pass: getenv("DB_PASS", ""),
		Host: getenv("DB_HOST", "127.0.0.1"),
		Port: getenv("DB_PORT", "3306"),
		Name: getenv("DB_NAME", "fintech"),
	}
	dsn := NewDSN(cfg)
	return InitDB(dsn)
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
