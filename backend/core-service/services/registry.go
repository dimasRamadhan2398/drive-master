package services

import (
	"core-service/pkg/kafka"
	"core-service/repositories"
)

type Registry struct {
	repoRegistry   repositories.IRepositoryRegistry
	eventPublisher *kafka.EventPublisher
	analyticsSvc   IAnalyticsService
	cacheSvc       ICacheService
}

// GetEventService implements [IServiceRegistry].
func (r *Registry) GetEventService() IEventService {
	return NewEventService(r.repoRegistry.GetEvent(), r.repoRegistry.GetCache())
}

// GetRegionService implements [IServiceRegistry].
func (r *Registry) GetRegionService() IRegionService {
	return NewRegionService(r.repoRegistry.GetRegion(), r.cacheSvc)
}

// GetCacheService implements [IServiceRegistry].
func (r *Registry) GetCacheService() ICacheService {
	return r.cacheSvc
}

// GetCarService implements [IServiceRegistry].
func (r *Registry) GetCarService() ICarService {
	return NewCarService(r.repoRegistry.GetCar(), r.eventPublisher)
}

// GetPackageService implements [IServiceRegistry].
func (r *Registry) GetPackageService() IPackageService {
	return NewPackageService(r.repoRegistry.GetPackage(), r.eventPublisher)
}

// GetArticleService implements [IServiceRegistry].
func (r *Registry) GetArticleService() IArticleService {
	return NewArticleService(r.repoRegistry.GetArticle(), r.repoRegistry.GetFAQ(), r.eventPublisher)
}

// GetAnalyticsService implements [IServiceRegistry].
func (r *Registry) GetAnalyticsService() IAnalyticsService {
	return NewAnalyticsService()
}

// GetSalesService implements [IServiceRegistry].
func (r *Registry) GetSalesService() ISalesService {
	return NewSalesService(r.repoRegistry.GetSales())
}

// GetGeneralSettingsService implements [IServiceRegistry].
func (r *Registry) GetGeneralSettingsService() IGeneralSettingsService {
	return NewGeneralSettingsService(r.repoRegistry.GetGeneralSettings())
}

// GetFAQService implements [IServiceRegistry].
func (r *Registry) GetFAQService() IFAQService {
	return NewFAQService(r.repoRegistry.GetFAQ())
}

type IServiceRegistry interface {
	GetEventService() IEventService
	GetRegionService() IRegionService
	GetCarService() ICarService
	GetPackageService() IPackageService
	GetArticleService() IArticleService
	GetAnalyticsService() IAnalyticsService
	GetCacheService() ICacheService
	GetSalesService() ISalesService
	GetGeneralSettingsService() IGeneralSettingsService
	GetFAQService() IFAQService
}

func NewServiceRegistry(repoRegistry repositories.IRepositoryRegistry, eventPublisher *kafka.EventPublisher) IServiceRegistry {
	return &Registry{
		repoRegistry:   repoRegistry,
		eventPublisher: eventPublisher,
		analyticsSvc:   NewAnalyticsService(),
		cacheSvc:       NewCacheService(repoRegistry.GetCache()),
	}
}
