package routes

import (
	"booking-service/controllers"
	"booking-service/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type BookingRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	authMiddleware middlewares.IAuthMiddleware
}

type IBookingRoute interface {
	Run()
}

func NewBookingRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup, authMiddleware middlewares.IAuthMiddleware) IBookingRoute {
	return &BookingRoute{controller: controller, group: group, authMiddleware: authMiddleware}
}

func (r *BookingRoute) Run() {
	bookings := r.group.Group("/bookings")
	{
		bookings.GET("", r.authMiddleware.Authenticate(),r.controller.GetBookingController().ListBookings)
		bookings.POST("", r.authMiddleware.Authenticate(), r.controller.GetBookingController().CreateBooking)
		bookings.GET("/:id", r.authMiddleware.Authenticate(), r.controller.GetBookingController().GetBooking)
		bookings.PUT("/:id", r.authMiddleware.Authenticate(), r.controller.GetBookingController().UpdateBooking)
		bookings.POST("/:id/cancel", r.authMiddleware.Authenticate(), r.controller.GetBookingController().CancelBooking)
		bookings.POST("/:id/confirm", r.authMiddleware.Authenticate(), r.controller.GetBookingController().ConfirmBooking)
		bookings.POST("/:id/complete", r.authMiddleware.Authenticate(), r.controller.GetBookingController().CompleteBooking)
	}
}