package repositories

import (
	"context"
	"time"

	"booking-service/models"
	"booking-service/models/dto"
	"booking-service/pkg/base"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IEnrollmentRepository interface {
	Create(ctx context.Context, enrollment *models.Enrollment) error
	CreateTx(tx *gorm.DB, enrollment *models.Enrollment) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Enrollment, error)
	FindByIDWithPreload(ctx context.Context, id uuid.UUID) (*models.Enrollment, error)
	Update(ctx context.Context, enrollment *models.Enrollment) error
	UpdateTx(tx *gorm.DB, enrollment *models.Enrollment) error
	Delete(ctx context.Context, enrollment *models.Enrollment) error
	FindAll(ctx context.Context) ([]models.Enrollment, error)
	FindAllPaginated(ctx context.Context, page, limit int) ([]models.Enrollment, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Enrollment, error)
	FindByUserIDPaginated(ctx context.Context, userID uuid.UUID, page, limit int) ([]models.Enrollment, error)
	FindByStatus(ctx context.Context, status models.EnrollmentStatus) ([]models.Enrollment, int64, error)
	FindByPackageID(ctx context.Context, packageID uuid.UUID) ([]models.Enrollment, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.EnrollmentStatus) error
	UpdateStatusTx(tx *gorm.DB, id uuid.UUID, status models.EnrollmentStatus) error
	MarkAsPaid(ctx context.Context, id uuid.UUID, paidAt time.Time, totalPrice float64) error
	AnonymizeByUserID(ctx context.Context, userID uuid.UUID, anonymizedAt time.Time) error
	CountAll(ctx context.Context) (int64, error)
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	CountByStatus(ctx context.Context, status models.EnrollmentStatus) (int64, error)
	Exists(ctx context.Context, condition any, args ...any) (bool, error)
	ToResponse(enrollment *models.Enrollment) dto.EnrollmentResponse
	ToListResponse(enrollments []models.Enrollment, total int64, page, limit int) dto.EnrollmentListResponse
}

type EnrollmentRepository struct {
	*base.BaseRepository
	db *gorm.DB
}

func NewEnrollmentRepository(db *gorm.DB) IEnrollmentRepository {
	return &EnrollmentRepository{BaseRepository: base.NewBaseRepository(db), db: db}
}

func (r *EnrollmentRepository) Create(ctx context.Context, enrollment *models.Enrollment) error {
	return r.BaseRepository.Create(enrollment)
}

func (r *EnrollmentRepository) CreateTx(tx *gorm.DB, enrollment *models.Enrollment) error {
	return r.BaseRepository.CreateTx(tx, enrollment)
}

func (r *EnrollmentRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	if err := r.BaseRepository.FindByIDWithPreload(&enrollment, id); err != nil {
		return nil, err
	}
	return &enrollment, nil
}

func (r *EnrollmentRepository) FindByIDWithPreload(ctx context.Context, id uuid.UUID) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	opts := base.NewQueryOptions().WithPreloads("Entitlements")
	if err := r.BaseRepository.FindMany(&models.Enrollment{}, &enrollment, opts); err != nil {
		return nil, err
	}
	return &enrollment, nil
}

func (r *EnrollmentRepository) Update(ctx context.Context, enrollment *models.Enrollment) error {
	return r.BaseRepository.Update(enrollment)
}

func (r *EnrollmentRepository) UpdateTx(tx *gorm.DB, enrollment *models.Enrollment) error {
	return tx.Save(enrollment).Error
}

func (r *EnrollmentRepository) Delete(ctx context.Context, enrollment *models.Enrollment) error {
	return r.BaseRepository.Delete(enrollment)
}

func (r *EnrollmentRepository) FindAll(ctx context.Context) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	opts := base.NewQueryOptions().
		WithPreloads("Entitlements").
		WithOrder("created_at DESC")
	if err := r.BaseRepository.FindMany(&models.Enrollment{}, &enrollments, opts); err != nil {
		return nil, err
	}
	return enrollments, nil
}

func (r *EnrollmentRepository) FindAllPaginated(ctx context.Context, page, limit int) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	offset := (page - 1) * limit
	opts := base.NewQueryOptions().
		WithPreloads("Entitlements").
		WithOrder("created_at DESC").
		WithPagination(offset, limit)
	if err := r.BaseRepository.FindMany(&models.Enrollment{}, &enrollments, opts); err != nil {
		return nil, err
	}
	return enrollments, nil
}

func (r *EnrollmentRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"user_id": userID}).
		WithPreloads("Entitlements").
		WithOrder("created_at DESC")
	if err := r.BaseRepository.FindMany(&models.Enrollment{}, &enrollments, opts); err != nil {
		return nil, err
	}
	return enrollments, nil
}

func (r *EnrollmentRepository) FindByUserIDPaginated(ctx context.Context, userID uuid.UUID, page, limit int) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	offset := (page - 1) * limit
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"user_id": userID}).
		WithPreloads("Entitlements").
		WithOrder("created_at DESC").
		WithPagination(offset, limit)
	if err := r.BaseRepository.FindMany(&models.Enrollment{}, &enrollments, opts); err != nil {
		return nil, err
	}
	return enrollments, nil
}

func (r *EnrollmentRepository) FindByStatus(ctx context.Context, status models.EnrollmentStatus) ([]models.Enrollment, int64, error) {
	var enrollments []models.Enrollment
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"status": status}).
		WithPreloads("Entitlements").
		WithOrder("created_at DESC")
	if err := r.BaseRepository.FindMany(&models.Enrollment{}, &enrollments, opts); err != nil {
		return nil, 0, err
	}
	return enrollments, int64(len(enrollments)), nil
}

func (r *EnrollmentRepository) FindByPackageID(ctx context.Context, packageID uuid.UUID) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"package_id": packageID}).
		WithPreloads("Entitlements").
		WithOrder("created_at DESC")
	if err := r.BaseRepository.FindMany(&models.Enrollment{}, &enrollments, opts); err != nil {
		return nil, err
	}
	return enrollments, nil
}

func (r *EnrollmentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.EnrollmentStatus) error {
	return r.BaseRepository.Exec(
		"UPDATE enrollments SET status = ?, updated_at = ? WHERE id = ?",
		status, time.Now(), id,
	)
}

func (r *EnrollmentRepository) UpdateStatusTx(tx *gorm.DB, id uuid.UUID, status models.EnrollmentStatus) error {
	return tx.Exec(
		"UPDATE enrollments SET status = ?, updated_at = ? WHERE id = ?",
		status, time.Now(), id,
	).Error
}

func (r *EnrollmentRepository) MarkAsPaid(ctx context.Context, id uuid.UUID, paidAt time.Time, totalPrice float64) error {
	return r.BaseRepository.Exec(
		"UPDATE enrollments SET status = 'paid', paid_at = ?, total_price = ?, updated_at = ? WHERE id = ?",
		paidAt, totalPrice, time.Now(), id,
	)
}

func (r *EnrollmentRepository) AnonymizeByUserID(ctx context.Context, userID uuid.UUID, anonymizedAt time.Time) error {
	return r.BaseRepository.Exec(
		"UPDATE enrollments SET anonymized_at = ? WHERE user_id = ?",
		anonymizedAt, userID,
	)
}

func (r *EnrollmentRepository) CountAll(ctx context.Context) (int64, error) {
	return r.BaseRepository.Count(&models.Enrollment{}, base.NewQueryOptions())
}

func (r *EnrollmentRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	opts := base.NewQueryOptions().WithWhere(map[string]any{"user_id": userID})
	return r.BaseRepository.Count(&models.Enrollment{}, opts)
}

func (r *EnrollmentRepository) CountByStatus(ctx context.Context, status models.EnrollmentStatus) (int64, error) {
	opts := base.NewQueryOptions().WithWhere(map[string]any{"status": status})
	return r.BaseRepository.Count(&models.Enrollment{}, opts)
}

func (r *EnrollmentRepository) Exists(ctx context.Context, condition any, args ...any) (bool, error) {
	return r.BaseRepository.Exists(&models.Enrollment{}, condition, args...)
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
		Pagination: dto.NewPaginationMeta(total, page, limit),
	}
}