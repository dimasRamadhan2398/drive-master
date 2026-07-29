package controllers

import (
	"net/http"
	"strconv"

	"user-service/models/dto"
	apperrors "user-service/pkg/errors"
	responseRes "user-service/pkg/response"
	"user-service/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EntitlementController struct {
	entitlementService services.IEntitlementService
}

type IEntitlementController interface {
	CreateEntitlement(ctx *gin.Context)
	GetEntitlement(ctx *gin.Context)
	GetEntitlementByID(ctx *gin.Context)
	UpdateEntitlement(ctx *gin.Context)
	DeleteEntitlement(ctx *gin.Context)
	ListEntitlements(ctx *gin.Context)
	UseSession(ctx *gin.Context)
	SyncEntitlement(ctx *gin.Context)
}

func NewEntitlementController(entitlementService services.IEntitlementService) IEntitlementController {
	return &EntitlementController{
		entitlementService: entitlementService,
	}
}

// @Summary Create Entitlement
// @Description Add a new entitlement for a member
// @Tags Members
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param request body dto.CreateEntitlementInput true "Entitlement data"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /members/{id}/entitlements [post]
func (c *EntitlementController) CreateEntitlement(ctx *gin.Context) {
	memberID, err := parseUUID(ctx, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	var input dto.CreateEntitlementInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	resp, err := c.entitlementService.CreateEntitlement(ctx.Request.Context(), memberID, input)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusCreated, "Entitlement created successfully", resp)
}

// @Summary Get Entitlement
// @Description Get a specific entitlement for a member
// @Tags Members
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param entId path string true "Entitlement ID (UUID)"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /members/{id}/entitlements/{entId} [get]
func (c *EntitlementController) GetEntitlement(ctx *gin.Context) {
	memberID, err := parseUUID(ctx, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	entitlementID, err := parseUUID(ctx, "entId")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	resp, err := c.entitlementService.GetEntitlement(ctx.Request.Context(), memberID, entitlementID)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Entitlement retrieved successfully", resp)
}

// @Summary Get Entitlement By ID
// @Description Get a specific entitlement by its ID (internal use)
// @Tags Members
// @Produce json
// @Param entId path string true "Entitlement ID (UUID)"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /entitlements/{entId} [get]
func (c *EntitlementController) GetEntitlementByID(ctx *gin.Context) {
	entitlementID, err := parseUUID(ctx, "entId")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	resp, err := c.entitlementService.GetEntitlementByID(ctx.Request.Context(), entitlementID)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Entitlement retrieved successfully", resp)
}

// @Summary Update Entitlement
// @Description Update an entitlement for a member
// @Tags Members
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param entId path string true "Entitlement ID (UUID)"
// @Param request body dto.UpdateEntitlementInput true "Entitlement data"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /members/{id}/entitlements/{entId} [put]
func (c *EntitlementController) UpdateEntitlement(ctx *gin.Context) {
	memberID, err := parseUUID(ctx, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	entitlementID, err := parseUUID(ctx, "entId")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	var input dto.UpdateEntitlementInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	resp, err := c.entitlementService.UpdateEntitlement(ctx.Request.Context(), memberID, entitlementID, input)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Entitlement updated successfully", resp)
}

// @Summary Delete Entitlement
// @Description Delete an entitlement for a member
// @Tags Members
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param entId path string true "Entitlement ID (UUID)"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /members/{id}/entitlements/{entId} [delete]
func (c *EntitlementController) DeleteEntitlement(ctx *gin.Context) {
	memberID, err := parseUUID(ctx, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	entitlementID, err := parseUUID(ctx, "entId")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	if err := c.entitlementService.DeleteEntitlement(ctx.Request.Context(), memberID, entitlementID); err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Entitlement deleted successfully", nil)
}

// @Summary List Entitlements
// @Description Get all entitlements for a member
// @Tags Members
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.Response
// @Router /members/{id}/entitlements [get]
func (c *EntitlementController) ListEntitlements(ctx *gin.Context) {
	memberID, err := parseUUID(ctx, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	resp, err := c.entitlementService.ListEntitlements(ctx.Request.Context(), memberID, page, limit)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Entitlements retrieved successfully", resp)
}

// @Summary Use Session
// @Description Use one session from an entitlement
// @Tags Members
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param entId path string true "Entitlement ID (UUID)"
// @Param request body dto.UseSessionInput true "Session usage data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /members/{id}/entitlements/{entId}/use-session [post]
func (c *EntitlementController) UseSession(ctx *gin.Context) {
	memberID, err := parseUUID(ctx, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	entitlementID, err := parseUUID(ctx, "entId")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	var input dto.UseSessionInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	resp, err := c.entitlementService.UseSession(ctx.Request.Context(), memberID, entitlementID, input)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Session used successfully", resp)
}

func (c *EntitlementController) SyncEntitlement(ctx *gin.Context) {
	var input struct {
		MemberID         string `json:"member_id"`
		BookingID         string `json:"booking_id"`
		PackageID         string `json:"package_id"`
		PackageName       string `json:"package_name"`
		TotalSessions     int    `json:"total_sessions"`
		IsNightSession   bool   `json:"is_night_session"`
		IsWeekendSession bool   `json:"is_weekend_session"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	memberID, err := uuid.Parse(input.MemberID)
	if err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	bookingID, err := uuid.Parse(input.BookingID)
	if err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	packageID, err := uuid.Parse(input.PackageID)
	if err != nil {
		packageID = uuid.Nil
	}

	totalSessions := input.TotalSessions
	if totalSessions <= 0 {
		totalSessions = 10
	}

	resp, err := c.entitlementService.SyncEntitlementFromBooking(
		ctx.Request.Context(),
		memberID,
		bookingID,
		packageID,
		input.PackageName,
		totalSessions,
		input.IsNightSession,
		input.IsWeekendSession,
	)

	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Entitlement synced successfully", resp)
}

// parseUUID parses a UUID from a path parameter
func parseUUID(ctx *gin.Context, param string) (uuid.UUID, error) {
	idStr := ctx.Param(param)
	return uuid.Parse(idStr)
}
