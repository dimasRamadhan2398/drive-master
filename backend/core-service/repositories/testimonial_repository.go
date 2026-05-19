package repositories

import (
	"context"
	"core-service/models"
	"core-service/pkg/base"

	"github.com/google/uuid"
)

// ITestimonialRepository defines testimonial operations
type ITestimonialRepository interface {
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

type TestimonialRepository struct {
	*base.BaseRepository
}

func NewTestimonialRepository(baseRepo *base.BaseRepository) ITestimonialRepository {
	return &TestimonialRepository{BaseRepository: baseRepo}
}

// CreateTestimonial creates a new testimonial
func (r *TestimonialRepository) CreateTestimonial(ctx context.Context, testimonial *models.Testimonial) error {
	return r.BaseRepository.Create(ctx, testimonial)
}

// GetTestimonialByID retrieves a testimonial by ID
func (r *TestimonialRepository) GetTestimonialByID(ctx context.Context, id uuid.UUID) (*models.Testimonial, error) {
	var testimonial models.Testimonial
	if err := r.BaseRepository.FindByID(ctx, &testimonial, id); err != nil {
		return nil, err
	}
	return &testimonial, nil
}

// GetAllTestimonials retrieves all testimonials
func (r *TestimonialRepository) GetAllTestimonials(ctx context.Context) ([]models.Testimonial, error) {
	var testimonials []models.Testimonial
	opts := base.NewQueryOptions()
	opts.Limit = 0 // No limit
	opts.Order = "created_at DESC"
	if err := r.BaseRepository.FindMany(ctx, &models.Testimonial{}, &testimonials, opts); err != nil {
		return nil, err
	}
	return testimonials, nil
}

// GetPublishedTestimonials retrieves all published testimonials
func (r *TestimonialRepository) GetPublishedTestimonials(ctx context.Context) ([]models.Testimonial, error) {
	var testimonials []models.Testimonial
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"status": models.TestimonialStatusPublished}).
		WithOrder("sort_order ASC, created_at DESC")
	if err := r.BaseRepository.FindMany(ctx, &models.Testimonial{}, &testimonials, opts); err != nil {
		return nil, err
	}
	return testimonials, nil
}

// GetFeaturedTestimonials retrieves featured published testimonials
func (r *TestimonialRepository) GetFeaturedTestimonials(ctx context.Context) ([]models.Testimonial, error) {
	var testimonials []models.Testimonial
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{
			"status":      models.TestimonialStatusPublished,
			"is_featured": true,
		}).
		WithOrder("sort_order ASC, created_at DESC")
	if err := r.BaseRepository.FindMany(ctx, &models.Testimonial{}, &testimonials, opts); err != nil {
		return nil, err
	}
	return testimonials, nil
}

// GetTestimonialsByUserID retrieves testimonials by user ID
func (r *TestimonialRepository) GetTestimonialsByUserID(ctx context.Context, userID uuid.UUID) ([]models.Testimonial, error) {
	var testimonials []models.Testimonial
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"user_id": userID}).
		WithOrder("created_at DESC")
	if err := r.BaseRepository.FindMany(ctx, &models.Testimonial{}, &testimonials, opts); err != nil {
		return nil, err
	}
	return testimonials, nil
}

// UpdateTestimonial updates a testimonial
func (r *TestimonialRepository) UpdateTestimonial(ctx context.Context, testimonial *models.Testimonial) error {
	return r.BaseRepository.Update(ctx, testimonial)
}

// DeleteTestimonial deletes a testimonial
func (r *TestimonialRepository) DeleteTestimonial(ctx context.Context, id uuid.UUID) error {
	var testimonial models.Testimonial
	if err := r.BaseRepository.FindByID(ctx, &testimonial, id); err != nil {
		return err
	}
	return r.BaseRepository.Delete(ctx, &testimonial)
}

// CountTestimonials returns the total number of testimonials
func (r *TestimonialRepository) CountTestimonials(ctx context.Context) (int64, error) {
	return r.BaseRepository.Count(ctx, &models.Testimonial{}, nil)
}