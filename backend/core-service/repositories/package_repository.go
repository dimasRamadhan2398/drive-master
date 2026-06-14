package repositories

import (
	"context"
	"core-service/models"
	"core-service/pkg/base"

	"github.com/google/uuid"
)

type IPackageRepository interface {
	Create(ctx context.Context, pkg *models.Package) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Package, error)
	FindByIDWithBenefits(ctx context.Context, id uuid.UUID) (*models.Package, error)
	FindAll(ctx context.Context) ([]models.Package, error)
	FindAllPaginated(ctx context.Context, page, limit int) ([]models.Package, int64, error)
	FindByType(ctx context.Context, packageType models.PackageType) ([]models.Package, error)
	FindByStatus(ctx context.Context, status models.PackageStatus) ([]models.Package, error)
	Update(ctx context.Context, pkg *models.Package) error
	Delete(ctx context.Context, pkg *models.Package) error
	Count(ctx context.Context) (int64, error)
	ToggleStatus(ctx context.Context, id uuid.UUID) (*models.Package, error)
}

type PackageRepository struct {
	*base.BaseRepository
}

func NewPackageRepository(baseRepo *base.BaseRepository) IPackageRepository {
	return &PackageRepository{BaseRepository: baseRepo}
}

// Create creates a new package
func (r *PackageRepository) Create(ctx context.Context, pkg *models.Package) error {
	return r.BaseRepository.Create(ctx, pkg)
}

// FindByID finds a package by ID
func (r *PackageRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Package, error) {
	var pkg models.Package
	if err := r.BaseRepository.FindByID(ctx, &pkg, id); err != nil {
		return nil, err
	}
	return &pkg, nil
}

// FindByIDWithBenefits finds a package by ID with benefits preloaded
func (r *PackageRepository) FindByIDWithBenefits(ctx context.Context, id uuid.UUID) (*models.Package, error) {
	var pkg models.Package
	opts := base.NewQueryOptions().WithPreloads("Benefits")
	if err := r.BaseRepository.FindOne(ctx, &pkg, "id = ?", id, opts); err != nil {
		return nil, err
	}
	return &pkg, nil
}

// FindAll retrieves all packages
func (r *PackageRepository) FindAll(ctx context.Context) ([]models.Package, error) {
	var packages []models.Package
	opts := base.NewQueryOptions()
	if err := r.BaseRepository.FindMany(ctx, &models.Package{}, &packages, opts); err != nil {
		return nil, err
	}
	return packages, nil
}

// FindAllPaginated retrieves packages with pagination
func (r *PackageRepository) FindAllPaginated(ctx context.Context, page, limit int) ([]models.Package, int64, error) {
	var packages []models.Package
	

	// Count total
	total, err := r.BaseRepository.Count(ctx, &models.Package{}, nil)
	if err != nil {
		return nil, 0, err
	}
	
	// Get paginated results
	offset := (page - 1) * limit
	opts := base.NewQueryOptions().
		WithPagination(offset, limit).
		WithOrder("created_at DESC")

	if err := r.BaseRepository.FindMany(ctx, &models.Package{}, &packages, opts); err != nil {
		return nil, 0, err
	}

	return packages, total, nil
}

// FindByType retrieves packages by type
func (r *PackageRepository) FindByType(ctx context.Context, packageType models.PackageType) ([]models.Package, error) {
	var packages []models.Package
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"package_type": packageType})
	if err := r.BaseRepository.FindMany(ctx, &models.Package{}, &packages, opts); err != nil {
		return nil, err
	}
	return packages, nil
}

// FindByStatus retrieves packages by status
func (r *PackageRepository) FindByStatus(ctx context.Context, status models.PackageStatus) ([]models.Package, error) {
	var packages []models.Package
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"status": status})
	if err := r.BaseRepository.FindMany(ctx, &models.Package{}, &packages, opts); err != nil {
		return nil, err
	}
	return packages, nil
}

// Update updates a package
func (r *PackageRepository) Update(ctx context.Context, pkg *models.Package) error {
	return r.BaseRepository.Update(ctx, pkg)
}

// Delete deletes a package
func (r *PackageRepository) Delete(ctx context.Context, pkg *models.Package) error {
	return r.BaseRepository.Delete(ctx, pkg)
}

// Count returns the total number of packages
func (r *PackageRepository) Count(ctx context.Context) (int64, error) {
	return r.BaseRepository.Count(ctx, &models.Package{}, nil)
}

// ToggleStatus toggles the package status between active and inactive
func (r *PackageRepository) ToggleStatus(ctx context.Context, id uuid.UUID) (*models.Package, error) {
	var pkg models.Package
	if err := r.BaseRepository.FindByID(ctx, &pkg, id); err != nil {
		return nil, err
	}

	// Toggle status
	if pkg.Status == models.PackageStatusActive {
		pkg.Status = models.PackageStatusInactive
	} else {
		pkg.Status = models.PackageStatusActive
	}

	if err := r.BaseRepository.Update(ctx, &pkg); err != nil {
		return nil, err
	}

	return &pkg, nil
}