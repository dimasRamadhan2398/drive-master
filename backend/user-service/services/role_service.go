package services

import (
	"user-service/models"
	"user-service/repositories"

	"github.com/google/uuid"
)

type IRoleService interface {
	GetRole(id uint) (*models.Role, error)
	GetAllRoles() ([]models.Role, error)
	UpdateUserRole(userID uuid.UUID, roleID uint) error
}

type RoleService struct {
	repo repositories.IRoleRepository
}

func NewRoleService(repo repositories.IRoleRepository) *RoleService {
	return &RoleService{repo: repo}
}

// GetRole retrieves a role by ID
func (s *RoleService) GetRole(id uint) (*models.Role, error) {
	return s.repo.FindRoleByID(id)
}

// GetAllRoles retrieves all roles
func (s *RoleService) GetAllRoles() ([]models.Role, error) {
	return s.repo.FindAllRoles()
}

// UpdateUserRole updates the user's role
func (s *RoleService) UpdateUserRole(userID uuid.UUID, roleID uint) error {
	return s.repo.UpdateUserRole(userID, roleID)
}