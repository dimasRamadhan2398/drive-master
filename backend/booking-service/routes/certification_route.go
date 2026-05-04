package routes

import (
	"booking-service/controllers"

	"github.com/gin-gonic/gin"
)

type CertificationRoute struct {
	controller controllers.IControllerRegistry
}

type ICertificationRoute interface {
	Run(group *gin.RouterGroup)
}

func (r *CertificationRoute) Run(group *gin.RouterGroup) {
	certifications := group.Group("/certifications")
	{
		certifications.GET("", r.controller.GetCertificationController().ListCertifications)
		certifications.POST("", r.controller.GetCertificationController().CreateCertification)
		certifications.GET("/:id", r.controller.GetCertificationController().GetCertification)
		certifications.PUT("/:id/status", r.controller.GetCertificationController().UpdateCertificationStatus)
		certifications.POST("/issue", r.controller.GetCertificationController().IssueCertification)
		certifications.POST("/revoke", r.controller.GetCertificationController().RevokeCertification)
		certifications.GET("/user/:userId", r.controller.GetCertificationController().GetUserCertifications)
		certifications.GET("/package/:packageId", r.controller.GetCertificationController().GetCertificationsByPackage)
	}
}