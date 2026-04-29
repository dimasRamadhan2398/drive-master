package repositories

import (
	"user-service/models"
	"user-service/pkg/base"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IWorkExperienceRepository interface {
	Create(workExp *models.WorkExperience) error
	Update(workExp *models.WorkExperience) error
	Delete(id uint) error
	FindByInstructorID(instructorID uuid.UUID) ([]models.WorkExperience, error)
}

type WorkExperienceRepository struct {
	*base.BaseRepository
}

func NewWorkExperienceRepository(db *gorm.DB) IWorkExperienceRepository {
	return &WorkExperienceRepository{BaseRepository: base.NewBaseRepository(db)}
}

func (r *WorkExperienceRepository) Create(workExp *models.WorkExperience) error {
	return r.BaseRepository.Create(workExp)
}

func (r *WorkExperienceRepository) Update(workExp *models.WorkExperience) error {
	return r.BaseRepository.Update(workExp)
}

func (r *WorkExperienceRepository) Delete(id uint) error {
	return r.BaseRepository.Delete(&models.WorkExperience{ID: id})
}

func (r *WorkExperienceRepository) FindByInstructorID(instructorID uuid.UUID) ([]models.WorkExperience, error) {
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



