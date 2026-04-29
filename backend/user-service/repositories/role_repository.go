package repositories

import (
	"errors"

	"user-service/models"
	"user-service/pkg/base"
	apperrors "user-service/pkg/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IRoleRepository interface {
	FindRoleByID(id uint) (*models.Role, error)
	FindAllRoles() ([]models.Role, error)
	UpdateUserRole(userID uuid.UUID, roleID uint) error
}

type RoleRepository struct {
	*base.BaseRepository
}

func NewRoleRepository(db *gorm.DB) IRoleRepository {
	return &RoleRepository{BaseRepository: base.NewBaseRepository(db)}
}

// FindRoleByID implements IRoleRepository
func (r *RoleRepository) FindRoleByID(id uint) (*models.Role, error) {
	var role models.Role
	if err := r.GetDB().First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &role, nil
}

// FindAllRoles implements IRoleRepository
func (r *RoleRepository) FindAllRoles() ([]models.Role, error) {
	var roles []models.Role
	if err := r.GetDB().Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// UpdateUserRole implements IRoleRepository
func (r *RoleRepository) UpdateUserRole(userID uuid.UUID, roleID uint) error {
	result := r.GetDB().Model(&models.User{}).Where("id = ?", userID).Update("role_id", roleID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}