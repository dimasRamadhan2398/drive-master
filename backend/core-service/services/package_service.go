package services

import (
	"context"
	"core-service/models"
	"core-service/pkg/kafka"
	"core-service/repositories"

	"github.com/google/uuid"
)

type IPackageService interface {
	CreatePackage(ctx context.Context, pkg *models.Package) error
	GetPackageByID(ctx context.Context, id uuid.UUID) (*models.Package, error)
	GetPackageByIDWithBenefits(ctx context.Context, id uuid.UUID) (*models.Package, error)
	GetAllPackages(ctx context.Context) ([]models.Package, error)
	GetAllPackagesPaginated(ctx context.Context, page, limit int) ([]models.Package, int64, error)
	GetPackagesByType(ctx context.Context, packageType models.PackageType) ([]models.Package, error)
	GetPackagesByStatus(ctx context.Context, status models.PackageStatus) ([]models.Package, error)
	UpdatePackage(ctx context.Context, pkg *models.Package) error
	DeletePackage(ctx context.Context, pkg *models.Package) error
	CountPackages(ctx context.Context) (int64, error)
	ToggleStatusPackage(ctx context.Context, id uuid.UUID) (*models.Package, error)
	IncrementStudentCount(ctx context.Context, id uuid.UUID) error
}

type PackageService struct {
	packageRepo   repositories.IPackageRepository
	eventPublisher *kafka.EventPublisher
}

func NewPackageService(packageRepo repositories.IPackageRepository, eventPublisher *kafka.EventPublisher) IPackageService {
	return &PackageService{
		packageRepo:   packageRepo,
		eventPublisher: eventPublisher,
	}
}

// CreatePackage creates a new package
func (s *PackageService) CreatePackage(ctx context.Context, pkg *models.Package) error {
	if err := s.packageRepo.Create(ctx, pkg); err != nil {
		return err
	}

	// Publish event (async to not block response)
	if s.eventPublisher != nil {
		go s.eventPublisher.PublishPackageCreated(context.Background(), pkg)
	}

	return nil
}

// GetPackageByID retrieves a package by ID
func (s *PackageService) GetPackageByID(ctx context.Context, id uuid.UUID) (*models.Package, error) {
	return s.packageRepo.FindByID(ctx, id)
}

// GetPackageByIDWithBenefits retrieves a package by ID with benefits preloaded
func (s *PackageService) GetPackageByIDWithBenefits(ctx context.Context, id uuid.UUID) (*models.Package, error) {
	return s.packageRepo.FindByIDWithBenefits(ctx, id)
}

// GetAllPackages retrieves all packages
func (s *PackageService) GetAllPackages(ctx context.Context) ([]models.Package, error) {
	return s.packageRepo.FindAll(ctx)
}

// GetAllPackagesPaginated retrieves packages with pagination
func (s *PackageService) GetAllPackagesPaginated(ctx context.Context, page, limit int) ([]models.Package, int64, error) {
	return s.packageRepo.FindAllPaginated(ctx, page, limit)
}

// GetPackagesByType retrieves packages by type
func (s *PackageService) GetPackagesByType(ctx context.Context, packageType models.PackageType) ([]models.Package, error) {
	return s.packageRepo.FindByType(ctx, packageType)
}

// GetPackagesByStatus retrieves packages by status
func (s *PackageService) GetPackagesByStatus(ctx context.Context, status models.PackageStatus) ([]models.Package, error) {
	return s.packageRepo.FindByStatus(ctx, status)
}

// UpdatePackage updates a package
func (s *PackageService) UpdatePackage(ctx context.Context, pkg *models.Package) error {
	if err := s.packageRepo.Update(ctx, pkg); err != nil {
		return err
	}

	// Publish event (async to not block response)
	if s.eventPublisher != nil {
		go s.eventPublisher.PublishPackageUpdated(context.Background(), pkg)
	}

	return nil
}

// DeletePackage deletes a package
func (s *PackageService) DeletePackage(ctx context.Context, pkg *models.Package) error {
	packageID := pkg.ID.String()

	if err := s.packageRepo.Delete(ctx, pkg); err != nil {
		return err
	}

	// Publish event (async to not block response)
	if s.eventPublisher != nil {
		go s.eventPublisher.PublishPackageDeleted(context.Background(), packageID)
	}

	return nil
}

// CountPackages returns the total number of packages
func (s *PackageService) CountPackages(ctx context.Context) (int64, error) {
	return s.packageRepo.Count(ctx)
}

// ToggleStatusPackage toggles the package status between active and inactive
func (s *PackageService) ToggleStatusPackage(ctx context.Context, id uuid.UUID) (*models.Package, error) {
	pkg, err := s.packageRepo.ToggleStatus(ctx, id)
	if err != nil {
		return nil, err
	}

	// Publish event (async to not block response)
	if s.eventPublisher != nil {
		go s.eventPublisher.PublishPackageUpdated(context.Background(), pkg)
	}

	return pkg, nil
}

// IncrementStudentCount increments the student count for a package
func (s *PackageService) IncrementStudentCount(ctx context.Context, id uuid.UUID) error {
	return s.packageRepo.IncrementStudentCount(ctx, id)
}