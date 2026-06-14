package dto

import (
	"time"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusPaid      PaymentStatus = "paid"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCancelled PaymentStatus = "cancelled"
	PaymentStatusExpired   PaymentStatus = "expired"
)

// PaymentMethod represents the payment method used
type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodEWallet     PaymentMethod = "e_wallet"
	PaymentMethodQris        PaymentMethod = "qris"
	PaymentMethodCStore      PaymentMethod = "cstore"
)

// Payment DTOs

type CreatePaymentRequest struct {
	EnrollmentID uint    `json:"enrollmentId" binding:"required"`
	UserID       uint    `json:"userId" binding:"required"`
	Amount       float64 `json:"amount" binding:"required"`
	PaymentMethod PaymentMethod `json:"paymentMethod" binding:"required"`
}

type PaymentResponse struct {
	ID            uint           `json:"id"`
	EnrollmentID  uint           `json:"enrollmentId"`
	UserID        uint           `json:"userId"`
	OrderID       string         `json:"orderId"`
	Amount        float64        `json:"amount"`
	PaymentMethod PaymentMethod  `json:"paymentMethod"`
	Status        PaymentStatus  `json:"status"`
	PaymentURL    string         `json:"paymentUrl,omitempty"`
	PaidAt        *time.Time     `json:"paidAt,omitempty"`
	ExpiresAt     *time.Time     `json:"expiresAt,omitempty"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
}

type PaymentListResponse struct {
	Data       []PaymentResponse `json:"data"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"totalPages"`
}

type PaymentCallbackRequest struct {
	OrderID     string `json:"orderId" form:"order_id"`
	StatusCode  string `json:"statusCode" form:"status_code"`
	TransactionStatus string `json:"transactionStatus" form:"transaction_status"`
	PaymentType string `json:"paymentType" form:"payment_type"`
	GrossAmount string `json:"grossAmount" form:"gross_amount"`
}

type PaymentDetailResponse struct {
	OrderID       string         `json:"orderId"`
	PaymentType   PaymentMethod  `json:"paymentType"`
	TransactionID string         `json:"transactionId"`
	GrossAmount   float64        `json:"grossAmount"`
	Status        PaymentStatus  `json:"status"`
	PaymentURL    string         `json:"paymentUrl"`
	ExpiryTime    *time.Time     `json:"expiryTime"`
}