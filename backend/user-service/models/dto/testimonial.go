package dto

import (
	"time"

	"user-service/models"

	"github.com/google/uuid"
)

// ========== TESTIMONIAL DTOs ==========

type CreateTestimonialRequest struct {
	UserID     uuid.UUID                `json:"userId" binding:"required"`
	UserName   string                   `json:"userName" binding:"required"`
	UserImage  string                   `json:"userImage"`
	UserRole   string                   `json:"userRole"`
	Content    string                   `json:"content" binding:"required"`
	Rating     float64                  `json:"rating" binding:"required,min=1,max=5"`
	Tags       string                   `json:"tags"`
	Status     models.TestimonialStatus `json:"status"`
	IsFeatured bool                     `json:"isFeatured"`
	AddedBy    uuid.UUID                `json:"addedBy" binding:"required"`
}

type UpdateTestimonialRequest struct {
	UserName   string                   `json:"userName"`
	UserImage  string                   `json:"userImage"`
	UserRole   string                   `json:"userRole"`
	Content    string                   `json:"content"`
	Rating     float64                  `json:"rating"`
	Tags       string                   `json:"tags"`
	Status     models.TestimonialStatus `json:"status"`
	IsFeatured bool                     `json:"isFeatured"`
	SortOrder  int                      `json:"sortOrder"`
}

type TestimonialResponse struct {
	ID         uuid.UUID                `json:"id"`
	UserID     uuid.UUID                `json:"userId"`
	UserName   string                   `json:"userName"`
	UserImage  string                   `json:"userImage"`
	UserRole   string                   `json:"userRole"`
	Content    string                   `json:"content"`
	Rating     float64                  `json:"rating"`
	Tags       string                   `json:"tags"`
	Status     models.TestimonialStatus `json:"status"`
	IsFeatured bool                     `json:"isFeatured"`
	AddedBy    uuid.UUID                `json:"addedBy"`
	AddedAt    time.Time                `json:"addedAt"`
	SortOrder  int                      `json:"sortOrder"`
	CreatedAt  time.Time                `json:"createdAt"`
	UpdatedAt  time.Time                `json:"updatedAt"`
}

type TestimonialListResponse struct {
	Testimonials []TestimonialResponse `json:"testimonials"`
	Total        int64                 `json:"total"`
}

type ToggleFeaturedRequest struct {
	IsFeatured *bool `json:"isFeatured"`
}
