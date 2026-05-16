package controllers

import (
	"net/http"
	"strconv"

	"booking-service/models/dto"
	"booking-service/services"

	"github.com/gin-gonic/gin"
)

// @title Booking Service API
// @version 1.0
// @description API for managing bookings, sessions, entitlements, and certifications
// @host localhost:8082
// @BasePath /api/v1

type BookingController struct {
	bookingService services.IBookingService
}

func NewBookingController(bookingService services.IBookingService) IBookingController {
	return &BookingController{bookingService: bookingService}
}

type IBookingController interface {
	CreateBooking(c *gin.Context)
	GetBooking(c *gin.Context)
	UpdateBooking(c *gin.Context)
	ListBookings(c *gin.Context)
	CancelBooking(c *gin.Context)
	ConfirmBooking(c *gin.Context)
	CompleteBooking(c *gin.Context)
}

// CreateBooking godoc
// @Summary Create a new booking
// @Description Creates a new booking with the provided details
// @Tags bookings
// @Accept json
// @Produce json
// @Param booking body dto.CreateBookingRequest true "Booking data"
// @Success 201 {object} dto.BookingResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /bookings [post]
func (c *BookingController) CreateBooking(ctx *gin.Context) {
	var req dto.CreateBookingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.bookingService.CreateBooking(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// GetBooking godoc
// @Summary Get a booking by ID
// @Description Retrieves a booking by its ID
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} dto.BookingResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /bookings/{id} [get]
func (c *BookingController) GetBooking(ctx *gin.Context) {
	id, err := getUintIDFromPath(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	resp, err := c.bookingService.GetBooking(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateBooking godoc
// @Summary Update a booking
// @Description Updates a booking with the provided details
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Param booking body dto.UpdateBookingRequest true "Booking data"
// @Success 200 {object} dto.BookingResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /bookings/{id} [put]
func (c *BookingController) UpdateBooking(ctx *gin.Context) {
	id, err := getUintIDFromPath(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	var req dto.UpdateBookingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.bookingService.UpdateBooking(ctx.Request.Context(), id, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// ListBookings godoc
// @Summary List all bookings
// @Description Retrieves a paginated list of bookings
// @Tags bookings
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} dto.BookingListResponse
// @Failure 500 {object} map[string]string
// @Router /bookings [get]
func (c *BookingController) ListBookings(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	resp, err := c.bookingService.ListBookings(ctx.Request.Context(), page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// CancelBooking godoc
// @Summary Cancel a booking
// @Description Cancels a booking by its ID
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /bookings/{id}/cancel [post]
func (c *BookingController) CancelBooking(ctx *gin.Context) {
	id, err := getUintIDFromPath(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	if err := c.bookingService.CancelBooking(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "booking cancelled"})
}

// ConfirmBooking godoc
// @Summary Confirm a booking
// @Description Confirms a booking by its ID
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /bookings/{id}/confirm [post]
func (c *BookingController) ConfirmBooking(ctx *gin.Context) {
	id, err := getUintIDFromPath(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	if err := c.bookingService.ConfirmBooking(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "booking confirmed"})
}

// CompleteBooking godoc
// @Summary Complete a booking
// @Description Marks a booking as completed
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /bookings/{id}/complete [post]
func (c *BookingController) CompleteBooking(ctx *gin.Context) {
	id, err := getUintIDFromPath(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	if err := c.bookingService.CompleteBooking(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "booking completed"})
}

// Helper to get uint ID from path
func getUintIDFromPath(c *gin.Context, paramName string) (uint, error) {
	idStr := c.Param(paramName)
	if idStr == "" {
		return 0, nil
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}