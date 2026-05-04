package routes

import (
	"user-service/controllers"
	"user-service/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type CoverageAreaRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type ICoverageAreaRoute interface {
	Run()
}

func NewCoverageAreaRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) ICoverageAreaRoute {
	return &CoverageAreaRoute{controller: controller, group: group}
}

func (r *CoverageAreaRoute) Run() {
	group := r.group.Group("/instructors")
	group.POST("/:id/coverage-areas", middlewares.Authenticate(), r.controller.GetCoverageAreaController().AddCoverageArea)
	group.GET("/:id/coverage-areas", middlewares.Authenticate(), r.controller.GetCoverageAreaController().GetCoverageAreas)
	group.DELETE("/:id/coverage-areas/:areaId", middlewares.Authenticate(), r.controller.GetCoverageAreaController().RemoveCoverageArea)
}
