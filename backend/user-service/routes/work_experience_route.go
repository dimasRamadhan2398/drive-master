package routes

import (
	"user-service/controllers"
	"user-service/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type WorkExperienceRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	authMiddleware middlewares.IAuthMiddleware
}

type IWorkExperienceRoute interface {
	Run()
}

func NewWorkExperienceRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup, authMiddleware middlewares.IAuthMiddleware) IWorkExperienceRoute {
	return &WorkExperienceRoute{controller: controller, group: group, authMiddleware: authMiddleware}
}

func (r *WorkExperienceRoute) Run() {
	group := r.group.Group("/instructors")
	group.POST("/:id/work-experiences", r.authMiddleware.Authenticate(), r.controller.GetWorkExperienceController().CreateWorkExperience)
	group.GET("/:id/work-experiences", r.authMiddleware.Authenticate(), r.controller.GetWorkExperienceController().GetWorkExperience)

	workExpGroup := r.group.Group("/work-experiences")
	workExpGroup.PUT("/:expId", r.authMiddleware.Authenticate(), r.controller.GetWorkExperienceController().UpdateWorkExperience)
	workExpGroup.DELETE("/:expId", r.authMiddleware.Authenticate(), r.controller.GetWorkExperienceController().DeleteWorkExperience)
}
