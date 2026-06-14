package controllers

import (
	"net/http"
	"strconv"

	"booking-service/models/dto"
	"booking-service/pkg/base"
	"booking-service/services"

	"github.com/gin-gonic/gin"
)

type ScheduleController struct {
	scheduleService services.IScheduleService
}

func NewScheduleController(scheduleService services.IScheduleService) IScheduleController {
	return &ScheduleController{scheduleService: scheduleService}
}

type IScheduleController interface {
	CreateSchedule(c *gin.Context)
	GetSchedule(c *gin.Context)
	UpdateSchedule(c *gin.Context)
	DeleteSchedule(c *gin.Context)
	ListSchedules(c *gin.Context)
	ListSchedulesFiltered(c *gin.Context)
	GetAvailableSchedules(c *gin.Context)
	BookSlot(c *gin.Context)
	CancelBooking(c *gin.Context)
}

// CreateSchedule godoc
// @Summary Create a new schedule slot
// @Description Creates a new schedule slot for instructor/car availability
// @Tags schedules
// @Accept json
// @Produce json
// @Param schedule body dto.CreateScheduleRequest true "Schedule data"
// @Success 201 {object} dto.ScheduleResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /schedules [post]
func (c *ScheduleController) CreateSchedule(ctx *gin.Context) {
	var req dto.CreateScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.scheduleService.CreateSchedule(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// GetSchedule godoc
// @Summary Get a schedule by ID
// @Description Retrieves a schedule slot by its ID
// @Tags schedules
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 200 {object} dto.ScheduleResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /schedules/{id} [get]
func (c *ScheduleController) GetSchedule(ctx *gin.Context) {
	id, err := base.GetUintIDFromPath(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule id"})
		return
	}

	resp, err := c.scheduleService.GetSchedule(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateSchedule godoc
// @Summary Update a schedule
// @Description Updates a schedule slot with the provided details
// @Tags schedules
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Param schedule body dto.UpdateScheduleRequest true "Schedule data"
// @Success 200 {object} dto.ScheduleResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /schedules/{id} [put]
func (c *ScheduleController) UpdateSchedule(ctx *gin.Context) {
	id, err := base.GetUintIDFromPath(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule id"})
		return
	}

	var req dto.UpdateScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.scheduleService.UpdateSchedule(ctx.Request.Context(), id, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// DeleteSchedule godoc
// @Summary Delete a schedule
// @Description Deletes a schedule slot by its ID (only available slots can be deleted)
// @Tags schedules
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /schedules/{id} [delete]
func (c *ScheduleController) DeleteSchedule(ctx *gin.Context) {
	id, err := base.GetUintIDFromPath(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule id"})
		return
	}

	if err := c.scheduleService.DeleteSchedule(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "schedule deleted"})
}

// ListSchedules godoc
// @Summary List all schedules
// @Description Retrieves a paginated list of schedules
// @Tags schedules
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} dto.ScheduleListResponse
// @Failure 500 {object} map[string]string
// @Router /schedules [get]
func (c *ScheduleController) ListSchedules(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	resp, err := c.scheduleService.ListSchedules(ctx.Request.Context(), page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// ListSchedulesFiltered godoc
// @Summary List schedules with filters
// @Description Retrieves a filtered list of schedules
// @Tags schedules
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param date query string false "Specific date (YYYY-MM-DD)"
// @Param startDate query string false "Start date (YYYY-MM-DD)"
// @Param endDate query string false "End date (YYYY-MM-DD)"
// @Param instructorId query int false "Filter by instructor"
// @Param carId query int false "Filter by car"
// @Param status query string false "Filter by status"
// @Success 200 {object} dto.ScheduleListResponse
// @Failure 500 {object} map[string]string
// @Router /schedules/filter [get]
func (c *ScheduleController) ListSchedulesFiltered(ctx *gin.Context) {
	params := dto.ScheduleFilterParams{
		ListParams: dto.ListParams{
			Page:  1,
			Limit: 10,
		},
	}

	if page, err := strconv.Atoi(ctx.DefaultQuery("page", "1")); err == nil {
		params.Page = page
	}
	if limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10")); err == nil {
		params.Limit = limit
	}

	params.Date = ctx.Query("date")
	params.StartDate = ctx.Query("startDate")
	params.EndDate = ctx.Query("endDate")

	params.InstructorID = ctx.Query("instructorId")
	if carID, err := strconv.Atoi(ctx.Query("carId")); err == nil {
		params.CarID = uint(carID)
	}

	params.Status = ctx.Query("status")

	resp, err := c.scheduleService.ListSchedulesFiltered(ctx.Request.Context(), params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetAvailableSchedules godoc
// @Summary Get available schedules
// @Description Retrieves available schedules within a date range
// @Tags schedules
// @Accept json
// @Produce json
// @Param startDate query string true "Start date (YYYY-MM-DD)"
// @Param endDate query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} dto.ScheduleListResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /schedules/available [get]
func (c *ScheduleController) GetAvailableSchedules(ctx *gin.Context) {
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	if startDate == "" || endDate == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "startDate and endDate are required"})
		return
	}

	resp, err := c.scheduleService.GetAvailableSchedules(ctx.Request.Context(), startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// BookSlot godoc
// @Summary Book a schedule slot
// @Description Books an available schedule slot for a user
// @Tags schedules
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Param booking body dto.BookSlotRequest true "Booking data"
// @Success 200 {object} dto.ScheduleResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /schedules/{id}/book [post]
func (c *ScheduleController) BookSlot(ctx *gin.Context) {
	id, err := base.GetUintIDFromPath(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule id"})
		return
	}

	var req dto.BookSlotRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.scheduleService.BookSlot(ctx.Request.Context(), id, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// CancelBooking godoc
// @Summary Cancel a booking
// @Description Cancels a schedule booking and releases the slot
// @Tags schedules
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /schedules/{id}/cancel [post]
func (c *ScheduleController) CancelBooking(ctx *gin.Context) {
	id, err := base.GetUintIDFromPath(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule id"})
		return
	}

	if err := c.scheduleService.CancelBooking(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "booking cancelled"})
}