package routes

import (
	"log"

	"payment-service/controllers"
	"payment-service/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type TransactionRoute struct {
	controller     controllers.ITransactionController
	group          *gin.RouterGroup
	publicGroup    *gin.RouterGroup
	authMiddleware middlewares.IAuthMiddleware
}

type ITransactionRoute interface {
	Run()
}

func NewTransactionRoute(
	controller controllers.ITransactionController,
	group *gin.RouterGroup,
	publicGroup *gin.RouterGroup,
	authMiddleware middlewares.IAuthMiddleware,
) ITransactionRoute {
	return &TransactionRoute{
		controller:     controller,
		group:          group,
		publicGroup:    publicGroup,
		authMiddleware: authMiddleware,
	}
}

func (t *TransactionRoute) Run() {
	group := t.group.Group("/transactions")
	log.Printf("[TransactionRoute] Registered routes under /api/v1/transactions:")
	log.Printf("  GET /api/v1/transactions")
	log.Printf("  GET /api/v1/transactions/:id")
	log.Printf("  GET /api/v1/transactions/payment/:paymentId")

	group.GET("", t.authMiddleware.Authenticate(), t.controller.ListTransactions)
	group.GET("/all", t.authMiddleware.Authenticate(), t.controller.ListTransactions)
	group.GET("/:id", t.authMiddleware.Authenticate(), t.controller.GetTransaction)
	group.GET("/payment/:paymentId", t.authMiddleware.Authenticate(), t.controller.GetTransactionsByPaymentID)

	// New checkout endpoints
	paymentsGroup := t.group.Group("/payments")
	paymentsGroup.POST("/transactions", t.authMiddleware.Authenticate(), t.controller.CreateTransaction)
	paymentsGroup.GET("", t.authMiddleware.Authenticate(), t.controller.ListPayments)
	// Use :id consistently and handle orderId vs paymentId in controller
	paymentsGroup.GET("/order/:orderId", t.authMiddleware.Authenticate(), t.controller.GetPaymentByOrderID)
	paymentsGroup.GET("/:id/details", t.authMiddleware.Authenticate(), t.controller.GetPaymentDetail)
	paymentsGroup.GET("/:id/status", t.authMiddleware.Authenticate(), t.controller.GetPaymentStatus)
	paymentsGroup.GET("/:id", t.authMiddleware.Authenticate(), t.controller.GetPayment)

	publicPaymentsGroup := t.publicGroup.Group("/payments")
	publicPaymentsGroup.POST("/callback", t.controller.Callback)
}