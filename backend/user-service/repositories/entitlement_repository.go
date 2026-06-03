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
	FindActiveByMemberIDs(ctx context.Context, memberIDs []uuid.UUID) (map[uuid.UUID]int, error) // Returns total remaining sessions per memberID
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

// FindActiveByMemberIDs calculates total remaining sessions from active entitlements for multiple members
// Returns a map of memberID -> total remaining sessions
func (r *EntitlementRepository) FindActiveByMemberIDs(ctx context.Context, memberIDs []uuid.UUID) (map[uuid.UUID]int, error) {
	result := make(map[uuid.UUID]int)

	if len(memberIDs) == 0 {
		return result, nil
	}

	// Build the query with IN clause
	// Use raw query to properly handle IN clause with UUIDs
	query := `
		SELECT member_id, SUM(remaining) as total_remaining
		FROM entitlements
		WHERE member_id = ANY($1) AND status = 'active'
		GROUP BY member_id
	`

	rows, err := r.DB.Raw(query, memberIDs).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var memberID uuid.UUID
		var totalRemaining int
		if err := rows.Scan(&memberID, &totalRemaining); err != nil {
			return nil, err
		}
		result[memberID] = totalRemaining
	}

	return result, nil
}