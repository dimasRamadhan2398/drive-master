package controllers

import (
	"net/http"

	"booking-service/services"

	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	sessionService services.ISessionService
	revenueService services.IRevenueService
}

func NewDashboardController(
	sessionService services.ISessionService,
	revenueService services.IRevenueService,
) IDashboardController {
	return &DashboardController{
		sessionService:  sessionService,
		revenueService: revenueService,
	}
}

type IDashboardController interface {
	GetSessionStats(ctx *gin.Context)
	GetRevenueStats(ctx *gin.Context)
}

func (c *DashboardController) GetSessionStats(ctx *gin.Context) {
	stats, err := c.sessionService.GetStats(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Session stats retrieved successfully",
		"data":    stats,
	})
}

func (c *DashboardController) GetRevenueStats(ctx *gin.Context) {
	// Get current month's revenue stats with growth from last month
	stats, err := c.revenueService.GetCurrentMonthRevenue(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Revenue stats retrieved successfully",
		"data":    stats,
	})
}

// StatsResponse for unified response format
type StatsResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewStatsResponse(data interface{}) StatsResponse {
	return StatsResponse{
		Success: true,
		Message: "Stats retrieved successfully",
		Data:    data,
	}
}

// SessionStatsResponse from booking-service
type SessionStatsResponse struct {
	TotalSessions     int64 `json:"totalSessions"`
	ActiveSessions    int64 `json:"activeSessions"`
	CompletedSessions int64 `json:"completedSessions"`
	PendingSessions   int64 `json:"pendingSessions"`
}

// CertificationStatsResponse from booking-service
type CertificationStatsResponse struct {
	TotalCertifications    int64 `json:"totalCertifications"`
	IssuedCertifications   int64 `json:"issuedCertifications"`
	ActiveCertifications   int64 `json:"activeCertifications"`
	RevokedCertifications  int64 `json:"revokedCertifications"`
}

// UserServiceStatsResponse from user-service
type UserServiceStatsResponse struct {
	TotalUsers          int64 `json:"totalUsers"`
	TotalMembers        int64 `json:"totalMembers"`
	TotalInstructors    int64 `json:"totalInstructors"`
	RecentRegistrations int64 `json:"recentRegistrations"`
}

// SalesStatsResponse from core-service
type SalesStatsResponse struct {
	TotalRevenue   int64  `json:"totalRevenue"`
	TotalSales     int64  `json:"totalSales"`
	TotalRefunds   int64  `json:"totalRefunds"`
	GrowthRate     int64  `json:"growthRate"`
	Currency       string `json:"currency"`
}