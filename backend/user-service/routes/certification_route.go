package routes

import (
	"user-service/controllers"
	"user-service/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type CertificationRoute struct {
	controller    controllers.IControllerRegistry
	group         *gin.RouterGroup
	authMiddleware middlewares.IAuthMiddleware
}

type ICertificationRoute interface {
	Run()
}

func NewCertificationRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup, authMiddleware middlewares.IAuthMiddleware) ICertificationRoute {
	return &CertificationRoute{controller: controller, group: group, authMiddleware: authMiddleware}
}

func (r *CertificationRoute) Run() {
	group := r.group.Group("/instructors")
	{
		// Certification routes
		group.POST("/:id/certifications", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().CreateCertification)
		group.GET("/:id/certifications", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().ListCertifications)
		group.GET("/:id/certifications/:certId", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().GetCertification)
		group.PUT("/:id/certifications/:certId", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().UpdateCertification)
		group.DELETE("/:id/certifications/:certId", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().DeleteCertification)
		group.POST("/:id/certifications/:certId/verify", r.authMiddleware.Authenticate(), r.controller.GetCertificationController().VerifyCertification)
	}
}