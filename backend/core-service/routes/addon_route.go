package routes

import (
	"core-service/controllers"

	"github.com/gin-gonic/gin"
)

type AddOnRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type IAddOnRoute interface {
	Run()
}

func NewAddOnRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) IAddOnRoute {
	return &AddOnRoute{
		controller: controller,
		group:      group,
	}
}

func (r *AddOnRoute) Run() {
	addons := r.group.Group("/addons")
	{
		addons.GET("/all", r.controller.GetAddOnController().GetAllAddOns)
		addons.GET("/:id", r.controller.GetAddOnController().GetAddOnByID)
		addons.POST("/create", r.controller.GetAddOnController().CreateAddOn)
		addons.PUT("/:id", r.controller.GetAddOnController().UpdateAddOn)
		addons.DELETE("/:id", r.controller.GetAddOnController().DeleteAddOn)
		addons.PUT("/toggle-status/:id", r.controller.GetAddOnController().ToggleStatusAddOn)
	}
}
