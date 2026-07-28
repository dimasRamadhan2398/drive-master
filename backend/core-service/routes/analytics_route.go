package routes

import (
	"core-service/controllers"

	"github.com/gin-gonic/gin"
)

type AnalyticsRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type IAnalyticsRoute interface {
	Run()
}

func NewAnalyticsRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) IAnalyticsRoute {
	return &AnalyticsRoute{
		controller: controller,
		group:      group,
	}
}

func (r *AnalyticsRoute) Run() {
	admin := r.group.Group("/admin/analytics")
	{
		admin.GET("/overview", r.controller.GetAnalyticsController().GetOverview)
		admin.GET("/funnel", r.controller.GetAnalyticsController().GetFunnel)
	}
}
