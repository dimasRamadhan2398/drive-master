package controllers

import (
	"core-service/models/dto"
	"core-service/pkg/response"
	"core-service/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/go-playground/validator/v10"
)

// FAQController handles FAQ-related HTTP requests
type FAQController struct {
	faqService services.IFAQService
}

// IF AQController defines the interface for FAQ controller
type IFAQController interface {
	GetAllFAQs(ctx *gin.Context)
	GetActiveFAQs(ctx *gin.Context)
	GetFAQByID(ctx *gin.Context)
	CreateFAQ(ctx *gin.Context)
	UpdateFAQ(ctx *gin.Context)
	ReorderFAQ(ctx *gin.Context)
	DeleteFAQ(ctx *gin.Context)
}

// NewFAQController creates a new FAQ controller
func NewFAQController(faqService services.IFAQService) IFAQController {
	return &FAQController{
		faqService: faqService,
	}
}

// GetAllFAQs handles GET /api/v1/faqs
// @Summary Get all FAQs
// @Description Retrieves all FAQs (including inactive)
// @Tags FAQs
// @Produce json
// @Success 200 {object} response.Response
// @Router /faqs [get]
func (c *FAQController) GetAllFAQs(ctx *gin.Context) {
	faqs, err := c.faqService.GetAllFAQs(ctx.Request.Context())
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch FAQs")
		return
	}

	response.OK(ctx, "FAQs fetched successfully", faqs)
}

// GetActiveFAQs handles GET /api/v1/faqs/active
// @Summary Get active FAQs
// @Description Retrieves all active FAQs
// @Tags FAQs
// @Produce json
// @Success 200 {object} response.Response
// @Router /faqs/active [get]
func (c *FAQController) GetActiveFAQs(ctx *gin.Context) {
	faqs, err := c.faqService.GetActiveFAQs(ctx.Request.Context())
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch FAQs")
		return
	}

	response.OK(ctx, "Active FAQs fetched successfully", faqs)
}

// GetFAQByID handles GET /api/v1/faqs/:id
// @Summary Get FAQ by ID
// @Description Retrieves a specific FAQ by ID
// @Tags FAQs
// @Produce json
// @Param id path string true "FAQ ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /faqs/{id} [get]
func (c *FAQController) GetFAQByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid FAQ ID format")
		return
	}

	faq, err := c.faqService.GetFAQByID(ctx.Request.Context(), id)
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch FAQ")
		return
	}

	if faq == nil {
		response.NotFound(ctx, "FAQ not found")
		return
	}

	response.OK(ctx, "FAQ fetched successfully", faq)
}

// CreateFAQ handles POST /api/v1/faqs
// @Summary Create a new FAQ
// @Description Creates a new FAQ
// @Tags FAQs
// @Accept json
// @Produce json
// @Param request body dto.CreateFAQRequest true "FAQ data"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /faqs [post]
func (c *FAQController) CreateFAQ(ctx *gin.Context) {
	var req dto.CreateFAQRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body: "+err.Error())
		return
	}

	if err := validator.New().Struct(req); err != nil {
		response.BadRequest(ctx, "Validation failed: "+err.Error())
		return
	}

	// Set default category if empty
	category := req.Category
	if category == "" {
		category = "general"
	}

	faq, err := c.faqService.CreateFAQ(ctx.Request.Context(), req.Question, req.Answer, category, req.Order)
	if err != nil {
		response.InternalServerError(ctx, "Failed to create FAQ")
		return
	}

	response.Created(ctx, "FAQ created successfully", faq)
}

// UpdateFAQ handles PUT /api/v1/faqs/:id
// @Summary Update an FAQ
// @Description Updates an existing FAQ
// @Tags FAQs
// @Accept json
// @Produce json
// @Param id path string true "FAQ ID"
// @Param request body dto.UpdateFAQRequest true "FAQ data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /faqs/{id} [put]
func (c *FAQController) UpdateFAQ(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid FAQ ID format")
		return
	}

	// Check if FAQ exists
	existingFAQ, err := c.faqService.GetFAQByID(ctx.Request.Context(), id)
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch FAQ")
		return
	}
	if existingFAQ == nil {
		response.NotFound(ctx, "FAQ not found")
		return
	}

	var req dto.UpdateFAQRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body: "+err.Error())
		return
	}

	// Use existing values as defaults
	question := existingFAQ.Question
	answer := existingFAQ.Answer
	category := existingFAQ.Category
	order := existingFAQ.Order
	isActive := existingFAQ.IsActive

	// Override with new values if provided
	if req.Question != "" {
		question = req.Question
	}
	if req.Answer != "" {
		answer = req.Answer
	}
	if req.Category != "" {
		category = req.Category
	}
	if req.Order != 0 {
		order = req.Order
	}
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	faq, err := c.faqService.UpdateFAQ(ctx.Request.Context(), id, question, answer, category, order, isActive)
	if err != nil {
		response.InternalServerError(ctx, "Failed to update FAQ")
		return
	}

	response.OK(ctx, "FAQ updated successfully", faq)
}

// ReorderFAQ handles PUT /api/v1/faqs/:id/reorder
// @Summary Reorder an FAQ
// @Description Changes the order of an FAQ and adjusts other FAQs accordingly
// @Tags FAQs
// @Accept json
// @Produce json
// @Param id path string true "FAQ ID"
// @Param request body dto.ReorderFAQRequest true "New order"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /faqs/{id}/reorder [put]
func (c *FAQController) ReorderFAQ(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid FAQ ID format")
		return
	}

	var req dto.ReorderFAQRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body: "+err.Error())
		return
	}

	if err := c.faqService.ReorderFAQ(ctx.Request.Context(), id, req.NewOrder); err != nil {
		response.InternalServerError(ctx, "Failed to reorder FAQ")
		return
	}

	// Fetch the updated FAQ to return
	faq, err := c.faqService.GetFAQByID(ctx.Request.Context(), id)
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch FAQ after reorder")
		return
	}

	response.OK(ctx, "FAQ reordered successfully", faq)
}

// DeleteFAQ handles DELETE /api/v1/faqs/:id
// @Summary Delete an FAQ
// @Description Soft-deletes an FAQ by ID
// @Tags FAQs
// @Produce json
// @Param id path string true "FAQ ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /faqs/{id} [delete]
func (c *FAQController) DeleteFAQ(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid FAQ ID format")
		return
	}

	// Check if FAQ exists
	existingFAQ, err := c.faqService.GetFAQByID(ctx.Request.Context(), id)
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch FAQ")
		return
	}
	if existingFAQ == nil {
		response.NotFound(ctx, "FAQ not found")
		return
	}

	if err := c.faqService.DeleteFAQ(ctx.Request.Context(), id); err != nil {
		response.InternalServerError(ctx, "Failed to delete FAQ")
		return
	}

	response.OK(ctx, "FAQ deleted successfully", nil)
}
