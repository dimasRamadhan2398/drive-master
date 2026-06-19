package services

import (
	"context"
	"time"
	"user-service/models/dto"
	"user-service/repositories"

	"user-service/pkg/constants"
)

type IDashboardService interface {
	GetStats(ctx context.Context) (*dto.DashboardStatsResponse, error)
}

type DashboardService struct {
	userRepo repositories.IUserRepository
	roleRepo repositories.IRoleRepository
}

func NewDashboardService(
	userRepo repositories.IUserRepository,
	roleRepo repositories.IRoleRepository,
) IDashboardService {
	return &DashboardService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (s *DashboardService) GetStats(ctx context.Context) (*dto.DashboardStatsResponse, error) {
	// Get total users
	totalUsers, err := s.userRepo.CountAll(ctx)
	if err != nil {
		totalUsers = 0
	}

	// Get member role
	memberRole, err := s.roleRepo.FindRoleByName(ctx, constants.Member)
	memberCount := int64(0)
	if err == nil && memberRole != nil {
		memberCount, _ = s.userRepo.CountByRoleID(ctx, memberRole.ID)
	}

	// Get instructor role
	instructorRole, err := s.roleRepo.FindRoleByName(ctx, constants.Instructor)
	instructorCount := int64(0)
	if err == nil && instructorRole != nil {
		instructorCount, _ = s.userRepo.CountByRoleID(ctx, instructorRole.ID)
	}

	// Get recent registrations (last 30 days)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	filters := &dto.RegistrationFilters{
		FromDate: &thirtyDaysAgo,
	}
	recentUsers, _ := s.userRepo.FindRecentRegistrations(ctx, 0, filters)
	recentCount := int64(len(recentUsers))

	return &dto.DashboardStatsResponse{
		TotalUsers:          totalUsers,
		TotalMembers:        memberCount,
		TotalInstructors:    instructorCount,
		RecentRegistrations: recentCount,
	}, nil
}
