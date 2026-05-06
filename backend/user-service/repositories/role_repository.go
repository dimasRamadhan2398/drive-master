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
	err := r.BaseRepository.FindMany(&models.Role{}, nil, base.NewQueryOptions().WithWhere(map[string]any{"id": roleID}))
	if err != nil {
		return apperrors.TranslateDBError(err)
	}
	return nil
}