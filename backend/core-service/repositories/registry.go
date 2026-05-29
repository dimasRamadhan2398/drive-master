package repositories

import (
	"core-service/pkg/base"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

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
	GetCar() ICarRepository
	GetPackage() IPackageRepository
	GetTestimonial() ITestimonialRepository
	GetArticle() IArticleRepository
}

// GetCache returns the cache repository
func (r *Registry) GetCache() ICacheRepository {
	if r.cacheClient == nil {
		return nil
	}
	return NewCacheRepository(r.cacheClient)
}

// GetRegion returns the region repository
func (r *Registry) GetRegion() IRegionRepository {
	return NewRegionRepository(r.baseRepo, r.GetCache())
}

// GetEvent returns the event repository
func (r *Registry) GetEvent() IEventRepository {
	return NewEventRepository(r.db)
}

// GetCar returns the car repository
func (r *Registry) GetCar() ICarRepository {
	return NewCarRepository(r.baseRepo)
}

// GetPackage returns the package repository
func (r *Registry) GetPackage() IPackageRepository {
	return NewPackageRepository(r.baseRepo)
}

// GetTestimonial returns the testimonial repository
func (r *Registry) GetTestimonial() ITestimonialRepository {
	return NewTestimonialRepository(r.baseRepo)
}

// GetArticle returns the article repository
func (r *Registry) GetArticle() IArticleRepository {
	return NewArticleRepository(r.baseRepo)
}