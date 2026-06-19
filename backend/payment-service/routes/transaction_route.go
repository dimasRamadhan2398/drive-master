package routes

import (
	"log"

	"payment-service/controllers"
	"payment-service/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type TransactionRoute struct {
	controller      controllers.ITransactionController
	group           *gin.RouterGroup
	authMiddleware  middlewares.IAuthMiddleware
}

type ITransactionRoute interface {
	Run()
}

func NewTransactionRoute(controller controllers.ITransactionController, group *gin.RouterGroup, authMiddleware middlewares.IAuthMiddleware) ITransactionRoute {
	return &TransactionRoute{
		controller:     controller,
		group:          group,
		authMiddleware: authMiddleware,
	}
}

func (t *TransactionRoute) Run() {
	group := t.group.Group("/transactions")
	log.Printf("[TransactionRoute] Registered routes under /api/v1/transactions:")
	log.Printf("  GET /api/v1/transactions")
	log.Printf("  GET /api/v1/transactions/:id")
	log.Printf("  GET /api/v1/transactions/payment/:paymentId")

	group.GET("/all", t.authMiddleware.Authenticate(), t.controller.ListTransactions)
	group.GET("/:id", t.authMiddleware.Authenticate(), t.controller.GetTransaction)
	group.GET("/payment/:paymentId", t.authMiddleware.Authenticate(), t.controller.GetTransactionsByPaymentID)
}