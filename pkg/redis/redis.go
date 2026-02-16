package redis

import (
	"context"
	"fmt"
	"time"

	"beeblog/pkg/config"
	"beeblog/pkg/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var RedisClient *redis.Client

// NewRedis 构造函数
func NewRedis(cfg config.RedisConfig, logger *zap.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect redis: %w", err)
	}

	logger.Info("Redis connection established", zap.String("addr", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)))
	return client, nil
}

// WireSet 声明本包的 ProviderSet
var WireSet = NewRedis

func Close() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}
