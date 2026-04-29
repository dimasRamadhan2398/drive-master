package routes

import (
	"user-service/controllers"
	"user-service/pkg/validations"
	"user-service/services"

	"github.com/gin-gonic/gin"
)

type Registry struct {
	controller    controllers.IControllerRegistry
	group         *gin.RouterGroup
	authMiddleware *middleware.AuthMiddleware
	client        clients.IClientRegistry
}

type IRegistry interface {
	Serve()
	SetController(controller controllers.IControllerRegistry)
}

func NewRouteRegistry(
	group *gin.RouterGroup,
	authMiddleware *middleware.AuthMiddleware,
) IRegistry {
	return &Registry{
		group:          group,
		authMiddleware: authMiddleware,
	}
}

func (r *Registry) Serve(controller controllers.IControllerRegistry) {
	r.auth()
	r.users()
	r.mfa()
	r.trustedDevices()
}

func (r *Registry) SetController(controller controllers.IControllerRegistry) {
	r.controller = controller
}

func NewControllerRegistry(service services.IServiceRegistry, validator *validations.Validator) controllers.IControllerRegistry {
	return controllers.NewControllerRegistry(service, validator)
}

func (r *Registry) auth() {
	auth := r.group.Group("/auth")
	{
		auth.POST("/login", c.Login)
		auth.POST("/register", c.Register)
		auth.POST("/forgot-password", c.ForgotPassword)
		auth.POST("/reset-password", c.ResetPassword)

		// Protected routes
		protected := auth.Group("")
		protected.Use(r.authMiddleware.Authenticate())
		{
			protected.POST("/logout", r.controller.GetAuth().Logout)
			protected.POST("/change-password", r.controller.GetAuth().ChangePassword)
			protected.GET("/me", r.controller.GetAuth().Me)
		}
	}
}