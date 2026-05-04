package routes

import (
	"user-service/controllers"
	"user-service/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type WorkExperienceRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type IWorkExperienceRoute interface {
	Run()
}

func NewWorkExperienceRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) IWorkExperienceRoute {
	return &WorkExperienceRoute{controller: controller, group: group}
}

func (r *WorkExperienceRoute) Run() {
	group := r.group.Group("/instructors")
	group.POST("/:id/work-experiences", middlewares.Authenticate(), r.controller.GetWorkExperienceController().CreateWorkExperience)
	group.GET("/:id/work-experiences", middlewares.Authenticate(), r.controller.GetWorkExperienceController().GetWorkExperience)

	workExpGroup := r.group.Group("/work-experiences")
	workExpGroup.PUT("/:expId", middlewares.Authenticate(), r.controller.GetWorkExperienceController().UpdateWorkExperience)
	workExpGroup.DELETE("/:expId", middlewares.Authenticate(), r.controller.GetWorkExperienceController().DeleteWorkExperience)
}
