package repositories

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type ICacheRepository interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type CacheRepository struct {
	client *redis.Client
}

// Delete implements [ICacheRepository].
func (c *CacheRepository) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Get implements [ICacheRepository].
func (c *CacheRepository) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Set implements [ICacheRepository].
func (c *CacheRepository) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, 24*time.Hour).Err()
}

func NewCacheRepository(client *redis.Client) ICacheRepository {
	return &CacheRepository{client: client}
}
