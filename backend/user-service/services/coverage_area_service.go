package services

import (
	"context"
	"user-service/models"
	"user-service/repositories"

	"github.com/google/uuid"
)

type ICoverageAreaService interface {
	AddCoverageArea(ctx context.Context, instructorID uuid.UUID, areaID uint) (*models.InstructorArea, error)
	RemoveCoverageArea(ctx context.Context, instructorID uuid.UUID, areaID uint) error
	GetCoverageAreas(ctx context.Context, instructorID uuid.UUID) ([]models.InstructorArea, error)
}

type CoverageAreaService struct {
	coverageAreaRepo repositories.ICoverageAreaRepository
}

func NewCoverageAreaService(
	coverageAreaRepo repositories.ICoverageAreaRepository,
) ICoverageAreaService {
	return &CoverageAreaService{
		coverageAreaRepo: coverageAreaRepo,
	}
}

// AddCoverageArea adds a coverage area to an instructor
func (s *CoverageAreaService) AddCoverageArea(ctx context.Context, instructorID uuid.UUID, areaID uint) (*models.InstructorArea, error) {
	area := &models.InstructorArea{
		InstructorID: instructorID,
		AreaID:       areaID,
	}

	if err := s.coverageAreaRepo.AddCoverageArea(ctx, area); err != nil {
		return nil, err
	}

	return area, nil
}

// RemoveCoverageArea removes a coverage area from an instructor
func (s *CoverageAreaService) RemoveCoverageArea(ctx context.Context, instructorID uuid.UUID, areaID uint) error {
	return s.coverageAreaRepo.RemoveCoverageArea(ctx, instructorID, areaID)
}

// GetCoverageAreas retrieves all coverage areas for an instructor
func (s *CoverageAreaService) GetCoverageAreas(ctx context.Context, instructorID uuid.UUID) ([]models.InstructorArea, error) {
	return s.coverageAreaRepo.FindCoverageAreasByInstructorID(ctx, instructorID)
}
