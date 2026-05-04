package repositories

import (
	"context"
	"time"
	"user-service/models"
	"user-service/models/dto"
	"user-service/pkg/base"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Create(ctx context.Context, user *dto.RegisterRequest) (*models.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByPhoneNumber(ctx context.Context, phoneNumber string) (bool, error)
	FindByRoleID(ctx context.Context, roleID uint) ([]models.User, error)
	FindAll(ctx context.Context) ([]models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, user *models.User) error
	GetAllWithProfiles(ctx context.Context) ([]models.User, error)
	FindByIDWithProfiles(ctx context.Context, id uuid.UUID) (*models.User, error)
	CountByRoleID(ctx context.Context, roleID uint) (int64, error)
	FindByRoleIDWithPagination(ctx context.Context, roleID uint, offset, limit int) ([]models.User, error)
}

type UserRepository struct {
	*base.BaseRepository
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{BaseRepository: base.NewBaseRepository(db)}
}

func (r *UserRepository) Create(ctx context.Context, user *dto.RegisterRequest) (*models.User, error) {
	t, err := time.Parse(user.DateOfBirth, "2006-01-02")
	if(err != nil) {
		t = time.Now()
	}
	userModel := models.User{
		Name:         user.Name,
		Username:     user.Username,
		Email:        user.Email,
		EmailAddress: user.Email,
		PhoneNumber:  user.PhoneNumber,
		PasswordHash: user.Password,
		DateOfBirth:  t,
		RoleID:       user.RoleID,
		IsActive:     true,
	}

	if err := r.BaseRepository.Create(&userModel); err != nil {
		return nil, err
	}

	return &userModel, nil
}


func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.BaseRepository.FindByID(&user, id); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.BaseRepository.FindOne(&user, "email = ?", email); err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername finds a user by username
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := r.BaseRepository.FindOne(&user, "username = ?", username); err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByPhoneNumber finds a user by phone number
func (r *UserRepository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error) {
	var user models.User
	if err := r.BaseRepository.FindOne(&user, "phone_number = ?", phoneNumber); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetAllWithProfiles(ctx context.Context) ([]models.User, error) {
	var users []models.User
	opts := base.NewQueryOptions().
		WithPreloads("Role", "MemberProfile", "InstructorProfile")
	if err := r.BaseRepository.FindMany(&models.User{}, &users, opts); err != nil {
		return nil, err
	}
	return users, nil
}

// FindByIDWithProfiles finds a user by UUID with their Role, MemberProfile, and InstructorProfile
func (r *UserRepository) FindByIDWithProfiles(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.BaseRepository.FindByIDWithPreload(&user, id, "Role", "MemberProfile", "InstructorProfile"); err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return r.BaseRepository.Exists(&models.User{}, "email = ?", email)
}

// ExistsByUsername checks if a user exists with the given username
func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	return r.BaseRepository.Exists(&models.User{}, "username = ?", username)
}

// ExistsByPhoneNumber checks if a user exists with the given phone number
func (r *UserRepository) ExistsByPhoneNumber(ctx context.Context, phoneNumber string) (bool, error) {
	return r.BaseRepository.Exists(&models.User{}, "phone_number = ?", phoneNumber)
}

func (r *UserRepository) FindByRoleID(ctx context.Context, roleID uint) ([]models.User, error) {
	var users []models.User
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"role_id": roleID})
	if err := r.BaseRepository.FindMany(&models.User{}, &users, opts); err != nil {
		return nil, err
	}
	return users, nil
}

// FindAll retrieves all users
func (r *UserRepository) FindAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	opts := base.NewQueryOptions()
	if err := r.BaseRepository.FindMany(&models.User{}, &users, opts); err != nil {
		return nil, err
	}
	return users, nil
}

// Update saves changes to an existing user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	return r.BaseRepository.Update(user)
}

// Delete soft-deletes a user
func (r *UserRepository) Delete(ctx context.Context, user *models.User) error {
	return r.BaseRepository.Delete(user)
}

// CountByRoleID counts the number of users with a specific role ID
func (r *UserRepository) CountByRoleID(ctx context.Context, roleID uint) (int64, error) {
	var count int64
	if err := r.DB.Model(&models.User{}).Where("role_id = ?", roleID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// FindByRoleIDWithPagination retrieves users with a specific role ID with pagination
func (r *UserRepository) FindByRoleIDWithPagination(ctx context.Context, roleID uint, offset, limit int) ([]models.User, error) {
	var users []models.User
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"role_id": roleID}).
		WithPagination(offset, limit).
		WithPreloads("Role", "MemberProfile", "InstructorProfile")
	if err := r.BaseRepository.FindMany(&models.User{}, &users, opts); err != nil {
		return nil, err
	}
	return users, nil
}
