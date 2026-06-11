package routes

import (
	"log"
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
	log.Printf("[MemberRoute] Registered routes under /api/v1/members:")
	log.Printf("  GET /api/v1/members/all")
	log.Printf("  GET /api/v1/members/recent")
	log.Printf("  GET /api/v1/members/search")
	log.Printf("  GET /api/v1/members/:userId/profile")
	log.Printf("  PUT /api/v1/members/:userId/profile")
	log.Printf("  GET /api/v1/members/:userId/certificates")
	log.Printf("  GET /api/v1/members/:userId/certificates/:certId/download")

	group.GET("/all", func(c *gin.Context) {
		log.Printf("[MemberRoute] GET /api/v1/members/all called - Path: %s, Query: %v", c.Request.URL.Path, c.Request.URL.Query())
		m.authMiddleware.Authenticate()(c)
		if c.IsAborted() {
			return
		}
		m.controller.GetMemberController().GetMemberLists(c)
	})
	group.GET("/recent", m.authMiddleware.AuthorizeAdmin(), m.controller.GetMemberController().FindRecentRegistrations)
	group.GET("/:userId/profile", m.authMiddleware.Authenticate(), m.controller.GetMemberController().GetMemberProfile)
	group.PUT("/:userId/profile", m.authMiddleware.Authenticate(), m.controller.GetMemberController().UpdateMemberProfile)
	group.GET("/search", m.authMiddleware.Authenticate(), m.controller.GetMemberController().SearchMembersWithPagination)
	group.GET("/:userId/certificates", m.authMiddleware.Authenticate(), m.controller.GetMemberController().GetMemberCertificates)
	group.GET("/:userId/certificates/:certId/download", m.authMiddleware.Authenticate(), m.controller.GetMemberController().DownloadMemberCertificate)

}