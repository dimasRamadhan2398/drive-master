package controllers

import (
	"net/http"
	"strconv"
	"user-service/models"
	"user-service/models/dto"
	apperrors "user-service/pkg/errors"
	responseRes "user-service/pkg/response"
	"user-service/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type WorkExperienceController struct {
	workExperienceService services.IWorkExperienceService
}

// @Summary Create Work Experience
// @Description Add work experience for an instructor
// @Tags Instructors
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param request body dto.CreateWorkExperienceRequest true "Work experience data"
// @Success 201 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Failure 404 {object} responseRes.Response
// @Router /instructors/{id}/work-experiences [post]
func (w *WorkExperienceController) CreateWorkExperience(ctx *gin.Context) {
	var instructorIDStr = ctx.Param("id")
	instructorID, err := uuid.Parse(instructorIDStr)
	if err != nil {
		responseRes.Error(ctx, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid instructor ID format", "")
		return
	}

	var input dto.CreateWorkExperienceRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		responseRes.Error(ctx, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity), err.Error(), "")
		return
	}

	result, err := w.workExperienceService.CreateWorkExperience(instructorID, input)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusCreated, "Work experience created successfully", result)
}

// @Summary Delete Work Experience
// @Description Delete a work experience
// @Tags Instructors
// @Produce json
// @Param expId path int true "Work Experience ID"
// @Success 200 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Failure 404 {object} responseRes.Response
// @Router /work-experiences/{expId} [delete]
func (w *WorkExperienceController) DeleteWorkExperience(ctx *gin.Context) {
	var expIdStr = ctx.Param("id")
	expIdNum, err := strconv.ParseUint(expIdStr, 10, 64)
	if err != nil {
		responseRes.Error(ctx, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid work experience ID format", "")
		return
	}

	err = w.workExperienceService.DeleteWorkExperience(uint(expIdNum))
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Work experience deleted successfully", nil)
}

// @Summary Get Work Experience
// @Description Get work experience by ID
// @Tags Instructors
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Failure 404 {object} responseRes.Response
// @Router /instructors/{id}/work-experiences [get]
func (w *WorkExperienceController) GetWorkExperience(ctx *gin.Context) {
	var expId = ctx.Param("id")
	if len(expId) < 10 {
		responseRes.Error(ctx, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid work experience ID format", "")
		return
	}

	expIdUuid, err := uuid.Parse(expId)
	if err != nil {
		responseRes.Error(ctx, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid work experience ID format", "")
		return
	}

	result, err := w.workExperienceService.GetWorkExperiences(expIdUuid)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Work experience retrieved successfully", result)
}

// UpdateWorkExperience implements [IWorkExperienceController].
// @Summary Update Work Experience
// @Description Update a work experience
// @Tags Instructors
// @Accept json
// @Produce json
// @Param expId path int true "Work Experience ID"
// @Param request body dto.CreateWorkExperienceRequest true "Update work experience data"
// @Success 200 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Failure 404 {object} responseRes.Response
// @Router /work-experiences/{expId} [put]
func (w *WorkExperienceController) UpdateWorkExperience(ctx *gin.Context) {
	expID, err := getUintIDFromPath(ctx, "expId")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	var input dto.UpdateWorkExperienceRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	// TODO: Get work experience by ID and update
	workExp := &models.WorkExperience{
		ID:           expID,
		CompanyName:  *input.CompanyName,
		Role:         *input.Role,
		StartDate:    *input.StartDate,
		EndDate:      input.EndDate,
		Description:  *input.Description,
	}

	if err := w.workExperienceService.UpdateWorkExperience(workExp); err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Work experience updated successfully", workExp)
}

type IWorkExperienceController interface {
	CreateWorkExperience(c *gin.Context)
	GetWorkExperience(c *gin.Context)
	UpdateWorkExperience(c *gin.Context)
	DeleteWorkExperience(c *gin.Context)
}

func NewWorkExperienceController(workExperienceService services.IWorkExperienceService) IWorkExperienceController {
	return &WorkExperienceController{
		workExperienceService: workExperienceService,
	}
}