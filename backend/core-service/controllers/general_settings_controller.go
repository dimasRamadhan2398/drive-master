package controllers

import (
	"core-service/models/dto"
	"core-service/pkg/response"
	"core-service/services"

	"github.com/gin-gonic/gin"
)

// GeneralSettingsController handles general settings-related HTTP requests
type GeneralSettingsController struct {
	service services.IGeneralSettingsService
}

// IGeneralSettingsController defines the interface for general settings controller
type IGeneralSettingsController interface {
	GetSettings(ctx *gin.Context)
	UpdateSettings(ctx *gin.Context)
}

// NewGeneralSettingsController creates a new general settings controller
func NewGeneralSettingsController(service services.IGeneralSettingsService) IGeneralSettingsController {
	return &GeneralSettingsController{
		service: service,
	}
}

// GetSettings handles GET /api/v1/general-settings
// @Summary Get general settings
// @Description Retrieves the general business settings
// @Tags General Settings
// @Produce json
// @Success 200 {object} response.Response
// @Router /general-settings [get]
func (c *GeneralSettingsController) GetSettings(ctx *gin.Context) {
	settings, err := c.service.GetSettings(ctx.Request.Context())
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch general settings: "+err.Error())
		return
	}

	response.OK(ctx, "General settings fetched successfully", dto.ToGeneralSettingsResponse(settings))
}

// UpdateSettings handles PUT /api/v1/general-settings
// @Summary Update general settings
// @Description Updates the general business settings
// @Tags General Settings
// @Accept json
// @Produce json
// @Param request body dto.UpdateGeneralSettingsRequest true "General settings data"
// @Success 200 {object} response.Response
// @Router /general-settings [put]
func (c *GeneralSettingsController) UpdateSettings(ctx *gin.Context) {
	var req dto.UpdateGeneralSettingsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body: "+err.Error())
		return
	}

	settings, err := c.service.UpdateSettings(ctx.Request.Context(), &req)
	if err != nil {
		response.InternalServerError(ctx, "Failed to update general settings: "+err.Error())
		return
	}

	response.OK(ctx, "General settings updated successfully", dto.ToGeneralSettingsResponse(settings))
}