package controllers

import "user-service/services"

type ControllerRegistry struct {
	service services.IServiceRegistry
}

// IControllerRegistry defines methods for getting controllers
type IControllerRegistry interface {
	GetUserController() *UserController
}

// NewControllerRegistry creates a new controller registry
func NewControllerRegistry(service services.IServiceRegistry) IControllerRegistry {
	return &ControllerRegistry{service: service}
}

// GetUserController returns the user controller
func (r *ControllerRegistry) GetUserController() *UserController {
	return NewUserController(
		r.service.GetUserService(),
		r.service.GetAuthService(),
		r.service.GetMemberService(),
		r.service.GetInstructorService(),
		r.service.GetRoleService(),
		r.service.GetEmailService(),
		r.service.GetMediaService(),
	)
}