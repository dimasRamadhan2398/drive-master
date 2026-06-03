package services

import (
	"context"
	"user-service/models"
	"user-service/models/dto"
	"user-service/pkg/base"
	"user-service/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	CreateUser(ctx context.Context, input dto.CreateUserRequest) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByIDWithProfiles(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetAllUsersWithProfiles(ctx context.Context) ([]models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetInstructorsWithPagination(ctx context.Context, page, limit int) (*dto.InstructorListResponse, error)
	GetMembersWithPagination(ctx context.Context, page, limit int) (*dto.MemberListResponse, error)
}

type UserService struct {
	*base.BaseService
	roleRepo repositories.IRoleRepository
	repo     repositories.IUserRepository
}

// GetMembersWithPagination implements [IUserService].

func NewUserService(repo repositories.IUserRepository, roleRepo repositories.IRoleRepository) IUserService {
	return &UserService{repo: repo, roleRepo: roleRepo, BaseService: base.NewBaseService()}
}

func (s *UserService) CreateUser(ctx context.Context, input dto.CreateUserRequest) (*models.User, error) {
	// Check if email already exists
	exists, err := s.repo.ExistsByEmail(ctx, input.EmailAddress)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	// Check if username already exists
	exists, err = s.repo.ExistsByUsername(ctx, input.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUsernameAlreadyExists
	}

	// Check if phone number already exists
	exists, err = s.repo.ExistsByPhoneNumber(ctx, input.PhoneNumber)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrPhoneNumberAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	registerReq := &dto.RegisterRequest{
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Username:    input.Username,
		Email:       input.EmailAddress,
		PhoneNumber: input.PhoneNumber,
		Password:    string(hashedPassword),
		RoleID:      input.RoleID,
	}

	createdUser, err := s.repo.Create(ctx, registerReq)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) GetUserByIDWithProfiles(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.repo.FindByIDWithProfiles(ctx, id)
}

func (s *UserService) GetAllUsersWithProfiles(ctx context.Context) ([]models.User, error) {
	return s.repo.GetAllWithProfiles(ctx)
}

func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	return s.repo.Update(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, user *models.User) error {
	return s.repo.Delete(ctx, user)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.FindByEmail(ctx, email)
}

// GetInstructorsWithPagination returns paginated list of instructors with their profiles
func (s *UserService) GetInstructorsWithPagination(ctx context.Context, page, limit int) (*dto.InstructorListResponse, error) {
	roleModel, err := s.roleRepo.FindRoleByName(ctx, "instructor")
	if err != nil {
		return nil, err
	}

	total, err := s.repo.CountByRoleID(ctx, roleModel.ID)
	if err != nil {
		return nil, err
	}

	// Calculate pagination
	offset := (page - 1) * limit

	// Get paginated users
	users, err := s.repo.FindByRoleIDWithPagination(ctx, roleModel.ID, offset, limit)
	if err != nil {
		return nil, err
	}

	// Convert to response DTOs
	data := make([]dto.UserWithProfileResponse, len(users))
	for i, user := range users {
		data[i] = dto.UserWithProfileResponse{
			GetUserResponse: dto.GetUserResponse{
				UserID:      user.ID,
				Email:       user.EmailAddress,
				Username:    user.Username,
				PhoneNumber: user.PhoneNumber,
				Image:       user.Image,
				DateOfBirth: user.DateOfBirth,
				Address:     user.Address,
				RoleID:      user.RoleID,
			},
		}
	}

	return &dto.InstructorListResponse{
		Data:       data,
		Pagination: dto.NewPaginationMeta(total, page, limit),
	}, nil
}

func (s *UserService) GetMembersWithPagination(ctx context.Context, page int, limit int) (*dto.MemberListResponse, error) {
	roleModel, err := s.roleRepo.FindRoleByName(ctx, "member")
	if err != nil {
		return nil, err
	}

	total, err := s.repo.CountByRoleID(ctx, roleModel.ID)
	if err != nil {
		return nil, err
	}

	// Calculate pagination
	offset := (page - 1) * limit

	// Get paginated users
	users, err := s.repo.FindByRoleIDWithPagination(ctx, roleModel.ID, offset, limit)
	if err != nil {
		return nil, err
	}

	// Convert to response DTOs
	data := make([]dto.UserWithProfileResponse, len(users))
	for i, user := range users {
		data[i] = dto.UserWithProfileResponse{
			GetUserResponse: dto.GetUserResponse{
				UserID:      user.ID,
				Email:       user.EmailAddress,
				Username:    user.Username,
				PhoneNumber: user.PhoneNumber,
				Image:       user.Image,
				DateOfBirth: user.DateOfBirth,
				Address:     user.Address,
				RoleID:      user.RoleID,
			},
		}
	}

	return &dto.MemberListResponse{
		Data:       data,
		Pagination: dto.NewPaginationMeta(total, page, limit),
	}, nil
}
// Error definitions
var (
	ErrEmailAlreadyExists       = &UserServiceError{Message: "email already exists"}
	ErrUsernameAlreadyExists    = &UserServiceError{Message: "username already exists"}
	ErrPhoneNumberAlreadyExists = &UserServiceError{Message: "phone number already exists"}
)

type UserServiceError struct {
	Message string
}

func (e *UserServiceError) Error() string {
	return e.Message
}
