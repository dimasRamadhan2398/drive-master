package services

import (
	"context"
	"user-service/models"
	"user-service/repositories"

	"github.com/google/uuid"
)

type IMemberService interface {
	GetMemberProfile(ctx context.Context, userID uuid.UUID) (*models.MemberProfile, error)
	CreateMemberProfile(ctx context.Context, userID uuid.UUID) (*models.MemberProfile, error)
	UpdateMemberProfile(ctx context.Context, profile *models.MemberProfile) error
	DeleteMemberProfile(ctx context.Context, userID uuid.UUID) error
}

type MemberService struct {
	repo repositories.IMemberRepository
}

func NewMemberService(repo repositories.IMemberRepository) *MemberService {
	return &MemberService{repo: repo}
}

// GetMemberProfile retrieves a member profile by user ID
func (s *MemberService) GetMemberProfile(ctx context.Context, userID uuid.UUID) (*models.MemberProfile, error) {
	return s.repo.FindByUserID(ctx, userID)
}

// CreateMemberProfile creates a new member profile for a user
func (s *MemberService) CreateMemberProfile(ctx context.Context, userID uuid.UUID) (*models.MemberProfile, error) {
	profile := &models.MemberProfile{
		UserID:            userID,
		SessionsCompleted: 0,
		TrainingTime:      0,
		AverageRating:     0,
	}
	if err := s.repo.Create(ctx, profile); err != nil {
		return nil, err
	}
	return profile, nil
}

// UpdateMemberProfile updates a member profile
func (s *MemberService) UpdateMemberProfile(ctx context.Context, profile *models.MemberProfile) error {
	return s.repo.Update(ctx, profile)
}

// DeleteMemberProfile deletes a member profile by user ID
func (s *MemberService) DeleteMemberProfile(ctx context.Context, userID uuid.UUID) error {
	return s.repo.Delete(ctx, userID)
}