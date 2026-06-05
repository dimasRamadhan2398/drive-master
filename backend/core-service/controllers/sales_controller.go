package controllers

import (
	"core-service/models"
	"core-service/models/dto"
	"core-service/pkg/response"
	"core-service/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SalesController struct {
	salesService services.ISalesService
}

type ISalesController interface {
	CreateSale(ctx *gin.Context)
	GetSaleByID(ctx *gin.Context)
	GetSaleByOrderNumber(ctx *gin.Context)
	ListSales(ctx *gin.Context)
	UpdateSaleStatus(ctx *gin.Context)
	DeleteSale(ctx *gin.Context)

	// Analytics
	GetOverview(ctx *gin.Context)
	GetMonthlySales(ctx *gin.Context)
	GetYearlySales(ctx *gin.Context)
	GetSalesTrend(ctx *gin.Context)
	GetSalesBySource(ctx *gin.Context)
	GetSalesByPackageType(ctx *gin.Context)
}

func NewSalesController(salesService services.ISalesService) ISalesController {
	return &SalesController{
		salesService: salesService,
	}
}

// CreateSale handles POST /api/v1/admin/sales
// @Summary Create a new sale
// @Description Creates a new sale/order
// @Tags Sales
// @Accept json
// @Produce json
// @Param request body dto.CreateSaleRequest true "Sale data"
// @Success 201 {object} response.Response
// @Router /admin/sales [post]
func (c *SalesController) CreateSale(ctx *gin.Context) {
	var req dto.CreateSaleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request: "+err.Error())
		return
	}

	sale, err := c.salesService.CreateSale(ctx, req)
	if err != nil {
		response.InternalServerError(ctx, "Failed to create sale: "+err.Error())
		return
	}

	response.Success(ctx, 201, "Sale created successfully", sale)
}

// GetSaleByID handles GET /api/v1/admin/sales/:id
// @Summary Get a sale by ID
// @Description Retrieves a sale by its ID
// @Tags Sales
// @Produce json
// @Param id path string true "Sale ID"
// @Success 200 {object} response.Response
// @Router /admin/sales/{id} [get]
func (c *SalesController) GetSaleByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(ctx, "Invalid sale ID")
		return
	}

	sale, err := c.salesService.GetSaleByID(ctx, id)
	if err != nil {
		response.InternalServerError(ctx, "Failed to get sale: "+err.Error())
		return
	}

	response.OK(ctx, "Sale retrieved successfully", sale)
}

// GetSaleByOrderNumber handles GET /api/v1/admin/sales/order/:orderNumber
// @Summary Get a sale by order number
// @Description Retrieves a sale by its order number
// @Tags Sales
// @Produce json
// @Param orderNumber path string true "Order Number"
// @Success 200 {object} response.Response
// @Router /admin/sales/order/{orderNumber} [get]
func (c *SalesController) GetSaleByOrderNumber(ctx *gin.Context) {
	orderNumber := ctx.Param("orderNumber")
	if orderNumber == "" {
		response.BadRequest(ctx, "Order number is required")
		return
	}

	sale, err := c.salesService.GetSaleByOrderNumber(ctx, orderNumber)
	if err != nil {
		response.InternalServerError(ctx, "Failed to get sale: "+err.Error())
		return
	}

	response.OK(ctx, "Sale retrieved successfully", sale)
}

// ListSales handles GET /api/v1/admin/sales
// @Summary List all sales
// @Description Retrieves a paginated list of sales with optional filters
// @Tags Sales
// @Produce json
// @Param userId query string false "Filter by user ID"
// @Param status query string false "Filter by status"
// @Param source query string false "Filter by source"
// @Param startDate query string false "Filter by start date (YYYY-MM-DD)"
// @Param endDate query string false "Filter by end date (YYYY-MM-DD)"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Success 200 {object} response.Response
// @Router /admin/sales [get]
func (c *SalesController) ListSales(ctx *gin.Context) {
	var req dto.ListSalesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "Invalid request: "+err.Error())
		return
	}

	result, err := c.salesService.ListSales(ctx, req)
	if err != nil {
		response.InternalServerError(ctx, "Failed to list sales: "+err.Error())
		return
	}

	response.OK(ctx, "Sales retrieved successfully", result)
}

// UpdateSaleStatus handles PATCH /api/v1/admin/sales/:id/status
// @Summary Update sale status
// @Description Updates the status of a sale
// @Tags Sales
// @Accept json
// @Produce json
// @Param id path string true "Sale ID"
// @Param request body dto.UpdateSaleStatusRequest true "Status update data"
// @Success 200 {object} response.Response
// @Router /admin/sales/{id}/status [patch]
func (c *SalesController) UpdateSaleStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(ctx, "Invalid sale ID")
		return
	}

	var req dto.UpdateSaleStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request: "+err.Error())
		return
	}

	if err := c.salesService.UpdateSaleStatus(ctx, id, req.Status, req.Notes); err != nil {
		response.InternalServerError(ctx, "Failed to update sale status: "+err.Error())
		return
	}

	response.OK(ctx, "Sale status updated successfully", nil)
}

// DeleteSale handles DELETE /api/v1/admin/sales/:id
// @Summary Delete a sale
// @Description Deletes a sale by its ID
// @Tags Sales
// @Produce json
// @Param id path string true "Sale ID"
// @Success 200 {object} response.Response
// @Router /admin/sales/{id} [delete]
func (c *SalesController) DeleteSale(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(ctx, "Invalid sale ID")
		return
	}

	if err := c.salesService.DeleteSale(ctx, id); err != nil {
		response.InternalServerError(ctx, "Failed to delete sale: "+err.Error())
		return
	}

	response.OK(ctx, "Sale deleted successfully", nil)
}

// GetOverview handles GET /api/v1/admin/sales/analytics/overview
// @Summary Get sales overview
// @Description Retrieves overview statistics for sales within a date range
// @Tags Sales Analytics
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} response.Response
// @Router /admin/sales/analytics/overview [get]
func (c *SalesController) GetOverview(ctx *gin.Context) {
	startDate := ctx.DefaultQuery("start_date", "30daysAgo")
	endDate := ctx.DefaultQuery("end_date", "today")

	data, err := c.salesService.GetOverview(ctx, startDate, endDate)
	if err != nil {
		response.InternalServerError(ctx, "Failed to get sales overview: "+err.Error())
		return
	}

	response.OK(ctx, "Sales overview retrieved successfully", data)
}

// GetMonthlySales handles GET /api/v1/admin/sales/analytics/monthly/:year/:month
// @Summary Get monthly sales data
// @Description Retrieves aggregated sales data for a specific month
// @Tags Sales Analytics
// @Produce json
// @Param year path int true "Year"
// @Param month path int true "Month"
// @Success 200 {object} response.Response
// @Router /admin/sales/analytics/monthly/{year}/{month} [get]
func (c *SalesController) GetMonthlySales(ctx *gin.Context) {
	yearStr := ctx.Param("year")
	monthStr := ctx.Param("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		response.BadRequest(ctx, "Invalid year")
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		response.BadRequest(ctx, "Invalid month")
		return
	}

	data, err := c.salesService.GetMonthlySales(ctx, year, month)
	if err != nil {
		response.InternalServerError(ctx, "Failed to get monthly sales: "+err.Error())
		return
	}

	response.OK(ctx, "Monthly sales retrieved successfully", data)
}

// GetYearlySales handles GET /api/v1/admin/sales/analytics/yearly/:year
// @Summary Get yearly sales data
// @Description Retrieves aggregated sales data for a specific year
// @Tags Sales Analytics
// @Produce json
// @Param year path int true "Year"
// @Success 200 {object} response.Response
// @Router /admin/sales/analytics/yearly/{year} [get]
func (c *SalesController) GetYearlySales(ctx *gin.Context) {
	yearStr := ctx.Param("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		response.BadRequest(ctx, "Invalid year")
		return
	}

	data, err := c.salesService.GetYearlySales(ctx, year)
	if err != nil {
		response.InternalServerError(ctx, "Failed to get yearly sales: "+err.Error())
		return
	}

	response.OK(ctx, "Yearly sales retrieved successfully", data)
}

// GetSalesTrend handles GET /api/v1/admin/sales/analytics/trend
// @Summary Get sales trend
// @Description Retrieves daily sales data for trend analysis
// @Tags Sales Analytics
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} response.Response
// @Router /admin/sales/analytics/trend [get]
func (c *SalesController) GetSalesTrend(ctx *gin.Context) {
	startDate := ctx.DefaultQuery("start_date", "30daysAgo")
	endDate := ctx.DefaultQuery("end_date", "today")

	data, err := c.salesService.GetSalesTrend(ctx, startDate, endDate)
	if err != nil {
		response.InternalServerError(ctx, "Failed to get sales trend: "+err.Error())
		return
	}

	response.OK(ctx, "Sales trend retrieved successfully", data)
}

// GetSalesBySource handles GET /api/v1/admin/sales/analytics/by-source
// @Summary Get sales by source
// @Description Retrieves sales breakdown by source (web, mobile, admin)
// @Tags Sales Analytics
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} response.Response
// @Router /admin/sales/analytics/by-source [get]
func (c *SalesController) GetSalesBySource(ctx *gin.Context) {
	startDate := ctx.DefaultQuery("start_date", "30daysAgo")
	endDate := ctx.DefaultQuery("end_date", "today")

	data, err := c.salesService.GetSalesBySource(ctx, startDate, endDate)
	if err != nil {
		response.InternalServerError(ctx, "Failed to get sales by source: "+err.Error())
		return
	}

	response.OK(ctx, "Sales by source retrieved successfully", data)
}

// GetSalesByPackageType handles GET /api/v1/admin/sales/analytics/by-package-type
// @Summary Get sales by package type
// @Description Retrieves sales breakdown by package type (bronze, silver, gold, platinum)
// @Tags Sales Analytics
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} response.Response
// @Router /admin/sales/analytics/by-package-type [get]
func (c *SalesController) GetSalesByPackageType(ctx *gin.Context) {
	startDate := ctx.DefaultQuery("start_date", "30daysAgo")
	endDate := ctx.DefaultQuery("end_date", "today")

	data, err := c.salesService.GetSalesByPackageType(ctx, startDate, endDate)
	if err != nil {
		response.InternalServerError(ctx, "Failed to get sales by package type: "+err.Error())
		return
	}

	response.OK(ctx, "Sales by package type retrieved successfully", data)
}

// Helper function to convert SaleStatus
func parseSaleStatus(s string) models.SaleStatus {
	return models.SaleStatus(s)
}
