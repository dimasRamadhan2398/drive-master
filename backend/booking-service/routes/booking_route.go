package routes

import (
	"booking-service/controllers"

	"github.com/gin-gonic/gin"
)

type BookingRoute struct {
	controller controllers.IControllerRegistry
}

type IBookingRoute interface {
	Run(group *gin.RouterGroup)
}

func (r *BookingRoute) Run(group *gin.RouterGroup) {
	bookings := group.Group("/bookings")
	{
		bookings.GET("", r.controller.GetBookingController().ListBookings)
		bookings.POST("", r.controller.GetBookingController().CreateBooking)
		bookings.GET("/:id", r.controller.GetBookingController().GetBooking)
		bookings.PUT("/:id", r.controller.GetBookingController().UpdateBooking)
		bookings.POST("/:id/cancel", r.controller.GetBookingController().CancelBooking)
		bookings.POST("/:id/confirm", r.controller.GetBookingController().ConfirmBooking)
		bookings.POST("/:id/complete", r.controller.GetBookingController().CompleteBooking)
	}
}