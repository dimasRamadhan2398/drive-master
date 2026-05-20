package repositories

import (
	"context"
	"user-service/models"
	"user-service/pkg/base"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ICoverageAreaRepository interface {
	FindCoverageAreasByInstructorID(ctx context.Context, instructorID uuid.UUID) ([]models.InstructorArea, error)
	FindByInstructorAndArea(ctx context.Context, instructorID uuid.UUID, areaType models.AreaType, areaID uint) (*models.InstructorArea, error)
	AddCoverageArea(ctx context.Context, area *models.InstructorArea) error
	RemoveCoverageArea(ctx context.Context, instructorID uuid.UUID, areaType models.AreaType, areaID uint) error
}

type CoverageAreaRepository struct {
	*base.BaseRepository
}

// AddCoverageArea implements [ICoverageAreaRepository].
func (c *CoverageAreaRepository) AddCoverageArea(ctx context.Context, area *models.InstructorArea) error {
	return c.BaseRepository.Create(area)
}

// FindCoverageAreasByInstructorID implements [ICoverageAreaRepository].
func (c *CoverageAreaRepository) FindCoverageAreasByInstructorID(ctx context.Context, instructorID uuid.UUID) ([]models.InstructorArea, error) {
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

// FindByInstructorAndArea finds a specific coverage area for an instructor
func (c *CoverageAreaRepository) FindByInstructorAndArea(ctx context.Context, instructorID uuid.UUID, areaType models.AreaType, areaID uint) (*models.InstructorArea, error) {
	var area models.InstructorArea
	if err := c.BaseRepository.FindWithOptions(&models.InstructorArea{}, &area, &base.QueryOptions{
		Where: map[string]interface{}{
			"instructor_id": instructorID,
			"area_type":      areaType,
			"area_id":        areaID,
		},
		Limit: 1,
	}); err != nil {
		return nil, err
	}
	return &area, nil
}

// RemoveCoverageArea implements [ICoverageAreaRepository].
func (c *CoverageAreaRepository) RemoveCoverageArea(ctx context.Context, instructorID uuid.UUID, areaType models.AreaType, areaID uint) error {
	area := &models.InstructorArea{
		InstructorID: instructorID,
		AreaType:     areaType,
		AreaID:       areaID,
	}
	return c.BaseRepository.Delete(area)
}

func NewCoverageArea(db *gorm.DB) ICoverageAreaRepository {
	return &CoverageAreaRepository{BaseRepository: base.NewBaseRepository(db)}
}