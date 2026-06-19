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
	}

	instructorGroup := r.group.Group("/instructors")
	{
		// Instructor Certification routes
		instructorGroup.POST("/:id/certifications", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().CreateCertification)
		instructorGroup.GET("/:id/certifications", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().ListCertifications)
		instructorGroup.GET("/:id/certifications/:certId", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().GetCertification)
		instructorGroup.PUT("/:id/certifications/:certId", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().UpdateCertification)
		instructorGroup.DELETE("/:id/certifications/:certId", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().DeleteCertification)
		instructorGroup.POST("/:id/certifications/:certId/verify", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().VerifyCertification)
	}
}
