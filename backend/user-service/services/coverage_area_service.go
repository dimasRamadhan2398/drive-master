package services

import (
	"user-service/models"
	"user-service/repositories"

	"github.com/google/uuid"
)

type ICoverageAreaService interface {
	AddCoverageArea(instructorID uuid.UUID, areaID uint) (*models.InstructorArea, error)
	RemoveCoverageArea(instructorID uuid.UUID, areaID uint) error
	GetCoverageAreas(instructorID uuid.UUID) ([]models.InstructorArea, error)
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
func (s *CoverageAreaService) AddCoverageArea(instructorID uuid.UUID, areaID uint) (*models.InstructorArea, error) {
	area := &models.InstructorArea{
		InstructorID: instructorID,
		AreaID:       areaID,
	}

	if err := s.coverageAreaRepo.AddCoverageArea(area); err != nil {
		return nil, err
	}

	return area, nil
}

// RemoveCoverageArea removes a coverage area from an instructor
func (s *CoverageAreaService) RemoveCoverageArea(instructorID uuid.UUID, areaID uint) error {
	return s.coverageAreaRepo.RemoveCoverageArea(instructorID, areaID)
}

// GetCoverageAreas retrieves all coverage areas for an instructor
func (s *CoverageAreaService) GetCoverageAreas(instructorID uuid.UUID) ([]models.InstructorArea, error) {
	return s.coverageAreaRepo.FindCoverageAreasByInstructorID(instructorID)
}
