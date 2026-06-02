package routes

import (
	"user-service/controllers"
	"user-service/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type MemberRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	authMiddleware middlewares.IAuthMiddleware
}

type IMemberRoute interface {
	Run()
}

func NewMemberRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup, authMiddleware middlewares.IAuthMiddleware) IMemberRoute {
	return &MemberRoute{controller: controller, group: group, authMiddleware: authMiddleware}
}

func (m *MemberRoute) Run() {
	group := m.group.Group("/members")
	group.GET("/", m.authMiddleware.Authenticate(), m.controller.GetMemberController().GetMemberLists)
	group.GET("/:userId/profile", m.authMiddleware.Authenticate(), m.controller.GetMemberController().GetMemberProfile)
	group.PUT("/:userId/profile", m.authMiddleware.Authenticate(), m.controller.GetMemberController().UpdateMemberProfile)
}