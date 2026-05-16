package routes

import (
	"user-service/controllers"
	"user-service/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type CoverageAreaRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	authMiddleware middlewares.IAuthMiddleware
}

type ICoverageAreaRoute interface {
	Run()
}

func NewCoverageAreaRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup, authMiddleware middlewares.IAuthMiddleware) ICoverageAreaRoute {
	return &CoverageAreaRoute{controller: controller, group: group, authMiddleware: authMiddleware}
}

func (r *CoverageAreaRoute) Run() {
	group := r.group.Group("/instructors")
	group.POST("/:id/coverage-areas", r.authMiddleware.Authenticate(), r.controller.GetCoverageAreaController().AddCoverageArea)
	group.GET("/:id/coverage-areas", r.authMiddleware.Authenticate(), r.controller.GetCoverageAreaController().GetCoverageAreas)
	group.DELETE("/:id/coverage-areas/:areaId", r.authMiddleware.Authenticate(), r.controller.GetCoverageAreaController().RemoveCoverageArea)
}
