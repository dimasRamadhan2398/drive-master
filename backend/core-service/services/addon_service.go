package services

import (
	"context"
	"core-service/models"
	"core-service/repositories"

	"github.com/google/uuid"
)

type IAddOnService interface {
	CreateAddOn(ctx context.Context, addon *models.AddOn) error
	GetAddOnByID(ctx context.Context, id uuid.UUID) (*models.AddOn, error)
	GetAllAddOns(ctx context.Context) ([]models.AddOn, error)
	GetAllAddOnsPaginated(ctx context.Context, page, limit int) ([]models.AddOn, int64, error)
	GetActiveAddOns(ctx context.Context) ([]models.AddOn, error)
	UpdateAddOn(ctx context.Context, addon *models.AddOn) error
	DeleteAddOn(ctx context.Context, addon *models.AddOn) error
	ToggleStatusAddOn(ctx context.Context, id uuid.UUID) (*models.AddOn, error)
}

type AddOnService struct {
	addOnRepo repositories.IAddOnRepository
}

func NewAddOnService(addOnRepo repositories.IAddOnRepository) IAddOnService {
	return &AddOnService{
		addOnRepo: addOnRepo,
	}
}

// CreateAddOn creates a new add-on
func (s *AddOnService) CreateAddOn(ctx context.Context, addon *models.AddOn) error {
	return s.addOnRepo.Create(ctx, addon)
}

// GetAddOnByID retrieves an add-on by ID
func (s *AddOnService) GetAddOnByID(ctx context.Context, id uuid.UUID) (*models.AddOn, error) {
	return s.addOnRepo.FindByID(ctx, id)
}

// GetAllAddOns retrieves all add-ons
func (s *AddOnService) GetAllAddOns(ctx context.Context) ([]models.AddOn, error) {
	return s.addOnRepo.FindAll(ctx)
}

// GetAllAddOnsPaginated retrieves add-ons with pagination
func (s *AddOnService) GetAllAddOnsPaginated(ctx context.Context, page, limit int) ([]models.AddOn, int64, error) {
	return s.addOnRepo.FindAllPaginated(ctx, page, limit)
}

// GetActiveAddOns retrieves all active add-ons
func (s *AddOnService) GetActiveAddOns(ctx context.Context) ([]models.AddOn, error) {
	return s.addOnRepo.FindByStatus(ctx, models.AddOnStatusActive)
}

// UpdateAddOn updates an add-on
func (s *AddOnService) UpdateAddOn(ctx context.Context, addon *models.AddOn) error {
	return s.addOnRepo.Update(ctx, addon)
}

// DeleteAddOn deletes an add-on
func (s *AddOnService) DeleteAddOn(ctx context.Context, addon *models.AddOn) error {
	return s.addOnRepo.Delete(ctx, addon)
}

// ToggleStatusAddOn toggles the add-on status between active and inactive
func (s *AddOnService) ToggleStatusAddOn(ctx context.Context, id uuid.UUID) (*models.AddOn, error) {
	return s.addOnRepo.ToggleStatus(ctx, id)
}
