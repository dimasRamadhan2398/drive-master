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
	GetCarController() ICarController
	GetPackageController() IPackageController
	GetTestimonialController() ITestimonialController
	GetArticleController() IArticleController
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

// GetCarController returns the car controller
func (r *Registry) GetCarController() ICarController {
	return NewCarController(r.svcRegistry.GetCarService())
}

// GetPackageController returns the package controller
func (r *Registry) GetPackageController() IPackageController {
	return NewPackageController(r.svcRegistry.GetPackageService())
}

// GetTestimonialController returns the testimonial controller
func (r *Registry) GetTestimonialController() ITestimonialController {
	return NewTestimonialController(r.svcRegistry.GetTestimonialService())
}

// GetArticleController returns the article controller
func (r *Registry) GetArticleController() IArticleController {
	return NewArticleController(r.svcRegistry.GetArticleService())
}

// GetRepositoryRegistry returns the repository registry (for dependency injection)
func NewRepositoryRegistry(db interface{}) repositories.IRepositoryRegistry {
	// Placeholder - will be set from main.go
	return nil
}