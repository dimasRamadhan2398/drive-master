package models

import (
	"time"

	"github.com/google/uuid"
)

// SaleStatus represents the status of a sale
type SaleStatus string

const (
	SaleStatusPending   SaleStatus = "pending"
	SaleStatusCompleted SaleStatus = "completed"
	SaleStatusRefunded  SaleStatus = "refunded"
	SaleStatusCancelled SaleStatus = "cancelled"
)

// Sale represents a sale/order record
type Sale struct {
	ID              uuid.UUID   `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrderNumber     string      `json:"orderNumber" gorm:"type:varchar(50);uniqueIndex;not null"`
	PaymentID       *uuid.UUID  `json:"paymentId" gorm:"type:uuid;index"`
	UserID          uuid.UUID   `json:"userId" gorm:"type:uuid;not null;index"`
	PackageID       *uuid.UUID  `json:"packageId" gorm:"type:uuid;index"`
	PackageName     string      `json:"packageName" gorm:"size:255"` // Denormalized for reporting
	PackageType     string      `json:"packageType" gorm:"size:50"`  // Denormalized for reporting
	Items           []SaleItem  `json:"items,omitempty" gorm:"foreignKey:SaleID"`
	TotalAmount     float64     `json:"totalAmount" gorm:"type:decimal(12,2);not null"`
	DiscountAmount  float64     `json:"discountAmount" gorm:"type:decimal(12,2);default:0"`
	FinalAmount     float64     `json:"finalAmount" gorm:"type:decimal(12,2);not null"`
	Status          SaleStatus  `json:"status" gorm:"type:varchar(20);default:'pending';index"`
	Source          string      `json:"source" gorm:"type:varchar(50)"` // web, mobile, admin
	PaymentMethod   string      `json:"paymentMethod" gorm:"type:varchar(50)"`
	Currency        string      `json:"currency" gorm:"type:varchar(3);default:'IDR'"`
	PaidAt          *time.Time `json:"paidAt"`
	RefundedAt      *time.Time `json:"refundedAt"`
	Notes           string      `json:"notes" gorm:"type:text"`
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       time.Time   `json:"updatedAt"`
}

// SaleItem represents individual items in a sale
type SaleItem struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SaleID      uuid.UUID `json:"saleId" gorm:"type:uuid;not null;index"`
	Sale        *Sale    `json:"sale,omitempty" gorm:"foreignKey:SaleID"`
	PackageID   uuid.UUID `json:"packageId" gorm:"type:uuid;not null"`
	PackageName string    `json:"packageName" gorm:"size:255"`
	Quantity    int       `json:"quantity" gorm:"default:1"`
	UnitPrice   float64   `json:"unitPrice" gorm:"type:decimal(12,2);not null"`
	Discount    float64   `json:"discount" gorm:"type:decimal(12,2);default:0"`
	Subtotal    float64   `json:"subtotal" gorm:"type:decimal(12,2);not null"`
	CreatedAt   time.Time `json:"createdAt"`
}

// MonthlySales represents aggregated monthly sales data
type MonthlySales struct {
	Year            int     `json:"year" gorm:"primaryKey"`
	Month           int     `json:"month" gorm:"primaryKey"`
	TotalSales      int64   `json:"totalSales"`
	TotalRevenue    float64 `json:"totalRevenue"`
	TotalDiscount   float64 `json:"totalDiscount"`
	TotalRefunds    float64 `json:"totalRefunds"`
	NetRevenue      float64 `json:"netRevenue"`
	AvgOrderValue   float64 `json:"avgOrderValue"`
	CanceledSales   int64   `json:"canceledSales"`
	PendingSales    int64   `json:"pendingSales"`
	CompletedSales  int64   `json:"completedSales"`
	SourceBreakdown string  `json:"sourceBreakdown" gorm:"type:text"` // JSON string of source -> revenue
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// DailySales represents aggregated daily sales data
type DailySales struct {
	Date           string  `json:"date"`
	TotalSales     int64   `json:"totalSales"`
	TotalRevenue   float64 `json:"totalRevenue"`
	TotalDiscount  float64 `json:"totalDiscount"`
	NetRevenue     float64 `json:"netRevenue"`
	CompletedSales int64   `json:"completedSales"`
	RefundedSales  int64   `json:"refundedSales"`
	PendingSales   int64   `json:"pendingSales"`
}

// SalesBySource represents sales breakdown by source
type SalesBySource struct {
	Source       string  `json:"source"`
	TotalSales   int64   `json:"totalSales"`
	TotalRevenue float64 `json:"totalRevenue"`
	Percentage   float64 `json:"percentage"`
}

// SalesByPackageType represents sales breakdown by package type
type SalesByPackageType struct {
	PackageType  string  `json:"packageType"`
	TotalSales   int64   `json:"totalSales"`
	TotalRevenue float64 `json:"totalRevenue"`
	Percentage   float64 `json:"percentage"`
}

// SalesOverviewStats represents overall sales statistics
type SalesOverviewStats struct {
	TotalRevenue    float64 `json:"totalRevenue"`
	TotalSales      int64   `json:"totalSales"`
	TotalRefunds    float64 `json:"totalRefunds"`
	NetRevenue      float64 `json:"netRevenue"`
	AvgOrderValue   float64 `json:"avgOrderValue"`
	CompletedSales  int64   `json:"completedSales"`
	PendingSales    int64   `json:"pendingSales"`
	CanceledSales   int64   `json:"canceledSales"`
	RefundedSales   int64   `json:"refundedSales"`
}
