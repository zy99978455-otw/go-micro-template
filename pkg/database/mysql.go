package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/zy99978455-otw/go-micro-template/pkg/config"
)

// NewMySQLClient 初始化 MySQL 连接
// 参数: cfg *config.Config (直接传入配置，不再读 global)
// 返回: *gorm.DB (实例), func() (清理函数), error
func NewMySQLClient(cfg *config.AppConfig) (*gorm.DB, func(), error) {
	
	// 1. 检查配置
	c := cfg.Mysql
	if c.Host == "" {
		// 允许空配置，返回 nil DB，不报错
		return nil, func() {}, nil
	}

	// 2. 组装 DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Name)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// 3. 尝试连接
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("MySQL connection failed: %w", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(c.MaxIdle)
	sqlDB.SetMaxOpenConns(c.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 4. 定义清理函数 (Cleanup)
	// 这个函数会在 main.go 退出时被调用
	cleanup := func() {
		fmt.Println(">>> Closing MySQL connection...")
		if sqlDB != nil {
			sqlDB.Close()
		}
	}
	
	fmt.Printf(">>> [Web2] MySQL Connected! Host: %s, DB: %s\n", c.Host, c.Name)

	return db, cleanup, nil
}
