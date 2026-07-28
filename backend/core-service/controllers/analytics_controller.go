package controllers

import (
	"core-service/pkg/response"
	"core-service/services"

	"github.com/gin-gonic/gin"
)

type AnalyticsController struct {
	analyticsService services.IAnalyticsService
}

type IAnalyticsController interface {
	GetOverview(ctx *gin.Context)
	GetFunnel(ctx *gin.Context)
}

func NewAnalyticsController(analyticsService services.IAnalyticsService) IAnalyticsController {
	return &AnalyticsController{
		analyticsService: analyticsService,
	}
}

// GetOverview handles GET /api/v1/admin/analytics/overview
// @Summary Get Google Analytics 4 overview data
// @Description Retrieves active users and page views over time
// @Tags Analytics
// @Produce json
// @Param start_date query string false "Start date (e.g. 30daysAgo)"
// @Param end_date query string false "End date (e.g. today)"
// @Success 200 {object} response.Response
// @Router /admin/analytics/overview [get]
func (c *AnalyticsController) GetOverview(ctx *gin.Context) {
	startDate := ctx.DefaultQuery("start_date", "30daysAgo")
	endDate := ctx.DefaultQuery("end_date", "today")

	data, err := c.analyticsService.GetOverviewReport(ctx, startDate, endDate)
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch GA4 overview report: "+err.Error())
		return
	}

	response.OK(ctx, "GA4 Overview report fetched successfully", data)
}

// GetFunnel handles GET /api/v1/admin/analytics/funnel
// @Summary Get Google Analytics 4 conversion funnel
// @Description Retrieves user progress from page_view to purchase
// @Tags Analytics
// @Produce json
// @Success 200 {object} response.Response
// @Router /admin/analytics/funnel [get]
func (c *AnalyticsController) GetFunnel(ctx *gin.Context) {
	data, err := c.analyticsService.GetFunnelReport(ctx)
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch GA4 funnel report: "+err.Error())
		return
	}

	response.OK(ctx, "GA4 Funnel report fetched successfully", data)
}
