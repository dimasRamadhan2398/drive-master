package repositories

import (
	"context"
	"errors"
	"time"
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
	FindByMemberID(ctx context.Context, memberID uuid.UUID, page, limit int) ([]models.Certification, int64, error)
	FindAllByMemberID(ctx context.Context, memberID uuid.UUID) ([]models.Certification, error)
	FindAll(ctx context.Context) ([]models.Certification, error)
	FindByMemberIDAndCertID(ctx context.Context, memberID, certID uuid.UUID) (*models.Certification, error)
	FindByInstructorAndID(ctx context.Context, instructorID, certID uuid.UUID) (*models.Certification, error)
	// FindByEntitlementID returns the certificate linked to the given entitlement, or nil if none exists yet.
	FindByEntitlementID(ctx context.Context, entitlementID uuid.UUID) (*models.Certification, error)
	CountByInstructorID(ctx context.Context, instructorID uuid.UUID) (int64, error)
	CountByMemberID(ctx context.Context, memberID uuid.UUID) (int64, error)
	CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error)
	GetStats(ctx context.Context) (*CertificateStats, error)
	GetStatsByDateRange(ctx context.Context, startDate, endDate time.Time) (*CertificateStats, error)
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

func (r *CertificationRepository) FindByMemberID(ctx context.Context, memberID uuid.UUID, page int, limit int) ([]models.Certification, int64, error) {
	var certs []models.Certification
	offset := (page - 1) * limit

	if err := r.BaseRepository.FindWithOptions(&models.Certification{}, &certs, &base.QueryOptions{
		Where: map[string]interface{}{
			"member_id": memberID,
		},
		Order:  "created_at DESC",
		Limit:  limit,
		Offset: offset,
	}); err != nil {
		return nil, 0, err
	}

	count, err := r.CountByMemberID(ctx, memberID)
	if err != nil {
		return nil, 0, err
	}

	return certs, count, nil
}

func (r *CertificationRepository) FindAllByMemberID(ctx context.Context, memberID uuid.UUID) ([]models.Certification, error) {
	var certs []models.Certification

	if err := r.BaseRepository.FindWithOptions(&models.Certification{}, &certs, &base.QueryOptions{
		Where: map[string]interface{}{
			"member_id": memberID,
		},
		Order: "created_at DESC",
	}); err != nil {
		return nil, err
	}

	return certs, nil
}

func (r *CertificationRepository) FindAll(ctx context.Context) ([]models.Certification, error) {
	var certs []models.Certification

	if err := r.BaseRepository.FindWithOptions(&models.Certification{}, &certs, &base.QueryOptions{
		Order: "created_at DESC",
	}); err != nil {
		return nil, err
	}

	return certs, nil
}

func (r *CertificationRepository) FindByMemberIDAndCertID(ctx context.Context, memberID, certID uuid.UUID) (*models.Certification, error) {
	var cert models.Certification
	if err := r.BaseRepository.FindWithOptions(&models.Certification{}, &cert, &base.QueryOptions{
		Where: map[string]interface{}{
			"id":        certID,
			"member_id": memberID,
		},
		Limit: 1,
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

// FindByEntitlementID returns the certificate linked to the given entitlement ID.
// Returns nil, nil when no certificate has been issued yet for that entitlement.
func (r *CertificationRepository) FindByEntitlementID(ctx context.Context, entitlementID uuid.UUID) (*models.Certification, error) {
	var cert models.Certification
	err := r.DB.Where("entitlement_id = ?", entitlementID).First(&cert).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
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

func (r *CertificationRepository) CountByMemberID(ctx context.Context, memberID uuid.UUID) (int64, error) {
	return r.BaseRepository.Count(&models.Certification{}, &base.QueryOptions{
		Where: map[string]interface{}{
			"member_id": memberID,
		},
	})
}

// CertificateStats holds certificate statistics
type CertificateStats struct {
	Total      int64
	Verified   int64
	Pending    int64
	Expired    int64
	Revoked    int64
}

func (r *CertificationRepository) CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.DB.Model(&models.Certification{}).
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Count(&count).Error
	return count, err
}

func (r *CertificationRepository) GetStats(ctx context.Context) (*CertificateStats, error) {
	stats := &CertificateStats{}

	// Get total certifications
	if err := r.DB.Model(&models.Certification{}).Count(&stats.Total).Error; err != nil {
		return nil, err
	}

	// Get verified certifications
	if err := r.DB.Model(&models.Certification{}).
		Where("status = ?", models.CertificationStatusVerified).
		Count(&stats.Verified).Error; err != nil {
		return nil, err
	}

	// Get pending certifications
	if err := r.DB.Model(&models.Certification{}).
		Where("status = ?", models.CertificationStatusPending).
		Count(&stats.Pending).Error; err != nil {
		return nil, err
	}

	// Get expired certifications
	if err := r.DB.Model(&models.Certification{}).
		Where("status = ?", models.CertificationStatusExpired).
		Count(&stats.Expired).Error; err != nil {
		return nil, err
	}

	// Get revoked certifications
	if err := r.DB.Model(&models.Certification{}).
		Where("status = ?", models.CertificationStatusRevoked).
		Count(&stats.Revoked).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *CertificationRepository) GetStatsByDateRange(ctx context.Context, startDate, endDate time.Time) (*CertificateStats, error) {
	stats := &CertificateStats{}

	// Get total certifications in date range
	if err := r.DB.Model(&models.Certification{}).
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Count(&stats.Total).Error; err != nil {
		return nil, err
	}

	// Get verified certifications in date range
	if err := r.DB.Model(&models.Certification{}).
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Where("status = ?", models.CertificationStatusVerified).
		Count(&stats.Verified).Error; err != nil {
		return nil, err
	}

	// Get pending certifications in date range
	if err := r.DB.Model(&models.Certification{}).
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Where("status = ?", models.CertificationStatusPending).
		Count(&stats.Pending).Error; err != nil {
		return nil, err
	}

	// Get expired certifications in date range
	if err := r.DB.Model(&models.Certification{}).
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Where("status = ?", models.CertificationStatusExpired).
		Count(&stats.Expired).Error; err != nil {
		return nil, err
	}

	// Get revoked certifications in date range
	if err := r.DB.Model(&models.Certification{}).
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Where("status = ?", models.CertificationStatusRevoked).
		Count(&stats.Revoked).Error; err != nil {
		return nil, err
	}

	return stats, nil
}