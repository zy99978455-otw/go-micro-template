package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zy99978455-otw/go-micro-template/pkg/config"
)

// NewRedisClient 初始化 Redis 连接
func NewRedisClient(cfg *config.AppConfig) (*redis.Client, func(), error) {
	
	// 1. 直接使用传入的配置
	c := cfg.Redis
	if c.Host == "" {
		// 空配置，返回 nil Client
		return nil, func() {}, nil
	}

	// 2. 初始化连接
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Host, c.Port),
		Password: c.Password,
		DB:       c.DB, 
	})

	// 3. 测试连接 (Ping)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		// 连接失败，返回错误
		return nil, nil, fmt.Errorf("Redis connection failed: %w", err)
	}

	// 4. 定义清理函数
	cleanup := func() {
		fmt.Println(">>> Closing Redis connection...")
		if rdb != nil {
			rdb.Close()
		}
	}

	fmt.Printf(">>> [Web2] Redis Connected! Addr: %s:%d\n", c.Host, c.Port)

	// 5. 返回实例
	return rdb, cleanup, nil
}
