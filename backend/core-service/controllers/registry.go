package controllers

import (
	"core-service/pkg/clients"
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
	GetAddOnController() IAddOnController
	GetArticleController() IArticleController
	GetAnalyticsController() IAnalyticsController
	GetSalesController() ISalesController
	GetGeneralSettingsController() IGeneralSettingsController
	GetFAQController() IFAQController
	GetTransactionController() ITransactionController
}

// NewControllerRegistry creates a new controller registry
func NewControllerRegistry(svcRegistry services.IServiceRegistry) IControllerRegistry {
	return &Registry{
		svcRegistry: svcRegistry,
	}
}

// NewControllerRegistryWithPayment creates a new controller registry with payment client
func NewControllerRegistryWithPayment(svcRegistry services.IServiceRegistry, paymentClient clients.ITransactionClient) IControllerRegistry {
	return &RegistryWithPayment{
		Registry: &Registry{
			svcRegistry: svcRegistry,
		},
		paymentClient: paymentClient,
	}
}

// RegistryWithPayment extends Registry with payment client
type RegistryWithPayment struct {
	*Registry
	paymentClient clients.ITransactionClient
}

// GetRegionController returns the region controller
func (r *Registry) GetRegionController() IRegionController {
	return NewRegionController(r.svcRegistry.GetRegionService())
}

// GetCarController returns the car controller
func (r *Registry) GetCarController() ICarController {
	return NewCarController(r.svcRegistry.GetCarService(), r.svcRegistry.GetMediaService())
}

// GetPackageController returns the package controller
func (r *Registry) GetPackageController() IPackageController {
	return NewPackageController(r.svcRegistry.GetPackageService())
}

// GetAddOnController returns the add-on controller
func (r *Registry) GetAddOnController() IAddOnController {
	return NewAddOnController(r.svcRegistry.GetAddOnService())
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

// GetTransactionController returns the transaction controller (requires payment client)
func (r *Registry) GetTransactionController() ITransactionController {
	return nil // Will be overridden by RegistryWithPayment
}

// GetTransactionController returns the transaction controller for RegistryWithPayment
func (r *RegistryWithPayment) GetTransactionController() ITransactionController {
	return NewTransactionController(r.paymentClient)
}

// GetRepositoryRegistry returns the repository registry (for dependency injection)
func NewRepositoryRegistry(db interface{}) repositories.IRepositoryRegistry {
	// Placeholder - will be set from main.go
	return nil
}