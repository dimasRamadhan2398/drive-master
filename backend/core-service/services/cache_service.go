package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"core-service/repositories"
)

type ICacheService interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

type CacheService struct {
	cacheRepo repositories.ICacheRepository
}

func NewCacheService(cacheRepo repositories.ICacheRepository) ICacheService {
	return &CacheService{
		cacheRepo: cacheRepo,
	}
}

// Get retrieves a value from cache and unmarshals it into dest
func (s *CacheService) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := s.cacheRepo.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

// Set marshals a value and stores it in cache
func (s *CacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	return s.cacheRepo.Set(ctx, key, data, ttl)
}

// Delete removes a key from cache
func (s *CacheService) Delete(ctx context.Context, key string) error {
	return s.cacheRepo.Delete(ctx, key)
}