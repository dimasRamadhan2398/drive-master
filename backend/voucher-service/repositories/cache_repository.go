package repositories

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{client: client}
}

func (r *CacheRepository) SetUserProfile(ctx context.Context, key, value string) error {
	return r.client.Set(ctx, key, value, 24*time.Hour).Err()
}
