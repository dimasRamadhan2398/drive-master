package controllers

import (
	"core-service/models/dto"
	"core-service/pkg/clients"
	"core-service/pkg/response"
	"fmt"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	paymentClient clients.ITransactionClient
}

type ITransactionController interface {
	GetRecentTransactions(ctx *gin.Context)
}

// NewTransactionController creates a new transaction controller
func NewTransactionController(paymentClient clients.ITransactionClient) ITransactionController {
	return &TransactionController{
		paymentClient: paymentClient,
	}
}

// GetRecentTransactions handles GET /api/v1/admin/sales/analytics/recent-transactions
// @Summary Get recent transactions from payment service
// @Description Retrieves recent payment transactions from the payment service
// @Tags Sales Analytics
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.Response
// @Router /admin/sales/analytics/recent-transactions [get]
func (c *TransactionController) GetRecentTransactions(ctx *gin.Context) {
	page := 1
	limit := 10

	// Parse pagination params
	if p := ctx.Query("page"); p != "" {
		if parsed, ok := fmt.Sscanf(p, "%d", &page); ok == nil && parsed > 0 {
			// OK
		}
	}
	if l := ctx.Query("limit"); l != "" {
		if parsed, ok := fmt.Sscanf(l, "%d", &limit); ok == nil && parsed > 0 {
			// OK
		}
	}

	// Default limits
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	data, err := c.paymentClient.GetRecentTransactions(ctx.Request.Context(), page, limit)
	if err != nil {
		response.InternalServerError(ctx, "Failed to get recent transactions: "+err.Error())
		return
	}

	// Convert to DTO
	transactions := make([]dto.RecentTransactionResponse, len(data.Data))
	for i, tx := range data.Data {
		transactions[i] = dto.RecentTransactionResponse{
			ID:            tx.ID,
			PaymentID:     tx.PaymentID,
			Type:          tx.Type,
			Status:        tx.Status,
			Amount:        tx.Amount,
			Currency:      tx.Currency,
			Gateway:       tx.Gateway,
			GatewayTxnID:  tx.GatewayTxnID,
			ErrorCode:     tx.ErrorCode,
			ErrorMessage:  tx.ErrorMessage,
			ProcessedAt:   tx.ProcessedAt,
			CreatedAt:     tx.CreatedAt,
		}
	}

	response.OK(ctx, "Recent transactions retrieved successfully", dto.RecentTransactionsResponse{
		Data:  transactions,
		Total: data.Total,
		Page:  data.Page,
		Limit: data.Limit,
	})
}
