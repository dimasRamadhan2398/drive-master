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
	ReorderFAQs(ctx context.Context, id uuid.UUID, newOrder int) error
	GetMaxOrder(ctx context.Context) (int, error)
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

// ReorderFAQs reorders an FAQ and adjusts other FAQs accordingly
// If moving down (e.g., order 4 to 2), items between new and old position shift down
// If moving up (e.g., order 2 to 4), items between new and old position shift up
func (r *FAQRepository) ReorderFAQs(ctx context.Context, id uuid.UUID, newOrder int) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get the FAQ being moved
		var faq models.FAQ
		if err := tx.Where("id = ?", id).First(&faq).Error; err != nil {
			return err
		}

		oldOrder := faq.Order

		// If order hasn't changed, nothing to do
		if oldOrder == newOrder {
			return nil
		}

		if newOrder > oldOrder {
			// Moving down: shift items between old and new position up by 1
			// e.g., moving from order 2 to 4: items at 2,3,4 -> 3,4,5
			if err := tx.Model(&models.FAQ{}).
				Where("\"order\" > ? AND \"order\" <= ?", oldOrder, newOrder).
				Update("\"order\"", gorm.Expr("\"order\" - 1")).Error; err != nil {
				return err
			}
		} else {
			// Moving up: shift items between new and old position down by 1
			// e.g., moving from order 4 to 2: items at 2,3,4 -> 1,2,3
			if err := tx.Model(&models.FAQ{}).
				Where("\"order\" >= ? AND \"order\" < ?", newOrder, oldOrder).
				Update("\"order\"", gorm.Expr("\"order\" + 1")).Error; err != nil {
				return err
			}
		}

		// Update the FAQ's new order
		faq.Order = newOrder
		if err := tx.Save(&faq).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetMaxOrder returns the maximum order value among all FAQs
func (r *FAQRepository) GetMaxOrder(ctx context.Context) (int, error) {
	var maxOrder int
	err := r.DB.WithContext(ctx).Model(&models.FAQ{}).Select("COALESCE(MAX(\"order\"), 0)").Scan(&maxOrder).Error
	return maxOrder, err
}