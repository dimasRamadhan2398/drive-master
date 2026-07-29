package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	paymentGateway    services.IPaymentGatewayService
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
	SimulatePayment(c *gin.Context)
	ListPayments(c *gin.Context)
}

func NewTransactionController(
	transactionRepo repositories.ITransactionRepository,
	paymentRepo repositories.IPaymentRepository,
	paymentMethodRepo repositories.IPaymentMethodRepository,
	paymentGateway services.IPaymentGatewayService,
	eventPublisher kafka.IEventPublisher,
) ITransactionController {
	return &TransactionController{
		transactionRepo:   transactionRepo,
		paymentRepo:       paymentRepo,
		paymentMethodRepo: paymentMethodRepo,
		paymentGateway:    paymentGateway,
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
	Enrollment    *struct {
		ID         string  `json:"id"`
		UserID     string  `json:"userId"`
		PackageID  string  `json:"packageId"`
		Status     string  `json:"status"`
		TotalPrice float64 `json:"totalPrice"`
		Price      float64 `json:"price"`
	} `json:"enrollment,omitempty"`
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
	bookingServiceURL := getEnv("BOOKING_SERVICE_URL", "http://127.0.0.1:8003")
	if bookingServiceURL == "http://booking-service:8003" {
		bookingServiceURL = "http://127.0.0.1:8003"
	}
	url := fmt.Sprintf("%s/api/v1/enrollments/%s", bookingServiceURL, req.EnrollmentID)

	httpClient := &http.Client{Timeout: 10 * time.Second}
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to build booking-service request: "+err.Error())
		return
	}
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		httpReq.Header.Set("Authorization", authHeader)
	}

	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		altURL := fmt.Sprintf("http://127.0.0.1:9003/api/v1/enrollments/%s", req.EnrollmentID)
		if altReq, altErr := http.NewRequest("GET", altURL, nil); altErr == nil {
			if authHeader != "" {
				altReq.Header.Set("Authorization", authHeader)
			}
			if altResp, altErr2 := httpClient.Do(altReq); altErr2 == nil {
				httpResp = altResp
				err = nil
			}
		}
	}

	type EnrollmentItem struct {
		ID         string  `json:"id"`
		UserID     string  `json:"userId"`
		PackageID  string  `json:"packageId"`
		Status     string  `json:"status"`
		TotalPrice float64 `json:"totalPrice"`
	}
	var enrollment EnrollmentItem

	if err != nil {
		if req.Enrollment != nil && req.Enrollment.ID != "" {
			price := req.Enrollment.TotalPrice
			if price == 0 {
				price = req.Enrollment.Price
			}
			if price == 0 {
				price = req.Amount
			}
			enrollment = EnrollmentItem{
				ID:         req.Enrollment.ID,
				UserID:     req.Enrollment.UserID,
				PackageID:  req.Enrollment.PackageID,
				Status:     req.Enrollment.Status,
				TotalPrice: price,
			}
		} else if req.Amount > 0 {
			enrollment = EnrollmentItem{
				ID:         req.EnrollmentID,
				UserID:     req.UserID,
				TotalPrice: req.Amount,
			}
		} else {
			response.Error(c, http.StatusInternalServerError, "Failed to call booking-service: "+err.Error())
			return
		}
	} else {
		defer httpResp.Body.Close()
		bodyBytes, readErr := io.ReadAll(httpResp.Body)

		if httpResp.StatusCode != http.StatusOK {
			errMsg := ""
			if readErr == nil && len(bodyBytes) > 0 {
				var errJson struct {
					Error   string `json:"error"`
					Message string `json:"message"`
				}
				if err := json.Unmarshal(bodyBytes, &errJson); err == nil {
					if errJson.Error != "" {
						errMsg = errJson.Error
					} else if errJson.Message != "" {
						errMsg = errJson.Message
					}
				}
				if errMsg == "" {
					errMsg = string(bodyBytes)
				}
			}
			if errMsg == "" {
				errMsg = fmt.Sprintf("HTTP status %d", httpResp.StatusCode)
			}

			if req.Enrollment != nil && req.Enrollment.ID != "" {
				price := req.Enrollment.TotalPrice
				if price == 0 {
					price = req.Enrollment.Price
				}
				if price == 0 {
					price = req.Amount
				}
				enrollment = EnrollmentItem{
					ID:         req.Enrollment.ID,
					UserID:     req.Enrollment.UserID,
					PackageID:  req.Enrollment.PackageID,
					Status:     req.Enrollment.Status,
					TotalPrice: price,
				}
			} else if req.Amount > 0 {
				enrollment = EnrollmentItem{
					ID:         req.EnrollmentID,
					UserID:     req.UserID,
					TotalPrice: req.Amount,
				}
			} else {
				response.Error(c, httpResp.StatusCode, "Booking service error: "+errMsg)
				return
			}
		} else {
			var parsed struct {
				Enrollment EnrollmentItem `json:"enrollment"`
				Data       struct {
					Enrollment EnrollmentItem `json:"enrollment"`
					ID         string         `json:"id"`
					UserID     string         `json:"userId"`
					PackageID  string         `json:"packageId"`
					Status     string         `json:"status"`
					TotalPrice float64        `json:"totalPrice"`
				} `json:"data"`
				ID         string  `json:"id"`
				UserID     string  `json:"userId"`
				PackageID  string  `json:"packageId"`
				Status     string  `json:"status"`
				TotalPrice float64 `json:"totalPrice"`
			}
			if err := json.Unmarshal(bodyBytes, &parsed); err == nil {
				if parsed.Enrollment.ID != "" {
					enrollment = parsed.Enrollment
				} else if parsed.Data.Enrollment.ID != "" {
					enrollment = parsed.Data.Enrollment
				} else if parsed.Data.ID != "" {
					enrollment = EnrollmentItem{
						ID:         parsed.Data.ID,
						UserID:     parsed.Data.UserID,
						PackageID:  parsed.Data.PackageID,
						Status:     parsed.Data.Status,
						TotalPrice: parsed.Data.TotalPrice,
					}
				} else if parsed.ID != "" {
					enrollment = EnrollmentItem{
						ID:         parsed.ID,
						UserID:     parsed.UserID,
						PackageID:  parsed.PackageID,
						Status:     parsed.Status,
						TotalPrice: parsed.TotalPrice,
					}
				}
			}
			if enrollment.TotalPrice == 0 {
				if req.Amount > 0 {
					enrollment.TotalPrice = req.Amount
				} else if req.Enrollment != nil && req.Enrollment.TotalPrice > 0 {
					enrollment.TotalPrice = req.Enrollment.TotalPrice
				} else if req.Enrollment != nil && req.Enrollment.Price > 0 {
					enrollment.TotalPrice = req.Enrollment.Price
				}
			}

			// Validate user ownership: prevent creating a payment for another user's enrollment
			if enrollment.UserID != "" && userUUID != uuid.Nil {
				if enrollmentUserUUID, err := uuid.Parse(enrollment.UserID); err == nil && enrollmentUserUUID != userUUID {
					response.Error(c, http.StatusBadRequest, fmt.Sprintf("Enrollment %s belongs to a different user (%s)", enrollment.ID, enrollment.UserID))
					return
				}
			}
		}
	}

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

	// Build metadata containing package & payment info for downstream use (e.g. sale creation)
	type paymentMeta struct {
		PackageID     string `json:"package_id"`
		PackageName   string `json:"package_name"`
		PaymentMethod string `json:"payment_method"`
	}
	metaBytes, _ := json.Marshal(paymentMeta{
		PackageID:     enrollment.PackageID,
		PackageName:   packageName,
		PaymentMethod: req.PaymentMethod,
	})
	metadataStr := string(metaBytes)

	nowTime := time.Now()

	bypass := os.Getenv("PAYMENT_BYPASS") == "true" || os.Getenv("PAYMENT_GATEWAY") == "bypass" || os.Getenv("PAYMENT_GATEWAY") == "BYPASS"

	var gatewayPaymentURL string = ""
	var gatewayTxnID string = orderID
	var gatewayName string = "bypass"
	var paymentStatus models.PaymentStatus = models.PaymentStatusSuccess
	var transactionStatus models.TransactionStatus = models.TransactionStatusSuccess
	var paidAt *time.Time = &nowTime
	var processedAt *time.Time = &nowTime

	if !bypass && t.paymentGateway != nil {
		checkoutResp, err := t.paymentGateway.CreateCheckout(orderID, enrollment.TotalPrice, packageName, customerName, customerEmail)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to initiate payment checkout: "+err.Error())
			return
		}
		gatewayPaymentURL = checkoutResp.RedirectURL
		gatewayTxnID = checkoutResp.Token
		gatewayName = t.paymentGateway.GetName()
		paymentStatus = models.PaymentStatusPending
		transactionStatus = models.TransactionStatusPending
		paidAt = nil
		processedAt = nil
	}

	// Create payment
	payment := &models.Payment{
		ID:                 uuid.New(),
		OrderID:            orderID,
		BookingID:          &enrollmentUUID,
		UserID:             userUUID,
		Amount:             enrollment.TotalPrice,
		Currency:           "IDR",
		Status:             paymentStatus,
		PaymentMethodID:    &pm.ID,
		Gateway:            gatewayName,
		GatewayOrderID:     orderID,
		GatewayPaymentURL:  gatewayPaymentURL,
		Metadata:           metadataStr,
		Description:        fmt.Sprintf("Pembelian %s", packageName),
		PaidAt:             paidAt,
		ExpiryTime:         func() *time.Time { t := time.Now().Add(24 * time.Hour); return &t }(),
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
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
		Status:          transactionStatus,
		Amount:          enrollment.TotalPrice,
		Currency:        "IDR",
		Gateway:         gatewayName,
		GatewayTxnID:    gatewayTxnID,
		GatewayResponse: "{}",
		PaymentMethodID: &pm.ID,
		ProcessedAt:     processedAt,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	if err := t.transactionRepo.Create(tx); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create transaction record: "+err.Error())
		return
	}

	if bypass {
		// Publish Kafka event transaction.paid immediately to update enrollment and grant user entitlement
		event := map[string]interface{}{
			"id":        uuid.New().String(),
			"type":      "transaction.paid",
			"timestamp": time.Now().Format(time.RFC3339),
			"user_id":   payment.UserID.String(),
			"data": map[string]interface{}{
				"enrollment_id":  payment.BookingID.String(),
				"total_price":    payment.Amount,
				"payment_method": req.PaymentMethod,
				"package_id":     enrollment.PackageID,
				"package_name":   packageName,
			},
			"success": true,
		}

		err = t.eventPublisher.Publish(payment.ID.String(), event)
		if err != nil {
			fmt.Printf("[CreateTransaction Bypass] Error publishing kafka event: %v\n", err)
		} else {
			fmt.Printf("[CreateTransaction Bypass] Successfully published transaction.paid event for enrollment %s\n", payment.BookingID.String())
		}

		response.Success(c, http.StatusCreated, "Payment transaction bypassed successfully", gin.H{
			"id":            payment.ID.String(),
			"enrollmentId":  req.EnrollmentID,
			"userId":        userUUID.String(),
			"orderId":       orderID,
			"amount":        payment.Amount,
			"paymentMethod": req.PaymentMethod,
			"status":        "success",
			"paymentUrl":    "",
			"createdAt":     payment.CreatedAt.Format(time.RFC3339),
			"updatedAt":     payment.UpdatedAt.Format(time.RFC3339),
		})
	} else {
		response.Success(c, http.StatusCreated, "Payment transaction initiated successfully", gin.H{
			"id":            payment.ID.String(),
			"enrollmentId":  req.EnrollmentID,
			"userId":        userUUID.String(),
			"orderId":       orderID,
			"amount":        payment.Amount,
			"paymentMethod": req.PaymentMethod,
			"status":        "pending",
			"paymentUrl":    gatewayPaymentURL,
			"createdAt":     payment.CreatedAt.Format(time.RFC3339),
			"updatedAt":     payment.UpdatedAt.Format(time.RFC3339),
		})
	}
}

func (t *TransactionController) Callback(c *gin.Context) {
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to read request body")
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody))

	headers := make(map[string]string)
	for k, v := range c.Request.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	notification, err := t.paymentGateway.VerifyNotification(headers, rawBody)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to verify notification: "+err.Error())
		return
	}

	// Look up payment by order id
	payment, err := t.paymentRepo.GetByOrderID(notification.OrderID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Payment record not found")
		return
	}

	status := notification.TransactionStatus
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
		// Parse metadata to extract package & payment method info stored at creation time
		var meta struct {
			PackageID   string `json:"package_id"`
			PackageName string `json:"package_name"`
			PaymentMethod string `json:"payment_method"`
		}
		if payment.Metadata != "" && payment.Metadata != "{}" {
			_ = json.Unmarshal([]byte(payment.Metadata), &meta)
		}

		event := map[string]interface{}{
			"id":        uuid.New().String(),
			"type":      "transaction.paid",
			"timestamp": time.Now().Format(time.RFC3339),
			"user_id":   payment.UserID.String(),
			"data": map[string]interface{}{
				"enrollment_id":  payment.BookingID.String(),
				"total_price":    payment.Amount,
				"payment_method": meta.PaymentMethod,
				"package_id":     meta.PackageID,
				"package_name":   meta.PackageName,
			},
			"success": true,
		}

		// Direct HTTP call to booking-service as instant fallback
		if payment.BookingID != nil {
			bookingServiceURL := getEnv("BOOKING_SERVICE_URL", "http://127.0.0.1:8003")
			if bookingServiceURL == "http://booking-service:8003" {
				bookingServiceURL = "http://127.0.0.1:8003"
			}
			payUrl := fmt.Sprintf("%s/api/v1/enrollments/%s/pay", bookingServiceURL, payment.BookingID.String())
			payBody, _ := json.Marshal(map[string]interface{}{
				"totalPrice": payment.Amount,
			})
			payReq, payErr := http.NewRequest("POST", payUrl, bytes.NewBuffer(payBody))
			if payErr == nil {
				payReq.Header.Set("Content-Type", "application/json")
				authHeader := c.GetHeader("Authorization")
				if authHeader != "" {
					payReq.Header.Set("Authorization", authHeader)
				}
				httpClient := &http.Client{Timeout: 10 * time.Second}
				payResp, payDoErr := httpClient.Do(payReq)
				if payDoErr == nil {
					payResp.Body.Close()
				}
			}
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
	orderID := c.Param("id")
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
	orderID := c.Param("id")
	payment, err := t.paymentRepo.GetByOrderID(orderID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Payment not found")
		return
	}

	// If the payment is still pending in our DB, check the gateway directly as a fallback
	if payment.Status == models.PaymentStatusPending && t.paymentGateway != nil {
		status, err := t.paymentGateway.GetTransactionStatus(orderID)
		if err == nil {
			if status != models.PaymentStatusPending {
				// Status changed! Update payment record
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
						txs[i].GatewayTxnID = orderID // CheckTransaction doesn't return full ID here, but orderID is safe fallback
						now := time.Now()
						txs[i].ProcessedAt = &now
						txs[i].UpdatedAt = now
						_ = t.transactionRepo.Update(&txs[i])
					}
				}

				// Publish Kafka event transaction.paid if payment is successful
				if status == models.PaymentStatusSuccess {
					// Parse metadata to extract package & payment method info stored at creation time
					var meta struct {
						PackageID   string `json:"package_id"`
						PackageName string `json:"package_name"`
						PaymentMethod string `json:"payment_method"`
					}
					if payment.Metadata != "" && payment.Metadata != "{}" {
						_ = json.Unmarshal([]byte(payment.Metadata), &meta)
					}

					event := map[string]interface{}{
						"id":        uuid.New().String(),
						"type":      "transaction.paid",
						"timestamp": time.Now().Format(time.RFC3339),
						"user_id":   payment.UserID.String(),
						"data": map[string]interface{}{
							"enrollment_id":  payment.BookingID.String(),
							"total_price":    payment.Amount,
							"payment_method": meta.PaymentMethod,
							"package_id":     meta.PackageID,
							"package_name":   meta.PackageName,
						},
						"success": true,
					}

					// Direct HTTP call to booking-service as instant fallback
					if payment.BookingID != nil {
						bookingServiceURL := getEnv("BOOKING_SERVICE_URL", "http://127.0.0.1:8003")
						if bookingServiceURL == "http://booking-service:8003" {
							bookingServiceURL = "http://127.0.0.1:8003"
						}
						payUrl := fmt.Sprintf("%s/api/v1/enrollments/%s/pay", bookingServiceURL, payment.BookingID.String())
						payBody, _ := json.Marshal(map[string]interface{}{
							"totalPrice": payment.Amount,
						})
						payReq, payErr := http.NewRequest("POST", payUrl, bytes.NewBuffer(payBody))
						if payErr == nil {
							payReq.Header.Set("Content-Type", "application/json")
							authHeader := c.GetHeader("Authorization")
							if authHeader != "" {
								payReq.Header.Set("Authorization", authHeader)
							}
							httpClient := &http.Client{Timeout: 10 * time.Second}
							payResp, payDoErr := httpClient.Do(payReq)
							if payDoErr == nil {
								payResp.Body.Close()
							}
						}
					}

					err := t.eventPublisher.Publish(payment.ID.String(), event)
					if err != nil {
						fmt.Printf("[Status Check Fallback] Error publishing kafka event: %v\n", err)
					} else {
						fmt.Printf("[Status Check Fallback] Successfully published transaction.paid event for enrollment %s\n", payment.BookingID.String())
					}
				}
			}
		}
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

func (t *TransactionController) SimulatePayment(c *gin.Context) {
	orderID := c.Param("orderId")
	if orderID == "" {
		orderID = c.Param("id")
	}
	if orderID == "" {
		var req struct {
			OrderID string `json:"orderId"`
		}
		_ = c.ShouldBindJSON(&req)
		orderID = req.OrderID
	}

	if orderID == "" {
		response.Error(c, http.StatusBadRequest, "Order ID is required")
		return
	}

	payment, err := t.paymentRepo.GetByOrderID(orderID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Payment record not found")
		return
	}

	if t.paymentGateway != nil {
		_ = t.paymentGateway.SimulatePayment(orderID, payment.Amount)
	}

	now := time.Now()
	payment.Status = models.PaymentStatusSuccess
	payment.PaidAt = &now
	payment.UpdatedAt = now
	_ = t.paymentRepo.Update(payment)

	txs, err := t.transactionRepo.GetByPaymentID(payment.ID)
	if err == nil && len(txs) > 0 {
		for i := range txs {
			txs[i].Status = models.TransactionStatusSuccess
			txs[i].GatewayTxnID = orderID
			txs[i].ProcessedAt = &now
			txs[i].UpdatedAt = now
			_ = t.transactionRepo.Update(&txs[i])
		}
	}

	if payment.BookingID != nil {
		bookingServiceURL := getEnv("BOOKING_SERVICE_URL", "http://127.0.0.1:8003")
		if bookingServiceURL == "http://booking-service:8003" {
			bookingServiceURL = "http://127.0.0.1:8003"
		}
		payUrl := fmt.Sprintf("%s/api/v1/enrollments/%s/pay", bookingServiceURL, payment.BookingID.String())
		payBody, _ := json.Marshal(map[string]interface{}{
			"totalPrice": payment.Amount,
		})
		payReq, payErr := http.NewRequest("POST", payUrl, bytes.NewBuffer(payBody))
		if payErr == nil {
			payReq.Header.Set("Content-Type", "application/json")
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				payReq.Header.Set("Authorization", authHeader)
			}
			httpClient := &http.Client{Timeout: 10 * time.Second}
			payResp, payDoErr := httpClient.Do(payReq)
			if payDoErr == nil {
				payResp.Body.Close()
			}
		}
	}

	var meta struct {
		PackageID     string `json:"package_id"`
		PackageName   string `json:"package_name"`
		PaymentMethod string `json:"payment_method"`
	}
	if payment.Metadata != "" && payment.Metadata != "{}" {
		_ = json.Unmarshal([]byte(payment.Metadata), &meta)
	}

	event := map[string]interface{}{
		"id":        uuid.New().String(),
		"type":      "transaction.paid",
		"timestamp": time.Now().Format(time.RFC3339),
		"user_id":   payment.UserID.String(),
		"data": map[string]interface{}{
			"enrollment_id":  func() string { if payment.BookingID != nil { return payment.BookingID.String() }; return "" }(),
			"total_price":    payment.Amount,
			"payment_method": meta.PaymentMethod,
			"package_id":     meta.PackageID,
			"package_name":   meta.PackageName,
		},
		"success": true,
	}
	_ = t.eventPublisher.Publish(payment.ID.String(), event)

	response.Success(c, http.StatusOK, "Payment simulated successfully", gin.H{
		"orderId": orderID,
		"status":  "paid",
	})
}