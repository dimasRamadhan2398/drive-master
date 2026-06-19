package controllers

import (
	"net/http"

	"user-service/models"
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
// @Description Add coverage area for an instructor by specifying area type and ID
// @Tags Instructors
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param request body dto.AddCoverageAreaInput true "Coverage area data (areaType: province/regency/district, areaId: ID from core-service)"
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

	areaType := models.AreaType(input.AreaType)

	if _, err := c.coverageAreaService.AddCoverageArea(ctx, userID, areaType, input.AreaID); err != nil {
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
// @Param areaType path string true "Area Type (province/regency/district)"
// @Param areaId path int true "Area ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /instructors/{id}/coverage-areas/{areaType}/{areaId} [delete]
func (c *CoverageAreaController) RemoveCoverageArea(ctx *gin.Context) {
	userID, err := getUserIDFromPath(ctx, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	areaType := models.AreaType(ctx.Param("areaType"))
	areaID, err := getUintIDFromPath(ctx, "areaId")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	if err := c.coverageAreaService.RemoveCoverageArea(ctx, userID, areaType, areaID); err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Coverage area removed successfully", nil)
}

// @Summary Get Coverage Areas
// @Description Get all coverage areas for an instructor with region details
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

	// Convert to response DTO
	response := make([]dto.CoverageAreaResponse, 0, len(areas))
	for _, area := range areas {
		response = append(response, dto.CoverageAreaResponse{
			InstructorID: area.InstructorID,
			AreaType:     string(area.AreaType),
			AreaID:       area.AreaID,
			AreaName:     area.AreaName,
		})
	}

	responseRes.Success(ctx, http.StatusOK, "Coverage areas retrieved successfully", response)
}

// func getUserIDFromPath(ctx *gin.Context, param string) (uuid.UUID, error) {
// 	idStr := ctx.Param(param)
// 	id, err := uuid.Parse(idStr)
// 	if err != nil {
// 		return uuid.Nil, err
// 	}
// 	return id, nil
// }

// func getUintIDFromPath(ctx *gin.Context, param string) (uint, error) {
// 	idStr := ctx.Param(param)
// 	var id uint
// 	if _, err := parseUint(idStr, &id); err != nil {
// 		return 0, err
// 	}
// 	return id, nil
// }

func parseUint(s string, result *uint) (bool, error) {
	var n uint
	for _, c := range s {
		if c < '0' || c > '9' {
			return false, nil
		}
		n = n*10 + uint(c-'0')
	}
	*result = n
	return true, nil
}
