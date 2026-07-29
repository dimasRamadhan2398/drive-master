package repositories

import (
	"context"
	"core-service/models"
	"core-service/pkg/base"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IContactRepository interface {
	Create(ctx context.Context, contact *models.ContactInquiry) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.ContactInquiry, error)
	FindAll(ctx context.Context) ([]models.ContactInquiry, error)
}

type ContactRepository struct {
	*base.BaseRepository
}

func NewContactRepository(baseRepo *base.BaseRepository) IContactRepository {
	return &ContactRepository{
		BaseRepository: baseRepo,
	}
}

func (r *ContactRepository) Create(ctx context.Context, contact *models.ContactInquiry) error {
	return r.BaseRepository.Create(ctx, contact)
}

func (r *ContactRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.ContactInquiry, error) {
	var contact models.ContactInquiry
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&contact).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &contact, nil
}

func (r *ContactRepository) FindAll(ctx context.Context) ([]models.ContactInquiry, error) {
	var contacts []models.ContactInquiry
	if err := r.DB.WithContext(ctx).Order("created_at DESC").Find(&contacts).Error; err != nil {
		return nil, err
	}
	return contacts, nil
}
