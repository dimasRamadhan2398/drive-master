package routes

import (
	"user-service/controllers"
	"user-service/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type Registry struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	authMiddleware middlewares.IAuthMiddleware
}

type IRouteRegister interface {
	Serve()
}

func NewRouteRegistry(controller controllers.IControllerRegistry, group *gin.RouterGroup, authMiddleware middlewares.IAuthMiddleware) IRouteRegister {
	return &Registry{controller: controller, group: group, authMiddleware: authMiddleware}
}

func (r *Registry) Serve() {
	r.GetUserRoute().Run()
	r.GetAuthRoute().Run()
	r.GetMemberRoute().Run()
	r.GetInstructorRoute().Run()
	r.GetCertificationRoute().Run()
	r.GetEntitlementRoute().Run()
	r.GetWorkExperienceRoute().Run()
	r.GetCoverageAreaRoute().Run()
	r.GetTestimonialRoute().Run()
	r.GetDashboardRoute().Run()
}

func (r *Registry) GetAuthRoute() IAuthRoute {
	return NewAuthRoute(r.controller, r.group, r.authMiddleware)
}

func (r *Registry) GetUserRoute() IUserRoute {
	return NewUserRoute(r.controller, r.group, r.authMiddleware)
}

func (r *Registry) GetMemberRoute() IMemberRoute {
	return NewMemberRoute(r.controller, r.group, r.authMiddleware)
}

func (r *Registry) GetInstructorRoute() IInstructorRoute {
	return NewInstructorRoute(r.controller, r.group, r.authMiddleware)
}

func (r *Registry) GetWorkExperienceRoute() IWorkExperienceRoute {
	return NewWorkExperienceRoute(r.controller, r.group, r.authMiddleware)
}

func (r *Registry) GetCoverageAreaRoute() ICoverageAreaRoute {
	return NewCoverageAreaRoute(r.controller, r.group, r.authMiddleware)
}

func (r *Registry) GetCertificationRoute() ICertificationRoute {
	return NewCertificationRoute(r.controller, r.group, r.authMiddleware)
}

func (r *Registry) GetEntitlementRoute() IEntitlementRoute {
	return NewEntitlementRoute(r.controller, r.group, r.authMiddleware)
}

func (r *Registry) GetTestimonialRoute() ITestimonialRoute {
	return NewTestimonialRoute(r.controller, r.group)
}

func (r *Registry) GetDashboardRoute() IDashboardRoute {
	return NewDashboardRoute(r.controller, r.group, r.authMiddleware)
}
