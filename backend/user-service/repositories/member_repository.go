package repositories

import (
	"context"
	"strings"
	"user-service/models"
	"user-service/pkg/base"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IMemberRepository interface {
	// Member profile
	FindByUserID(ctx context.Context, userID uuid.UUID) (*models.MemberProfile, error)
	FindMany(ctx context.Context, opts *base.QueryOptions) ([]models.MemberProfile, error)
	Create(ctx context.Context, profile *models.MemberProfile) error
	Update(ctx context.Context, profile *models.MemberProfile) error
	Delete(ctx context.Context, userID uuid.UUID) error
	// Pagination support
	FindByRoleIDWithPagination(ctx context.Context, roleID uint, offset, limit int, query string) ([]models.User, error)
	CountByRoleID(ctx context.Context, roleID uint) (int64, error)	
	// Update total available sessions
	UpdateTotalAvailableSessions(ctx context.Context, userID uuid.UUID, total int) error
}

type MemberRepository struct {
	*base.BaseRepository
	db *gorm.DB
}

// FindByUserID implements [IMemberRepository].
func (m *MemberRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*models.MemberProfile, error) {
	var profile models.MemberProfile
	if err := m.BaseRepository.FindOne(&profile, "user_id = ?", userID); err != nil {
		return nil, err
	}
	return &profile, nil
}

// FindMany implements [IMemberRepository].
// Returns a slice of member profiles based on the provided query options.
func (m *MemberRepository) FindMany(ctx context.Context, opts *base.QueryOptions) ([]models.MemberProfile, error) {
	var profiles []models.MemberProfile
	if err := m.BaseRepository.FindMany(&models.MemberProfile{}, &profiles, opts); err != nil {
		return nil, err
	}
	return profiles, nil
}

// Create implements [IMemberRepository].
func (m *MemberRepository) Create(ctx context.Context, profile *models.MemberProfile) error {
	return m.BaseRepository.Create(profile)
}

// Update implements [IMemberRepository].
// Subtle: this method shadows the method (*BaseRepository).Update of MemberRepository.BaseRepository.
func (m *MemberRepository) Update(ctx context.Context, profile *models.MemberProfile) error {
	return m.BaseRepository.Update(profile)
}

// Delete implements [IMemberRepository].
// Deletes a member profile by user ID.
func (m *MemberRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	return m.BaseRepository.Delete(&models.MemberProfile{UserID: userID})
}

// CountByRoleID counts the number of users with a specific role ID
func (m *MemberRepository) CountByRoleID(ctx context.Context, roleID uint) (int64, error) {
	var count int64
	if err := m.DB.Model(&models.User{}).Where("role_id = ?", roleID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// FindByRoleIDWithPagination retrieves users with a specific role ID with pagination
func (m *MemberRepository) FindByRoleIDWithPagination(ctx context.Context, roleID uint, offset, limit int, query string) ([]models.User, error) {
	var users []models.User

	// Convert query to lowercase for case-insensitive search
	lowerQuery := "%" + strings.ToLower(query) + "%"

	// Use raw SQL with LOWER() for case-insensitive search
	sqlQuery := `
		SELECT id, first_name, last_name, username, password_hash,
		       email_address, phone_number, image, date_of_birth,
		       address, is_active, is_verified, role_id,
		       created_at, updated_at
		FROM users
		WHERE role_id = ? AND (LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ? OR LOWER(username) LIKE ?)
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	if err := m.DB.Raw(sqlQuery, roleID, lowerQuery, lowerQuery, lowerQuery, limit, offset).Scan(&users).Error; err != nil {
		return nil, err
	}

	// Preload Role and MemberProfile separately without overwriting user fields
	for i := range users {
		// Only load the relation data, don't reload the whole user
		var role models.Role
		if err := m.DB.First(&role, "id = ?", users[i].RoleID).Error; err == nil {
			users[i].Role = role
		}
		var memberProfile models.MemberProfile
		if err := m.DB.First(&memberProfile, "user_id = ?", users[i].ID).Error; err == nil {
			users[i].MemberProfile = &memberProfile
		}
	}

	return users, nil
}

// UpdateTotalAvailableSessions updates the total available sessions for a member profile
func (m *MemberRepository) UpdateTotalAvailableSessions(ctx context.Context, userID uuid.UUID, total int) error {
	return m.BaseRepository.Exec("UPDATE member_profiles SET total_available_sessions = ?, updated_at = NOW() WHERE user_id = ?", total, userID)
}

func NewMemberRepository(db *gorm.DB) IMemberRepository {
	return &MemberRepository{BaseRepository: base.NewBaseRepository(db), db: db}
}