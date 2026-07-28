package models

import (
	"time"

	"github.com/google/uuid"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSuccess   PaymentStatus = "success"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusExpired   PaymentStatus = "expired"
	PaymentStatusCancelled PaymentStatus = "cancelled"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypeCharge    TransactionType = "charge"
	TransactionTypeRefund    TransactionType = "refund"
	TransactionTypeReversal  TransactionType = "reversal"
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusSuccess   TransactionStatus = "success"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusReversed  TransactionStatus = "reversed"
)

// Payment represents a payment record
type Payment struct {
	ID                 uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrderID            string         `json:"orderId" gorm:"type:varchar(50);uniqueIndex;not null"`
	BookingID         *uuid.UUID     `json:"bookingId" gorm:"type:uuid;index"`
	UserID            uuid.UUID      `json:"userId" gorm:"type:uuid;not null;index"`
	Amount            float64        `json:"amount" gorm:"type:decimal(12,2);not null"`
	Currency          string         `json:"currency" gorm:"type:varchar(3);default:'IDR'"`
	Status            PaymentStatus  `json:"status" gorm:"type:varchar(20);default:'pending';index"`
	PaymentMethodID   *uint          `json:"paymentMethodId" gorm:"index"`
	PaymentMethod     *PaymentMethod `json:"paymentMethod,omitempty" gorm:"foreignKey:PaymentMethodID"`
	Gateway           string         `json:"gateway" gorm:"type:varchar(50)"` // midtrans, doku, pakasir, etc.
	GatewayOrderID    string         `json:"gatewayOrderId" gorm:"type:varchar(100)"`
	GatewayPaymentURL string         `json:"gatewayPaymentUrl" gorm:"type:varchar(500)"`
	VaNumber          string         `json:"vaNumber" gorm:"type:varchar(100)"`
	QrCodeURL         string         `json:"qrCodeUrl" gorm:"type:varchar(500)"`
	Metadata          string         `json:"metadata" gorm:"type:jsonb"` // JSON metadata
	Description       string         `json:"description" gorm:"type:text"`
	ExpiryTime        *time.Time     `json:"expiryTime"`
	PaidAt            *time.Time     `json:"paidAt"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         time.Time      `json:"updatedAt"`

	// Relations
	Transactions []Transaction `json:"transactions,omitempty" gorm:"foreignKey:PaymentID"`
}

// Transaction represents a payment transaction record
type Transaction struct {
	ID              uuid.UUID        `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PaymentID       uuid.UUID        `json:"paymentId" gorm:"type:uuid;not null;index"`
	Payment         *Payment         `json:"payment,omitempty" gorm:"foreignKey:PaymentID"`
	Type            TransactionType  `json:"type" gorm:"type:varchar(20);not null"`
	Status          TransactionStatus `json:"status" gorm:"type:varchar(20);default:'pending';index"`
	Amount          float64          `json:"amount" gorm:"type:decimal(12,2);not null"`
	Currency        string           `json:"currency" gorm:"type:varchar(3);default:'IDR'"`
	Gateway         string           `json:"gateway" gorm:"type:varchar(50)"` // midtrans, doku, pakasir
	GatewayTxnID    string           `json:"gatewayTxnId" gorm:"type:varchar(100)"`
	GatewayResponse string           `json:"gatewayResponse" gorm:"type:jsonb"`
	ErrorCode       string           `json:"errorCode" gorm:"type:varchar(50)"`
	ErrorMessage    string           `json:"errorMessage" gorm:"type:text"`
	PaymentMethodID *uint            `json:"paymentMethodId" gorm:"index"`
	PaymentMethod   *PaymentMethod   `json:"paymentMethod,omitempty" gorm:"foreignKey:PaymentMethodID"`
	ProcessedAt     *time.Time       `json:"processedAt"`
	CreatedAt       time.Time        `json:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt"`
}

// PaymentMethod represents available payment methods
type PaymentMethod struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Code        string    `json:"code" gorm:"type:varchar(50);uniqueIndex;not null"`
	Name        string    `json:"name" gorm:"type:varchar(100);not null"`
	Description string    `json:"description" gorm:"type:text"`
	Icon        string    `json:"icon" gorm:"type:varchar(255)"`
	IsActive    bool      `json:"isActive" gorm:"default:true"`
	SortOrder   int       `json:"sortOrder" gorm:"default:0"`
	Gateway     string    `json:"gateway" gorm:"type:varchar(50)"` // Which gateway supports this
	MinAmount   float64   `json:"minAmount" gorm:"type:decimal(12,2);default:0"`
	MaxAmount   float64   `json:"maxAmount" gorm:"type:decimal(12,2);default:0"`
	FeeType     string    `json:"feeType" gorm:"type:varchar(20)"` // fixed, percentage, both
	FeeAmount   float64   `json:"feeAmount" gorm:"type:decimal(12,2);default:0"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Refund represents a refund request
type Refund struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PaymentID       uuid.UUID  `json:"paymentId" gorm:"type:uuid;not null;index"`
	Payment         *Payment   `json:"payment,omitempty" gorm:"foreignKey:PaymentID"`
	TransactionID   uuid.UUID  `json:"transactionId" gorm:"type:uuid;not null"`
	Amount          float64    `json:"amount" gorm:"type:decimal(12,2);not null"`
	Reason          string     `json:"reason" gorm:"type:text"`
	Status          string     `json:"status" gorm:"type:varchar(20);default:'pending'"`
	GatewayRefundID string     `json:"gatewayRefundId" gorm:"type:varchar(100)"`
	ProcessedAt     *time.Time `json:"processedAt"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}