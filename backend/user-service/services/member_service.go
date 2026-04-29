package services

import (
	"user-service/models"
	"user-service/repositories"

	"github.com/google/uuid"
)

type IMemberService interface {
	GetMemberProfile(userID uuid.UUID) (*models.MemberProfile, error)
	CreateMemberProfile(userID uuid.UUID) (*models.MemberProfile, error)
	UpdateMemberProfile(profile *models.MemberProfile) error
	DeleteMemberProfile(userID uuid.UUID) error
}

type MemberService struct {
	repo repositories.IMemberRepository
}

func NewMemberService(repo repositories.IMemberRepository) *MemberService {
	return &MemberService{repo: repo}
}

// GetMemberProfile retrieves a member profile by user ID
func (s *MemberService) GetMemberProfile(userID uuid.UUID) (*models.MemberProfile, error) {
	return s.repo.FindByUserID(userID)
}

// CreateMemberProfile creates a new member profile for a user
func (s *MemberService) CreateMemberProfile(userID uuid.UUID) (*models.MemberProfile, error) {
	profile := &models.MemberProfile{
		UserID:            userID,
		SessionsCompleted: 0,
		TrainingTime:      0,
		AverageRating:     0,
	}
	if err := s.repo.Create(profile); err != nil {
		return nil, err
	}
	return profile, nil
}

// UpdateMemberProfile updates a member profile
func (s *MemberService) UpdateMemberProfile(profile *models.MemberProfile) error {
	return s.repo.Update(profile)
}

// DeleteMemberProfile deletes a member profile by user ID
func (s *MemberService) DeleteMemberProfile(userID uuid.UUID) error {
	return s.repo.Delete(userID)
}