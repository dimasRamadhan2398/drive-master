package routes

import (
	"booking-service/controllers"
	"booking-service/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type EnrollmentRoute struct {
	controller     controllers.IControllerRegistry
	group         *gin.RouterGroup
	authMiddleware middlewares.IAuthMiddleware
}

type IEnrollmentRoute interface {
	Run()
}

func NewEnrollmentRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup, authMiddleware middlewares.IAuthMiddleware) IEnrollmentRoute {
	return &EnrollmentRoute{controller: controller, group: group, authMiddleware: authMiddleware}
}

func (r *EnrollmentRoute) Run() {
	enrollments := r.group.Group("/enrollments")
	{
		enrollments.GET("", r.authMiddleware.Authenticate(), r.controller.GetEnrollmentController().ListEnrollments)
		enrollments.POST("", r.authMiddleware.Authenticate(), r.controller.GetEnrollmentController().CreateEnrollment)
		enrollments.GET("/:id", r.authMiddleware.Authenticate(), r.controller.GetEnrollmentController().GetEnrollment)
		enrollments.PUT("/:id", r.authMiddleware.Authenticate(), r.controller.GetEnrollmentController().UpdateEnrollment)
		enrollments.POST("/:id/cancel", r.authMiddleware.Authenticate(), r.controller.GetEnrollmentController().CancelEnrollment)
		enrollments.POST("/:id/pay", r.authMiddleware.Authenticate(), r.controller.GetEnrollmentController().MarkAsPaid)
		enrollments.GET("/user/:userId", r.authMiddleware.Authenticate(), r.controller.GetEnrollmentController().ListUserEnrollments)
	}
}
