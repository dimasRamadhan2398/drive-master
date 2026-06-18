package routes

import (
	"log"
	"booking-service/controllers"
	"booking-service/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type DashboardRoute struct {
	controller     controllers.IControllerRegistry
	group          *gin.RouterGroup
	authMiddleware middlewares.IAuthMiddleware
}

type IDashboardRoute interface {
	Run()
}

func NewDashboardRoute(
	controller controllers.IControllerRegistry,
	group *gin.RouterGroup,
	authMiddleware middlewares.IAuthMiddleware,
) IDashboardRoute {
	return &DashboardRoute{
		controller:     controller,
		group:          group,
		authMiddleware: authMiddleware,
	}
}

func (r *DashboardRoute) Run() {
	group := r.group.Group("/dashboard")
	log.Printf("[DashboardRoute] Registered routes under /api/v1/dashboard:")
	log.Printf("  GET /api/v1/dashboard/sessions/stats")
	log.Printf("  GET /api/v1/dashboard/revenue/stats")

	group.GET("/sessions/stats", r.authMiddleware.Authenticate(), r.controller.GetDashboardController().GetSessionStats)
	group.GET("/revenue/stats", r.authMiddleware.Authenticate(), r.controller.GetDashboardController().GetRevenueStats)
}