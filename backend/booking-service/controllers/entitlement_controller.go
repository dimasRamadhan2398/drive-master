package controllers

import (
	"net/http"
	"strconv"

	"booking-service/models/dto"
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

func (c *EntitlementController) GetEntitlement(ctx *gin.Context) {
	id, err := getUintIDFromPath(ctx, "id")
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

func (c *EntitlementController) UpdateEntitlement(ctx *gin.Context) {
	id, err := getUintIDFromPath(ctx, "id")
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

func (c *EntitlementController) DeleteEntitlement(ctx *gin.Context) {
	id, err := getUintIDFromPath(ctx, "id")
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

func (c *EntitlementController) GetUserEntitlements(ctx *gin.Context) {
	userID, err := getUintIDFromPath(ctx, "userId")
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