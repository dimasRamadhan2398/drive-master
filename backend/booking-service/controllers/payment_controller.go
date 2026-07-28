package controllers

import (
	"net/http"
	"strconv"

	"booking-service/models/dto"
	"booking-service/pkg/base"
	"booking-service/services"
	"booking-service/pkg/response"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	paymentService services.IPaymentService
}

func NewPaymentController(paymentService services.IPaymentService) *PaymentController {
	return &PaymentController{
		paymentService: paymentService,
	}
}

// CreatePayment godoc
// @Summary Create a new payment
// @Description Creates a new payment for an enrollment
// @Tags payments
// @Accept json
// @Produce json
// @Param payment body dto.CreatePaymentRequest true "Payment details"
// @Success 201 {object} dto.PaymentResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payments [post]
func (c *PaymentController) CreatePayment(ctx *gin.Context) {
	var req dto.CreatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.paymentService.CreatePayment(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response.Created(ctx, "Payment created successfully", resp)
}

// GetPayment godoc
// @Summary Get payment by ID
// @Description Retrieves a payment by its ID
// @Tags payments
// @Accept json
// @Produce json
// @Param id path int true "Payment ID"
// @Success 200 {object} dto.PaymentResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payments/{id} [get]
func (c *PaymentController) GetPayment(ctx *gin.Context) {
	id, err := base.GetUintIDFromPath(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.paymentService.GetPayment(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response.OK(ctx, "Payment retrieved successfully", resp)
}

// GetPaymentByOrderID godoc
// @Summary Get payment by order ID
// @Description Retrieves a payment by its Midtrans order ID
// @Tags payments
// @Accept json
// @Produce json
// @Param orderId path string true "Order ID"
// @Success 200 {object} dto.PaymentResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payments/order/{orderId} [get]
func (c *PaymentController) GetPaymentByOrderID(ctx *gin.Context) {
	orderID := ctx.Param("orderId")
	if orderID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "order ID is required"})
		return
	}

	resp, err := c.paymentService.GetPaymentByOrderID(ctx.Request.Context(), orderID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response.OK(ctx, "Payment retrieved successfully", resp)
}

// GetPaymentDetail godoc
// @Summary Get payment detail by order ID
// @Description Retrieves detailed payment information by order ID
// @Tags payments
// @Accept json
// @Produce json
// @Param orderId path string true "Order ID"
// @Success 200 {object} dto.PaymentDetailResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payments/order/{orderId}/detail [get]
func (c *PaymentController) GetPaymentDetail(ctx *gin.Context) {
	orderID := ctx.Param("orderId")
	if orderID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "order ID is required"})
		return
	}

	resp, err := c.paymentService.GetPaymentDetail(ctx.Request.Context(), orderID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response.OK(ctx, "Payment detail retrieved successfully", resp)
}

// ListPayments godoc
// @Summary List all payments
// @Description Retrieves a paginated list of all payments
// @Tags payments
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} dto.PaymentListResponse
// @Failure 500 {object} map[string]string
// @Router /payments [get]
func (c *PaymentController) ListPayments(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	resp, err := c.paymentService.ListPayments(ctx.Request.Context(), page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response.OK(ctx, "Payments retrieved successfully", resp)
}

// ListUserPayments godoc
// @Summary List user payments
// @Description Retrieves payments for a specific user
// @Tags payments
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} dto.PaymentListResponse
// @Failure 500 {object} map[string]string
// @Router /payments/user/{userId} [get]
func (c *PaymentController) ListUserPayments(ctx *gin.Context) {
	userID, err := base.GetUUIDIDFromPath(ctx, "userId")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	resp, err := c.paymentService.ListUserPayments(ctx.Request.Context(), userID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response.OK(ctx, "User payments retrieved successfully", resp)
}

// HandleCallback godoc
// @Summary Handle Midtrans callback
// @Description Handles payment callback notifications from Midtrans
// @Tags payments
// @Accept x-www-form-urlencoded
// @Produce json
// @Param callback body dto.PaymentCallbackRequest true "Midtrans callback"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payments/callback [post]
func (c *PaymentController) HandleCallback(ctx *gin.Context) {
	var callback dto.PaymentCallbackRequest
	if err := ctx.ShouldBind(&callback); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.paymentService.HandleCallback(ctx.Request.Context(), callback); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response.OK(ctx, "Callback processed successfully", gin.H{"status": "success"})
}

// CancelPayment godoc
// @Summary Cancel a payment
// @Description Cancels a pending payment
// @Tags payments
// @Accept json
// @Produce json
// @Param orderId path string true "Order ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payments/order/{orderId}/cancel [post]
func (c *PaymentController) CancelPayment(ctx *gin.Context) {
	orderID := ctx.Param("orderId")
	if orderID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "order ID is required"})
		return
	}

	if err := c.paymentService.CancelPayment(ctx.Request.Context(), orderID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response.OK(ctx, "Payment cancelled successfully", gin.H{"status": "cancelled"})
}

// IPaymentController interface for controller registry
type IPaymentController interface {
	CreatePayment(ctx *gin.Context)
	GetPayment(ctx *gin.Context)
	GetPaymentByOrderID(ctx *gin.Context)
	GetPaymentDetail(ctx *gin.Context)
	ListPayments(ctx *gin.Context)
	ListUserPayments(ctx *gin.Context)
	HandleCallback(ctx *gin.Context)
	CancelPayment(ctx *gin.Context)
}