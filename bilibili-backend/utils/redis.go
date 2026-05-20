package utils

import (
	"context"
	"fmt"
	"log"

	"bilibili-backend/config"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// InitRedis 初始化 Redis 连接
func InitRedis() error {
	cfg := config.C.Redis
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}

	log.Println("[Redis] Connected successfully")
	return nil
}
