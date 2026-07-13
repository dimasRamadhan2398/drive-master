package models

import (
	"time"

	"github.com/google/uuid"
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed   TransactionStatus = "failed"
	TransactionStatusRefunded TransactionStatus = "refunded"
)

// TransactionItemType represents the type of item in a transaction
type TransactionItemType string

const (
	TransactionItemTypePackage TransactionItemType = "package"
	TransactionItemTypeAddOn  TransactionItemType = "addon"
)

// Transaction represents a financial transaction record that captures all line items
// This is the source of truth for what was purchased and how much it cost
type Transaction struct {
	ID            uuid.UUID         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	EnrollmentID  uuid.UUID        `json:"enrollmentId" gorm:"type:uuid;not null;uniqueIndex"` // One transaction per enrollment
	UserID        uuid.UUID        `json:"userId" gorm:"type:uuid;not null;index"`
	BasePrice     float64         `json:"basePrice" gorm:"type:decimal(10,2);default:0"`      // Package price
	AddOnsTotal   float64         `json:"addOnsTotal" gorm:"type:decimal(10,2);default:0"`     // Total of all add-ons
	TotalAmount   float64         `json:"totalAmount" gorm:"type:decimal(10,2);not null"`     // Final amount to pay
	Status        TransactionStatus `json:"status" gorm:"type:varchar(20);default:'pending'"`
	PaidAt        *time.Time       `json:"paidAt"`     // When payment was confirmed
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`

	// Relationships
	Items []TransactionItem `json:"items" gorm:"foreignKey:TransactionID"`
}

// TransactionItem represents a single line item in a transaction
type TransactionItem struct {
	ID            uuid.UUID             `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TransactionID uuid.UUID             `json:"transactionId" gorm:"type:uuid;not null;index"`
	ItemType     TransactionItemType   `json:"itemType" gorm:"type:varchar(20);not null"`      // "package" or "addon"
	ItemID       uuid.UUID             `json:"itemId" gorm:"type:uuid;not null"`               // Package or AddOn ID
	ItemName     string                `json:"itemName" gorm:"size:255;not null"`              // Name for display
	Quantity     int                   `json:"quantity" gorm:"default:1"`                     // Quantity (usually 1 for packages)
	UnitPrice    float64               `json:"unitPrice" gorm:"type:decimal(10,2);not null"`   // Price per unit
	Subtotal     float64               `json:"subtotal" gorm:"type:decimal(10,2);not null"`    // UnitPrice * Quantity
	Sessions     int                   `json:"sessions" gorm:"default:0"`                     // Sessions included (for add-ons)
	CreatedAt    time.Time            `json:"createdAt"`
}
