package services

import (
	"core-service/pkg/kafka"
	"core-service/repositories"
)

type Registry struct {
	repoRegistry   repositories.IRepositoryRegistry
	eventPublisher *kafka.EventPublisher
}

// GetEventService implements [IServiceRegistry].
func (r *Registry) GetEventService() IEventService {
	return NewEventService(r.repoRegistry.GetEvent(), r.repoRegistry.GetCache())
}

// GetRegionService implements [IServiceRegistry].
func (r *Registry) GetRegionService() IRegionService {
	return NewRegionService(r.repoRegistry.GetRegion())
}

// GetCarService implements [IServiceRegistry].
func (r *Registry) GetCarService() ICarService {
	return NewCarService(r.repoRegistry.GetCar(), r.eventPublisher)
}

// GetPackageService implements [IServiceRegistry].
func (r *Registry) GetPackageService() IPackageService {
	return NewPackageService(r.repoRegistry.GetPackage(), r.eventPublisher)
}

// GetTestimonialService implements [IServiceRegistry].
func (r *Registry) GetTestimonialService() ITestimonialService {
	return NewTestimonialService(r.repoRegistry.GetTestimonial())
}

// GetArticleService implements [IServiceRegistry].
func (r *Registry) GetArticleService() IArticleService {
	return NewArticleService(r.repoRegistry.GetArticle())
}

type IServiceRegistry interface {
	GetEventService() IEventService
	GetRegionService() IRegionService
	GetCarService() ICarService
	GetPackageService() IPackageService
	GetTestimonialService() ITestimonialService
	GetArticleService() IArticleService
}

func NewServiceRegistry(repoRegistry repositories.IRepositoryRegistry, eventPublisher *kafka.EventPublisher) IServiceRegistry {
	return &Registry{
		repoRegistry:   repoRegistry,
		eventPublisher: eventPublisher,
	}
}
