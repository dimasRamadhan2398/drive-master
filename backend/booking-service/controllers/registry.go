package controllers

import (
	"booking-service/services"
)

type Registry struct {
	service services.IServiceRegistry
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

func (r *Registry) GetEnrollmentController() IEnrollmentController {
	return NewEnrollmentController(
		r.service.GetEnrollmentService(),
	)
}

func (r *Registry) GetScheduleController() IScheduleController {
	return NewScheduleController(
		r.service.GetScheduleService(),
	)
}

func (r *Registry) GetDashboardController() IDashboardController {
	return NewDashboardController(
		r.service.GetSessionService(),
		r.service.GetRevenueService(),
	)
}

func (r *Registry) GetPaymentController() IPaymentController {
	return NewPaymentController(
		r.service.GetPaymentService(),
	)
}

// IControllerRegistry defines methods for getting controllers
type IControllerRegistry interface {
	GetSessionController() ISessionController
	GetEntitlementController() IEntitlementController
	GetEnrollmentController() IEnrollmentController
	GetScheduleController() IScheduleController
	GetDashboardController() IDashboardController
	GetPaymentController() IPaymentController
}

// NewControllerRegistry creates a new controller registry
func NewControllerRegistry(service services.IServiceRegistry) IControllerRegistry {
	return &Registry{service: service}
}