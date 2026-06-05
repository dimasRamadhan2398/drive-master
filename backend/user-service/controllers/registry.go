package controllers

import "user-service/services"

type ControllerRegistry struct {
	service services.IServiceRegistry
}

// GetAuthController implements [IControllerRegistry].
func (r *ControllerRegistry) GetAuthController() IAuthController {
	return NewAuthController(
		r.service.GetUserService(),
		r.service.GetAuthService(),
		r.service.GetEmailService(),
		r.service.GetRoleService(),
	)
}

// GetInstructorController implements [IControllerRegistry].
func (r *ControllerRegistry) GetInstructorController() IInstructorController {
	return NewInstructorController(
		r.service.GetUserService(),
		r.service.GetAuthService(),
		r.service.GetMemberService(),
		r.service.GetInstructorService(),
		r.service.GetRoleService(),
		r.service.GetEmailService(),
		r.service.GetMediaService(),
	)
}

// GetMemberController implements [IControllerRegistry].
func (r *ControllerRegistry) GetMemberController() IMemberController {
	return NewMemberController(
		r.service.GetUserService(),
		r.service.GetAuthService(),
		r.service.GetMemberService(),
		r.service.GetRoleService(),
		r.service.GetEmailService(),
		r.service.GetMediaService(),
	)
}

// GetRoleController implements [IControllerRegistry].
func (r *ControllerRegistry) GetRoleController() IRoleController {
	return NewRoleController(r.service.GetRoleService())
}

// GetRoleController implements [IControllerRegistry].
func (r *ControllerRegistry) GetWorkExperienceController() IWorkExperienceController {
	return NewWorkExperienceController(r.service.GetWorkExperienceService())
}

// GetCoverageAreaController implements [IControllerRegistry].
func (r *ControllerRegistry) GetCoverageAreaController() ICoverageAreaController {
	return NewCoverageAreaController(r.service.GetCoverageAreaService())
}

// GetCertificationController implements [IControllerRegistry].
func (r *ControllerRegistry) GetCertificationController() ICertificationController {
	return NewCertificationController(
		r.service.GetCertificationService(),
	)
}

// GetEntitlementController implements [IControllerRegistry].
func (r *ControllerRegistry) GetEntitlementController() IEntitlementController {
	return NewEntitlementController(
		r.service.GetEntitlementService(),
	)
}

// GetTestimonialController implements [IControllerRegistry].
func (r *ControllerRegistry) GetTestimonialController() ITestimonialController {
	return NewTestimonialController(r.service.GetTestimonialService())
}

// IControllerRegistry defines methods for getting controllers
type IControllerRegistry interface {
	GetUserController() IUserController
	GetAuthController() IAuthController
	GetRoleController() IRoleController
	GetInstructorController() IInstructorController
	GetMemberController() IMemberController
	GetWorkExperienceController() IWorkExperienceController
	GetCoverageAreaController() ICoverageAreaController
	GetCertificationController() ICertificationController
	GetEntitlementController() IEntitlementController
	GetTestimonialController() ITestimonialController
}

// NewControllerRegistry creates a new controller registry
func NewControllerRegistry(service services.IServiceRegistry) IControllerRegistry {
	return &ControllerRegistry{service: service}
}

// GetUserController returns the user controller
func (r *ControllerRegistry) GetUserController() IUserController {
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
