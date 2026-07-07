package repositories

import (
	"context"
	"user-service/models"
	"user-service/pkg/base"
	apperrors "user-service/pkg/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IInstructorRepository interface {
	FindInstructorProfileByUserID(ctx context.Context, userID uuid.UUID) (*models.InstructorProfile, error)
	UpdateInstructorProfile(ctx context.Context, profile *models.InstructorProfile) error
	UpdateInstructorPhotoURL(ctx context.Context, userID uuid.UUID, photoURL string) error
	CreateInstructorProfile(ctx context.Context, profile *models.InstructorProfile) error
	CreateInstructorProfileTx(tx *gorm.DB, profile *models.InstructorProfile) error
	DeleteInstructorProfile(ctx context.Context, instructorID uuid.UUID) error
}

type InstructorRepository struct {
	*base.BaseRepository
}

// DeleteInstructorProfile implements [IInstructorRepository].
func (i *InstructorRepository) DeleteInstructorProfile(ctx context.Context, instructorID uuid.UUID) error {
	if err := i.BaseRepository.DB.Where("user_id = ?", instructorID).Delete(&models.InstructorProfile{}).Error; err != nil {
		return apperrors.ErrDatabase
	}
	return nil
}

// CreateInstructorProfile implements [IInstructorRepository].
func (i *InstructorRepository) CreateInstructorProfile(ctx context.Context, profile *models.InstructorProfile) error {
	return i.BaseRepository.Create(profile)
}

// CreateInstructorProfileTx creates an instructor profile within a transaction
func (i *InstructorRepository) CreateInstructorProfileTx(tx *gorm.DB, profile *models.InstructorProfile) error {
	return i.BaseRepository.CreateTx(tx, profile)
}

// FindInstructorProfileByUserID implements [IInstructorRepository].
func (i *InstructorRepository) FindInstructorProfileByUserID(ctx context.Context, userID uuid.UUID) (*models.InstructorProfile, error) {
	var profile models.InstructorProfile
	if err := i.BaseRepository.FindOne(&profile, "user_id = ?", userID); err != nil {
		return nil, err
	}
	return &profile, nil
}

// UpdateInstructorProfile implements [IInstructorRepository].
func (i *InstructorRepository) UpdateInstructorProfile(ctx context.Context, profile *models.InstructorProfile) error {
	if profile.ID == uuid.Nil && profile.UserID != uuid.Nil {
		var existing models.InstructorProfile
		if err := i.BaseRepository.FindOne(&existing, "user_id = ?", profile.UserID); err == nil {
			profile.ID = existing.ID
		}
	}
	return i.BaseRepository.Update(profile)
}

// UpdateInstructorPhotoURL updates only the photo URL field for an instructor profile
func (i *InstructorRepository) UpdateInstructorPhotoURL(ctx context.Context, userID uuid.UUID, photoURL string) error {
	return i.BaseRepository.UpdateByField(
		&models.InstructorProfile{},
		map[string]interface{}{"photo_url": photoURL},
		"user_id = ?",
		userID,
	)
}

func NewInstructorRepository(db *gorm.DB) IInstructorRepository {
	return &InstructorRepository{BaseRepository: base.NewBaseRepository(db)}
}
