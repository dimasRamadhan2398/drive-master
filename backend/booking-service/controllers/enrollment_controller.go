package controllers

import (
	"net/http"
	"strconv"

	"booking-service/models/dto"
	"booking-service/pkg/base"
	"booking-service/services"

	"github.com/gin-gonic/gin"
)

// @title Booking Service API
// @version 1.0
// @description API for managing bookings, sessions, entitlements, and certifications
// @host localhost:8082
// @BasePath /api/v1

type EnrollmentController struct {
	enrollmentService services.IEnrollmentService
}

func NewEnrollmentController(enrollmentService services.IEnrollmentService) IEnrollmentController {
	return &EnrollmentController{enrollmentService: enrollmentService}
}

type IEnrollmentController interface {
	CreateEnrollment(c *gin.Context)
	GetEnrollment(c *gin.Context)
	UpdateEnrollment(c *gin.Context)
	CancelEnrollment(c *gin.Context)
	MarkAsPaid(c *gin.Context)
	ListEnrollments(c *gin.Context)
	ListUserEnrollments(c *gin.Context)
}

// CreateEnrollment godoc
// @Summary Create a new enrollment
// @Description Creates a new enrollment for a user package purchase
// @Tags enrollments
// @Accept json
// @Produce json
// @Param enrollment body dto.CreateEnrollmentRequest true "Enrollment data"
// @Success 201 {object} dto.EnrollmentResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /enrollments [post]
func (c *EnrollmentController) CreateEnrollment(ctx *gin.Context) {
	var req dto.CreateEnrollmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.enrollmentService.CreateEnrollment(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// GetEnrollment godoc
// @Summary Get an enrollment by ID
// @Description Retrieves an enrollment by its ID
// @Tags enrollments
// @Accept json
// @Produce json
// @Param id path string true "Enrollment ID"
// @Success 200 {object} dto.EnrollmentResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /enrollments/{id} [get]
func (c *EnrollmentController) GetEnrollment(ctx *gin.Context) {
	id, err := base.GetUUIDIDFromPath(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid enrollment id"})
		return
	}

	resp, err := c.enrollmentService.GetEnrollment(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateEnrollment godoc
// @Summary Update an enrollment
// @Description Updates an enrollment with the provided details
// @Tags enrollments
// @Accept json
// @Produce json
// @Param id path string true "Enrollment ID"
// @Param enrollment body dto.UpdateEnrollmentRequest true "Enrollment data"
// @Success 200 {object} dto.EnrollmentResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /enrollments/{id} [put]
func (c *EnrollmentController) UpdateEnrollment(ctx *gin.Context) {
	id, err := base.GetUUIDIDFromPath(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid enrollment id"})
		return
	}

	var req dto.UpdateEnrollmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.enrollmentService.UpdateEnrollment(ctx.Request.Context(), id, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// CancelEnrollment godoc
// @Summary Cancel an enrollment
// @Description Cancels an enrollment by its ID
// @Tags enrollments
// @Accept json
// @Produce json
// @Param id path string true "Enrollment ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /enrollments/{id}/cancel [post]
func (c *EnrollmentController) CancelEnrollment(ctx *gin.Context) {
	id, err := base.GetUUIDIDFromPath(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid enrollment id"})
		return
	}

	if err := c.enrollmentService.CancelEnrollment(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "enrollment cancelled"})
}

// MarkAsPaid godoc
// @Summary Mark enrollment as paid
// @Description Marks an enrollment as paid and activates it
// @Tags enrollments
// @Accept json
// @Produce json
// @Param id path string true "Enrollment ID"
// @Param body body map[string]interface{} true "Payment details"
// @Success 200 {object} dto.EnrollmentResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /enrollments/{id}/pay [post]
func (c *EnrollmentController) MarkAsPaid(ctx *gin.Context) {
	id, err := base.GetUUIDIDFromPath(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid enrollment id"})
		return
	}

	var req struct {
		TotalPrice float64 `json:"totalPrice" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.enrollmentService.MarkAsPaid(ctx.Request.Context(), id, req.TotalPrice)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// ListEnrollments godoc
// @Summary List all enrollments
// @Description Retrieves a paginated list of enrollments
// @Tags enrollments
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} dto.EnrollmentListResponse
// @Failure 500 {object} map[string]string
// @Router /enrollments [get]
func (c *EnrollmentController) ListEnrollments(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	resp, err := c.enrollmentService.ListEnrollments(ctx.Request.Context(), page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// ListUserEnrollments godoc
// @Summary List user enrollments
// @Description Retrieves enrollments for a specific user
// @Tags enrollments
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} dto.EnrollmentListResponse
// @Failure 500 {object} map[string]string
// @Router /enrollments/user/{userId} [get]
func (c *EnrollmentController) ListUserEnrollments(ctx *gin.Context) {
	userID, err := base.GetUUIDIDFromPath(ctx, "userId")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	resp, err := c.enrollmentService.ListUserEnrollments(ctx.Request.Context(), userID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}