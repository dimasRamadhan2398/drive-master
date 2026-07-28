package routes

import (
	"booking-service/controllers"
	"booking-service/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type PaymentRoute struct {
	controllerRegistry controllers.IControllerRegistry
	group             *gin.RouterGroup
	authMiddleware     middlewares.IAuthMiddleware
}

func NewPaymentRoute(
	controllerRegistry controllers.IControllerRegistry,
	group *gin.RouterGroup,
	authMiddleware middlewares.IAuthMiddleware,
) IPaymentRoute {
	return &PaymentRoute{
		controllerRegistry: controllerRegistry,
		group:              group,
		authMiddleware:     authMiddleware,
	}
}

func (r *PaymentRoute) Run() {
	payments := r.group.Group("/payments")
	{
		// Public routes (for Midtrans callback - no auth required)
		payments.POST("/callback", r.controllerRegistry.GetPaymentController().HandleCallback)

		// Protected routes
		payments.Use(r.authMiddleware.Authenticate())

		// Create payment
		payments.POST("", r.controllerRegistry.GetPaymentController().CreatePayment)

		// List payments
		payments.GET("", r.controllerRegistry.GetPaymentController().ListPayments)

		// Get payment by ID
		payments.GET("/:id", r.controllerRegistry.GetPaymentController().GetPayment)

		// Get payment by order ID
		payments.GET("/order/:orderId", r.controllerRegistry.GetPaymentController().GetPaymentByOrderID)

		// Get payment detail by order ID
		payments.GET("/order/:orderId/detail", r.controllerRegistry.GetPaymentController().GetPaymentDetail)

		// List user payments
		payments.GET("/user/:userId", r.controllerRegistry.GetPaymentController().ListUserPayments)

		// Cancel payment
		payments.POST("/order/:orderId/cancel", r.controllerRegistry.GetPaymentController().CancelPayment)
	}
}

// IPaymentRoute interface for route registry
type IPaymentRoute interface {
	Run()
}