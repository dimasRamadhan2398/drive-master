package routes

import (
	"booking-service/controllers"
	"booking-service/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type CertificationRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	authMiddleware middlewares.IAuthMiddleware
}

type ICertificationRoute interface {
	Run()
}

func NewCertificationRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup, authMiddleware middlewares.IAuthMiddleware) ICertificationRoute {
	return &BookingRoute{controller: controller, group: group, authMiddleware: authMiddleware}
}

func (r *CertificationRoute) Run() {
	certifications := r.group.Group("/certifications")
	{
		certifications.GET("/all", r.controller.GetCertificationController().ListCertifications)
		certifications.POST("/create", r.controller.GetCertificationController().CreateCertification)
		certifications.GET("/:id", r.controller.GetCertificationController().GetCertification)
		certifications.PUT("/:id/status", r.controller.GetCertificationController().UpdateCertificationStatus)
		certifications.POST("/issue", r.controller.GetCertificationController().IssueCertification)
		certifications.POST("/revoke", r.controller.GetCertificationController().RevokeCertification)
		certifications.GET("/user/:userId", r.controller.GetCertificationController().GetUserCertifications)
		certifications.GET("/package/:packageId", r.controller.GetCertificationController().GetCertificationsByPackage)
	}
}