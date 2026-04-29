package routes

import (
	"user-service/controllers"
	"user-service/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type InstructorRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type IInstructorRoute interface {
	Run()
}

func NewInstructorRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) IInstructorRoute {
	return &InstructorRoute{controller: controller, group: group}
}

func (u *InstructorRoute) Run() {
	group := u.group.Group("/instructors")
	group.GET("/login", middlewares.Authenticate(), u.)
	group.POST("/register", u.controller.GetUserController().Register)
	group.PUT("/forgot-password",  u.controller.GetUserController().ForgotPassword)
	group.PUT("/reset-password",  u.controller.GetUserController().ResetPassword)
}
