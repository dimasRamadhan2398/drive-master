package services

import (
	"core-service/repositories"
)

type Registry struct {
	repoRegistry repositories.IRepositoryRegistry
}

// GetEventService implements [IServiceRegistry].
func (r *Registry) GetEventService() IEventService {
	return NewEventService(r.repoRegistry.GetEvent(), r.repoRegistry.GetCache());
}

// GetRegionService implements [IServiceRegistry].
func (r *Registry) GetRegionService() IRegionService {
	return NewRegionService(r.repoRegistry.GetRegion())
}

type IServiceRegistry interface {
	GetEventService() IEventService
	GetRegionService() IRegionService
}

func NewServiceRegistry(repoRegistry repositories.IRepositoryRegistry) IServiceRegistry {
	return &Registry{repoRegistry: repoRegistry}
}
