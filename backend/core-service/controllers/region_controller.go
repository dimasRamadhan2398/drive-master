package controllers

import (
	"core-service/services"
	"core-service/pkg/response"

	"github.com/gin-gonic/gin"
)

// RegionController handles region-related HTTP requests
type RegionController struct {
	regionService services.IRegionService
}

// IRegionController defines the interface for region controller
type IRegionController interface {
	GetAllProvinces(ctx *gin.Context)
	GetRegenciesByProvince(ctx *gin.Context)
	GetDistrictsByRegency(ctx *gin.Context)
}

// NewRegionController creates a new region controller
func NewRegionController(regionService services.IRegionService) IRegionController {
	return &RegionController{
		regionService: regionService,
	}
}

// GetAllProvinces handles GET /api/v1/regions/provinces
// @Summary Get all provinces
// @Description Retrieves all provinces
// @Tags Regions
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/regions/provinces [get]
func (c *RegionController) GetAllProvinces(ctx *gin.Context) {
	provinces, err := c.regionService.GetAllProvinces()
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch provinces")
		return
	}

	response.OK(ctx, "Provinces fetched successfully", provinces)
}

// GetRegenciesByProvince handles GET /api/v1/regions/provinces/:province/regencies
// @Summary Get regencies by province
// @Description Retrieves all regencies for a specific province
// @Tags Regions
// @Produce json
// @Param province path string true "Province ID"
// @Success 200 {object} response.Response
// @Router /api/v1/regions/provinces/{province}/regencies [get]
func (c *RegionController) GetRegenciesByProvince(ctx *gin.Context) {
	province := ctx.Param("province")
	if province == "" {
		response.BadRequest(ctx, "Province ID is required")
		return
	}

	regencies, err := c.regionService.GetRegenciesByProvince(province)
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch regencies")
		return
	}

	response.OK(ctx, "Regencies fetched successfully", regencies)
}

// GetDistrictsByRegency handles GET /api/v1/regions/regencies/:regency/districts
// @Summary Get districts by regency
// @Description Retrieves all districts for a specific regency
// @Tags Regions
// @Produce json
// @Param province query string true "Province ID"
// @Param regency path string true "Regency ID"
// @Success 200 {object} response.Response
// @Router /api/v1/regions/regencies/{regency}/districts [get]
func (c *RegionController) GetDistrictsByRegency(ctx *gin.Context) {
	province := ctx.Query("province")
	regency := ctx.Param("regency")

	if regency == "" {
		response.BadRequest(ctx, "Regency ID is required")
		return
	}

	if province == "" {
		response.BadRequest(ctx, "Province ID is required as query parameter")
		return
	}

	districts, err := c.regionService.GetDistrictsByRegency(province, regency)
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch districts")
		return
	}

	response.OK(ctx, "Districts fetched successfully", districts)
}