package services

import (
	"context"
	"core-service/models"
	"core-service/models/dto"
	"core-service/pkg/base"
	"core-service/repositories"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ISalesService interface {
	// CRUD operations
	CreateSale(ctx context.Context, req dto.CreateSaleRequest) (*models.Sale, error)
	GetSaleByID(ctx context.Context, id uuid.UUID) (*models.Sale, error)
	GetSaleByOrderNumber(ctx context.Context, orderNumber string) (*models.Sale, error)
	ListSales(ctx context.Context, req dto.ListSalesRequest) (*dto.ListSalesResponse, error)
	UpdateSaleStatus(ctx context.Context, id uuid.UUID, status models.SaleStatus, notes string) error
	DeleteSale(ctx context.Context, id uuid.UUID) error

	// Analytics operations
	GetOverview(ctx context.Context, startDate, endDate string) (*dto.SalesOverviewResponse, error)
	GetMonthlySales(ctx context.Context, year, month int) (*dto.MonthlySalesResponse, error)
	GetYearlySales(ctx context.Context, year int) ([]dto.MonthlySalesResponse, error)
	GetSalesTrend(ctx context.Context, startDate, endDate string) (*dto.SalesTrendResponse, error)
	GetSalesBySource(ctx context.Context, startDate, endDate string) ([]dto.SalesBySourceResponse, error)
	GetSalesByPackageType(ctx context.Context, startDate, endDate string) ([]dto.SalesByPackageTypeResponse, error)

	// Helpers
	GenerateOrderNumber() string
}

type SalesService struct {
	salesRepo    repositories.ISalesRepository
}

func NewSalesService(salesRepo repositories.ISalesRepository) ISalesService {
	return &SalesService{
		salesRepo: salesRepo,
	}
}

// CreateSale creates a new sale
func (s *SalesService) CreateSale(ctx context.Context, req dto.CreateSaleRequest) (*models.Sale, error) {
	// Calculate totals
	var totalAmount, totalDiscount float64
	var items []models.SaleItem
	var primaryPackageID *uuid.UUID
	var primaryPackageName string
	var primaryPackageType string

	for _, item := range req.Items {
		subtotal := (item.UnitPrice * float64(item.Quantity)) - item.Discount
		totalAmount += item.UnitPrice * float64(item.Quantity)
		totalDiscount += item.Discount

		saleItem := models.SaleItem{
			ID:          uuid.New(),
			PackageID:   item.PackageID,
			PackageName: item.PackageName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Discount:    item.Discount,
			Subtotal:    subtotal,
		}
		items = append(items, saleItem)

		// Use first item as primary
		if primaryPackageID == nil {
			pid := item.PackageID
			primaryPackageID = &pid
			primaryPackageName = item.PackageName
		}
	}

	finalAmount := totalAmount - totalDiscount

	sale := &models.Sale{
		ID:             uuid.New(),
		OrderNumber:    s.GenerateOrderNumber(),
		UserID:         req.UserID,
		PackageID:      primaryPackageID,
		PackageName:    primaryPackageName,
		PackageType:    primaryPackageType,
		TotalAmount:    totalAmount,
		DiscountAmount: totalDiscount,
		FinalAmount:    finalAmount,
		Status:         models.SaleStatusPending,
		Source:         req.Source,
		PaymentMethod:  req.PaymentMethod,
		Currency:       "IDR",
		Notes:          req.Notes,
	}

	if err := s.salesRepo.CreateWithItems(ctx, sale, items); err != nil {
		return nil, err
	}

	sale.Items = items
	return sale, nil
}

// GetSaleByID retrieves a sale by ID
func (s *SalesService) GetSaleByID(ctx context.Context, id uuid.UUID) (*models.Sale, error) {
	return s.salesRepo.FindByIDWithItems(ctx, id)
}

// GetSaleByOrderNumber retrieves a sale by order number
func (s *SalesService) GetSaleByOrderNumber(ctx context.Context, orderNumber string) (*models.Sale, error) {
	return s.salesRepo.FindByOrderNumber(ctx, orderNumber)
}

// ListSales retrieves sales with filtering and pagination
func (s *SalesService) ListSales(ctx context.Context, req dto.ListSalesRequest) (*dto.ListSalesResponse, error) {
	where := make(map[string]any)

	if req.UserID != nil {
		where["user_id"] = req.UserID
	}
	if req.Status != "" {
		where["status"] = req.Status
	}
	if req.Source != "" {
		where["source"] = req.Source
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	opts := base.NewQueryOptions().
		WithWhere(where).
		WithPagination((page-1)*pageSize, pageSize).
		WithOrder("created_at DESC").
		WithPreloads("Items")

	sales, total, err := s.salesRepo.FindAll(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Convert to response DTOs
	data := make([]dto.SaleResponse, len(sales))
	for i, sale := range sales {
		data[i] = s.toSaleResponse(&sale)
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &dto.ListSalesResponse{
		Data: data,
		Pagination: dto.PaginationDTO{
			Page:       page,
			Limit:      pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// UpdateSaleStatus updates the status of a sale
func (s *SalesService) UpdateSaleStatus(ctx context.Context, id uuid.UUID, status models.SaleStatus, notes string) error {
	sale, err := s.salesRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	sale.Status = status
	if notes != "" {
		sale.Notes = notes
	}

	// Update timestamps based on status
	now := time.Now()
	switch status {
	case models.SaleStatusCompleted:
		sale.PaidAt = &now
	case models.SaleStatusRefunded:
		sale.RefundedAt = &now
	}

	return s.salesRepo.Update(ctx, sale)
}

// DeleteSale deletes a sale
func (s *SalesService) DeleteSale(ctx context.Context, id uuid.UUID) error {
	sale, err := s.salesRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return s.salesRepo.Delete(ctx, sale)
}

// GetOverview retrieves sales overview statistics
func (s *SalesService) GetOverview(ctx context.Context, startDate, endDate string) (*dto.SalesOverviewResponse, error) {
	stats, err := s.salesRepo.GetOverviewStats(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Calculate growth rate (compare with previous period)
	growthRate := s.calculateGrowthRate(ctx, startDate, endDate)

	return &dto.SalesOverviewResponse{
		TotalRevenue:   stats.TotalRevenue,
		TotalSales:     stats.TotalSales,
		TotalRefunds:   stats.TotalRefunds,
		NetRevenue:     stats.NetRevenue,
		AvgOrderValue:  stats.AvgOrderValue,
		CompletedSales: stats.CompletedSales,
		PendingSales:   stats.PendingSales,
		CanceledSales:  stats.CanceledSales,
		RefundedSales:  stats.RefundedSales,
		GrowthRate:     growthRate,
	}, nil
}

// GetMonthlySales retrieves aggregated monthly sales data
func (s *SalesService) GetMonthlySales(ctx context.Context, year, month int) (*dto.MonthlySalesResponse, error) {
	data, err := s.salesRepo.GetMonthlySales(ctx, year, month)
	if err != nil {
		return nil, err
	}

	// Get source breakdown
	sourceData, _ := s.salesRepo.GetSalesBySource(ctx, fmt.Sprintf("%d-%02d-01", year, month), fmt.Sprintf("%d-%02d-31", year, month))

	sourceBreakdown := make(map[string]float64)
	for _, s := range sourceData {
		sourceBreakdown[s.Source] = s.TotalRevenue
	}

	return &dto.MonthlySalesResponse{
		Year:            data.Year,
		Month:           data.Month,
		TotalSales:      data.TotalSales,
		TotalRevenue:    data.TotalRevenue,
		TotalDiscount:   data.TotalDiscount,
		TotalRefunds:    data.TotalRefunds,
		NetRevenue:      data.NetRevenue,
		AvgOrderValue:   data.AvgOrderValue,
		CanceledSales:   data.CanceledSales,
		PendingSales:    data.PendingSales,
		CompletedSales:  data.CompletedSales,
		SourceBreakdown: sourceBreakdown,
	}, nil
}

// GetYearlySales retrieves aggregated monthly sales data for a year
func (s *SalesService) GetYearlySales(ctx context.Context, year int) ([]dto.MonthlySalesResponse, error) {
	data, err := s.salesRepo.GetMonthlySalesYear(ctx, year)
	if err != nil {
		return nil, err
	}

	results := make([]dto.MonthlySalesResponse, len(data))
	for i, month := range data {
		results[i] = dto.MonthlySalesResponse{
			Year:            month.Year,
			Month:           month.Month,
			TotalSales:      month.TotalSales,
			TotalRevenue:    month.TotalRevenue,
			TotalDiscount:   month.TotalDiscount,
			TotalRefunds:    month.TotalRefunds,
			NetRevenue:      month.NetRevenue,
			AvgOrderValue:   month.AvgOrderValue,
			CanceledSales:   month.CanceledSales,
			PendingSales:    month.PendingSales,
			CompletedSales:  month.CompletedSales,
			SourceBreakdown: make(map[string]float64),
		}
	}

	return results, nil
}

// GetSalesTrend retrieves daily sales data for trend analysis
func (s *SalesService) GetSalesTrend(ctx context.Context, startDate, endDate string) (*dto.SalesTrendResponse, error) {
	daily, err := s.salesRepo.GetSalesTrend(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	data := make([]dto.DailySalesResponse, len(daily))
	for i, day := range daily {
		data[i] = dto.DailySalesResponse{
			Date:           day.Date,
			TotalSales:     day.TotalSales,
			TotalRevenue:   day.TotalRevenue,
			TotalDiscount:  day.TotalDiscount,
			NetRevenue:     day.NetRevenue,
			CompletedSales: day.CompletedSales,
			RefundedSales:  day.RefundedSales,
			PendingSales:   day.PendingSales,
		}
	}

	return &dto.SalesTrendResponse{
		Period:    "daily",
		Data:      data,
		StartDate: startDate,
		EndDate:   endDate,
	}, nil
}

// GetSalesBySource retrieves sales breakdown by source
func (s *SalesService) GetSalesBySource(ctx context.Context, startDate, endDate string) ([]dto.SalesBySourceResponse, error) {
	data, err := s.salesRepo.GetSalesBySource(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Calculate total for percentage calculation
	var totalRevenue float64
	for _, item := range data {
		totalRevenue += item.TotalRevenue
	}

	results := make([]dto.SalesBySourceResponse, len(data))
	for i, item := range data {
		percentage := 0.0
		if totalRevenue > 0 {
			percentage = (item.TotalRevenue / totalRevenue) * 100
		}

		results[i] = dto.SalesBySourceResponse{
			Source:       item.Source,
			TotalSales:   item.TotalSales,
			TotalRevenue: item.TotalRevenue,
			Percentage:   percentage,
		}
	}

	return results, nil
}

// GetSalesByPackageType retrieves sales breakdown by package type
func (s *SalesService) GetSalesByPackageType(ctx context.Context, startDate, endDate string) ([]dto.SalesByPackageTypeResponse, error) {
	data, err := s.salesRepo.GetSalesByPackageType(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Calculate total for percentage calculation
	var totalRevenue float64
	for _, item := range data {
		totalRevenue += item.TotalRevenue
	}

	results := make([]dto.SalesByPackageTypeResponse, len(data))
	for i, item := range data {
		percentage := 0.0
		if totalRevenue > 0 {
			percentage = (item.TotalRevenue / totalRevenue) * 100
		}

		results[i] = dto.SalesByPackageTypeResponse{
			PackageType:  item.PackageType,
			TotalSales:   item.TotalSales,
			TotalRevenue: item.TotalRevenue,
			Percentage:   percentage,
		}
	}

	return results, nil
}

// toSaleResponse converts a Sale model to SaleResponse DTO
func (s *SalesService) toSaleResponse(sale *models.Sale) dto.SaleResponse {
	items := make([]dto.SaleItemResponse, len(sale.Items))
	for i, item := range sale.Items {
		items[i] = dto.SaleItemResponse{
			ID:          item.ID,
			SaleID:      item.SaleID,
			PackageID:   item.PackageID,
			PackageName: item.PackageName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Discount:    item.Discount,
			Subtotal:    item.Subtotal,
		}
	}

	return dto.SaleResponse{
		ID:             sale.ID,
		OrderNumber:    sale.OrderNumber,
		PaymentID:      sale.PaymentID,
		UserID:         sale.UserID,
		PackageID:      sale.PackageID,
		PackageName:    sale.PackageName,
		PackageType:    sale.PackageType,
		Items:          items,
		TotalAmount:    sale.TotalAmount,
		DiscountAmount: sale.DiscountAmount,
		FinalAmount:    sale.FinalAmount,
		Status:         sale.Status,
		Source:         sale.Source,
		PaymentMethod:  sale.PaymentMethod,
		Currency:       sale.Currency,
		PaidAt:         sale.PaidAt,
		RefundedAt:     sale.RefundedAt,
		Notes:          sale.Notes,
		CreatedAt:      sale.CreatedAt,
		UpdatedAt:      sale.UpdatedAt,
	}
}

// calculateGrowthRate calculates the growth rate compared to the previous period
func (s *SalesService) calculateGrowthRate(ctx context.Context, startDate, endDate string) float64 {
	// Parse dates
	start, err1 := time.Parse("2006-01-02", startDate)
	end, err2 := time.Parse("2006-01-02", endDate)
	if err1 != nil || err2 != nil {
		return 0
	}

	// Calculate previous period
	periodDays := int(end.Sub(start).Hours() / 24)
	prevEnd := start.AddDate(0, 0, -1)
	prevStart := prevEnd.AddDate(0, 0, -periodDays)

	currentRevenue, _ := s.salesRepo.SumRevenue(ctx, startDate, endDate)
	prevRevenue, _ := s.salesRepo.SumRevenue(ctx, prevStart.Format("2006-01-02"), prevEnd.Format("2006-01-02"))

	if prevRevenue == 0 {
		return 0
	}

	return ((currentRevenue - prevRevenue) / prevRevenue) * 100
}

// GenerateOrderNumber generates a unique order number
func (s *SalesService) GenerateOrderNumber() string {
	return fmt.Sprintf("ORD-%s-%d", time.Now().Format("20060102150405"), time.Now().Nanosecond()%1000)
}
