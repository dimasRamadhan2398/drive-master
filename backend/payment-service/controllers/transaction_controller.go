package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"payment-service/models"
	"payment-service/models/dto"
	"payment-service/pkg/kafka"
	"payment-service/pkg/response"
	"payment-service/repositories"
	"payment-service/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionController struct {
	transactionRepo   repositories.ITransactionRepository
	paymentRepo       repositories.IPaymentRepository
	paymentMethodRepo repositories.IPaymentMethodRepository
	midtransSvc       services.IMidtransService
	eventPublisher    kafka.IEventPublisher
}

type ITransactionController interface {
	GetTransaction(c *gin.Context)
	GetTransactionsByPaymentID(c *gin.Context)
	ListTransactions(c *gin.Context)
	CreateTransaction(c *gin.Context)
	Callback(c *gin.Context)
	GetPayment(c *gin.Context)
	GetPaymentByOrderID(c *gin.Context)
	GetPaymentDetail(c *gin.Context)
	GetPaymentStatus(c *gin.Context)
	ListPayments(c *gin.Context)
}

func NewTransactionController(
	transactionRepo repositories.ITransactionRepository,
	paymentRepo repositories.IPaymentRepository,
	paymentMethodRepo repositories.IPaymentMethodRepository,
	midtransSvc services.IMidtransService,
	eventPublisher kafka.IEventPublisher,
) ITransactionController {
	return &TransactionController{
		transactionRepo:   transactionRepo,
		paymentRepo:       paymentRepo,
		paymentMethodRepo: paymentMethodRepo,
		midtransSvc:       midtransSvc,
		eventPublisher:    eventPublisher,
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

type CreateTransactionRequest struct {
	EnrollmentID  string  `json:"enrollmentId" binding:"required"`
	UserID        string  `json:"userId"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"paymentMethod" binding:"required"`
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func (t *TransactionController) CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	enrollmentUUID, err := uuid.Parse(req.EnrollmentID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid enrollment ID format")
		return
	}

	userCtxID, exists := c.Get("user_id")
	var userUUID uuid.UUID
	if exists {
		if uStr, ok := userCtxID.(string); ok {
			userUUID, _ = uuid.Parse(uStr)
		}
	}
	if userUUID == uuid.Nil && req.UserID != "" {
		userUUID, _ = uuid.Parse(req.UserID)
	}
	if userUUID == uuid.Nil {
		response.Error(c, http.StatusUnauthorized, "User context not found")
		return
	}

	code := req.PaymentMethod
	if code == "e_wallet" {
		code = "ewallet"
	}
	pm, err := t.paymentMethodRepo.GetByCode(code)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid or unsupported payment method")
		return
	}

	// Call booking-service
	bookingServiceURL := getEnv("BOOKING_SERVICE_URL", "http://booking-service:8003")
	url := fmt.Sprintf("%s/api/v1/enrollments/%s", bookingServiceURL, req.EnrollmentID)

	httpClient := &http.Client{Timeout: 10 * time.Second}
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to build booking-service request")
		return
	}
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		httpReq.Header.Set("Authorization", authHeader)
	}
	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to call booking-service: "+err.Error())
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		response.Error(c, httpResp.StatusCode, "Booking service returned error status")
		return
	}

	var enrollmentResp struct {
		Enrollment struct {
			ID         string  `json:"id"`
			UserID     string  `json:"userId"`
			PackageID  string  `json:"packageId"`
			Status     string  `json:"status"`
			TotalPrice float64 `json:"totalPrice"`
		} `json:"enrollment"`
	}
	if err := json.NewDecoder(httpResp.Body).Decode(&enrollmentResp); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to decode enrollment details: "+err.Error())
		return
	}
	enrollment := enrollmentResp.Enrollment

	// Call core-service
	packageName := "Driving Package"
	coreServiceURL := getEnv("CORE_SERVICE_URL", "http://core-service:8002")
	if enrollment.PackageID != "" {
		pkgUrl := fmt.Sprintf("%s/api/v1/packages/%s", coreServiceURL, enrollment.PackageID)
		pkgReq, err := http.NewRequest("GET", pkgUrl, nil)
		if err == nil {
			if authHeader != "" {
				pkgReq.Header.Set("Authorization", authHeader)
			}
			pkgResp, err := httpClient.Do(pkgReq)
			if err == nil {
				defer pkgResp.Body.Close()
				if pkgResp.StatusCode == http.StatusOK {
					var pkgData struct {
						Success bool `json:"success"`
						Data    struct {
							Name string `json:"name"`
						} `json:"data"`
					}
					if err := json.NewDecoder(pkgResp.Body).Decode(&pkgData); err == nil && pkgData.Success {
						packageName = pkgData.Data.Name
					}
				}
			}
		}
	}

	// Call user-service
	customerName := "Customer"
	customerEmail := "customer@example.com"
	userServiceURL := getEnv("USER_SERVICE_URL", "http://user-service:8001")
	if enrollment.UserID != "" {
		userUrl := fmt.Sprintf("%s/api/v1/users/%s", userServiceURL, enrollment.UserID)
		userReq, err := http.NewRequest("GET", userUrl, nil)
		if err == nil {
			if authHeader != "" {
				userReq.Header.Set("Authorization", authHeader)
			}
			userResp, err := httpClient.Do(userReq)
			if err == nil {
				defer userResp.Body.Close()
				if userResp.StatusCode == http.StatusOK {
					var userEnvelope struct {
						Success bool `json:"success"`
						Data    struct {
							Email       string `json:"email"`
							FirstName   string `json:"firstName"`
							LastName    string `json:"lastName"`
						} `json:"data"`
					}
					if err := json.NewDecoder(userResp.Body).Decode(&userEnvelope); err == nil && userEnvelope.Success {
						customerEmail = userEnvelope.Data.Email
						customerName = fmt.Sprintf("%s %s", userEnvelope.Data.FirstName, userEnvelope.Data.LastName)
						if customerName == " " || customerName == "" {
							customerName = "Customer"
						}
					}
				}
			}
		}
	}

	orderID := fmt.Sprintf("ORD-%s-%d", time.Now().Format("20060102150405"), uuid.New().ID())
	snapResp, err := t.midtransSvc.CreateSnapTransaction(orderID, enrollment.TotalPrice, packageName, customerName, customerEmail)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to initiate Midtrans checkout: "+err.Error())
		return
	}

	// Create payment
	payment := &models.Payment{
		ID:                 uuid.New(),
		OrderID:            orderID,
		BookingID:         &enrollmentUUID,
		UserID:            userUUID,
		Amount:            enrollment.TotalPrice,
		Currency:          "IDR",
		Status:            models.PaymentStatusPending,
		PaymentMethodID:   &pm.ID,
		Gateway:           "midtrans",
		GatewayOrderID:    orderID,
		GatewayPaymentURL: snapResp.RedirectURL,
		Description:       fmt.Sprintf("Pembelian %s", packageName),
		ExpiryTime:        func() *time.Time { t := time.Now().Add(24 * time.Hour); return &t }(),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	if err := t.paymentRepo.Create(payment); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create payment record: "+err.Error())
		return
	}

	// Create transaction
	tx := &models.Transaction{
		ID:              uuid.New(),
		PaymentID:       payment.ID,
		Type:            models.TransactionTypeCharge,
		Status:          models.TransactionStatusPending,
		Amount:          enrollment.TotalPrice,
		Currency:        "IDR",
		Gateway:         "midtrans",
		GatewayTxnID:    snapResp.Token,
		PaymentMethodID: &pm.ID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	if err := t.transactionRepo.Create(tx); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create transaction record: "+err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Payment transaction initiated successfully", gin.H{
		"id":            payment.ID.String(),
		"enrollmentId":  req.EnrollmentID,
		"userId":        userUUID.String(),
		"orderId":       orderID,
		"amount":        payment.Amount,
		"paymentMethod": req.PaymentMethod,
		"status":        "pending",
		"paymentUrl":    snapResp.RedirectURL,
		"createdAt":     payment.CreatedAt.Format(time.RFC3339),
		"updatedAt":     payment.UpdatedAt.Format(time.RFC3339),
	})
}

func (t *TransactionController) Callback(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid payload: "+err.Error())
		return
	}

	notification, err := t.midtransSvc.ParseNotification(payload)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to parse Midtrans notification: "+err.Error())
		return
	}

	// Look up payment by order id
	payment, err := t.paymentRepo.GetByOrderID(notification.OrderID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Payment record not found")
		return
	}

	status := services.MapMidtransStatusToPaymentStatus(notification.TransactionStatus)
	payment.Status = status
	if status == models.PaymentStatusSuccess {
		now := time.Now()
		payment.PaidAt = &now
	}
	payment.UpdatedAt = time.Now()
	_ = t.paymentRepo.Update(payment)

	// Update associated transactions
	txs, err := t.transactionRepo.GetByPaymentID(payment.ID)
	if err == nil && len(txs) > 0 {
		for i := range txs {
			txs[i].Status = models.TransactionStatus(status)
			txs[i].GatewayTxnID = notification.TransactionID
			now := time.Now()
			txs[i].ProcessedAt = &now
			txs[i].UpdatedAt = now
			_ = t.transactionRepo.Update(&txs[i])
		}
	}

	// Publish Kafka event transaction.paid if payment is successful
	if status == models.PaymentStatusSuccess {
		event := map[string]interface{}{
			"id":        uuid.New().String(),
			"type":      "transaction.paid",
			"timestamp": time.Now().Format(time.RFC3339),
			"user_id":   payment.UserID.String(),
			"data": map[string]interface{}{
				"enrollment_id":     payment.BookingID.String(),
				"total_price":       payment.Amount,
				"payment_method_id": payment.PaymentMethodID,
			},
			"success": true,
		}

		err := t.eventPublisher.Publish(payment.ID.String(), event)
		if err != nil {
			fmt.Printf("[Callback] Error publishing kafka event: %v\n", err)
		} else {
			fmt.Printf("[Callback] Successfully published transaction.paid event for enrollment %s\n", payment.BookingID.String())
		}
	}

	response.Success(c, http.StatusOK, "Callback handled successfully", nil)
}

func (t *TransactionController) GetPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid payment ID format")
		return
	}

	payment, err := t.paymentRepo.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Payment not found")
		return
	}

	response.Success(c, http.StatusOK, "Payment retrieved successfully", payment)
}

func (t *TransactionController) GetPaymentByOrderID(c *gin.Context) {
	orderID := c.Param("orderId")
	payment, err := t.paymentRepo.GetByOrderID(orderID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Payment not found")
		return
	}

	response.Success(c, http.StatusOK, "Payment retrieved successfully", payment)
}

func (t *TransactionController) GetPaymentDetail(c *gin.Context) {
	orderID := c.Param("orderId")
	payment, err := t.paymentRepo.GetByOrderID(orderID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Payment not found")
		return
	}

	response.Success(c, http.StatusOK, "Payment detail retrieved successfully", gin.H{
		"orderId":       payment.OrderID,
		"paymentType":   payment.PaymentMethod.Code,
		"transactionId": payment.GatewayOrderID,
		"grossAmount":   payment.Amount,
		"status":        payment.Status,
		"paymentUrl":    payment.GatewayPaymentURL,
		"expiryTime":    payment.ExpiryTime,
	})
}

func (t *TransactionController) GetPaymentStatus(c *gin.Context) {
	orderID := c.Param("orderId")
	payment, err := t.paymentRepo.GetByOrderID(orderID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Payment not found")
		return
	}

	response.Success(c, http.StatusOK, "Payment status retrieved successfully", gin.H{
		"status": payment.Status,
	})
}

func (t *TransactionController) ListPayments(c *gin.Context) {
	var query dto.ListPaymentsRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 10
	}
	if query.PageSize > 100 {
		query.PageSize = 100
	}

	filter := repositories.ListPaymentsFilter{
		UserID: query.UserID,
		Status: query.Status,
	}

	payments, total, err := t.paymentRepo.List(filter, query.Page, query.PageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve payments: "+err.Error())
		return
	}

	var data []dto.GetPaymentResponse
	for _, p := range payments {
		var transactionsDTO []dto.TransactionDTO
		for _, tx := range p.Transactions {
			transactionsDTO = append(transactionsDTO, dto.TransactionDTO{
				ID:           tx.ID,
				PaymentID:    tx.PaymentID,
				Type:         tx.Type,
				Status:       tx.Status,
				Amount:       tx.Amount,
				Currency:     tx.Currency,
				Gateway:      tx.Gateway,
				GatewayTxnID: tx.GatewayTxnID,
				ErrorCode:    tx.ErrorCode,
				ErrorMessage: tx.ErrorMessage,
				ProcessedAt:  tx.ProcessedAt,
				CreatedAt:    tx.CreatedAt,
			})
		}

		var pmDTO *dto.PaymentMethodDTO
		if p.PaymentMethod != nil {
			pmDTO = &dto.PaymentMethodDTO{
				ID:          p.PaymentMethod.ID,
				Code:        p.PaymentMethod.Code,
				Name:        p.PaymentMethod.Name,
				Description: p.PaymentMethod.Description,
				Icon:        p.PaymentMethod.Icon,
				IsActive:    p.PaymentMethod.IsActive,
				Gateway:     p.PaymentMethod.Gateway,
				FeeType:     p.PaymentMethod.FeeType,
				FeeAmount:   p.PaymentMethod.FeeAmount,
			}
		}

		data = append(data, dto.GetPaymentResponse{
			ID:                 p.ID,
			OrderID:            p.OrderID,
			BookingID:          p.BookingID,
			UserID:             p.UserID,
			Amount:             p.Amount,
			Currency:           p.Currency,
			Status:             p.Status,
			PaymentMethod:      pmDTO,
			Gateway:            p.Gateway,
			GatewayOrderID:     p.GatewayOrderID,
			GatewayPaymentURL:  p.GatewayPaymentURL,
			VaNumber:           p.VaNumber,
			QrCodeURL:          p.QrCodeURL,
			Description:        p.Description,
			ExpiryTime:         p.ExpiryTime,
			PaidAt:             p.PaidAt,
			CreatedAt:          p.CreatedAt,
			UpdatedAt:          p.UpdatedAt,
			Transactions:       transactionsDTO,
		})
	}

	totalPages := int((total + int64(query.PageSize) - 1) / int64(query.PageSize))

	response.Success(c, http.StatusOK, "Payments retrieved successfully", gin.H{
		"data": data,
		"pagination": gin.H{
			"page":       query.Page,
			"limit":      query.PageSize,
			"total":      total,
			"totalPages": totalPages,
		},
	})
}