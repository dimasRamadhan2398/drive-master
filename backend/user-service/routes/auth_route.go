package routes

import (
	"user-service/controllers"

	"github.com/gin-gonic/gin"
)

type AuthRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type IAuthRoute interface {
	Run()
}

func NewAuthRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) IAuthRoute {
	return &AuthRoute{controller: controller, group: group}
}

func (u *AuthRoute) Run() {
	group := u.group.Group("/auth")
	group.POST("/login", u.controller.GetAuthController().Login)
	group.POST("/register", u.controller.GetAuthController().Register)
	group.PUT("/forgot-password",  u.controller.GetAuthController().ForgotPassword)
	group.PUT("/reset-password",  u.controller.GetAuthController().ResetPassword)
}

