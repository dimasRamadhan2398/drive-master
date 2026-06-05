package repositories

import (
	"context"
	"core-service/models"
	"core-service/pkg/base"

	"github.com/google/uuid"
)

// IGeneralSettingsRepository defines the interface for general settings repository
type IGeneralSettingsRepository interface {
	Create(ctx context.Context, settings *models.GeneralSettings) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.GeneralSettings, error)
	FindFirst(ctx context.Context) (*models.GeneralSettings, error)
	Update(ctx context.Context, settings *models.GeneralSettings) error
	Delete(ctx context.Context, settings *models.GeneralSettings) error
}

// GeneralSettingsRepository implements IGeneralSettingsRepository
type GeneralSettingsRepository struct {
	*base.BaseRepository
}

func NewGeneralSettingsRepository(baseRepo *base.BaseRepository) IGeneralSettingsRepository {
	return &GeneralSettingsRepository{
		BaseRepository: baseRepo,
	}
}

// Create creates a new general settings record
func (r *GeneralSettingsRepository) Create(ctx context.Context, settings *models.GeneralSettings) error {
	return r.BaseRepository.Create(ctx, settings)
}

// FindByID finds general settings by ID
func (r *GeneralSettingsRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.GeneralSettings, error) {
	var settings models.GeneralSettings
	if err := r.BaseRepository.FindByID(ctx, &settings, id); err != nil {
		return nil, err
	}
	return &settings, nil
}

// FindFirst finds the first general settings record (singleton pattern)
func (r *GeneralSettingsRepository) FindFirst(ctx context.Context) (*models.GeneralSettings, error) {
	var settings models.GeneralSettings
	if err := r.DB.WithContext(ctx).First(&settings).Error; err != nil {
		return nil, err
	}
	return &settings, nil
}

// Update updates general settings
func (r *GeneralSettingsRepository) Update(ctx context.Context, settings *models.GeneralSettings) error {
	return r.BaseRepository.Update(ctx, settings)
}

// Delete deletes general settings (rarely used, but included for completeness)
func (r *GeneralSettingsRepository) Delete(ctx context.Context, settings *models.GeneralSettings) error {
	return r.BaseRepository.Delete(ctx, settings)
}