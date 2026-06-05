package repositories

import (
	"context"
	"user-service/models"
	"user-service/pkg/base"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ICertificationRepository interface {
	Create(ctx context.Context, cert *models.Certification) error
	Update(ctx context.Context, cert *models.Certification) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Certification, error)
	FindByInstructorID(ctx context.Context, instructorID uuid.UUID, page, limit int) ([]models.Certification, int64, error)
	FindByMemberID(ctx context.Context, memberID uuid.UUID, page, limit int) (*models.Certification, error)
	FindByInstructorAndID(ctx context.Context, instructorID, certID uuid.UUID) (*models.Certification, error)
	CountByInstructorID(ctx context.Context, instructorID uuid.UUID) (int64, error)
}

type CertificationRepository struct {
	*base.BaseRepository
}

func NewCertificationRepository(db *gorm.DB) ICertificationRepository {
	return &CertificationRepository{BaseRepository: base.NewBaseRepository(db)}
}

func (r *CertificationRepository) Create(ctx context.Context, cert *models.Certification) error {
	return r.BaseRepository.Create(cert)
}

func (r *CertificationRepository) Update(ctx context.Context, cert *models.Certification) error {
	return r.BaseRepository.Update(cert)
}

func (r *CertificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.BaseRepository.Delete(&models.Certification{ID: id})
}

func (r *CertificationRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Certification, error) {
	var cert models.Certification
	if err := r.BaseRepository.FindByID(&cert, id); err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *CertificationRepository) FindByMemberID(ctx context.Context, memberID uuid.UUID, page int, limit int) (*models.Certification, error) {
	var cert models.Certification
	offset := (page - 1) * limit

	if err := r.BaseRepository.FindOne(&models.Certification{}, &cert, &base.QueryOptions{
		Where: map[string]interface{}{
			"memberId": memberID,
		},
		Order:  "created_at DESC",
		Limit:  limit,
		Offset: offset,
	}); err != nil {
		return nil, err
	}

	return &cert, nil
}


func (r *CertificationRepository) FindByInstructorID(ctx context.Context, instructorID uuid.UUID, page, limit int) ([]models.Certification, int64, error) {
	var certs []models.Certification
	offset := (page - 1) * limit

	if err := r.BaseRepository.FindWithOptions(&models.Certification{}, &certs, &base.QueryOptions{
		Where: map[string]interface{}{
			"instructorId": instructorID,
		},
		Order:  "created_at DESC",
		Limit:  limit,
		Offset: offset,
	}); err != nil {
		return nil, 0, err
	}

	count, err := r.CountByInstructorID(ctx, instructorID)
	if err != nil {
		return nil, 0, err
	}

	return certs, count, nil
}

func (r *CertificationRepository) FindByInstructorAndID(ctx context.Context, instructorID, certID uuid.UUID) (*models.Certification, error) {
	var cert models.Certification
	if err := r.BaseRepository.FindWithOptions(&models.Certification{}, &cert, &base.QueryOptions{
		Where: map[string]interface{}{
			"id":            certID,
			"instructor_id": instructorID,
		},
		Limit: 1,
	}); err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *CertificationRepository) CountByInstructorID(ctx context.Context, instructorID uuid.UUID) (int64, error) {
	return r.BaseRepository.Count(&models.Certification{}, &base.QueryOptions{
		Where: map[string]interface{}{
			"instructor_id": instructorID,
		},
	})
}
