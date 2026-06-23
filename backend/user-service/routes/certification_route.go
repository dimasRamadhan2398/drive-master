package routes

import (
	"user-service/controllers"
	"user-service/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type CertificationRoute struct {
	controller     controllers.IControllerRegistry
	group          *gin.RouterGroup
	authMiddleware middlewares.IAuthMiddleware
}

type ICertificationRoute interface {
	Run()
}

func NewCertificationRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup, authMiddleware middlewares.IAuthMiddleware) ICertificationRoute {
	return &CertificationRoute{controller: controller, group: group, authMiddleware: authMiddleware}
}

func (r *CertificationRoute) Run() {
	group := r.group.Group("/certificates")
	{
		// Certificate stats
		group.GET("/stats", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().GetCertificateStats)
		group.POST("/issue", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().IssueCertification)
		group.POST("/revoke", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().RevokeCertification)
		group.GET("/user/:userId", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().GetUserCertifications)
	}

	memberGroup := r.group.Group("/members")
	{
		// Member Certification routes
		memberGroup.POST("/:id/certifications", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().CreateCertification)
		memberGroup.GET("/:id/certifications", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().ListCertifications)
		memberGroup.GET("/:id/certifications/:certId", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().GetCertification)
		memberGroup.PUT("/:id/certifications/:certId", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().UpdateCertification)
		memberGroup.DELETE("/:id/certifications/:certId", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().DeleteCertification)
		memberGroup.POST("/:id/certifications/:certId/verify", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().VerifyCertification)
	}
}
