package repositories

import (
	"context"
	"core-service/models"
	"core-service/pkg/base"

	"github.com/google/uuid"
)

type IAddOnRepository interface {
	Create(ctx context.Context, addon *models.AddOn) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.AddOn, error)
	FindAll(ctx context.Context) ([]models.AddOn, error)
	FindAllPaginated(ctx context.Context, page, limit int) ([]models.AddOn, int64, error)
	FindByStatus(ctx context.Context, status models.AddOnStatus) ([]models.AddOn, error)
	Update(ctx context.Context, addon *models.AddOn) error
	Delete(ctx context.Context, addon *models.AddOn) error
	Count(ctx context.Context) (int64, error)
	ToggleStatus(ctx context.Context, id uuid.UUID) (*models.AddOn, error)
}

type AddOnRepository struct {
	*base.BaseRepository
}

func NewAddOnRepository(baseRepo *base.BaseRepository) IAddOnRepository {
	return &AddOnRepository{BaseRepository: baseRepo}
}

// Create creates a new add-on
func (r *AddOnRepository) Create(ctx context.Context, addon *models.AddOn) error {
	return r.BaseRepository.Create(ctx, addon)
}

// FindByID finds an add-on by ID
func (r *AddOnRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.AddOn, error) {
	var addon models.AddOn
	if err := r.BaseRepository.FindByID(ctx, &addon, id); err != nil {
		return nil, err
	}
	return &addon, nil
}

// FindAll retrieves all add-ons
func (r *AddOnRepository) FindAll(ctx context.Context) ([]models.AddOn, error) {
	var addons []models.AddOn
	opts := base.NewQueryOptions().WithOrder("sort_order ASC, created_at DESC")
	if err := r.BaseRepository.FindMany(ctx, &models.AddOn{}, &addons, opts); err != nil {
		return nil, err
	}
	return addons, nil
}

// FindAllPaginated retrieves add-ons with pagination
func (r *AddOnRepository) FindAllPaginated(ctx context.Context, page, limit int) ([]models.AddOn, int64, error) {
	var addons []models.AddOn

	// Count total
	total, err := r.BaseRepository.Count(ctx, &models.AddOn{}, nil)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	opts := base.NewQueryOptions().
		WithPagination(offset, limit).
		WithOrder("sort_order ASC, created_at DESC")

	if err := r.BaseRepository.FindMany(ctx, &models.AddOn{}, &addons, opts); err != nil {
		return nil, 0, err
	}

	return addons, total, nil
}

// FindByStatus retrieves add-ons by status
func (r *AddOnRepository) FindByStatus(ctx context.Context, status models.AddOnStatus) ([]models.AddOn, error) {
	var addons []models.AddOn
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"status": status}).
		WithOrder("sort_order ASC, created_at DESC")
	if err := r.BaseRepository.FindMany(ctx, &models.AddOn{}, &addons, opts); err != nil {
		return nil, err
	}
	return addons, nil
}

// Update updates an add-on
func (r *AddOnRepository) Update(ctx context.Context, addon *models.AddOn) error {
	return r.BaseRepository.Update(ctx, addon)
}

// Delete deletes an add-on
func (r *AddOnRepository) Delete(ctx context.Context, addon *models.AddOn) error {
	return r.BaseRepository.Delete(ctx, addon)
}

// Count returns the total number of add-ons
func (r *AddOnRepository) Count(ctx context.Context) (int64, error) {
	return r.BaseRepository.Count(ctx, &models.AddOn{}, nil)
}

// ToggleStatus toggles the add-on status between active and inactive
func (r *AddOnRepository) ToggleStatus(ctx context.Context, id uuid.UUID) (*models.AddOn, error) {
	var addon models.AddOn
	if err := r.BaseRepository.FindByID(ctx, &addon, id); err != nil {
		return nil, err
	}

	// Toggle status
	if addon.Status == models.AddOnStatusActive {
		addon.Status = models.AddOnStatusInactive
	} else {
		addon.Status = models.AddOnStatusActive
	}

	if err := r.BaseRepository.Update(ctx, &addon); err != nil {
		return nil, err
	}

	return &addon, nil
}
