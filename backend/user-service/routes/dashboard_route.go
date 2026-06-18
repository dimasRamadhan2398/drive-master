package routes

import (
	"log"
	"user-service/controllers"
	"user-service/pkg/middlewares"

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
	log.Printf("  GET /api/v1/dashboard/stats")
	log.Printf("  GET /api/v1/dashboard/certifications/stats")

	group.GET("/stats", r.authMiddleware.Authenticate(), r.controller.GetDashboardController().GetStats)
	group.GET("/certifications/stats", r.authMiddleware.Authenticate(), r.controller.GetDashboardController().GetCertificationStats)
}
