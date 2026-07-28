package dto

import (
	model "payment-service/models"
	"time"

	"github.com/google/uuid"
)

// CreatePaymentRequest is the request DTO for creating a payment
type CreatePaymentRequest struct {
	BookingID       uuid.UUID `json:"bookingId" binding:"required"`
	UserID          uuid.UUID `json:"userId" binding:"required"`
	Amount          float64   `json:"amount" binding:"required,gt=0"`
	Currency        string    `json:"currency" binding:"required"`
	PaymentMethodID uint      `json:"paymentMethodId" binding:"required"`
	Description     string    `json:"description"`
	Metadata        string    `json:"metadata"`
}

// CreatePaymentResponse is the response DTO for creating a payment
type CreatePaymentResponse struct {
	ID                 uuid.UUID `json:"id"`
	OrderID            string    `json:"orderId"`
	Amount             float64   `json:"amount"`
	Currency           string    `json:"currency"`
	Status             string    `json:"status"`
	PaymentMethod      string    `json:"paymentMethod"`
	Gateway            string    `json:"gateway"`
	GatewayPaymentURL  string    `json:"gatewayPaymentUrl"`
	VaNumber           string    `json:"vaNumber,omitempty"`
	QrCodeURL          string    `json:"qrCodeUrl,omitempty"`
	ExpiryTime         time.Time `json:"expiryTime"`
}

// GetPaymentResponse is the response DTO for getting a payment
type GetPaymentResponse struct {
	ID                 uuid.UUID         `json:"id"`
	OrderID            string            `json:"orderId"`
	BookingID          *uuid.UUID        `json:"bookingId"`
	UserID             uuid.UUID         `json:"userId"`
	Amount             float64           `json:"amount"`
	Currency           string            `json:"currency"`
	Status             model.PaymentStatus     `json:"status"`
	PaymentMethod      *PaymentMethodDTO `json:"paymentMethod,omitempty"`
	Gateway            string             `json:"gateway"`
	GatewayOrderID     string            `json:"gatewayOrderId"`
	GatewayPaymentURL  string            `json:"gatewayPaymentUrl"`
	VaNumber           string             `json:"vaNumber"`
	QrCodeURL          string             `json:"qrCodeUrl"`
	Description       string             `json:"description"`
	ExpiryTime         *time.Time        `json:"expiryTime"`
	PaidAt             *time.Time        `json:"paidAt"`
	CreatedAt          time.Time         `json:"createdAt"`
	UpdatedAt          time.Time         `json:"updatedAt"`
	Transactions       []TransactionDTO  `json:"transactions,omitempty"`
}

// PaymentMethodDTO is the DTO for payment method
type PaymentMethodDTO struct {
	ID          uint    `json:"id"`
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Icon        string   `json:"icon"`
	IsActive    bool     `json:"isActive"`
	Gateway     string   `json:"gateway"`
	FeeType     string   `json:"feeType"`
	FeeAmount   float64  `json:"feeAmount"`
}

// TransactionDTO is the DTO for transaction
type TransactionDTO struct {
	ID              uuid.UUID        `json:"id"`
	PaymentID       uuid.UUID        `json:"paymentId"`
	Type            model.TransactionType  `json:"type"`
	Status          model.TransactionStatus `json:"status"`
	Amount          float64          `json:"amount"`
	Currency        string           `json:"currency"`
	Gateway         string           `json:"gateway"`
	GatewayTxnID    string           `json:"gatewayTxnId"`
	ErrorCode       string           `json:"errorCode,omitempty"`
	ErrorMessage    string           `json:"errorMessage,omitempty"`
	ProcessedAt     *time.Time       `json:"processedAt"`
	CreatedAt       time.Time        `json:"createdAt"`
}

// ListPaymentsRequest is the request DTO for listing payments
type ListPaymentsRequest struct {
	UserID   *uuid.UUID     `form:"userId"`
	BookingID *uuid.UUID    `form:"bookingId"`
	Status   model.PaymentStatus  `form:"status"`
	Page     int            `form:"page,default=1"`
	PageSize int            `form:"pageSize,default=20"`
}

// ListPaymentsResponse is the response DTO for listing payments
type ListPaymentsResponse struct {
	Data       []GetPaymentResponse `json:"data"`
	Pagination PaginationDTO         `json:"pagination"`
}

// PaginationDTO is the DTO for pagination
type PaginationDTO struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total 	   int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

// PaymentCallbackRequest is the request DTO for payment gateway callback
type PaymentCallbackRequest struct {
	OrderID     string `json:"order_id"`
	Status      string `json:"status"`
	Gateway     string `json:"gateway"`
	PaymentType string `json:"payment_type"`
	Amount      string `json:"amount"`
	Timestamp   string `json:"timestamp"`
	Signature   string `json:"signature_key"`
}

// RefundRequest is the request DTO for creating a refund
type RefundRequest struct {
	PaymentID     uuid.UUID `json:"paymentId" binding:"required"`
	TransactionID uuid.UUID `json:"transactionId" binding:"required"`
	Amount        float64   `json:"amount" binding:"required,gt=0"`
	Reason        string    `json:"reason" binding:"required"`
}

// RefundResponse is the response DTO for refund
type RefundResponse struct {
	ID              uuid.UUID  `json:"id"`
	PaymentID       uuid.UUID  `json:"paymentId"`
	TransactionID   uuid.UUID  `json:"transactionId"`
	Amount          float64    `json:"amount"`
	Status          string     `json:"status"`
	GatewayRefundID string     `json:"gatewayRefundId"`
	ProcessedAt     *time.Time `json:"processedAt"`
	CreatedAt       time.Time  `json:"createdAt"`
}

// CancelPaymentRequest is the request DTO for cancelling a payment
type CancelPaymentRequest struct {
	Reason string `json:"reason"`
}

// WebhookEvent represents a webhook event from payment gateway
type WebhookEvent struct {
	EventType   string    `json:"event_type"`
	OrderID     string    `json:"order_id"`
	PaymentID   string    `json:"payment_id"`
	Status      string    `json:"status"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	Gateway     string    `json:"gateway"`
	Timestamp   time.Time `json:"timestamp"`
	RawResponse string    `json:"raw_response"`
}