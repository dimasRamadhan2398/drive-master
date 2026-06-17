package services

import (
	"context"
	"user-service/models"
	"user-service/models/dto"
	"user-service/repositories"

	"github.com/google/uuid"
)

type IMemberService interface {
	GetMemberProfile(ctx context.Context, userID uuid.UUID) (*models.MemberProfile, error)
	CreateMemberProfile(ctx context.Context, userID uuid.UUID) (*models.MemberProfile, error)
	UpdateMemberProfile(ctx context.Context, profile *models.MemberProfile) error
	DeleteMemberProfile(ctx context.Context, userID uuid.UUID) error
	GetMembersWithPagination(ctx context.Context, page, limit int, searchQuery *string) (*dto.MemberListResponse, error)
}

type MemberService struct {
	repo           repositories.IMemberRepository
	roleRepo       repositories.IRoleRepository
	entitlementSvc IEntitlementService
}

func NewMemberService(repo repositories.IMemberRepository, roleRepo repositories.IRoleRepository, entitlementSvc IEntitlementService) IMemberService {
	return &MemberService{
		repo:           repo,
		roleRepo:       roleRepo,
		entitlementSvc: entitlementSvc,
	}
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

// GetMembersWithPagination returns paginated list of members with their profiles
// Calculates totalAvailableSessions dynamically from active entitlements
func (s *MemberService) GetMembersWithPagination(ctx context.Context, page, limit int, searchQuery *string) (*dto.MemberListResponse, error) {
	// Find the "member" role
	roleModel, err := s.roleRepo.FindRoleByName(ctx, "member")
	if err != nil {
		return nil, err
	}

	// Get total count
	total, err := s.repo.CountByRoleID(ctx, roleModel.ID)
	if err != nil {
		return nil, err
	}

	// Calculate pagination offset
	offset := (page - 1) * limit

	// Determine search query - use empty string if nil
	query := ""
	if searchQuery != nil {
		query = *searchQuery
	}

	// Get paginated users
	users, err := s.repo.FindByRoleIDWithPagination(ctx, roleModel.ID, offset, limit, query)
	if err != nil {
		return nil, err
	}

	// Collect member IDs for bulk entitlements query
	memberIDs := make([]uuid.UUID, 0, len(users))
	for _, user := range users {
		if user.MemberProfile != nil {
			memberIDs = append(memberIDs, user.ID)
		}
	}

	// Get active entitlements for all members
	entitlements, err := s.entitlementSvc.FindSessionsStatsByMemberIDs(ctx, memberIDs)
	if err != nil {
		// Log but don't fail - just set available sessions to 0
		entitlements = make(map[uuid.UUID][]models.Entitlement)
	}

	// Convert to response DTOs
	data := make([]dto.UserWithProfileResponse, len(users))
	for i, user := range users {
		var memberProfile *dto.MemberProfileResponse
		if user.MemberProfile != nil {
			// Calculate total remaining sessions from active entitlements
			totalRemaining := 0
			if ents, ok := entitlements[user.ID]; ok {
				for _, ent := range ents {
					totalRemaining += ent.Remaining
				}
			}

			memberProfile = &dto.MemberProfileResponse{
				UserID:                 user.MemberProfile.UserID,
				SessionsCompleted:      user.MemberProfile.SessionsCompleted,
				TrainingTime:           user.MemberProfile.TrainingTime,
				AverageRating:          user.MemberProfile.AverageRating,
				TotalAvailableSessions: totalRemaining,
				Entitlements:           entitlements[user.ID],
			}
		}

		data[i] = dto.UserWithProfileResponse{
			GetUserResponse: dto.GetUserResponse{
				UserID:      user.ID,
				Email:       user.EmailAddress,
				Username:    user.Username,
				FirstName:   user.FirstName,
				LastName:    user.LastName,
				PhoneNumber: user.PhoneNumber,
				Image:       user.Image,
				DateOfBirth: user.DateOfBirth,
				Address:     user.Address,
				RoleID:      user.RoleID,
				Role: dto.RoleResponse{
					ID:   user.RoleID,
					Name: user.Role.Name,
				},
			},
			MemberProfile: memberProfile,
		}
	}

	return &dto.MemberListResponse{
		Data:       data,
		Pagination: dto.NewPaginationMeta(total, page, limit),
	}, nil
}
