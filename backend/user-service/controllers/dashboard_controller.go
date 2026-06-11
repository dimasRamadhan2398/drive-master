package controllers

import (
	"net/http"

	responseRes "user-service/pkg/response"
	"user-service/services"

	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	dashboardService services.IDashboardService
}

func NewDashboardController(dashboardService services.IDashboardService) IDashboardController {
	return &DashboardController{
		dashboardService: dashboardService,
	}
}

type IDashboardController interface {
	GetStats(ctx *gin.Context)
}

func (c *DashboardController) GetStats(ctx *gin.Context) {
	stats, err := c.dashboardService.GetStats(ctx.Request.Context())
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Dashboard stats retrieved successfully", stats)
}