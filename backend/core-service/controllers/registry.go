package controllers

import (
	"core-service/repositories"
	"core-service/services"
)

type Registry struct {
	svcRegistry services.IServiceRegistry
}

type IControllerRegistry interface {
	GetRegionController() IRegionController
}

// NewControllerRegistry creates a new controller registry
func NewControllerRegistry(svcRegistry services.IServiceRegistry) IControllerRegistry {
	return &Registry{
		svcRegistry: svcRegistry,
	}
}

// GetRegionController returns the region controller
func (r *Registry) GetRegionController() IRegionController {
	return NewRegionController(r.svcRegistry.GetRegionService())
}

// GetRepositoryRegistry returns the repository registry (for dependency injection)
func NewRepositoryRegistry(db interface{}) repositories.IRepositoryRegistry {
	// Placeholder - will be set from main.go
	return nil
}