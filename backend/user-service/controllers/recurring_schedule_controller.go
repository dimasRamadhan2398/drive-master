package controllers

import (
	"net/http"

	"user-service/models/dto"
	apperrors "user-service/pkg/errors"
	responseRes "user-service/pkg/response"
	"user-service/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RecurringScheduleController struct {
	recurringScheduleService services.IRecurringScheduleService
}

type IRecurringScheduleController interface {
	CreateRecurringSchedule(ctx *gin.Context)
	BulkCreateRecurringSchedules(ctx *gin.Context)
	GetRecurringSchedules(ctx *gin.Context)
	GetRecurringScheduleByID(ctx *gin.Context)
	UpdateRecurringSchedule(ctx *gin.Context)
	DeleteRecurringSchedule(ctx *gin.Context)
	DeleteAllRecurringSchedules(ctx *gin.Context)
}

func NewRecurringScheduleController(
	recurringScheduleService services.IRecurringScheduleService,
) IRecurringScheduleController {
	return&RecurringScheduleController{
		recurringScheduleService: recurringScheduleService,
	}
}

// @Summary Create Recurring Schedule
// @Description Create a single recurring time slot for an instructor
// @Tags Recurring Schedules
// @Accept json
// @Produce json
// @Param id path string true "Instructor User ID (UUID)"
// @Param request body dto.CreateRecurringScheduleRequest true "Recurring schedule data"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /instructors/{id}/recurring-schedules [post]
func (c *RecurringScheduleController) CreateRecurringSchedule(ctx *gin.Context) {
	instructorID, err := getUserIDFromPath(ctx, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	var req dto.CreateRecurringScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	result, err := c.recurringScheduleService.CreateRecurringSchedule(ctx.Request.Context(), instructorID, req)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusCreated, "Recurring schedule created successfully", result)
}

// @Summary Bulk Create Recurring Schedules
// @Description Create multiple recurring time slots at once for an instructor
// @Description Example: Create Mon-Fri 09:00-10:00 and 13:00-14:00
// @Tags Recurring Schedules
// @Accept json
// @Produce json
// @Param id path string true "Instructor User ID (UUID)"
// @Param request body dto.BulkCreateRecurringScheduleRequest true "Bulk recurring schedule data"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /instructors/{id}/recurring-schedules/bulk [post]
func (c *RecurringScheduleController) BulkCreateRecurringSchedules(ctx *gin.Context) {
	instructorID, err := getUserIDFromPath(ctx, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	var req dto.BulkCreateRecurringScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	result, err := c.recurringScheduleService.BulkCreateRecurringSchedules(ctx.Request.Context(), instructorID, req)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusCreated, "Recurring schedules created successfully", result)
}

// @Summary Get Recurring Schedules
// @Description Get all recurring schedules for an instructor
// @Tags Recurring Schedules
// @Produce json
// @Param id path string true "Instructor User ID (UUID)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /instructors/{id}/recurring-schedules [get]
func (c *RecurringScheduleController) GetRecurringSchedules(ctx *gin.Context) {
	instructorID, err := getUserIDFromPath(ctx, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	result, err := c.recurringScheduleService.GetRecurringSchedules(ctx.Request.Context(), instructorID)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Recurring schedules retrieved successfully", result)
}

// @Summary Get Recurring Schedule by ID
// @Description Get a single recurring schedule by ID
// @Tags Recurring Schedules
// @Produce json
// @Param id path string true "Instructor User ID (UUID)"
// @Param scheduleId path string true "Recurring Schedule ID (UUID)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /instructors/{id}/recurring-schedules/{scheduleId} [get]
func (c *RecurringScheduleController) GetRecurringScheduleByID(ctx *gin.Context) {
	scheduleID, err := uuid.Parse(ctx.Param("scheduleId"))
	if err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	result, err := c.recurringScheduleService.GetRecurringScheduleByID(ctx.Request.Context(), scheduleID)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Recurring schedule retrieved successfully", result)
}

// @Summary Update Recurring Schedule
// @Description Update a recurring schedule
// @Tags Recurring Schedules
// @Accept json
// @Produce json
// @Param id path string true "Instructor User ID (UUID)"
// @Param scheduleId path string true "Recurring Schedule ID (UUID)"
// @Param request body dto.UpdateRecurringScheduleRequest true "Update recurring schedule data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /instructors/{id}/recurring-schedules/{scheduleId} [put]
func (c *RecurringScheduleController) UpdateRecurringSchedule(ctx *gin.Context) {
	scheduleID, err := uuid.Parse(ctx.Param("scheduleId"))
	if err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	var req dto.UpdateRecurringScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	result, err := c.recurringScheduleService.UpdateRecurringSchedule(ctx.Request.Context(), scheduleID, req)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Recurring schedule updated successfully", result)
}

// @Summary Delete Recurring Schedule
// @Description Delete a single recurring schedule
// @Tags Recurring Schedules
// @Produce json
// @Param id path string true "Instructor User ID (UUID)"
// @Param scheduleId path string true "Recurring Schedule ID (UUID)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /instructors/{id}/recurring-schedules/{scheduleId} [delete]
func (c *RecurringScheduleController) DeleteRecurringSchedule(ctx *gin.Context) {
	scheduleID, err := uuid.Parse(ctx.Param("scheduleId"))
	if err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	if err := c.recurringScheduleService.DeleteRecurringSchedule(ctx.Request.Context(), scheduleID); err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Recurring schedule deleted successfully", nil)
}

// @Summary Delete All Recurring Schedules
// @Description Delete all recurring schedules for an instructor
// @Tags Recurring Schedules
// @Produce json
// @Param id path string true "Instructor User ID (UUID)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /instructors/{id}/recurring-schedules [delete]
func (c *RecurringScheduleController) DeleteAllRecurringSchedules(ctx *gin.Context) {
	instructorID, err := getUserIDFromPath(ctx, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	if err := c.recurringScheduleService.DeleteAllRecurringSchedules(ctx.Request.Context(), instructorID); err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "All recurring schedules deleted successfully", nil)
}
