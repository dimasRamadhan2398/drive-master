package controllers

import (
	"user-service/models"
	"user-service/models/dto"
	responseRes "user-service/pkg/response"
	"user-service/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// TestimonialController handles testimonial-related HTTP requests
type TestimonialController struct {
	testimonialService services.ITestimonialService
}

// ITestimonialController defines the interface for testimonial controller
type ITestimonialController interface {
	GetAllTestimonials(ctx *gin.Context)
	GetTestimonialByID(ctx *gin.Context)
	GetPublishedTestimonials(ctx *gin.Context)
	GetFeaturedTestimonials(ctx *gin.Context)
	GetTestimonialsByUserID(ctx *gin.Context)
	CreateTestimonial(ctx *gin.Context)
	UpdateTestimonial(ctx *gin.Context)
	DeleteTestimonial(ctx *gin.Context)
	RemoveFromFeatured(ctx *gin.Context)
}

// NewTestimonialController creates a new testimonial controller
func NewTestimonialController(testimonialService services.ITestimonialService) ITestimonialController {
	return &TestimonialController{
		testimonialService: testimonialService,
	}
}

// GetAllTestimonials handles GET /api/v1/testimonials
// @Summary Get all testimonials
// @Description Retrieves all testimonials (admin)
// @Tags Testimonials
// @Produce json
// @Success 200 {object} response.Response
// @Router /testimonials [get]
func (c *TestimonialController) GetAllTestimonials(ctx *gin.Context) {
	testimonials, err := c.testimonialService.GetAllTestimonials(ctx.Request.Context())
	if err != nil {
		responseRes.InternalServerError(ctx, "Failed to fetch testimonials")
		return
	}

	responseRes.OK(ctx, "Testimonials fetched successfully", testimonials)
}

// GetTestimonialByID handles GET /api/v1/testimonials/:id
// @Summary Get testimonial by ID
// @Description Retrieves a specific testimonial by ID
// @Tags Testimonials
// @Produce json
// @Param id path string true "Testimonial ID"
// @Success 200 {object} response.Response
// @Router /testimonials/{id} [get]
func (c *TestimonialController) GetTestimonialByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		responseRes.BadRequest(ctx, "Invalid testimonial ID format")
		return
	}

	testimonial, err := c.testimonialService.GetTestimonialByID(ctx.Request.Context(), id)
	if err != nil {
		responseRes.InternalServerError(ctx, "Failed to fetch testimonial")
		return
	}

	if testimonial == nil {
		responseRes.NotFound(ctx, "Testimonial not found")
		return
	}

	responseRes.OK(ctx, "Testimonial fetched successfully", testimonial)
}

// GetPublishedTestimonials handles GET /api/v1/testimonials/published
// @Summary Get published testimonials
// @Description Retrieves all published testimonials
// @Tags Testimonials
// @Produce json
// @Success 200 {object} response.Response
// @Router /testimonials/published [get]
func (c *TestimonialController) GetPublishedTestimonials(ctx *gin.Context) {
	testimonials, err := c.testimonialService.GetPublishedTestimonials(ctx.Request.Context())
	if err != nil {
		responseRes.InternalServerError(ctx, "Failed to fetch published testimonials")
		return
	}

	responseRes.OK(ctx, "Published testimonials fetched successfully", testimonials)
}

// GetFeaturedTestimonials handles GET /api/v1/testimonials/featured
// @Summary Get featured testimonials
// @Description Retrieves featured testimonials
// @Tags Testimonials
// @Produce json
// @Success 200 {object} response.Response
// @Router /testimonials/featured [get]
func (c *TestimonialController) GetFeaturedTestimonials(ctx *gin.Context) {
	testimonials, err := c.testimonialService.GetFeaturedTestimonials(ctx.Request.Context())
	if err != nil {
		responseRes.InternalServerError(ctx, "Failed to fetch featured testimonials")
		return
	}

	responseRes.OK(ctx, "Featured testimonials fetched successfully", testimonials)
}

// GetTestimonialsByUserID handles GET /api/v1/testimonials/user/:userId
// @Summary Get testimonials by user ID
// @Description Retrieves testimonials for a specific user
// @Tags Testimonials
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} response.Response
// @Router /testimonials/user/{userId} [get]
func (c *TestimonialController) GetTestimonialsByUserID(ctx *gin.Context) {
	userIDParam := ctx.Param("userId")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		responseRes.BadRequest(ctx, "Invalid user ID format")
		return
	}

	testimonials, err := c.testimonialService.GetTestimonialsByUserID(ctx.Request.Context(), userID)
	if err != nil {
		responseRes.InternalServerError(ctx, "Failed to fetch user testimonials")
		return
	}

	responseRes.OK(ctx, "User testimonials fetched successfully", testimonials)
}

// CreateTestimonial handles POST /api/v1/testimonials
// @Summary Create a new testimonial
// @Description Creates a new testimonial
// @Tags Testimonials
// @Accept json
// @Produce json
// @Param request body dto.CreateTestimonialRequest true "Testimonial data"
// @Success 201 {object} response.Response
// @Router /testimonials [post]
func (c *TestimonialController) CreateTestimonial(ctx *gin.Context) {
	var req dto.CreateTestimonialRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responseRes.BadRequest(ctx, "Invalid request body: "+err.Error())
		return
	}

	if err := validator.New().Struct(req); err != nil {
		responseRes.BadRequest(ctx, "Validation failed: "+err.Error())
		return
	}

	// Set default status if not provided
	status := req.Status
	if status == "" {
		status = models.TestimonialStatusDraft
	}

	testimonial := &models.Testimonial{
		UserID:     req.UserID,
		UserName:   req.UserName,
		UserImage:  req.UserImage,
		UserRole:   req.UserRole,
		Content:    req.Content,
		Rating:     req.Rating,
		Tags:       req.Tags,
		Status:     status,
		IsFeatured: req.IsFeatured,
		AddedBy:    req.AddedBy,
	}

	if err := c.testimonialService.CreateTestimonial(ctx.Request.Context(), testimonial); err != nil {
		responseRes.InternalServerError(ctx, "Failed to create testimonial")
		return
	}

	responseRes.Created(ctx, "Testimonial created successfully", testimonial)
}

// UpdateTestimonial handles PUT /api/v1/testimonials/:id
// @Summary Update a testimonial
// @Description Updates an existing testimonial
// @Tags Testimonials
// @Accept json
// @Produce json
// @Param id path string true "Testimonial ID"
// @Param request body dto.UpdateTestimonialRequest true "Testimonial data"
// @Success 200 {object} response.Response
// @Router /testimonials/{id} [put]
func (c *TestimonialController) UpdateTestimonial(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		responseRes.BadRequest(ctx, "Invalid testimonial ID format")
		return
	}

	// Get existing testimonial
	testimonial, err := c.testimonialService.GetTestimonialByID(ctx.Request.Context(), id)
	if err != nil {
		responseRes.InternalServerError(ctx, "Failed to fetch testimonial")
		return
	}

	if testimonial == nil {
		responseRes.NotFound(ctx, "Testimonial not found")
		return
	}

	var req dto.UpdateTestimonialRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responseRes.BadRequest(ctx, "Invalid request body: "+err.Error())
		return
	}

	// Update fields
	if req.UserName != "" {
		testimonial.UserName = req.UserName
	}
	if req.UserImage != "" {
		testimonial.UserImage = req.UserImage
	}
	if req.UserRole != "" {
		testimonial.UserRole = req.UserRole
	}
	if req.Content != "" {
		testimonial.Content = req.Content
	}
	if req.Rating != 0 {
		testimonial.Rating = req.Rating
	}
	if req.Tags != "" {
		testimonial.Tags = req.Tags
	}
	if req.Status != "" {
		testimonial.Status = req.Status
	}
	testimonial.IsFeatured = req.IsFeatured
	if req.SortOrder != 0 {
		testimonial.SortOrder = req.SortOrder
	}

	if err := c.testimonialService.UpdateTestimonial(ctx.Request.Context(), testimonial); err != nil {
		responseRes.InternalServerError(ctx, "Failed to update testimonial")
		return
	}

	responseRes.OK(ctx, "Testimonial updated successfully", testimonial)
}

// DeleteTestimonial handles DELETE /api/v1/testimonials/:id
// @Summary Delete a testimonial
// @Description Deletes a testimonial by ID
// @Tags Testimonials
// @Produce json
// @Param id path string true "Testimonial ID"
// @Success 200 {object} response.Response
// @Router /testimonials/{id} [delete]
func (c *TestimonialController) DeleteTestimonial(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		responseRes.BadRequest(ctx, "Invalid testimonial ID format")
		return
	}

	testimonial, err := c.testimonialService.GetTestimonialByID(ctx.Request.Context(), id)
	if err != nil {
		responseRes.InternalServerError(ctx, "Failed to fetch testimonial")
		return
	}

	if testimonial == nil {
		responseRes.NotFound(ctx, "Testimonial not found")
		return
	}

	if err := c.testimonialService.DeleteTestimonial(ctx.Request.Context(), id); err != nil {
		responseRes.InternalServerError(ctx, "Failed to delete testimonial")
		return
	}

	responseRes.OK(ctx, "Testimonial deleted successfully", nil)
}

// RemoveFromFeatured handles PUT /api/v1/testimonials/:id/unfeature
// @Summary Remove testimonial from featured
// @Description Removes the featured status from a testimonial
// @Tags Testimonials
// @Produce json
// @Param id path string true "Testimonial ID"
// @Success 200 {object} response.Response
// @Router /testimonials/{id}/unfeature [put]
func (c *TestimonialController) RemoveFromFeatured(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		responseRes.BadRequest(ctx, "Invalid testimonial ID format")
		return
	}

	testimonial, err := c.testimonialService.GetTestimonialByID(ctx.Request.Context(), id)
	if err != nil {
		responseRes.InternalServerError(ctx, "Failed to fetch testimonial")
		return
	}

	if testimonial == nil {
		responseRes.NotFound(ctx, "Testimonial not found")
		return
	}

	if err := c.testimonialService.RemoveFromFeatured(ctx.Request.Context(), id); err != nil {
		responseRes.InternalServerError(ctx, "Failed to remove from featured")
		return
	}

	responseRes.OK(ctx, "Testimonial removed from featured", nil)
}