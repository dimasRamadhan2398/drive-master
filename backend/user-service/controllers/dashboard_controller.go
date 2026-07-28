package controllers

import (
	"net/http"

	responseRes "user-service/pkg/response"
	"user-service/services"

	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	dashboardService     services.IDashboardService
	certificationService services.ICertificationService
}

func NewDashboardController(
	dashboardService services.IDashboardService,
	certificationService services.ICertificationService,
) IDashboardController {
	return &DashboardController{
		dashboardService:     dashboardService,
		certificationService: certificationService,
	}
}

type IDashboardController interface {
	GetStats(ctx *gin.Context)
	GetCertificationStats(ctx *gin.Context)
}

func (c *DashboardController) GetStats(ctx *gin.Context) {
	stats, err := c.dashboardService.GetStats(ctx.Request.Context())
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Dashboard stats retrieved successfully", stats)
}

func (c *DashboardController) GetCertificationStats(ctx *gin.Context) {
	stats, err := c.certificationService.GetCertificateStats(ctx.Request.Context())
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Certification stats retrieved successfully", stats)
}
