package repositories

import (
	"context"
	"strconv"
	"time"

	"booking-service/models"
	"booking-service/models/dto"
	"booking-service/pkg/base"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IEntitlementRepository interface {
	Create(ctx context.Context, entitlement *models.UserEntitlement) error
	CreateTx(tx *gorm.DB, entitlement *models.UserEntitlement) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.UserEntitlement, error)
	Update(ctx context.Context, entitlement *models.UserEntitlement) error
	Delete(ctx context.Context, entitlement *models.UserEntitlement) error
	FindAll(ctx context.Context) ([]models.UserEntitlement, error)
	FindByEnrollmentID(ctx context.Context, enrollmentID uuid.UUID) ([]models.UserEntitlement, error)
	FindByUserID(ctx context.Context, userID uint) ([]models.UserEntitlement, error)
	FindBySourceType(ctx context.Context, sourceType string) ([]models.UserEntitlement, error)
	FindActiveByUserID(ctx context.Context, userID uint) ([]models.UserEntitlement, error)
	UpdateUsedSessions(ctx context.Context, id uuid.UUID, usedSessions int) error
	AnonymizeByUserID(ctx context.Context, userID uuid.UUID, anonymizedAt time.Time) error
	CountAll(ctx context.Context) (int64, error)
	CountByEnrollmentID(ctx context.Context, enrollmentID uuid.UUID) (int64, error)
	Exists(ctx context.Context, condition any, args ...any) (bool, error)
	ToResponse(entitlement *models.UserEntitlement) dto.EntitlementResponse
	ToListResponse(entitlements []models.UserEntitlement, total int64, page, limit int) dto.EntitlementListResponse
}

type EntitlementRepository struct {
	*base.BaseRepository
	db *gorm.DB
}

func NewEntitlementRepository(db *gorm.DB) IEntitlementRepository {
	return &EntitlementRepository{BaseRepository: base.NewBaseRepository(db), db: db}
}

func (r *EntitlementRepository) Create(ctx context.Context, entitlement *models.UserEntitlement) error {
	return r.BaseRepository.Create(entitlement)
}

func (r *EntitlementRepository) CreateTx(tx *gorm.DB, entitlement *models.UserEntitlement) error {
	return r.BaseRepository.CreateTx(tx, entitlement)
}

func (r *EntitlementRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.UserEntitlement, error) {
	var entitlement models.UserEntitlement
	if err := r.BaseRepository.FindByIDWithPreload(&entitlement, id); err != nil {
		return nil, err
	}
	return &entitlement, nil
}

func (r *EntitlementRepository) Update(ctx context.Context, entitlement *models.UserEntitlement) error {
	return r.BaseRepository.Update(entitlement)
}

func (r *EntitlementRepository) Delete(ctx context.Context, entitlement *models.UserEntitlement) error {
	return r.BaseRepository.Delete(entitlement)
}

func (r *EntitlementRepository) FindAll(ctx context.Context) ([]models.UserEntitlement, error) {
	var entitlements []models.UserEntitlement
	opts := base.NewQueryOptions().WithOrder("created_at DESC")
	if err := r.BaseRepository.FindMany(&models.UserEntitlement{}, &entitlements, opts); err != nil {
		return nil, err
	}
	return entitlements, nil
}

func (r *EntitlementRepository) FindByEnrollmentID(ctx context.Context, enrollmentID uuid.UUID) ([]models.UserEntitlement, error) {
	var entitlements []models.UserEntitlement
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"enrollment_id": enrollmentID}).
		WithOrder("created_at DESC")
	if err := r.BaseRepository.FindMany(&models.UserEntitlement{}, &entitlements, opts); err != nil {
		return nil, err
	}
	return entitlements, nil
}

func (r *EntitlementRepository) FindByUserID(ctx context.Context, userID uint) ([]models.UserEntitlement, error) {
	var entitlements []models.UserEntitlement
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"user_id": userID}).
		WithOrder("created_at DESC")
	if err := r.BaseRepository.FindMany(&models.UserEntitlement{}, &entitlements, opts); err != nil {
		return nil, err
	}
	return entitlements, nil
}

func (r *EntitlementRepository) FindBySourceType(ctx context.Context, sourceType string) ([]models.UserEntitlement, error) {
	var entitlements []models.UserEntitlement
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"source_type": sourceType}).
		WithOrder("created_at DESC")
	if err := r.BaseRepository.FindMany(&models.UserEntitlement{}, &entitlements, opts); err != nil {
		return nil, err
	}
	return entitlements, nil
}

func (r *EntitlementRepository) FindActiveByUserID(ctx context.Context, userID uint) ([]models.UserEntitlement, error) {
	var entitlements []models.UserEntitlement
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{
			"user_id":        userID,
			"used_sessions < total_sessions": nil,
		}).
		WithOrder("created_at DESC")
	if err := r.BaseRepository.FindMany(&models.UserEntitlement{}, &entitlements, opts); err != nil {
		return nil, err
	}
	return entitlements, nil
}

func (r *EntitlementRepository) UpdateUsedSessions(ctx context.Context, id uuid.UUID, usedSessions int) error {
	return r.BaseRepository.Exec(
		"UPDATE user_entitlements SET used_sessions = ?, updated_at = ? WHERE id = ?",
		usedSessions, time.Now(), id,
	)
}

func (r *EntitlementRepository) AnonymizeByUserID(ctx context.Context, userID uuid.UUID, anonymizedAt time.Time) error {
	// Convert UUID to uint (this assumes the UUID can be parsed as a number)
	userIDUint, err := strconv.ParseUint(userID.String()[0:8], 16, 64)
	if err != nil {
		userIDUint = 0
	}

	return r.BaseRepository.Exec(
		"UPDATE user_entitlements SET anonymized_at = ? WHERE user_id = ?",
		anonymizedAt, userIDUint,
	)
}

func (r *EntitlementRepository) CountAll(ctx context.Context) (int64, error) {
	return r.BaseRepository.Count(&models.UserEntitlement{}, base.NewQueryOptions())
}

func (r *EntitlementRepository) CountByEnrollmentID(ctx context.Context, enrollmentID uuid.UUID) (int64, error) {
	opts := base.NewQueryOptions().WithWhere(map[string]any{"enrollment_id": enrollmentID})
	return r.BaseRepository.Count(&models.UserEntitlement{}, opts)
}

func (r *EntitlementRepository) Exists(ctx context.Context, condition any, args ...any) (bool, error) {
	return r.BaseRepository.Exists(&models.UserEntitlement{}, condition, args...)
}

// ToResponse converts a UserEntitlement model to EntitlementResponse DTO
func (r *EntitlementRepository) ToResponse(entitlement *models.UserEntitlement) dto.EntitlementResponse {
	return dto.EntitlementResponse{
		ID:                entitlement.ID,
		UserID:            entitlement.UserID,
		SourceType:        entitlement.SourceType,
		SourceID:          entitlement.SourceID,
		TotalSessions:     entitlement.TotalSessions,
		SessionsRemaining: entitlement.TotalSessions - entitlement.UsedSessions,
		ExpiresAt:         entitlement.ExpiresAt,
		CreatedAt:         entitlement.CreatedAt,
		UpdatedAt:         entitlement.UpdatedAt,
	}
}

// ToListResponse converts a slice of UserEntitlements to EntitlementListResponse DTO
func (r *EntitlementRepository) ToListResponse(entitlements []models.UserEntitlement, total int64, page, limit int) dto.EntitlementListResponse {
	items := make([]dto.EntitlementResponse, len(entitlements))
	for i, e := range entitlements {
		items[i] = r.ToResponse(&e)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return dto.EntitlementListResponse{
		Data:       items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}