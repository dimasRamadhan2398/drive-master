package routes

import (
	"user-service/controllers"
	"user-service/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type AuthRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	authMiddleware middlewares.IAuthMiddleware
}

type IAuthRoute interface {
	Run()
}

func NewAuthRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup, authMiddleware middlewares.IAuthMiddleware) IAuthRoute {
	return &AuthRoute{controller: controller, group: group, authMiddleware: authMiddleware}
}

func (u *AuthRoute) Run() {
	auth := u.group.Group("/auth")
	auth.POST("/login", u.controller.GetAuthController().Login)
	auth.POST("/register", u.controller.GetAuthController().Register)
	auth.POST("/forgot-password", u.controller.GetAuthController().ResetPassword)
	auth.POST("/confirm-reset-password", u.controller.GetAuthController().ConfirmResetPassword)
	auth.POST("/verify-otp", u.controller.GetAuthController().VerifyOTP)
	auth.POST("/resend-otp", u.controller.GetAuthController().ResendOTP)
	auth.POST("/refresh", u.controller.GetAuthController().RefreshToken)
}

