package routes

import (
	"booking-service/controllers"
	"booking-service/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type SessionRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	authMiddleware middlewares.IAuthMiddleware
}

type ISessionRoute interface {
	Run()
}

func NewSessionRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup, authMiddleware middlewares.IAuthMiddleware) ISessionRoute {
	return &SessionRoute{controller: controller, group: group, authMiddleware: authMiddleware}
}

func (r *SessionRoute) Run() {
	sessions := r.group.Group("/sessions")
	{
		sessions.GET("", r.controller.GetSessionController().ListSessions)
		sessions.POST("", r.controller.GetSessionController().CreateSession)
		sessions.GET("/:id", r.controller.GetSessionController().GetSession)
		sessions.POST("/:id/start", r.controller.GetSessionController().StartSession)
		sessions.POST("/:id/complete", r.controller.GetSessionController().CompleteSession)
		sessions.POST("/:id/cancel", r.controller.GetSessionController().CancelSession)
	}
}