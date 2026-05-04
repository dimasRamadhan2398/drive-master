package controllers

import (
	"net/http"
	"strconv"

	"booking-service/models/dto"
	"booking-service/services"

	"github.com/gin-gonic/gin"
)

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