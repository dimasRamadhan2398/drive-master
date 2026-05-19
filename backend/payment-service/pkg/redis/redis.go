package redis

import (
	"context"
	"fmt"

	"payment-service/pkg/config"

	"github.com/redis/go-redis/v9"
)

type Client = redis.Client

func NewRedisConnection(cfg *config.RedisConfig) (*Client, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}