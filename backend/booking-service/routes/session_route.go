package routes

import (
	"booking-service/controllers"

	"github.com/gin-gonic/gin"
)

type SessionRoute struct {
	controller controllers.IControllerRegistry
}

type ISessionRoute interface {
	Run(group *gin.RouterGroup)
}

func (r *SessionRoute) Run(group *gin.RouterGroup) {
	sessions := group.Group("/sessions")
	{
		sessions.GET("", r.controller.GetSessionController().ListSessions)
		sessions.POST("", r.controller.GetSessionController().CreateSession)
		sessions.GET("/:id", r.controller.GetSessionController().GetSession)
	}
}