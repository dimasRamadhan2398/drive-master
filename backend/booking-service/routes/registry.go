package routes

import (
	"booking-service/controllers"
	"booking-service/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type Registry struct {
	controller     controllers.IControllerRegistry
	group         *gin.RouterGroup
	authMiddleware middlewares.IAuthMiddleware
}

type IRouteRegister interface {
	Serve()
}

func NewRouteRegistry(controller controllers.IControllerRegistry, group *gin.RouterGroup, authMiddleware middlewares.IAuthMiddleware) IRouteRegister {
	return &Registry{controller: controller, group: group, authMiddleware: authMiddleware}
}

func (r *Registry) Serve() {
	r.GetSessionRoute().Run()
	r.GetEnrollmentRoute().Run()
	r.GetScheduleRoute().Run()
	r.GetDashboardRoute().Run()
	r.GetPaymentRoute().Run()
}

func (r *Registry) GetSessionRoute() ISessionRoute {
	return NewSessionRoute(r.controller, r.group, r.authMiddleware)
}

func (r *Registry) GetEnrollmentRoute() IEnrollmentRoute {
	return NewEnrollmentRoute(r.controller, r.group, r.authMiddleware)
}

func (r *Registry) GetScheduleRoute() IScheduleRoute {
	return NewScheduleRoute(r.controller, r.group, r.authMiddleware)
}

func (r *Registry) GetDashboardRoute() IDashboardRoute {
	return NewDashboardRoute(r.controller, r.group, r.authMiddleware)
}

func (r *Registry) GetPaymentRoute() IPaymentRoute {
	return NewPaymentRoute(r.controller, r.group, r.authMiddleware)
}