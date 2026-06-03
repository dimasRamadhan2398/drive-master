package routes

import (
	"core-service/controllers"

	"github.com/gin-gonic/gin"
)

type SalesRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type ISalesRoute interface {
	Run()
}

func NewSalesRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) ISalesRoute {
	return &SalesRoute{
		controller: controller,
		group:      group,
	}
}

func (r *SalesRoute) Run() {
	sales := r.group.Group("/admin/sales")
	{
		// CRUD endpoints
		sales.GET("", r.controller.GetSalesController().ListSales)
		sales.POST("", r.controller.GetSalesController().CreateSale)
		sales.GET("/:id", r.controller.GetSalesController().GetSaleByID)
		sales.DELETE("/:id", r.controller.GetSalesController().DeleteSale)
		sales.PATCH("/:id/status", r.controller.GetSalesController().UpdateSaleStatus)
		sales.GET("/order/:orderNumber", r.controller.GetSalesController().GetSaleByOrderNumber)

		// Analytics endpoints
		analytics := sales.Group("/analytics")
		{
			analytics.GET("/overview", r.controller.GetSalesController().GetOverview)
			analytics.GET("/trend", r.controller.GetSalesController().GetSalesTrend)
			analytics.GET("/by-source", r.controller.GetSalesController().GetSalesBySource)
			analytics.GET("/by-package-type", r.controller.GetSalesController().GetSalesByPackageType)
			analytics.GET("/monthly/:year/:month", r.controller.GetSalesController().GetMonthlySales)
			analytics.GET("/yearly/:year", r.controller.GetSalesController().GetYearlySales)
		}
	}
}
