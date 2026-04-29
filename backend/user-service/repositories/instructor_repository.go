package repositories

import (
	"user-service/models"
	"user-service/pkg/base"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IInstructorRepository interface {
	FindInstructorProfileByUserID(userID uuid.UUID) (*models.InstructorProfile, error)
	UpdateInstructorProfile(profile *models.InstructorProfile) error
	CreateInstructorProfile(profile *models.InstructorProfile) error
	DeleteInstructorProfile(instructorID uuid.UUID) error
}

type InstructorRepository struct {
	*base.BaseRepository
}

// DeleteInstructorProfile implements [IInstructorRepository].
func (i *InstructorRepository) DeleteInstructorProfile(instructorID uuid.UUID) error {
	return i.BaseRepository.Delete(&models.InstructorProfile{UserID: instructorID})
}

// CreateInstructorProfile implements [IInstructorRepository].
func (i *InstructorRepository) CreateInstructorProfile(profile *models.InstructorProfile) error {
	return i.BaseRepository.Create(profile)
}

// FindInstructorProfileByUserID implements [IInstructorRepository].
func (i *InstructorRepository) FindInstructorProfileByUserID(userID uuid.UUID) (*models.InstructorProfile, error) {
	var profile models.InstructorProfile
	if err := i.BaseRepository.FindOne(&profile, "user_id = ?", userID); err != nil {
		return nil, err
	}
	return &profile, nil
}

// UpdateInstructorProfile implements [IInstructorRepository].
func (i *InstructorRepository) UpdateInstructorProfile(profile *models.InstructorProfile) error {
	return i.BaseRepository.Update(profile)
}

func NewInstructorRepository(db *gorm.DB) IInstructorRepository {
	return &InstructorRepository{BaseRepository: base.NewBaseRepository(db)}
}
