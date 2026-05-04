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
	auth := u.group.Group("/auth")
	auth.POST("/login", u.controller.GetAuthController().Login)
	auth.POST("/register", u.controller.GetAuthController().Register)
	auth.POST("/forgot-password", u.controller.GetAuthController().ResetPassword)
	auth.POST("/confirm-reset-password", u.controller.GetAuthController().ConfirmResetPassword)
	auth.POST("/verify-otp", u.controller.GetAuthController().VerifyOTP)
	auth.POST("/resend-otp", u.controller.GetAuthController().ResendOTP)
}

