package repositories

import (
	"context"
	"time"

	"booking-service/models"
	"booking-service/models/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EnrollmentRepository handles enrollment database operations
type EnrollmentRepository struct {
	db *gorm.DB
}

func NewEnrollmentRepository(db *gorm.DB) *EnrollmentRepository {
	return &EnrollmentRepository{db: db}
}

func (r *EnrollmentRepository) Create(ctx context.Context, enrollment *models.Enrollment) error {
	return r.db.WithContext(ctx).Create(enrollment).Error
}

func (r *EnrollmentRepository) GetByID(ctx context.Context, id uint) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	if err := r.db.WithContext(ctx).
		Preload("Entitlements").
		First(&enrollment, id).Error; err != nil {
		return nil, err
	}
	return &enrollment, nil
}

func (r *EnrollmentRepository) GetByIDWithPreload(ctx context.Context, id uint) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	if err := r.db.WithContext(ctx).
		Preload("Entitlements").
		First(&enrollment, id).Error; err != nil {
		return nil, err
	}
	return &enrollment, nil
}

func (r *EnrollmentRepository) Update(ctx context.Context, enrollment *models.Enrollment) error {
	return r.db.WithContext(ctx).Save(enrollment).Error
}

func (r *EnrollmentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Enrollment{}, id).Error
}

func (r *EnrollmentRepository) List(ctx context.Context, page, limit int) ([]models.Enrollment, int64, error) {
	var enrollments []models.Enrollment
	var total int64

	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&models.Enrollment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Preload("Entitlements").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&enrollments).Error; err != nil {
		return nil, 0, err
	}

	return enrollments, total, nil
}

func (r *EnrollmentRepository) GetByUserID(ctx context.Context, userID uint, page, limit int) ([]models.Enrollment, int64, error) {
	var enrollments []models.Enrollment
	var total int64

	offset := (page - 1) * limit
	query := r.db.WithContext(ctx).Model(&models.Enrollment{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("Entitlements").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&enrollments).Error; err != nil {
		return nil, 0, err
	}

	return enrollments, total, nil
}

func (r *EnrollmentRepository) GetByStatus(ctx context.Context, status models.EnrollmentStatus, page, limit int) ([]models.Enrollment, int64, error) {
	var enrollments []models.Enrollment
	var total int64

	offset := (page - 1) * limit
	query := r.db.WithContext(ctx).Model(&models.Enrollment{}).Where("status = ?", status)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("Entitlements").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&enrollments).Error; err != nil {
		return nil, 0, err
	}

	return enrollments, total, nil
}

func (r *EnrollmentRepository) UpdateStatus(ctx context.Context, id uint, status models.EnrollmentStatus) error {
	return r.db.WithContext(ctx).
		Model(&models.Enrollment{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *EnrollmentRepository) MarkAsPaid(ctx context.Context, id uint, paidAt time.Time, totalPrice float64) error {
	return r.db.WithContext(ctx).
		Model(&models.Enrollment{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     models.EnrollmentStatusPaid,
			"paid_at":    paidAt,
			"total_price": totalPrice,
		}).Error
}

// ToResponse converts an Enrollment model to EnrollmentResponse DTO
func (r *EnrollmentRepository) ToResponse(enrollment *models.Enrollment) dto.EnrollmentResponse {
	return dto.EnrollmentResponse{
		ID:         enrollment.ID,
		UserID:     enrollment.UserID,
		PackageID:  enrollment.PackageID,
		Status:     string(enrollment.Status),
		TotalPrice: enrollment.TotalPrice,
		PaidAt:     enrollment.PaidAt,
		ExpiresAt:  enrollment.ExpiresAt,
		CreatedAt:  enrollment.CreatedAt,
		UpdatedAt:  enrollment.UpdatedAt,
	}
}

// ToListResponse converts a slice of Enrollments to EnrollmentListResponse DTO
func (r *EnrollmentRepository) ToListResponse(enrollments []models.Enrollment, total int64, page, limit int) dto.EnrollmentListResponse {
	items := make([]dto.EnrollmentResponse, len(enrollments))
	for i, e := range enrollments {
		items[i] = r.ToResponse(&e)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return dto.EnrollmentListResponse{
		Data:       items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}

// AnonymizeByUserID marks all enrollments for a user as anonymized
func (r *EnrollmentRepository) AnonymizeByUserID(ctx context.Context, userID uuid.UUID, anonymizedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&models.Enrollment{}).
		Where("user_id = ?", userID).
		Update("anonymized_at", anonymizedAt).Error
}