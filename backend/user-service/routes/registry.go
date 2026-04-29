package routes

import (
	"user-service/controllers"
	routes "user-service/routes/user"

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
	return &Registry{controller: controller, group: group}
}

func (r *Registry) Serve() {
	r.GetUserRoute().Run()
}

func (r *Registry) GetAuthRoute() routes.IUserRoute {
	return routes.authRoute(r.controller, r.group)
}

func (r *Registry) GetUserRoute() routes.IUserRoute {
	return routes.NewUserRoute(r.controller, r.group)
}
