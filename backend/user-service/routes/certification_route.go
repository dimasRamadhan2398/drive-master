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
	certGroup := r.group.Group("/certificates")
	{
		// Admin routes - Issue and manage certificates
		certGroup.GET("/stats", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().GetCertificateStats)
		certGroup.POST("", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().IssueCertificate)
		certGroup.GET("/member/:memberId", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().GetMemberCertificates)
		certGroup.DELETE("/:id", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().RevokeCertificate)
		certGroup.GET("", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().GetAllCertificates)

		// Member routes - View and download their certificates
		certGroup.GET("/:id", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().GetCertificate)
		certGroup.GET("/:id/pdf", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().DownloadCertificatePDF)
	}
}
