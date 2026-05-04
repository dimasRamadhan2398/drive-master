package controllers

import (
	"booking-service/services"
)

type Registry struct {
	service services.IServiceRegistry
}

func (r *Registry) GetBookingController() IBookingController {
	return NewBookingController(
		r.service.GetBookingService(),
	)
}

func (r *Registry) GetSessionController() ISessionController {
	return NewSessionController(
		r.service.GetSessionService(),
	)
}

func (r *Registry) GetEntitlementController() IEntitlementController {
	return NewEntitlementController(
		r.service.GetEntitlementService(),
	)
}

func (r *Registry) GetCertificationController() ICertificationController {
	return NewCertificationController(
		r.service.GetCertificationService(),
	)
}

// IControllerRegistry defines methods for getting controllers
type IControllerRegistry interface {
	GetBookingController() IBookingController
	GetSessionController() ISessionController
	GetEntitlementController() IEntitlementController
	GetCertificationController() ICertificationController
}

// NewControllerRegistry creates a new controller registry
func NewControllerRegistry(service services.IServiceRegistry) IControllerRegistry {
	return &Registry{service: service}
}