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
	r.GetCarRoute().Run()
	r.GetPackageRoute().Run()
	r.GetAddOnRoute().Run()
	r.GetArticleRoute().Run()
	r.GetAnalyticsRoute().Run()
	r.GetSalesRoute().Run()
	r.GetGeneralSettingsRoute().Run()
}

func (r *Registry) GetRegionRoute() IRegionRoute {
	return NewRegionRoute(r.controller, r.group)
}

func (r *Registry) GetCarRoute() ICarRoute {
	return NewCarRoute(r.controller, r.group)
}

func (r *Registry) GetPackageRoute() IPackageRoute {
	return NewPackageRoute(r.controller, r.group)
}

func (r *Registry) GetAddOnRoute() IAddOnRoute {
	return NewAddOnRoute(r.controller, r.group)
}

func (r *Registry) GetArticleRoute() IArticleRoute {
	return NewArticleRoute(r.controller, r.group)
}

func (r *Registry) GetAnalyticsRoute() IAnalyticsRoute {
	return NewAnalyticsRoute(r.controller, r.group)
}

func (r *Registry) GetSalesRoute() ISalesRoute {
	return NewSalesRoute(r.controller, r.group)
}

func (r *Registry) GetGeneralSettingsRoute() IGeneralSettingsRoute {
	return NewGeneralSettingsRoute(r.controller, r.group)
}