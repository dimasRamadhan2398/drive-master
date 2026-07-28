package services

import (
	"context"
	"time"

	cClient "booking-service/clients/core"
)

// IRevenueService defines the interface for revenue operations
type IRevenueService interface {
	GetRevenueStats(ctx context.Context, startDate, endDate string) (*RevenueStatsResponse, error)
	GetCurrentMonthRevenue(ctx context.Context) (*RevenueStatsResponse, error)
}

// RevenueService implements IRevenueService
type RevenueService struct {
	coreClient cClient.ICoreClient
}

// RevenueStatsResponse is the response for revenue stats
type RevenueStatsResponse struct {
	TotalRevenue     float64 `json:"totalRevenue"`
	TotalSales       int64   `json:"totalSales"`
	TotalRefunds     float64 `json:"totalRefunds"`
	NetRevenue       float64 `json:"netRevenue"`
	AvgOrderValue    float64 `json:"avgOrderValue"`
	CompletedSales   int64   `json:"completedSales"`
	PendingSales     int64   `json:"pendingSales"`
	CanceledSales    int64   `json:"canceledSales"`
	RefundedSales    int64   `json:"refundedSales"`
	GrowthRate       float64 `json:"growthRate"`
	GrowthPercentage float64 `json:"growthPercentage"`
}

// NewRevenueService creates a new revenue service
func NewRevenueService(coreClient cClient.ICoreClient) IRevenueService {
	return &RevenueService{
		coreClient: coreClient,
	}
}

// GetRevenueStats retrieves revenue stats for a given date range
func (s *RevenueService) GetRevenueStats(ctx context.Context, startDate, endDate string) (*RevenueStatsResponse, error) {
	overview, err := s.coreClient.GetSalesOverview(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return &RevenueStatsResponse{
		TotalRevenue:     overview.TotalRevenue,
		TotalSales:       overview.TotalSales,
		TotalRefunds:     overview.TotalRefunds,
		NetRevenue:       overview.NetRevenue,
		AvgOrderValue:    overview.AvgOrderValue,
		CompletedSales:   overview.CompletedSales,
		PendingSales:     overview.PendingSales,
		CanceledSales:    overview.CanceledSales,
		RefundedSales:    overview.RefundedSales,
		GrowthRate:       overview.GrowthRate,
		GrowthPercentage: overview.GrowthRate, // Alias for convenience
	}, nil
}

// GetCurrentMonthRevenue retrieves revenue stats for the current month with growth from last month
func (s *RevenueService) GetCurrentMonthRevenue(ctx context.Context) (*RevenueStatsResponse, error) {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	// Current month date range
	currentStart := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	currentEnd := currentStart.AddDate(0, 1, -1) // Last day of month

	// Format dates as YYYY-MM-DD
	startDate := currentStart.Format("2006-01-02")
	endDate := currentEnd.Format("2006-01-02")

	return s.GetRevenueStats(ctx, startDate, endDate)
}
