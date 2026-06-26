package controllers

import (
	"core-service/models"
	"core-service/models/dto"
	"core-service/pkg/response"
	"core-service/services"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// AddOnController handles add-on related HTTP requests
type AddOnController struct {
	addOnService services.IAddOnService
}

// IAddOnController defines the interface for add-on controller
type IAddOnController interface {
	GetAllAddOns(ctx *gin.Context)
	GetAddOnByID(ctx *gin.Context)
	CreateAddOn(ctx *gin.Context)
	UpdateAddOn(ctx *gin.Context)
	DeleteAddOn(ctx *gin.Context)
	ToggleStatusAddOn(ctx *gin.Context)
}

// NewAddOnController creates a new add-on controller
func NewAddOnController(addOnService services.IAddOnService) IAddOnController {
	return &AddOnController{
		addOnService: addOnService,
	}
}

// GetAllAddOns handles GET /api/v1/addons
// @Summary Get all add-ons
// @Description Retrieves all add-ons with pagination
// @Tags AddOns
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.Response
// @Router /addons [get]
func (c *AddOnController) GetAllAddOns(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	addons, total, err := c.addOnService.GetAllAddOnsPaginated(ctx.Request.Context(), page, limit)
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch add-ons")
		return
	}

	// Build response items
	var items []dto.AddOnResponse
	for _, addon := range addons {
		items = append(items, dto.AddOnResponse{
			ID:          addon.ID,
			Title:       addon.Title,
			Description: addon.Description,
			Price:       addon.Price,
			Sessions:    addon.Sessions,
			Status:      addon.Status,
			ImageURL:    addon.ImageURL,
			SortOrder:   addon.SortOrder,
			CreatedAt:   addon.CreatedAt,
			UpdatedAt:   addon.UpdatedAt,
		})
	}

	pagination := dto.NewPaginationMeta(total, page, limit)
	response.Paginated(ctx, http.StatusOK, "Add-ons fetched successfully", items, &pagination)
}

// GetAddOnByID handles GET /api/v1/addons/:id
// @Summary Get add-on by ID
// @Description Retrieves a specific add-on by ID
// @Tags AddOns
// @Produce json
// @Param id path string true "AddOn ID"
// @Success 200 {object} response.Response
// @Router /addons/{id} [get]
func (c *AddOnController) GetAddOnByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid add-on ID format")
		return
	}

	addon, err := c.addOnService.GetAddOnByID(ctx.Request.Context(), id)
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch add-on")
		return
	}

	if addon == nil {
		response.NotFound(ctx, "Add-on not found")
		return
	}

	resp := dto.AddOnResponse{
		ID:          addon.ID,
		Title:       addon.Title,
		Description: addon.Description,
		Price:       addon.Price,
		Sessions:    addon.Sessions,
		Status:      addon.Status,
		ImageURL:    addon.ImageURL,
		SortOrder:   addon.SortOrder,
		CreatedAt:   addon.CreatedAt,
		UpdatedAt:   addon.UpdatedAt,
	}

	response.OK(ctx, "Add-on fetched successfully", resp)
}

// CreateAddOn handles POST /api/v1/addons
// @Summary Create a new add-on
// @Description Creates a new add-on
// @Tags AddOns
// @Accept json
// @Produce json
// @Param request body dto.CreateAddOnRequest true "AddOn data"
// @Success 201 {object} response.Response
// @Router /addons [post]
func (c *AddOnController) CreateAddOn(ctx *gin.Context) {
	var req dto.CreateAddOnRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body: "+err.Error())
		return
	}

	if err := validator.New().Struct(req); err != nil {
		response.BadRequest(ctx, "Validation failed: "+err.Error())
		return
	}

	req.SetDefaults()

	addon := &models.AddOn{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Sessions:    req.Sessions,
		ImageURL:    req.ImageURL,
		SortOrder:   req.SortOrder,
		Status:      models.AddOnStatusActive,
	}

	if err := c.addOnService.CreateAddOn(ctx.Request.Context(), addon); err != nil {
		response.InternalServerError(ctx, "Failed to create add-on")
		return
	}

	resp := dto.AddOnResponse{
		ID:          addon.ID,
		Title:       addon.Title,
		Description: addon.Description,
		Price:       addon.Price,
		Sessions:    addon.Sessions,
		Status:      addon.Status,
		ImageURL:    addon.ImageURL,
		SortOrder:   addon.SortOrder,
		CreatedAt:   addon.CreatedAt,
		UpdatedAt:   addon.UpdatedAt,
	}

	response.Created(ctx, "Add-on created successfully", resp)
}

// UpdateAddOn handles PUT /api/v1/addons/:id
// @Summary Update an add-on
// @Description Updates an existing add-on
// @Tags AddOns
// @Accept json
// @Produce json
// @Param id path string true "AddOn ID"
// @Param request body dto.UpdateAddOnRequest true "AddOn data"
// @Success 200 {object} response.Response
// @Router /addons/{id} [put]
func (c *AddOnController) UpdateAddOn(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid add-on ID format")
		return
	}

	// Get existing add-on
	addon, err := c.addOnService.GetAddOnByID(ctx.Request.Context(), id)
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch add-on")
		return
	}

	if addon == nil {
		response.NotFound(ctx, "Add-on not found")
		return
	}

	var req dto.UpdateAddOnRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body: "+err.Error())
		return
	}

	// Update fields
	if req.Title != "" {
		addon.Title = req.Title
	}
	if req.Description != "" {
		addon.Description = req.Description
	}
	if req.Price != 0 {
		addon.Price = req.Price
	}
	if req.Sessions != 0 {
		addon.Sessions = req.Sessions
	}
	if req.Status != "" {
		addon.Status = req.Status
	}
	if req.ImageURL != "" {
		addon.ImageURL = req.ImageURL
	}
	if req.SortOrder != 0 {
		addon.SortOrder = req.SortOrder
	}

	if err := c.addOnService.UpdateAddOn(ctx.Request.Context(), addon); err != nil {
		response.InternalServerError(ctx, "Failed to update add-on")
		return
	}

	resp := dto.AddOnResponse{
		ID:          addon.ID,
		Title:       addon.Title,
		Description: addon.Description,
		Price:       addon.Price,
		Sessions:    addon.Sessions,
		Status:      addon.Status,
		ImageURL:    addon.ImageURL,
		SortOrder:   addon.SortOrder,
		CreatedAt:   addon.CreatedAt,
		UpdatedAt:   addon.UpdatedAt,
	}

	response.OK(ctx, "Add-on updated successfully", resp)
}

// DeleteAddOn handles DELETE /api/v1/addons/:id
// @Summary Delete an add-on
// @Description Deletes an add-on by ID
// @Tags AddOns
// @Produce json
// @Param id path string true "AddOn ID"
// @Success 200 {object} response.Response
// @Router /addons/{id} [delete]
func (c *AddOnController) DeleteAddOn(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid add-on ID format")
		return
	}

	addon, err := c.addOnService.GetAddOnByID(ctx.Request.Context(), id)
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch add-on")
		return
	}

	if addon == nil {
		response.NotFound(ctx, "Add-on not found")
		return
	}

	if err := c.addOnService.DeleteAddOn(ctx.Request.Context(), addon); err != nil {
		response.InternalServerError(ctx, "Failed to delete add-on")
		return
	}

	response.OK(ctx, "Add-on deleted successfully", nil)
}

// ToggleStatusAddOn handles PUT /api/v1/addons/toggle-status/:id
// @Summary Toggle add-on status
// @Description Toggles an add-on status between active and inactive
// @Tags AddOns
// @Produce json
// @Param id path string true "AddOn ID"
// @Success 200 {object} response.Response
// @Router /addons/toggle-status/{id} [put]
func (c *AddOnController) ToggleStatusAddOn(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid add-on ID format")
		return
	}

	addon, err := c.addOnService.ToggleStatusAddOn(ctx.Request.Context(), id)
	if err != nil {
		response.InternalServerError(ctx, "Failed to toggle add-on status")
		return
	}

	resp := dto.AddOnResponse{
		ID:          addon.ID,
		Title:       addon.Title,
		Description: addon.Description,
		Price:       addon.Price,
		Sessions:    addon.Sessions,
		Status:      addon.Status,
		ImageURL:    addon.ImageURL,
		SortOrder:   addon.SortOrder,
		CreatedAt:   addon.CreatedAt,
		UpdatedAt:   addon.UpdatedAt,
	}

	response.OK(ctx, "Add-on status toggled successfully", resp)
}
