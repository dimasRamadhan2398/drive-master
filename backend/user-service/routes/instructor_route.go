package routes

import (
	"log"
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
	
	routes := []string{
		"POST   /api/v1/instructors/register",
		"POST   /api/v1/instructors/new",
		"GET    /api/v1/instructors/all",
		"GET    /api/v1/instructors/:id/profile",
		"PUT    /api/v1/instructors/:id/profile",
		"DELETE /api/v1/instructors/:id",
		"POST   /api/v1/instructors/:id/media/upload",
		"POST   /api/v1/instructors/:id/media/upload-base64",
		"DELETE /api/v1/instructors/:id/media",
		"GET    /api/v1/instructors/:id/media/metadata",
	}	

	log.Printf("[InstructorRoute] Registered routes under /api/v1/instructors:")
	for _, route := range routes {
	log.Printf("  %s", route)
	}

	group.POST("/register", u.controller.GetInstructorController().RegisterInstructor)
	group.POST("/new", u.controller.GetInstructorController().CreateInstructorProfile)
	group.GET("/all", u.controller.GetInstructorController().GetInstructorLists)
	group.GET("/:id/profile", u.controller.GetInstructorController().GetInstructorProfile)
	group.PUT("/:id/profile", u.authMiddleware.Authenticate(), u.controller.GetInstructorController().UpdateInstructorProfile)
	group.DELETE("/:id", u.authMiddleware.Authenticate(), u.controller.GetInstructorController().DeleteInstructor)

	// Media routes with instructor user ID
	group.POST("/:id/media/upload", u.authMiddleware.Authenticate(), u.controller.GetInstructorController().UploadProfilePic)
	group.POST("/:id/media/upload-base64", u.authMiddleware.Authenticate(), u.controller.GetInstructorController().UploadBase64Media)
	group.DELETE("/:id/media", u.authMiddleware.Authenticate(), u.controller.GetInstructorController().DeleteProfilePic)
	group.GET("/:id/media/metadata", u.controller.GetInstructorController().GetMediaMetadata)
}
