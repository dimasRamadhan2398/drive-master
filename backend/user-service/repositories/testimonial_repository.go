package repositories

import (
	"context"
	"user-service/models"
	"user-service/pkg/base"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	RemoveFromFeatured(ctx context.Context, id uuid.UUID) error
}

type TestimonialRepository struct {
	*base.BaseRepository
}

func NewTestimonialRepository(db *gorm.DB) ITestimonialRepository {
	return &TestimonialRepository{BaseRepository: base.NewBaseRepository(db)}
}

// CreateTestimonial creates a new testimonial
func (r *TestimonialRepository) CreateTestimonial(ctx context.Context, testimonial *models.Testimonial) error {
	return r.BaseRepository.Create(testimonial)
}

// GetTestimonialByID retrieves a testimonial by ID
func (r *TestimonialRepository) GetTestimonialByID(ctx context.Context, id uuid.UUID) (*models.Testimonial, error) {
	var testimonial models.Testimonial
	if err := r.BaseRepository.FindByID(&testimonial, id); err != nil {
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
	if err := r.BaseRepository.FindMany( &models.Testimonial{}, &testimonials, opts); err != nil {
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
	if err := r.BaseRepository.FindMany( &models.Testimonial{}, &testimonials, opts); err != nil {
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
	if err := r.BaseRepository.FindMany( &models.Testimonial{}, &testimonials, opts); err != nil {
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
	if err := r.BaseRepository.FindMany( &models.Testimonial{}, &testimonials, opts); err != nil {
		return nil, err
	}
	return testimonials, nil
}

// UpdateTestimonial updates a testimonial
func (r *TestimonialRepository) UpdateTestimonial(ctx context.Context, testimonial *models.Testimonial) error {
	return r.BaseRepository.Update( testimonial)
}

// DeleteTestimonial deletes a testimonial
func (r *TestimonialRepository) DeleteTestimonial(ctx context.Context, id uuid.UUID) error {
	var testimonial models.Testimonial
	if err := r.BaseRepository.FindByID( &testimonial, id); err != nil {
		return err
	}
	return r.BaseRepository.Delete( &testimonial)
}

// CountTestimonials returns the total number of testimonials
func (r *TestimonialRepository) CountTestimonials(ctx context.Context) (int64, error) {
	return r.BaseRepository.Count( &models.Testimonial{}, nil)
}

// RemoveFromFeatured removes the featured status from a testimonial
func (r *TestimonialRepository) RemoveFromFeatured(ctx context.Context, id uuid.UUID) error {
	return r.BaseRepository.DB.Model(&models.Testimonial{}).Where("id = ?", id).Update("is_featured", false).Error
}