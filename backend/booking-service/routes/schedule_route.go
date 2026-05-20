package routes

import (
	"booking-service/controllers"
	"booking-service/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type ScheduleRoute struct {
	controller     controllers.IControllerRegistry
	group         *gin.RouterGroup
	authMiddleware middlewares.IAuthMiddleware
}

type IScheduleRoute interface {
	Run()
}

func NewScheduleRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup, authMiddleware middlewares.IAuthMiddleware) IScheduleRoute {
	return &ScheduleRoute{controller: controller, group: group, authMiddleware: authMiddleware}
}

func (r *ScheduleRoute) Run() {
	schedules := r.group.Group("/schedules")
	{
		schedules.GET("", r.authMiddleware.Authenticate(), r.controller.GetScheduleController().ListSchedules)
		schedules.POST("", r.authMiddleware.Authenticate(), r.controller.GetScheduleController().CreateSchedule)
		schedules.GET("/filter", r.authMiddleware.Authenticate(), r.controller.GetScheduleController().ListSchedulesFiltered)
		schedules.GET("/available", r.authMiddleware.Authenticate(), r.controller.GetScheduleController().GetAvailableSchedules)
		schedules.GET("/:id", r.authMiddleware.Authenticate(), r.controller.GetScheduleController().GetSchedule)
		schedules.PUT("/:id", r.authMiddleware.Authenticate(), r.controller.GetScheduleController().UpdateSchedule)
		schedules.DELETE("/:id", r.authMiddleware.Authenticate(), r.controller.GetScheduleController().DeleteSchedule)
		schedules.POST("/:id/book", r.authMiddleware.Authenticate(), r.controller.GetScheduleController().BookSlot)
		schedules.POST("/:id/cancel", r.authMiddleware.Authenticate(), r.controller.GetScheduleController().CancelBooking)
	}
}