package repositories

import (
	"context"

	"booking-service/models"
	"booking-service/models/dto"
	"booking-service/pkg/base"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserCertificationRepository interface {
	Create(ctx context.Context, uc *models.UserCertification) error
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.UserCertification, error)
	FindByUserAndCertification(ctx context.Context, userID, certificationID uuid.UUID) (*models.UserCertification, error)
	Delete(ctx context.Context, userID, certificationID uuid.UUID) error
	Exists(ctx context.Context, userID, certificationID uuid.UUID) (bool, error)
	ToResponse(uc *models.UserCertification, certResponse dto.CertificationResponse) dto.UserCertificationResponse
}

type UserCertificationRepository struct {
	*base.BaseRepository
	db *gorm.DB
}

func NewUserCertificationRepository(db *gorm.DB) IUserCertificationRepository {
	return &UserCertificationRepository{BaseRepository: base.NewBaseRepository(db), db: db}
}

func (r *UserCertificationRepository) Create(ctx context.Context, uc *models.UserCertification) error {
	return r.BaseRepository.Create(uc)
}

func (r *UserCertificationRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.UserCertification, error) {
	var userCerts []models.UserCertification
	opts := base.NewQueryOptions().
		WithPreloads("Certification").
		WithWhere(map[string]any{"user_id": userID})
	if err := r.BaseRepository.FindMany(&models.UserCertification{}, &userCerts, opts); err != nil {
		return nil, err
	}
	return userCerts, nil
}

func (r *UserCertificationRepository) FindByUserAndCertification(ctx context.Context, userID, certificationID uuid.UUID) (*models.UserCertification, error) {
	var uc models.UserCertification
	opts := base.NewQueryOptions().
		WithPreloads("Certification").
		WithWhere(map[string]any{"user_id": userID, "certification_id": certificationID})
	if err := r.BaseRepository.FindMany(&models.UserCertification{}, &uc, opts); err != nil {
		return nil, err
	}
	return &uc, nil
}

func (r *UserCertificationRepository) Delete(ctx context.Context, userID, certificationID uuid.UUID) error {
	return r.BaseRepository.Exec(
		"DELETE FROM user_certifications WHERE user_id = ? AND certification_id = ?",
		userID, certificationID,
	)
}

func (r *UserCertificationRepository) Exists(ctx context.Context, userID, certificationID uuid.UUID) (bool, error) {
	return r.BaseRepository.Exists(&models.UserCertification{}, "user_id = ? AND certification_id = ?", userID, certificationID)
}

// ToResponse converts a UserCertification model to UserCertificationResponse DTO
func (r *UserCertificationRepository) ToResponse(uc *models.UserCertification, certResponse dto.CertificationResponse) dto.UserCertificationResponse {
	return dto.UserCertificationResponse{
		UserID:          uc.UserID,
		CertificationID: uc.CertificationID,
		IssuedAt:        uc.IssuedAt,
		Certification:   certResponse,
	}
}