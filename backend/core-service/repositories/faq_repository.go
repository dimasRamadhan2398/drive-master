package repositories

import (
	"context"
	"core-service/models"
	"core-service/pkg/base"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// IFAQRepository defines the interface for FAQ repository
type IFAQRepository interface {
	Create(ctx context.Context, faq *models.FAQ) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.FAQ, error)
	FindAll(ctx context.Context) ([]models.FAQ, error)
	FindActive(ctx context.Context) ([]models.FAQ, error)
	Update(ctx context.Context, faq *models.FAQ) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteSoft(ctx context.Context, faq *models.FAQ) error
}

// FAQRepository implements IFAQRepository
type FAQRepository struct {
	*base.BaseRepository
}

// NewFAQRepository creates a new FAQ repository
func NewFAQRepository(baseRepo *base.BaseRepository) IFAQRepository {
	return &FAQRepository{
		BaseRepository: baseRepo,
	}
}

// Create creates a new FAQ
func (r *FAQRepository) Create(ctx context.Context, faq *models.FAQ) error {
	return r.BaseRepository.Create(ctx, faq)
}

// FindByID finds an FAQ by ID
func (r *FAQRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.FAQ, error) {
	var faq models.FAQ
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&faq).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &faq, nil
}

// FindAll finds all FAQs (including inactive)
func (r *FAQRepository) FindAll(ctx context.Context) ([]models.FAQ, error) {
	var faqs []models.FAQ
	err := r.DB.WithContext(ctx).Order("\"order\" ASC, created_at DESC").Find(&faqs).Error
	if err != nil {
		return nil, err
	}
	return faqs, nil
}

// FindActive finds all active FAQs
func (r *FAQRepository) FindActive(ctx context.Context) ([]models.FAQ, error) {
	var faqs []models.FAQ
	err := r.DB.WithContext(ctx).
		Where("is_active = ?", true).
		Order("\"order\" ASC, created_at DESC").
		Find(&faqs).Error
	if err != nil {
		return nil, err
	}
	return faqs, nil
}

// Update updates an FAQ
func (r *FAQRepository) Update(ctx context.Context, faq *models.FAQ) error {
	return r.BaseRepository.Update(ctx, faq)
}

// Delete permanently deletes an FAQ
func (r *FAQRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Where("id = ?", id).Delete(&models.FAQ{}).Error
}

// DeleteSoft soft-deletes an FAQ
func (r *FAQRepository) DeleteSoft(ctx context.Context, faq *models.FAQ) error {
	return r.BaseRepository.Delete(ctx, faq)
}