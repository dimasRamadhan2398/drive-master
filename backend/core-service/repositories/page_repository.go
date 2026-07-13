package repositories

import (
	"context"
	"core-service/models"
	"core-service/pkg/base"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IPageRepository interface {
	Create(ctx context.Context, page *models.Page) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Page, error)
	FindBySlug(ctx context.Context, slug string) (*models.Page, error)
	FindAll(ctx context.Context) ([]models.Page, error)
	Update(ctx context.Context, page *models.Page) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type PageRepository struct {
	*base.BaseRepository
}

func NewPageRepository(baseRepo *base.BaseRepository) IPageRepository {
	return &PageRepository{
		BaseRepository: baseRepo,
	}
}

func (r *PageRepository) Create(ctx context.Context, page *models.Page) error {
	return r.BaseRepository.Create(ctx, page)
}

func (r *PageRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Page, error) {
	var page models.Page
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&page).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &page, nil
}

func (r *PageRepository) FindBySlug(ctx context.Context, slug string) (*models.Page, error) {
	var page models.Page
	if err := r.DB.WithContext(ctx).Where("slug = ?", slug).First(&page).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &page, nil
}

func (r *PageRepository) FindAll(ctx context.Context) ([]models.Page, error) {
	var pages []models.Page
	if err := r.DB.WithContext(ctx).Order("created_at DESC").Find(&pages).Error; err != nil {
		return nil, err
	}
	return pages, nil
}

func (r *PageRepository) Update(ctx context.Context, page *models.Page) error {
	return r.BaseRepository.Update(ctx, page)
}

func (r *PageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Delete(&models.Page{}, "id = ?", id).Error
}
