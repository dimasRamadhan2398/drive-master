package routes

import (
	"core-service/controllers"

	"github.com/gin-gonic/gin"
)

// GeneralSettingsRoute handles general settings routes
type GeneralSettingsRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

// IGeneralSettingsRoute defines the interface for general settings route
type IGeneralSettingsRoute interface {
	Run()
}

func NewGeneralSettingsRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) IGeneralSettingsRoute {
	return &GeneralSettingsRoute{
		controller: controller,
		group:      group,
	}
}

// Run registers all general settings routes
func (r *GeneralSettingsRoute) Run() {
	settings := r.group.Group("/general-settings")
	{
		settings.GET("", r.controller.GetGeneralSettingsController().GetSettings)
		settings.PUT("", r.controller.GetGeneralSettingsController().UpdateSettings)
	}
}