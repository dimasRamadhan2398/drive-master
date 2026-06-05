package services

import (
	"context"
	"core-service/models"
	"core-service/models/dto"
	"core-service/repositories"

	"github.com/google/uuid"
)

// IGeneralSettingsService defines the interface for general settings service
type IGeneralSettingsService interface {
	GetSettings(ctx context.Context) (*models.GeneralSettings, error)
	UpdateSettings(ctx context.Context, req *dto.UpdateGeneralSettingsRequest) (*models.GeneralSettings, error)
}

// GeneralSettingsService implements IGeneralSettingsService
type GeneralSettingsService struct {
	repo repositories.IGeneralSettingsRepository
}

func NewGeneralSettingsService(repo repositories.IGeneralSettingsRepository) IGeneralSettingsService {
	return &GeneralSettingsService{
		repo: repo,
	}
}

// GetSettings retrieves the general settings (singleton pattern)
func (s *GeneralSettingsService) GetSettings(ctx context.Context) (*models.GeneralSettings, error) {
	settings, err := s.repo.FindFirst(ctx)
	if err != nil {
		// If no settings exist, create default settings
		if err.Error() == "record not found" {
			defaultSettings := &models.GeneralSettings{
				BusinessName: "Drive",
				NotifyEmail:  true,
			}
			if createErr := s.repo.Create(ctx, defaultSettings); createErr != nil {
				return nil, createErr
			}
			return defaultSettings, nil
		}
		return nil, err
	}
	return settings, nil
}

// UpdateSettings updates the general settings
func (s *GeneralSettingsService) UpdateSettings(ctx context.Context, req *dto.UpdateGeneralSettingsRequest) (*models.GeneralSettings, error) {
	// Get current settings or create new one
	settings, err := s.repo.FindFirst(ctx)
	if err != nil {
		if err.Error() == "record not found" {
			// Create new settings with defaults
			settings = &models.GeneralSettings{
				NotifyEmail: true,
			}
		} else {
			return nil, err
		}
	}

	// Update fields if provided
	if req.BusinessName != nil {
		settings.BusinessName = *req.BusinessName
	}
	if req.Email != nil {
		settings.Email = *req.Email
	}
	if req.Phone != nil {
		settings.Phone = *req.Phone
	}
	if req.Fax != nil {
		settings.Fax = *req.Fax
	}
	if req.WhatsApp != nil {
		settings.WhatsApp = *req.WhatsApp
	}
	if req.Address != nil {
		settings.Address = *req.Address
	}
	if req.HoursMonFri != nil {
		settings.HoursMonFri = *req.HoursMonFri
	}
	if req.HoursSatSun != nil {
		settings.HoursSatSun = *req.HoursSatSun
	}
	if req.HoursNightShift != nil {
		settings.HoursNightShift = *req.HoursNightShift
	}
	if req.PromoEndDate != nil {
		settings.PromoEndDate = req.PromoEndDate
	}
	if req.NotifyEmail != nil {
		settings.NotifyEmail = *req.NotifyEmail
	}
	if req.NotifySMS != nil {
		settings.NotifySMS = *req.NotifySMS
	}
	if req.NotifyWhatsApp != nil {
		settings.NotifyWhatsApp = *req.NotifyWhatsApp
	}

	// Check if this is a new record
	isNew := settings.ID == uuid.Nil
	if isNew {
		if err := s.repo.Create(ctx, settings); err != nil {
			return nil, err
		}
	} else {
		if err := s.repo.Update(ctx, settings); err != nil {
			return nil, err
		}
	}

	return settings, nil
}