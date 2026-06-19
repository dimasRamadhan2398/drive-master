package repositories

import (
	"context"
	"errors"

	"user-service/models"
	"user-service/pkg/base"
	apperrors "user-service/pkg/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IRoleRepository interface {
	FindRoleByID(ctx context.Context, id uint) (*models.Role, error)
	FindRoleByName(ctx context.Context, name string) (*models.Role, error)
	FindAllRoles(ctx context.Context) ([]models.Role, error)
	UpdateUserRole(ctx context.Context, userID uuid.UUID, roleID uint) error
}

type RoleRepository struct {
	*base.BaseRepository
}

func NewRoleRepository(db *gorm.DB) IRoleRepository {
	return &RoleRepository{BaseRepository: base.NewBaseRepository(db)}
}

// FindRoleByID implements IRoleRepository
func (r *RoleRepository) FindRoleByID(ctx context.Context, id uint) (*models.Role, error) {
	var role models.Role
	if err := r.BaseRepository.FindByID(&role, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.TranslateDBError(err)
	}
	return &role, nil
}

func (r *RoleRepository) FindRoleByName(ctx context.Context, name string) (*models.Role, error) {
	var roles []models.Role
	if err := r.BaseRepository.FindMany(&models.Role{}, &roles, base.NewQueryOptions().WithWhere(map[string]any{"name": name}).WithPagination(0, 1)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.TranslateDBError(err)
	}
	if len(roles) == 0 {
		return nil, apperrors.ErrNotFound
	}
	return &roles[0], nil
}

// FindAllRoles implements IRoleRepository
func (r *RoleRepository) FindAllRoles(ctx context.Context) ([]models.Role, error) {
	var roles []models.Role
	if err := r.BaseRepository.FindMany(&models.Role{}, &roles, base.NewQueryOptions()); err != nil {
		return nil, apperrors.TranslateDBError(err)
	}
	return roles, nil
}

// UpdateUserRole implements IRoleRepository
func (r *RoleRepository) UpdateUserRole(ctx context.Context, userID uuid.UUID, roleID uint) error {
	var roles []models.Role
	err := r.BaseRepository.FindMany(&models.Role{}, &roles, base.NewQueryOptions().WithWhere(map[string]any{"id": roleID}).WithPagination(0, 1))
	if err != nil {
		return apperrors.TranslateDBError(err)
	}
	if len(roles) == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}
