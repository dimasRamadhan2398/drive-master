package routes

import (
	"core-service/controllers"

	"github.com/gin-gonic/gin"
)

type Registry struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type IRouteRegister interface {
	Serve()
}

func NewRouteRegistry(controller controllers.IControllerRegistry, group *gin.RouterGroup) IRouteRegister {
	return &Registry{
		controller: controller,
		group:      group,
	}
}

func (r *Registry) Serve() {
	r.GetRegionRoute().Run()
}

func (r *Registry) GetRegionRoute() IRegionRoute {
	return NewRegionRoute(r.controller, r.group)
}