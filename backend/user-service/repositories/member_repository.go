package repositories

import (
	"user-service/models"
	"user-service/pkg/base"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IMemberRepository interface {
	// Member profile
	FindByUserID(userID uuid.UUID) (*models.MemberProfile, error)
	FindMany(opts *base.QueryOptions) ([]models.MemberProfile, error)
	Create(profile *models.MemberProfile) error
	Update(profile *models.MemberProfile) error
	Delete(userID uuid.UUID) error
}

type MemberRepository struct {
	*base.BaseRepository
}

// FindByUserID implements [IMemberRepository].
func (m *MemberRepository) FindByUserID(userID uuid.UUID) (*models.MemberProfile, error) {
	var profile models.MemberProfile
	if err := m.BaseRepository.FindOne(&profile, "user_id = ?", userID); err != nil {
		return nil, err
	}
	return &profile, nil
}

// FindMany implements [IMemberRepository].
// Returns a slice of member profiles based on the provided query options.
func (m *MemberRepository) FindMany(opts *base.QueryOptions) ([]models.MemberProfile, error) {
	var profiles []models.MemberProfile
	if err := m.BaseRepository.FindMany(&models.MemberProfile{}, &profiles, opts); err != nil {
		return nil, err
	}
	return profiles, nil
}

// Create implements [IMemberRepository].
func (m *MemberRepository) Create(profile *models.MemberProfile) error {
	return m.BaseRepository.Create(profile)
}

// Update implements [IMemberRepository].
// Subtle: this method shadows the method (*BaseRepository).Update of MemberRepository.BaseRepository.
func (m *MemberRepository) Update(profile *models.MemberProfile) error {
	return m.BaseRepository.Update(profile)
}

// Delete implements [IMemberRepository].
// Deletes a member profile by user ID.
func (m *MemberRepository) Delete(userID uuid.UUID) error {
	return m.BaseRepository.Delete(&models.MemberProfile{UserID: userID})
}

func NewMemberRepository(db *gorm.DB) IMemberRepository {
	return &MemberRepository{BaseRepository: base.NewBaseRepository(db)}
}
