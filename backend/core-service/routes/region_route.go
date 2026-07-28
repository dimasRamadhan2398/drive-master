package routes

import (
	"core-service/controllers"

	"github.com/gin-gonic/gin"
)

type RegionRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type IRegionRoute interface {
	Run()
}

func NewRegionRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) IRegionRoute {
	return &RegionRoute{
		controller: controller,
		group:      group,
	}
}

func (r *RegionRoute) Run() {
	regions := r.group.Group("/regions")
	{
		regions.GET("/provinces", r.controller.GetRegionController().GetAllProvinces)
		regions.GET("/provinces/:province/regencies", r.controller.GetRegionController().GetRegenciesByProvince)
		regions.GET("/regencies/:regency/districts", r.controller.GetRegionController().GetDistrictsByRegency)
	}
}