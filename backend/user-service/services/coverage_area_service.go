package services

import (
	"context"
	"fmt"

	"user-service/models"
	"user-service/repositories"

	"github.com/google/uuid"
)

// ICoverageAreaService interface
type ICoverageAreaService interface {
	AddCoverageArea(ctx context.Context, instructorID uuid.UUID, areaType models.AreaType, areaID uint) (*models.InstructorArea, error)
	RemoveCoverageArea(ctx context.Context, instructorID uuid.UUID, areaType models.AreaType, areaID uint) error
	GetCoverageAreas(ctx context.Context, instructorID uuid.UUID) ([]models.InstructorAreaWithDetails, error)
}

// CoverageAreaService handles coverage area business logic
type CoverageAreaService struct {
	coverageAreaRepo repositories.ICoverageAreaRepository
	regionService    IRegionService
}

// NewCoverageAreaService creates a new coverage area service
func NewCoverageAreaService(
	coverageAreaRepo repositories.ICoverageAreaRepository,
	regionService IRegionService,
) ICoverageAreaService {
	return &CoverageAreaService{
		coverageAreaRepo: coverageAreaRepo,
		regionService:    regionService,
	}
}

// AddCoverageArea adds a coverage area to an instructor
func (s *CoverageAreaService) AddCoverageArea(ctx context.Context, instructorID uuid.UUID, areaType models.AreaType, areaID uint) (*models.InstructorArea, error) {
	// Validate the area exists in core-service
	areaName, err := s.regionService.FindAreaByID(ctx, areaType, areaID)
	if err != nil {
		return nil, fmt.Errorf("invalid area: %w", err)
	}

	// Verify area doesn't already exist for this instructor
	existing, err := s.coverageAreaRepo.FindByInstructorAndArea(ctx, instructorID, areaType, areaID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("coverage area already exists")
	}

	area := &models.InstructorArea{
		InstructorID: instructorID,
		AreaType:     areaType,
		AreaID:       areaID,
	}

	if err := s.coverageAreaRepo.AddCoverageArea(ctx, area); err != nil {
		return nil, err
	}

	return area, nil
}

// RemoveCoverageArea removes a coverage area from an instructor
func (s *CoverageAreaService) RemoveCoverageArea(ctx context.Context, instructorID uuid.UUID, areaType models.AreaType, areaID uint) error {
	return s.coverageAreaRepo.RemoveCoverageArea(ctx, instructorID, areaType, areaID)
}

// GetCoverageAreas retrieves all coverage areas for an instructor with region details
func (s *CoverageAreaService) GetCoverageAreas(ctx context.Context, instructorID uuid.UUID) ([]models.InstructorAreaWithDetails, error) {
	areas, err := s.coverageAreaRepo.FindCoverageAreasByInstructorID(ctx, instructorID)
	if err != nil {
		return nil, err
	}

	// Enrich with area names from core-service
	result := make([]models.InstructorAreaWithDetails, 0, len(areas))
	for _, area := range areas {
		areaName, err := s.regionService.FindAreaByID(ctx, area.AreaType, area.AreaID)
		if err != nil {
			areaName = fmt.Sprintf("Unknown (%d)", area.AreaID)
		}

		result = append(result, models.InstructorAreaWithDetails{
			InstructorID: area.InstructorID,
			AreaType:     area.AreaType,
			AreaID:       area.AreaID,
			AreaName:     areaName,
		})
	}

	return result, nil
}