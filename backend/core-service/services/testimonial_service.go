package services

import (
	"context"
	"core-service/models"
	"core-service/repositories"

	"github.com/google/uuid"
)

type ITestimonialService interface {
	CreateTestimonial(ctx context.Context, testimonial *models.Testimonial) error
	GetTestimonialByID(ctx context.Context, id uuid.UUID) (*models.Testimonial, error)
	GetAllTestimonials(ctx context.Context) ([]models.Testimonial, error)
	GetPublishedTestimonials(ctx context.Context) ([]models.Testimonial, error)
	GetFeaturedTestimonials(ctx context.Context) ([]models.Testimonial, error)
	GetTestimonialsByUserID(ctx context.Context, userID uuid.UUID) ([]models.Testimonial, error)
	UpdateTestimonial(ctx context.Context, testimonial *models.Testimonial) error
	DeleteTestimonial(ctx context.Context, id uuid.UUID) error
	CountTestimonials(ctx context.Context) (int64, error)
}

type TestimonialService struct {
	testimonialRepo repositories.ITestimonialRepository
}

func NewTestimonialService(testimonialRepo repositories.ITestimonialRepository) ITestimonialService {
	return &TestimonialService{
		testimonialRepo: testimonialRepo,
	}
}

// CreateTestimonial creates a new testimonial
func (s *TestimonialService) CreateTestimonial(ctx context.Context, testimonial *models.Testimonial) error {
	return s.testimonialRepo.CreateTestimonial(ctx, testimonial)
}

// GetTestimonialByID retrieves a testimonial by ID
func (s *TestimonialService) GetTestimonialByID(ctx context.Context, id uuid.UUID) (*models.Testimonial, error) {
	return s.testimonialRepo.GetTestimonialByID(ctx, id)
}

// GetAllTestimonials retrieves all testimonials
func (s *TestimonialService) GetAllTestimonials(ctx context.Context) ([]models.Testimonial, error) {
	return s.testimonialRepo.GetAllTestimonials(ctx)
}

// GetPublishedTestimonials retrieves all published testimonials
func (s *TestimonialService) GetPublishedTestimonials(ctx context.Context) ([]models.Testimonial, error) {
	return s.testimonialRepo.GetPublishedTestimonials(ctx)
}

// GetFeaturedTestimonials retrieves featured testimonials
func (s *TestimonialService) GetFeaturedTestimonials(ctx context.Context) ([]models.Testimonial, error) {
	return s.testimonialRepo.GetFeaturedTestimonials(ctx)
}

// GetTestimonialsByUserID retrieves testimonials by user ID
func (s *TestimonialService) GetTestimonialsByUserID(ctx context.Context, userID uuid.UUID) ([]models.Testimonial, error) {
	return s.testimonialRepo.GetTestimonialsByUserID(ctx, userID)
}

// UpdateTestimonial updates a testimonial
func (s *TestimonialService) UpdateTestimonial(ctx context.Context, testimonial *models.Testimonial) error {
	return s.testimonialRepo.UpdateTestimonial(ctx, testimonial)
}

// DeleteTestimonial deletes a testimonial
func (s *TestimonialService) DeleteTestimonial(ctx context.Context, id uuid.UUID) error {
	return s.testimonialRepo.DeleteTestimonial(ctx, id)
}

// CountTestimonials returns the total number of testimonials
func (s *TestimonialService) CountTestimonials(ctx context.Context) (int64, error) {
	return s.testimonialRepo.CountTestimonials(ctx)
}