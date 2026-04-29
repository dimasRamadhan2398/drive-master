package repositories

import (
	"user-service/models"
	"user-service/pkg/base"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ICoverageAreaRepository interface {
	FindCoverageAreasByInstructorID(instructorID uuid.UUID) ([]models.InstructorArea, error)
	AddCoverageArea(area *models.InstructorArea) error
	RemoveCoverageArea(instructorID uuid.UUID, areaID uint) error
}

type CoverageAreaRepository struct {
	*base.BaseRepository
}

// AddCoverageArea implements [ICoverageAreaRepository].
func (c *CoverageAreaRepository) AddCoverageArea(area *models.InstructorArea) error {
	return c.BaseRepository.Create(area)
}

// FindCoverageAreasByInstructorID implements [ICoverageAreaRepository].
func (c *CoverageAreaRepository) FindCoverageAreasByInstructorID(instructorID uuid.UUID) ([]models.InstructorArea, error) {
	var areas []models.InstructorArea
	if err := c.BaseRepository.FindWithOptions(&models.InstructorArea{}, &areas, &base.QueryOptions{
		Where: map[string]interface{}{
			"instructor_id": instructorID,
		},
	}); err != nil {
		return nil, err
	}
	return areas, nil
}

// RemoveCoverageArea implements [ICoverageAreaRepository].
func (c *CoverageAreaRepository) RemoveCoverageArea(instructorID uuid.UUID, areaID uint) error {
	err := c.BaseRepository.Delete(&models.InstructorArea{InstructorID: instructorID})
	if err != nil {
		return err
	}
	return nil
}

func NewCoverageArea(db *gorm.DB) ICoverageAreaRepository {
	return &CoverageAreaRepository{BaseRepository: base.NewBaseRepository(db)}
}
