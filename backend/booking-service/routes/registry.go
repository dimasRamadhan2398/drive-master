package routes

import (
	"booking-service/controllers"

	"github.com/gin-gonic/gin"
)

type Registry struct {
	controller controllers.IControllerRegistry
	engine     *gin.Engine
}

type IRouteRegister interface {
	Serve()
	GetBookingRoute() IBookingRoute
	GetSessionRoute() ISessionRoute
	GetEntitlementRoute() IEntitlementRoute
	GetCertificationRoute() ICertificationRoute
}

func NewRouteRegistry(controller controllers.IControllerRegistry, engine *gin.Engine) IRouteRegister {
	return &Registry{controller: controller, engine: engine}
}

func (r *Registry) Serve() {
	group := r.engine.Group("/api/v1")
	r.GetBookingRoute().Run(group)
	r.GetSessionRoute().Run(group)
	r.GetEntitlementRoute().Run(group)
	r.GetCertificationRoute().Run(group)
}

func (r *Registry) GetBookingRoute() IBookingRoute {
	return &BookingRoute{controller: r.controller}
}

func (r *Registry) GetSessionRoute() ISessionRoute {
	return &SessionRoute{controller: r.controller}
}

func (r *Registry) GetEntitlementRoute() IEntitlementRoute {
	return &EntitlementRoute{controller: r.controller}
}

func (r *Registry) GetCertificationRoute() ICertificationRoute {
	return &CertificationRoute{controller: r.controller}
}