package services

import (
	"user-service/models"
	"user-service/models/dto"
	"user-service/pkg/base"
	"user-service/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	CreateUser(input dto.CreateUserRequest) (*models.User, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetUserByIDWithProfiles(id uuid.UUID) (*models.User, error)
	GetAllUsersWithProfiles() ([]models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
}

type UserService struct {
	*base.BaseService
	repo repositories.IUserRepository
}

func NewUserService(repo repositories.IUserRepository) *UserService {
	return &UserService{repo: repo, BaseService: base.NewBaseService()}
}

func (s *UserService) CreateUser(input dto.CreateUserRequest) (*models.User, error) {
	// Check if email already exists
	exists, err := s.repo.ExistsByEmail(input.EmailAddress)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	// Check if username already exists
	exists, err = s.repo.ExistsByUsername(input.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUsernameAlreadyExists
	}

	// Check if phone number already exists
	exists, err = s.repo.ExistsByPhoneNumber(input.PhoneNumber)
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

	// Create user
	user := &models.User{
		Username:     input.Username,
		PasswordHash: string(hashedPassword),
		Email:        input.EmailAddress,
		EmailAddress: input.EmailAddress,
		PhoneNumber:  input.PhoneNumber,
		Image:        input.Image,
		DateOfBirth:  input.DateOfBirth,
		RoleID:       input.RoleID,
	}

	// Set optional fields
	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Address != "" {
		user.Address = input.Address
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByID(id uuid.UUID) (*models.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) GetUserByIDWithProfiles(id uuid.UUID) (*models.User, error) {
	return s.repo.FindByIDWithProfiles(id)
}

func (s *UserService) GetAllUsersWithProfiles() ([]models.User, error) {
	return s.repo.GetAllWithProfiles()
}

func (s *UserService) UpdateUser(user *models.User) error {
	return s.repo.Update(user)
}

func (s *UserService) DeleteUser(user *models.User) error {
	return s.repo.Delete(user)
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.repo.FindByEmail(email)
}

// Error definitions
var (
	ErrEmailAlreadyExists     = &UserServiceError{Message: "email already exists"}
	ErrUsernameAlreadyExists  = &UserServiceError{Message: "username already exists"}
	ErrPhoneNumberAlreadyExists = &UserServiceError{Message: "phone number already exists"}
)

type UserServiceError struct {
	Message string
}

func (e *UserServiceError) Error() string {
	return e.Message
}