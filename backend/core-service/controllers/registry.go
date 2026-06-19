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
	GetArticleController() IArticleController
	GetAnalyticsController() IAnalyticsController
	GetSalesController() ISalesController
	GetGeneralSettingsController() IGeneralSettingsController
	GetFAQController() IFAQController
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

// GetArticleController returns the article controller
func (r *Registry) GetArticleController() IArticleController {
	return NewArticleController(r.svcRegistry.GetArticleService())
}

// GetAnalyticsController returns the analytics controller
func (r *Registry) GetAnalyticsController() IAnalyticsController {
	return NewAnalyticsController(r.svcRegistry.GetAnalyticsService())
}

// GetSalesController returns the sales controller
func (r *Registry) GetSalesController() ISalesController {
	return NewSalesController(r.svcRegistry.GetSalesService())
}

// GetGeneralSettingsController returns the general settings controller
func (r *Registry) GetGeneralSettingsController() IGeneralSettingsController {
	return NewGeneralSettingsController(r.svcRegistry.GetGeneralSettingsService())
}

// GetFAQController returns the FAQ controller
func (r *Registry) GetFAQController() IFAQController {
	return NewFAQController(r.svcRegistry.GetFAQService())
}

// GetRepositoryRegistry returns the repository registry (for dependency injection)
func NewRepositoryRegistry(db interface{}) repositories.IRepositoryRegistry {
	// Placeholder - will be set from main.go
	return nil
}