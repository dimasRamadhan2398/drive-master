package controllers

import (
	"net/http"
	"strconv"

	"booking-service/models/dto"
	"booking-service/pkg/base"
	"booking-service/services"

	"github.com/gin-gonic/gin"
)

type EntitlementController struct {
	entitlementService services.IEntitlementService
}

func NewEntitlementController(entitlementService services.IEntitlementService) IEntitlementController {
	return &EntitlementController{entitlementService: entitlementService}
}

type IEntitlementController interface {
	CreateEntitlement(c *gin.Context)
	GetEntitlement(c *gin.Context)
	UpdateEntitlement(c *gin.Context)
	DeleteEntitlement(c *gin.Context)
	ListEntitlements(c *gin.Context)
	GetUserEntitlements(c *gin.Context)
}

// CreateEntitlement godoc
// @Summary Create a new entitlement
// @Description Creates a new entitlement with the provided details
// @Tags entitlements
// @Accept json
// @Produce json
// @Param entitlement body dto.CreateEntitlementRequest true "Entitlement data"
// @Success 201 {object} dto.EntitlementResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /entitlements [post]
func (c *EntitlementController) CreateEntitlement(ctx *gin.Context) {
	var req dto.CreateEntitlementRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.entitlementService.CreateEntitlement(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// GetEntitlement godoc
// @Summary Get an entitlement by ID
// @Description Retrieves an entitlement by its ID
// @Tags entitlements
// @Accept json
// @Produce json
// @Param id path int true "Entitlement ID"
// @Success 200 {object} dto.EntitlementResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /entitlements/{id} [get]
func (c *EntitlementController) GetEntitlement(ctx *gin.Context) {
	id, err := base.GetUintIDFromPath(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid entitlement id"})
		return
	}

	resp, err := c.entitlementService.GetEntitlement(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateEntitlement godoc
// @Summary Update an entitlement
// @Description Updates an entitlement with the provided details
// @Tags entitlements
// @Accept json
// @Produce json
// @Param id path int true "Entitlement ID"
// @Param entitlement body dto.UpdateEntitlementRequest true "Entitlement data"
// @Success 200 {object} dto.EntitlementResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /entitlements/{id} [put]
func (c *EntitlementController) UpdateEntitlement(ctx *gin.Context) {
	id, err := base.GetUintIDFromPath(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid entitlement id"})
		return
	}

	var req dto.UpdateEntitlementRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.entitlementService.UpdateEntitlement(ctx.Request.Context(), id, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// DeleteEntitlement godoc
// @Summary Delete an entitlement
// @Description Deletes an entitlement by its ID
// @Tags entitlements
// @Accept json
// @Produce json
// @Param id path int true "Entitlement ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /entitlements/{id} [delete]
func (c *EntitlementController) DeleteEntitlement(ctx *gin.Context) {
	id, err := base.GetUintIDFromPath(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid entitlement id"})
		return
	}

	if err := c.entitlementService.DeleteEntitlement(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "entitlement deleted"})
}

// ListEntitlements godoc
// @Summary List all entitlements
// @Description Retrieves a paginated list of entitlements
// @Tags entitlements
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} dto.EntitlementListResponse
// @Failure 500 {object} map[string]string
// @Router /entitlements [get]
func (c *EntitlementController) ListEntitlements(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	resp, err := c.entitlementService.ListEntitlements(ctx.Request.Context(), page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetUserEntitlements godoc
// @Summary Get entitlements by user ID
// @Description Retrieves all entitlements for a specific user
// @Tags entitlements
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Success 200 {object} dto.EntitlementListResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /entitlements/user/{userId} [get]
func (c *EntitlementController) GetUserEntitlements(ctx *gin.Context) {
	userID, err := base.GetUintIDFromPath(ctx, "userId")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	resp, err := c.entitlementService.GetUserEntitlements(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}