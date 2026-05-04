package controllers

import (
	"net/http"
	"user-service/models/dto"
	apperrors "user-service/pkg/errors"
	responseRes "user-service/pkg/response"
	"user-service/services"

	"github.com/gin-gonic/gin"
)

type CoverageAreaController struct {
	coverageAreaService services.ICoverageAreaService
}

type ICoverageAreaController interface {
	AddCoverageArea(ctx *gin.Context)
	RemoveCoverageArea(ctx *gin.Context)
	GetCoverageAreas(ctx *gin.Context)
}

func NewCoverageAreaController(
	coverageAreaService services.ICoverageAreaService,
) ICoverageAreaController {
	return &CoverageAreaController{
		coverageAreaService: coverageAreaService,
	}
}

// @Summary Add Coverage Area
// @Description Add coverage area for an instructor
// @Tags Instructors
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param request body dto.AddCoverageAreaInput true "Coverage area data"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /instructors/{id}/coverage-areas [post]
func (c *CoverageAreaController) AddCoverageArea(ctx *gin.Context) {
	userID, err := getUserIDFromPath(ctx, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	var input dto.AddCoverageAreaInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	// TODO: Look up area ID by name from location service
	// For now, use a placeholder area ID
	areaID := uint(1)

	if _, err := c.coverageAreaService.AddCoverageArea(ctx, userID, areaID); err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusCreated, "Coverage area added successfully", nil)
}

// @Summary Remove Coverage Area
// @Description Remove coverage area for an instructor
// @Tags Instructors
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param areaId path int true "Coverage Area ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /instructors/{id}/coverage-areas/{areaId} [delete]
func (c *CoverageAreaController) RemoveCoverageArea(ctx *gin.Context) {
	userID, err := getUserIDFromPath(ctx, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	areaID, err := getUintIDFromPath(ctx, "areaId")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	if err := c.coverageAreaService.RemoveCoverageArea(ctx, userID, areaID); err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Coverage area removed successfully", nil)
}

// @Summary Get Coverage Areas
// @Description Get all coverage areas for an instructor
// @Tags Instructors
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /instructors/{id}/coverage-areas [get]
func (c *CoverageAreaController) GetCoverageAreas(ctx *gin.Context) {
	userID, err := getUserIDFromPath(ctx, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	areas, err := c.coverageAreaService.GetCoverageAreas(ctx, userID)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Coverage areas retrieved successfully", areas)
}
