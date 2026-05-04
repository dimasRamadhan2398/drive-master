package routes

import (
	"user-service/controllers"

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
	r.GetAuthRoute().Run()
	r.GetInstructorRoute().Run()
	r.GetWorkExperienceRoute().Run()
	r.GetCoverageAreaRoute().Run()
}

func (r *Registry) GetAuthRoute() IAuthRoute {
	return NewAuthRoute(r.controller, r.group)
}

func (r *Registry) GetUserRoute() IUserRoute {
	return NewUserRoute(r.controller, r.group)
}

func (r *Registry) GetInstructorRoute() IInstructorRoute {
	return NewInstructorRoute(r.controller, r.group)
}

func (r *Registry) GetWorkExperienceRoute() IWorkExperienceRoute {
	return NewWorkExperienceRoute(r.controller, r.group)
}

func (r *Registry) GetCoverageAreaRoute() ICoverageAreaRoute {
	return NewCoverageAreaRoute(r.controller, r.group)
}
