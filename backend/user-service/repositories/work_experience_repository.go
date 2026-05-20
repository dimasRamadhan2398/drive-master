package repositories

import (
	"context"
	"user-service/models"
	"user-service/pkg/base"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IWorkExperienceRepository interface {
	Create(ctx context.Context, workExp *models.WorkExperience) error
	CreateBatch(ctx context.Context, workExps []models.WorkExperience) error
	Update(ctx context.Context, workExp *models.WorkExperience) error
	Delete(ctx context.Context, id uint) error
	FindByInstructorID(ctx context.Context, instructorID uuid.UUID) ([]models.WorkExperience, error)
}

type WorkExperienceRepository struct {
	*base.BaseRepository
}

// CreateBatch implements [IWorkExperienceRepository].
func (r *WorkExperienceRepository) CreateBatch(ctx context.Context, workExps []models.WorkExperience) error {
	// Convert []models.WorkExperience to []interface{} for CreateInBatches
	interfaceSlice := make([]interface{}, len(workExps))
	for i, v := range workExps {
		interfaceSlice[i] = v
	}
	return r.BaseRepository.CreateInBatches(interfaceSlice)
}

func NewWorkExperienceRepository(db *gorm.DB) IWorkExperienceRepository {
	return &WorkExperienceRepository{BaseRepository: base.NewBaseRepository(db)}
}

func (r *WorkExperienceRepository) Create(ctx context.Context, workExp *models.WorkExperience) error {
	return r.BaseRepository.Create(workExp)
}

func (r *WorkExperienceRepository) Update(ctx context.Context, workExp *models.WorkExperience) error {
	return r.BaseRepository.Update(workExp)
}

func (r *WorkExperienceRepository) Delete(ctx context.Context, id uint) error {
	return r.BaseRepository.Delete(&models.WorkExperience{ID: id})
}

func (r *WorkExperienceRepository) FindByInstructorID(ctx context.Context, instructorID uuid.UUID) ([]models.WorkExperience, error) {
	var workExps []models.WorkExperience
	if err := r.BaseRepository.FindWithOptions(&models.WorkExperience{}, &workExps, &base.QueryOptions{
		Where: map[string]interface{}{
			"instructor_id": instructorID,
		},
		Order: "start_date DESC",
	}); err != nil {
		return nil, err
	}
	return workExps, nil
}