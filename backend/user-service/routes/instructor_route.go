package routes

import (
	"user-service/controllers"
	"user-service/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type InstructorRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	authMiddleware middlewares.IAuthMiddleware
}

type IInstructorRoute interface {
	Run()
}

func NewInstructorRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup, authMiddleware middlewares.IAuthMiddleware) IInstructorRoute {
	return &InstructorRoute{controller: controller, group: group, authMiddleware: authMiddleware}
}

func (u *InstructorRoute) Run() {
	group := u.group.Group("/instructors")
	group.GET("/", u.controller.GetInstructorController().GetInstructorLists)
	group.GET("/:id/profile", u.controller.GetInstructorController().GetInstructorProfile)
	group.PUT("/:id/profile", u.authMiddleware.Authenticate(), u.controller.GetInstructorController().UpdateInstructorProfile)

	// Media routes with instructor user ID
	group.POST("/:id/media/upload", u.authMiddleware.Authenticate(), u.controller.GetInstructorController().UploadProfilePic)
	group.POST("/:id/media/upload-base64", u.authMiddleware.Authenticate(), u.controller.GetInstructorController().UploadBase64Media)
	group.DELETE("/:id/media", u.authMiddleware.Authenticate(), u.controller.GetInstructorController().DeleteProfilePic)
	group.GET("/:id/media/metadata", u.controller.GetInstructorController().GetMediaMetadata)
}
