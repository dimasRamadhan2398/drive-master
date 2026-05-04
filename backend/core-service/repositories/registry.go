package repositories

import (
	"context"

	"core-service/pkg/base"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ICacheRepository defines cache operations
type ICacheRepository interface {
	SetUserProfile(ctx context.Context, key, value string) error
}

// IEventRepository defines event persistence operations
type IEventRepository interface {
	SaveProcessedEvent(eventType, payload string) error
}

// Registry implements IRepositoryRegistry
type Registry struct {
	db          *gorm.DB
	baseRepo    *base.BaseRepository
	cacheClient *redis.Client
}

func NewRepositoryRegistry(db *gorm.DB) *Registry {
	return &Registry{
		db:       db,
		baseRepo: base.NewBaseRepository(db),
	}
}

// SetCacheClient sets the redis client for cache repository
func (r *Registry) SetCacheClient(client *redis.Client) {
	r.cacheClient = client
}

// IRepositoryRegistry defines the interface for repository registry
type IRepositoryRegistry interface {
	GetRegion() IRegionRepository
	GetCache() ICacheRepository
	GetEvent() IEventRepository
}

// GetRegion returns the region repository
func (r *Registry) GetRegion() IRegionRepository {
	return NewRegionRepository(r.baseRepo)
}

// GetCache returns the cache repository
func (r *Registry) GetCache() ICacheRepository {
	if r.cacheClient == nil {
		return nil
	}
	return NewCacheRepository(r.cacheClient)
}

// GetEvent returns the event repository
func (r *Registry) GetEvent() IEventRepository {
	return NewEventRepository(r.db)
}