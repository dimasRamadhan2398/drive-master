package dto

import (
	"core-service/models"
	"time"

	"github.com/google/uuid"
)

// ========== SALES DTOs ==========

// CreateSaleRequest is the request DTO for creating a sale
type CreateSaleRequest struct {
	UserID        uuid.UUID          `json:"userId" binding:"required"`
	PackageID     *uuid.UUID         `json:"packageId"`
	Items         []CreateSaleItemRequest `json:"items" binding:"required,min=1"`
	PaymentMethod string             `json:"paymentMethod" binding:"required"`
	Source        string             `json:"source"`
	Notes         string             `json:"notes"`
}

// CreateSaleItemRequest is the request DTO for a sale item
type CreateSaleItemRequest struct {
	PackageID   uuid.UUID `json:"packageId" binding:"required"`
	PackageName string    `json:"packageName"`
	Quantity    int       `json:"quantity" binding:"required,min=1"`
	UnitPrice   float64   `json:"unitPrice" binding:"required,gt=0"`
	Discount    float64   `json:"discount"`
}

// SaleResponse is the response DTO for a sale
type SaleResponse struct {
	ID              uuid.UUID             `json:"id"`
	OrderNumber     string                `json:"orderNumber"`
	PaymentID       *uuid.UUID            `json:"paymentId"`
	UserID          uuid.UUID             `json:"userId"`
	PackageID       *uuid.UUID            `json:"packageId"`
	PackageName     string                `json:"packageName"`
	PackageType     string                `json:"packageType"`
	Items           []SaleItemResponse     `json:"items,omitempty"`
	TotalAmount     float64               `json:"totalAmount"`
	DiscountAmount  float64               `json:"discountAmount"`
	FinalAmount     float64               `json:"finalAmount"`
	Status          models.SaleStatus     `json:"status"`
	Source          string                `json:"source"`
	PaymentMethod   string               `json:"paymentMethod"`
	Currency        string                `json:"currency"`
	PaidAt          *time.Time           `json:"paidAt"`
	RefundedAt      *time.Time           `json:"refundedAt"`
	Notes           string                `json:"notes"`
	CreatedAt       time.Time             `json:"createdAt"`
	UpdatedAt       time.Time             `json:"updatedAt"`
}

// SaleItemResponse is the response DTO for a sale item
type SaleItemResponse struct {
	ID          uuid.UUID `json:"id"`
	SaleID      uuid.UUID `json:"saleId"`
	PackageID   uuid.UUID `json:"packageId"`
	PackageName string    `json:"packageName"`
	Quantity    int       `json:"quantity"`
	UnitPrice   float64   `json:"unitPrice"`
	Discount    float64   `json:"discount"`
	Subtotal    float64   `json:"subtotal"`
}

// ListSalesRequest is the request DTO for listing sales
type ListSalesRequest struct {
	UserID       *uuid.UUID           `form:"userId"`
	Status       models.SaleStatus    `form:"status"`
	Source       string               `form:"source"`
	StartDate    string               `form:"startDate"`
	EndDate      string               `form:"endDate"`
	Page         int                  `form:"page,default=1"`
	PageSize     int                  `form:"pageSize,default=20"`
}

// ListSalesResponse is the response DTO for listing sales
type ListSalesResponse struct {
	Data       []SaleResponse `json:"data"`
	Pagination PaginationDTO  `json:"pagination"`
}

// PaginationDTO is the DTO for pagination
type PaginationDTO struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

// ========== SALES ANALYTICS DTOs ==========

// SalesOverviewResponse is the response DTO for sales overview
type SalesOverviewResponse struct {
	TotalRevenue     float64             `json:"totalRevenue"`
	TotalSales       int64               `json:"totalSales"`
	TotalRefunds     float64             `json:"totalRefunds"`
	NetRevenue       float64             `json:"netRevenue"`
	AvgOrderValue    float64             `json:"avgOrderValue"`
	CompletedSales   int64               `json:"completedSales"`
	PendingSales     int64               `json:"pendingSales"`
	CanceledSales    int64               `json:"canceledSales"`
	RefundedSales    int64               `json:"refundedSales"`
	GrowthRate       float64             `json:"growthRate"`
}

// MonthlySalesResponse is the response DTO for monthly sales
type MonthlySalesResponse struct {
	Year             int                    `json:"year"`
	Month            int                    `json:"month"`
	TotalSales       int64                  `json:"totalSales"`
	TotalRevenue     float64                `json:"totalRevenue"`
	TotalDiscount    float64                `json:"totalDiscount"`
	TotalRefunds     float64                `json:"totalRefunds"`
	NetRevenue       float64                `json:"netRevenue"`
	AvgOrderValue    float64                `json:"avgOrderValue"`
	CanceledSales    int64                  `json:"canceledSales"`
	PendingSales     int64                  `json:"pendingSales"`
	CompletedSales   int64                  `json:"completedSales"`
	SourceBreakdown  map[string]float64     `json:"sourceBreakdown"`
}

// DailySalesResponse is the response DTO for daily sales
type DailySalesResponse struct {
	Date           string `json:"date"`
	TotalSales     int64  `json:"totalSales"`
	TotalRevenue   float64 `json:"totalRevenue"`
	TotalDiscount  float64 `json:"totalDiscount"`
	NetRevenue     float64 `json:"netRevenue"`
	CompletedSales int64  `json:"completedSales"`
	RefundedSales  int64  `json:"refundedSales"`
	PendingSales   int64  `json:"pendingSales"`
}

// SalesTrendResponse is the response DTO for sales trend
type SalesTrendResponse struct {
	Period      string                 `json:"period"` // "daily", "weekly", "monthly"
	Data        []DailySalesResponse  `json:"data"`
	StartDate   string                 `json:"startDate"`
	EndDate     string                 `json:"endDate"`
}

// SalesBySourceResponse is the response DTO for sales by source
type SalesBySourceResponse struct {
	Source       string  `json:"source"`
	TotalSales   int64   `json:"totalSales"`
	TotalRevenue float64 `json:"totalRevenue"`
	Percentage   float64 `json:"percentage"`
}

// SalesByPackageTypeResponse is the response DTO for sales by package type
type SalesByPackageTypeResponse struct {
	PackageType  string  `json:"packageType"`
	TotalSales   int64   `json:"totalSales"`
	TotalRevenue float64 `json:"totalRevenue"`
	Percentage   float64 `json:"percentage"`
}

// UpdateSaleStatusRequest is the request DTO for updating sale status
type UpdateSaleStatusRequest struct {
	Status models.SaleStatus `json:"status" binding:"required"`
	Notes  string            `json:"notes"`
}
