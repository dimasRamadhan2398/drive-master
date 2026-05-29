package services

import (
	"context"
	"core-service/models"
	"core-service/pkg/kafka"
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
	PublishTestimonial(ctx context.Context, id uuid.UUID, publishedBy string) error
	ArchiveTestimonial(ctx context.Context, id uuid.UUID, archivedBy string) error
	CountTestimonials(ctx context.Context) (int64, error)
}

type TestimonialService struct {
	testimonialRepo repositories.ITestimonialRepository
	eventPublisher  *kafka.EventPublisher
}

func NewTestimonialService(testimonialRepo repositories.ITestimonialRepository, eventPublisher *kafka.EventPublisher) ITestimonialService {
	return &TestimonialService{
		testimonialRepo: testimonialRepo,
		eventPublisher:  eventPublisher,
	}
}

// CreateTestimonial creates a new testimonial
func (s *TestimonialService) CreateTestimonial(ctx context.Context, testimonial *models.Testimonial) error {
	if err := s.testimonialRepo.CreateTestimonial(ctx, testimonial); err != nil {
		return err
	}

	// Publish event (async to not block response)
	if s.eventPublisher != nil {
		go s.eventPublisher.PublishTestimonialCreated(context.Background(), testimonial)
	}

	return nil
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
	if err := s.testimonialRepo.UpdateTestimonial(ctx, testimonial); err != nil {
		return err
	}

	// Publish event (async to not block response)
	if s.eventPublisher != nil {
		go s.eventPublisher.PublishTestimonialUpdated(context.Background(), testimonial)
	}

	return nil
}

// DeleteTestimonial deletes a testimonial
func (s *TestimonialService) DeleteTestimonial(ctx context.Context, id uuid.UUID) error {
	testimonialID := id.String()

	if err := s.testimonialRepo.DeleteTestimonial(ctx, id); err != nil {
		return err
	}

	// Publish event (async to not block response)
	if s.eventPublisher != nil {
		go s.eventPublisher.PublishTestimonialDeleted(context.Background(), testimonialID)
	}

	return nil
}

// PublishTestimonial publishes a testimonial
func (s *TestimonialService) PublishTestimonial(ctx context.Context, id uuid.UUID, publishedBy string) error {
	testimonial, err := s.testimonialRepo.GetTestimonialByID(ctx, id)
	if err != nil {
		return err
	}

	testimonial.Status = models.TestimonialStatusPublished
	if err := s.testimonialRepo.UpdateTestimonial(ctx, testimonial); err != nil {
		return err
	}

	// Publish event (async to not block response)
	if s.eventPublisher != nil {
		go s.eventPublisher.PublishTestimonialPublished(context.Background(), id.String(), publishedBy)
	}

	return nil
}

// ArchiveTestimonial archives a testimonial
func (s *TestimonialService) ArchiveTestimonial(ctx context.Context, id uuid.UUID, archivedBy string) error {
	testimonial, err := s.testimonialRepo.GetTestimonialByID(ctx, id)
	if err != nil {
		return err
	}

	testimonial.Status = models.TestimonialStatusArchived
	if err := s.testimonialRepo.UpdateTestimonial(ctx, testimonial); err != nil {
		return err
	}

	// Publish event (async to not block response)
	if s.eventPublisher != nil {
		go s.eventPublisher.PublishTestimonialArchived(context.Background(), id.String(), archivedBy)
	}

	return nil
}

// CountTestimonials returns the total number of testimonials
func (s *TestimonialService) CountTestimonials(ctx context.Context) (int64, error) {
	return s.testimonialRepo.CountTestimonials(ctx)
}
