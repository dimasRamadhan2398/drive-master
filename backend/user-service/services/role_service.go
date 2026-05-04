package services

import (
	"context"
	"user-service/models"
	"user-service/repositories"

	"github.com/google/uuid"
)

type IRoleService interface {
	GetRole(ctx context.Context, id uint) (*models.Role, error)
	FindAllRoles(ctx context.Context) ([]models.Role, error)
	UpdateUserRole(ctx context.Context, userID uuid.UUID, roleID uint) error
}

type RoleService struct {
	repo repositories.IRoleRepository
}

func NewRoleService(repo repositories.IRoleRepository) IRoleService {
	return &RoleService{repo: repo}
}

// GetRole retrieves a role by ID
func (s *RoleService) GetRole(ctx context.Context, id uint) (*models.Role, error) {
	return s.repo.FindRoleByID(ctx, id)
}

// GetAllRoles retrieves all roles
func (s *RoleService) FindAllRoles(ctx context.Context) ([]models.Role, error) {
	return s.repo.FindAllRoles(ctx)
}

// UpdateUserRole updates the user's role
func (s *RoleService) UpdateUserRole(ctx context.Context, userID uuid.UUID, roleID uint) error {
	return s.repo.UpdateUserRole(ctx, userID, roleID)
}