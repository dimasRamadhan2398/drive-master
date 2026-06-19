package controllers

import (
	"net/http"

	"payment-service/models/dto"
	"payment-service/pkg/response"
	"payment-service/repositories"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionController struct {
	transactionRepo repositories.ITransactionRepository
}

type ITransactionController interface {
	GetTransaction(c *gin.Context)
	GetTransactionsByPaymentID(c *gin.Context)
	ListTransactions(c *gin.Context)
}

func NewTransactionController(transactionRepo repositories.ITransactionRepository) ITransactionController {
	return &TransactionController{
		transactionRepo: transactionRepo,
	}
}

// GetTransaction handles GET /api/v1/transactions/:id
// @Summary Get transaction by ID
// @Description Get a single transaction by its UUID
// @Tags Transactions
// @Produce json
// @Param id path string true "Transaction ID (UUID)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /transactions/{id} [get]
func (t *TransactionController) GetTransaction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid transaction ID format")
		return
	}

	tx, err := t.transactionRepo.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Transaction not found")
		return
	}

	response.Success(c, http.StatusOK, "Transaction retrieved successfully", tx)
}

// GetTransactionsByPaymentID handles GET /api/v1/transactions/payment/:paymentId
// @Summary Get transactions by payment ID
// @Description Get all transactions associated with a specific payment
// @Tags Transactions
// @Produce json
// @Param paymentId path string true "Payment ID (UUID)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /transactions/payment/{paymentId} [get]
func (t *TransactionController) GetTransactionsByPaymentID(c *gin.Context) {
	paymentIDStr := c.Param("paymentId")
	paymentID, err := uuid.Parse(paymentIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid payment ID format")
		return
	}

	txs, err := t.transactionRepo.GetByPaymentID(paymentID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve transactions")
		return
	}

	response.Success(c, http.StatusOK, "Transactions retrieved successfully", txs)
}

// ListTransactions handles GET /api/v1/transactions
// @Summary List all transactions with pagination
// @Description Get a paginated list of all transactions
// @Tags Transactions
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10, max: 100)"
// @Success 200 {object} response.Response
// @Router /transactions [get]
func (t *TransactionController) ListTransactions(c *gin.Context) {
	var query dto.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	// Set defaults
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 10
	}
	if query.Limit > 100 {
		query.Limit = 100
	}

	offset := (query.Page - 1) * query.Limit

	txs, err := t.transactionRepo.ListWithPagination(offset, query.Limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve transactions")
		return
	}

	total, err := t.transactionRepo.CountAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to count transactions")
		return
	}

	response.Success(c, http.StatusOK, "Transactions retrieved successfully", gin.H{
		"data":  txs,
		"total": total,
		"page":  query.Page,
		"limit": query.Limit,
	})
}