package repositories

import (
	"context"
	"user-service/models"
	"user-service/pkg/base"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IEntitlementRepository interface {
	Create(ctx context.Context, entitlement *models.Entitlement) error
	Update(ctx context.Context, entitlement *models.Entitlement) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Entitlement, error)
	FindByMemberID(ctx context.Context, memberID uuid.UUID, page, limit int) ([]models.Entitlement, int64, error)
	FindByMemberAndID(ctx context.Context, memberID, entitlementID uuid.UUID) (*models.Entitlement, error)
	FindByBookingID(ctx context.Context, bookingID uuid.UUID) (*models.Entitlement, error)
	FindActiveByMemberIDs(ctx context.Context, memberIDs []uuid.UUID) (map[uuid.UUID][]models.Entitlement, error)
	CountByMemberID(ctx context.Context, memberID uuid.UUID) (int64, error)
	DecrementRemaining(ctx context.Context, id uuid.UUID) error
}

type EntitlementRepository struct {
	*base.BaseRepository
}

func NewEntitlementRepository(db *gorm.DB) IEntitlementRepository {
	return &EntitlementRepository{BaseRepository: base.NewBaseRepository(db)}
}

func (r *EntitlementRepository) Create(ctx context.Context, entitlement *models.Entitlement) error {
	return r.BaseRepository.Create(entitlement)
}

func (r *EntitlementRepository) Update(ctx context.Context, entitlement *models.Entitlement) error {
	return r.BaseRepository.Update(entitlement)
}

func (r *EntitlementRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.BaseRepository.Delete(&models.Entitlement{ID: id})
}

func (r *EntitlementRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Entitlement, error) {
	var entitlement models.Entitlement
	if err := r.BaseRepository.FindByID(&entitlement, id); err != nil {
		return nil, err
	}
	return &entitlement, nil
}

func (r *EntitlementRepository) FindByMemberID(ctx context.Context, memberID uuid.UUID, page, limit int) ([]models.Entitlement, int64, error) {
	var entitlements []models.Entitlement
	offset := (page - 1) * limit

	if err := r.BaseRepository.FindWithOptions(&models.Entitlement{}, &entitlements, &base.QueryOptions{
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

	return entitlements, count, nil
}

func (r *EntitlementRepository) FindByMemberAndID(ctx context.Context, memberID, entitlementID uuid.UUID) (*models.Entitlement, error) {
	var entitlement models.Entitlement
	if err := r.BaseRepository.FindWithOptions(&models.Entitlement{}, &entitlement, &base.QueryOptions{
		Where: map[string]interface{}{
			"id":        entitlementID,
			"member_id": memberID,
		},
		Limit: 1,
	}); err != nil {
		return nil, err
	}
	return &entitlement, nil
}

func (r *EntitlementRepository) FindByBookingID(ctx context.Context, bookingID uuid.UUID) (*models.Entitlement, error) {
	var entitlement models.Entitlement
	if err := r.BaseRepository.FindWithOptions(&models.Entitlement{}, &entitlement, &base.QueryOptions{
		Where: map[string]interface{}{
			"booking_id": bookingID,
		},
		Limit: 1,
	}); err != nil {
		return nil, err
	}
	return &entitlement, nil
}

func (r *EntitlementRepository) CountByMemberID(ctx context.Context, memberID uuid.UUID) (int64, error) {
	return r.BaseRepository.Count(&models.Entitlement{}, &base.QueryOptions{
		Where: map[string]interface{}{
			"member_id": memberID,
		},
	})
}

func (r *EntitlementRepository) DecrementRemaining(ctx context.Context, id uuid.UUID) error {
	// Use raw SQL for atomic decrement
	return r.BaseRepository.Exec("UPDATE entitlements SET remaining = remaining - 1, used_sessions = used_sessions + 1 WHERE id = ? AND remaining > 0", id)
}


// FindActiveByMemberIDs returns active entitlements for multiple members
func (r *EntitlementRepository) FindActiveByMemberIDs(ctx context.Context, memberIDs []uuid.UUID) (map[uuid.UUID][]models.Entitlement, error) {
	result := make(map[uuid.UUID][]models.Entitlement)

	if len(memberIDs) == 0 {
		return result, nil
	}

	// Initialize empty slices for all member IDs
	for _, id := range memberIDs {
		result[id] = []models.Entitlement{}
	}

	var entitlements []models.Entitlement
	if err := r.BaseRepository.FindWithOptions(&models.Entitlement{}, &entitlements, &base.QueryOptions{
		Where: map[string]interface{}{
			"member_id IN": memberIDs,
			"status":       models.EntitlementStatusActive,
		},
	}); err != nil {
		return nil, err
	}

	// Group entitlements by member ID
	for _, ent := range entitlements {
		result[ent.MemberID] = append(result[ent.MemberID], ent)
	}

	return result, nil
}