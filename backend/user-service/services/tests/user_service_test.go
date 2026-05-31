package services_test

import (
	"context"
	"testing"
	"time"

	"user-service/models"
	"user-service/models/dto"
	apperrors "user-service/pkg/errors"
	"user-service/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ==================== Mock Repositories ====================

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *dto.RegisterRequest) (*models.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error) {
	args := m.Called(ctx, phoneNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) ExistsByPhoneNumber(ctx context.Context, phoneNumber string) (bool, error) {
	args := m.Called(ctx, phoneNumber)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) FindByRoleID(ctx context.Context, roleID uint) ([]models.User, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) FindAll(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetAllWithProfiles(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) FindByIDWithProfiles(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) CountByRoleID(ctx context.Context, roleID uint) (int64, error) {
	args := m.Called(ctx, roleID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) FindByRoleIDWithPagination(ctx context.Context, roleID uint, offset, limit int) ([]models.User, error) {
	args := m.Called(ctx, roleID, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) FindRoleByID(ctx context.Context, id uint) (*models.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleRepository) FindRoleByName(ctx context.Context, name string) (*models.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleRepository) FindAllRoles(ctx context.Context) ([]models.Role, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Role), args.Error(1)
}

func (m *MockRoleRepository) UpdateUserRole(ctx context.Context, userID uuid.UUID, roleID uint) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

// ==================== Test Setup ====================

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)
	return db
}

func createMockUser() *models.User {
	return &models.User{
		ID:           uuid.New(),
		Username:     "testuser",
		PasswordHash: "$2a$10$hashedpassword",
		EmailAddress: "test@example.com",
		PhoneNumber:  "+1234567890",
		RoleID:       1,
		IsActive:     true,
		IsVerified:   false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func createMockRole(id uint, name string) *models.Role {
	return &models.Role{
		ID:   id,
		Name: name,
	}
}

func createMockUserService(mockUserRepo *MockUserRepository, mockRoleRepo *MockRoleRepository) services.IUserService {
	return services.NewUserService(mockUserRepo, mockRoleRepo)
}

// ==================== CreateUser Tests ====================

func TestCreateUser_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	input := dto.CreateUserRequest{
		FirstName:    "John",
		LastName:     "Doe",
		Username:     "johndoe",
		Password:     "password123",
		EmailAddress: "johndoe@example.com",
		PhoneNumber:  "081234567890",
		DateOfBirth:  time.Date(1999, 8, 22, 0, 0, 0, 0, time.UTC),
		Image:        "https://example.com/image.jpg",
		RoleID:       1,
	}

	mockUserRepo.On("ExistsByEmail", ctx, input.EmailAddress).Return(false, nil)
	mockUserRepo.On("ExistsByUsername", ctx, input.Username).Return(false, nil)
	mockUserRepo.On("ExistsByPhoneNumber", ctx, input.PhoneNumber).Return(false, nil)
	mockUserRepo.On("Create", ctx, mock.AnythingOfType("*dto.RegisterRequest")).Return(&models.User{
		ID:           uuid.New(),
		Username:     input.Username,
		EmailAddress: input.EmailAddress,
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		PhoneNumber:  input.PhoneNumber,
		RoleID:       input.RoleID,
		IsActive:     true,
	}, nil)

	user, err := userService.CreateUser(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, input.Username, user.Username)
	assert.Equal(t, input.EmailAddress, user.EmailAddress)
	mockUserRepo.AssertExpectations(t)
}

func TestCreateUser_EmailAlreadyExists(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	input := dto.CreateUserRequest{
		FirstName:    "John",
		LastName:     "Doe",
		Username:     "johndoe",
		Password:     "password123",
		EmailAddress: "johndoe@example.com",
		PhoneNumber:  "081234567890",
		DateOfBirth:  time.Date(1999, 8, 22, 0, 0, 0, 0, time.UTC),
		Image:        "https://example.com/image.jpg",
		RoleID:       1,
	}

	mockUserRepo.On("ExistsByEmail", ctx, input.EmailAddress).Return(true, nil)

	user, err := userService.CreateUser(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, services.ErrEmailAlreadyExists, err)
	mockUserRepo.AssertExpectations(t)
}

func TestCreateUser_UsernameAlreadyExists(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	input := dto.CreateUserRequest{
		FirstName:    "John",
		LastName:     "Doe",
		Username:     "johndoe",
		Password:     "password123",
		EmailAddress: "johndoe@example.com",
		PhoneNumber:  "081234567890",
		DateOfBirth:  time.Date(1999, 8, 22, 0, 0, 0, 0, time.UTC),
		Image:        "https://example.com/image.jpg",
		RoleID:       1,
	}

	mockUserRepo.On("ExistsByEmail", ctx, input.EmailAddress).Return(false, nil)
	mockUserRepo.On("ExistsByUsername", ctx, input.Username).Return(true, nil)

	user, err := userService.CreateUser(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, services.ErrUsernameAlreadyExists, err)
	mockUserRepo.AssertExpectations(t)
}

func TestCreateUser_PhoneNumberAlreadyExists(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	input := dto.CreateUserRequest{
		FirstName:    "John",
		LastName:     "Doe",
		Username:     "johndoe",
		Password:     "password123",
		EmailAddress: "johndoe@example.com",
		PhoneNumber:  "081234567890",
		DateOfBirth:  time.Date(1999, 8, 22, 0, 0, 0, 0, time.UTC),
		Image:        "https://example.com/image.jpg",
		RoleID:       1,
	}

	mockUserRepo.On("ExistsByEmail", ctx, input.EmailAddress).Return(false, nil)
	mockUserRepo.On("ExistsByUsername", ctx, input.Username).Return(false, nil)
	mockUserRepo.On("ExistsByPhoneNumber", ctx, input.PhoneNumber).Return(true, nil)

	user, err := userService.CreateUser(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, services.ErrPhoneNumberAlreadyExists, err)
	mockUserRepo.AssertExpectations(t)
}

func TestCreateUser_ExistsByEmailError(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	input := dto.CreateUserRequest{
		FirstName:    "John",
		LastName:     "Doe",
		Username:     "johndoe",
		Password:     "password123",
		EmailAddress: "johndoe@example.com",
		PhoneNumber:  "081234567890",
		DateOfBirth:  time.Date(1999, 8, 22, 0, 0, 0, 0, time.UTC),
		Image:        "https://example.com/image.jpg",
		RoleID:       1,
	}

	mockUserRepo.On("ExistsByEmail", ctx, input.EmailAddress).Return(false, assert.AnError)

	user, err := userService.CreateUser(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockUserRepo.AssertExpectations(t)
}

func TestCreateUser_ExistsByUsernameError(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	input := dto.CreateUserRequest{
		FirstName:    "John",
		LastName:     "Doe",
		Username:     "johndoe",
		Password:     "password123",
		EmailAddress: "johndoe@example.com",
		PhoneNumber:  "081234567890",
		DateOfBirth:  time.Date(1999, 8, 22, 0, 0, 0, 0, time.UTC),
		Image:        "https://example.com/image.jpg",
		RoleID:       1,
	}

	mockUserRepo.On("ExistsByEmail", ctx, input.EmailAddress).Return(false, nil)
	mockUserRepo.On("ExistsByUsername", ctx, input.Username).Return(false, assert.AnError)

	user, err := userService.CreateUser(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockUserRepo.AssertExpectations(t)
}

func TestCreateUser_ExistsByPhoneNumberError(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	input := dto.CreateUserRequest{
		FirstName:    "John",
		LastName:     "Doe",
		Username:     "johndoe",
		Password:     "password123",
		EmailAddress: "johndoe@example.com",
		PhoneNumber:  "081234567890",
		DateOfBirth:  time.Date(1999, 8, 22, 0, 0, 0, 0, time.UTC),
		Image:        "https://example.com/image.jpg",
		RoleID:       1,
	}

	mockUserRepo.On("ExistsByEmail", ctx, input.EmailAddress).Return(false, nil)
	mockUserRepo.On("ExistsByUsername", ctx, input.Username).Return(false, nil)
	mockUserRepo.On("ExistsByPhoneNumber", ctx, input.PhoneNumber).Return(false, assert.AnError)

	user, err := userService.CreateUser(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockUserRepo.AssertExpectations(t)
}

func TestCreateUser_CreateError(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	input := dto.CreateUserRequest{
		FirstName:    "John",
		LastName:     "Doe",
		Username:     "johndoe",
		Password:     "password123",
		EmailAddress: "johndoe@example.com",
		PhoneNumber:  "081234567890",
		DateOfBirth:  time.Date(1999, 8, 22, 0, 0, 0, 0, time.UTC),
		Image:        "https://example.com/image.jpg",
		RoleID:       1,
	}

	mockUserRepo.On("ExistsByEmail", ctx, input.EmailAddress).Return(false, nil)
	mockUserRepo.On("ExistsByUsername", ctx, input.Username).Return(false, nil)
	mockUserRepo.On("ExistsByPhoneNumber", ctx, input.PhoneNumber).Return(false, nil)
	mockUserRepo.On("Create", ctx, mock.AnythingOfType("*dto.RegisterRequest")).Return(nil, assert.AnError)

	user, err := userService.CreateUser(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockUserRepo.AssertExpectations(t)
}

// ==================== GetUserByID Tests ====================

func TestGetUserByID_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	userID := uuid.New()
	expectedUser := createMockUser()
	expectedUser.ID = userID

	mockUserRepo.On("FindByID", ctx, userID).Return(expectedUser, nil)

	user, err := userService.GetUserByID(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	mockUserRepo.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	userID := uuid.New()

	mockUserRepo.On("FindByID", ctx, userID).Return(nil, apperrors.ErrNotFound)

	user, err := userService.GetUserByID(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockUserRepo.AssertExpectations(t)
}

func TestGetUserByID_Error(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	userID := uuid.New()

	mockUserRepo.On("FindByID", ctx, userID).Return(nil, assert.AnError)

	user, err := userService.GetUserByID(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockUserRepo.AssertExpectations(t)
}

// ==================== GetUserByIDWithProfiles Tests ====================

func TestGetUserByIDWithProfiles_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	userID := uuid.New()
	expectedUser := createMockUser()
	expectedUser.ID = userID
	expectedUser.Role = *createMockRole(1, "member")

	mockUserRepo.On("FindByIDWithProfiles", ctx, userID).Return(expectedUser, nil)

	user, err := userService.GetUserByIDWithProfiles(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	assert.NotNil(t, user.Role)
	mockUserRepo.AssertExpectations(t)
}

func TestGetUserByIDWithProfiles_NotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	userID := uuid.New()

	mockUserRepo.On("FindByIDWithProfiles", ctx, userID).Return(nil, apperrors.ErrNotFound)

	user, err := userService.GetUserByIDWithProfiles(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockUserRepo.AssertExpectations(t)
}

func TestGetUserByIDWithProfiles_Error(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	userID := uuid.New()

	mockUserRepo.On("FindByIDWithProfiles", ctx, userID).Return(nil, assert.AnError)

	user, err := userService.GetUserByIDWithProfiles(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockUserRepo.AssertExpectations(t)
}

// ==================== GetAllUsersWithProfiles Tests ====================

func TestGetAllUsersWithProfiles_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	expectedUsers := []models.User{
		*createMockUser(),
		*createMockUser(),
	}
	expectedUsers[0].Username = "user1"
	expectedUsers[1].Username = "user2"

	mockUserRepo.On("GetAllWithProfiles", ctx).Return(expectedUsers, nil)

	users, err := userService.GetAllUsersWithProfiles(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 2)
	mockUserRepo.AssertExpectations(t)
}

func TestGetAllUsersWithProfiles_EmptyList(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()

	mockUserRepo.On("GetAllWithProfiles", ctx).Return([]models.User{}, nil)

	users, err := userService.GetAllUsersWithProfiles(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 0)
	mockUserRepo.AssertExpectations(t)
}

func TestGetAllUsersWithProfiles_Error(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()

	mockUserRepo.On("GetAllWithProfiles", ctx).Return(nil, assert.AnError)

	users, err := userService.GetAllUsersWithProfiles(ctx)

	assert.Error(t, err)
	assert.Nil(t, users)
	mockUserRepo.AssertExpectations(t)
}

// ==================== UpdateUser Tests ====================

func TestUpdateUser_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	user := createMockUser()

	mockUserRepo.On("Update", ctx, user).Return(nil)

	err := userService.UpdateUser(ctx, user)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestUpdateUser_Error(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	user := createMockUser()

	mockUserRepo.On("Update", ctx, user).Return(assert.AnError)

	err := userService.UpdateUser(ctx, user)

	assert.Error(t, err)
	mockUserRepo.AssertExpectations(t)
}

// ==================== DeleteUser Tests ====================

func TestDeleteUser_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	user := createMockUser()

	mockUserRepo.On("Delete", ctx, user).Return(nil)

	err := userService.DeleteUser(ctx, user)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestDeleteUser_Error(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	user := createMockUser()

	mockUserRepo.On("Delete", ctx, user).Return(assert.AnError)

	err := userService.DeleteUser(ctx, user)

	assert.Error(t, err)
	mockUserRepo.AssertExpectations(t)
}

// ==================== GetUserByEmail Tests ====================

func TestGetUserByEmail_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	email := "test@example.com"
	expectedUser := createMockUser()

	mockUserRepo.On("FindByEmail", ctx, email).Return(expectedUser, nil)

	user, err := userService.GetUserByEmail(ctx, email)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.EmailAddress)
	mockUserRepo.AssertExpectations(t)
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	email := "notfound@example.com"

	mockUserRepo.On("FindByEmail", ctx, email).Return(nil, apperrors.ErrNotFound)

	user, err := userService.GetUserByEmail(ctx, email)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockUserRepo.AssertExpectations(t)
}

func TestGetUserByEmail_Error(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	email := "test@example.com"

	mockUserRepo.On("FindByEmail", ctx, email).Return(nil, assert.AnError)

	user, err := userService.GetUserByEmail(ctx, email)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockUserRepo.AssertExpectations(t)
}

// ==================== GetInstructorsWithPagination Tests ====================

func TestGetInstructorsWithPagination_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	page := 1
	limit := 10
	instructorRole := createMockRole(2, "instructor")
	instructors := []models.User{
		*createMockUser(),
		*createMockUser(),
	}
	instructors[0].RoleID = 2
	instructors[1].RoleID = 2

	mockRoleRepo.On("FindRoleByName", ctx, "instructor").Return(instructorRole, nil)
	mockUserRepo.On("CountByRoleID", ctx, instructorRole.ID).Return(int64(2), nil)
	mockUserRepo.On("FindByRoleIDWithPagination", ctx, instructorRole.ID, 0, limit).Return(instructors, nil)

	response, err := userService.GetInstructorsWithPagination(ctx, page, limit)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, int64(2), response.Total)
	assert.Equal(t, page, response.Page)
	assert.Equal(t, limit, response.Limit)
	assert.Len(t, response.Data, 2)
	mockRoleRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestGetInstructorsWithPagination_RoleNotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	page := 1
	limit := 10

	mockRoleRepo.On("FindRoleByName", ctx, "instructor").Return(nil, apperrors.ErrNotFound)

	response, err := userService.GetInstructorsWithPagination(ctx, page, limit)

	assert.Error(t, err)
	assert.Nil(t, response)
	mockRoleRepo.AssertExpectations(t)
}

func TestGetInstructorsWithPagination_CountError(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	page := 1
	limit := 10
	instructorRole := createMockRole(2, "instructor")

	mockRoleRepo.On("FindRoleByName", ctx, "instructor").Return(instructorRole, nil)
	mockUserRepo.On("CountByRoleID", ctx, instructorRole.ID).Return(int64(0), assert.AnError)

	response, err := userService.GetInstructorsWithPagination(ctx, page, limit)

	assert.Error(t, err)
	assert.Nil(t, response)
	mockRoleRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestGetInstructorsWithPagination_FindError(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	page := 1
	limit := 10
	instructorRole := createMockRole(2, "instructor")

	mockRoleRepo.On("FindRoleByName", ctx, "instructor").Return(instructorRole, nil)
	mockUserRepo.On("CountByRoleID", ctx, instructorRole.ID).Return(int64(2), nil)
	mockUserRepo.On("FindByRoleIDWithPagination", ctx, instructorRole.ID, 0, limit).Return(nil, assert.AnError)

	response, err := userService.GetInstructorsWithPagination(ctx, page, limit)

	assert.Error(t, err)
	assert.Nil(t, response)
	mockRoleRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestGetInstructorsWithPagination_EmptyList(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	page := 1
	limit := 10
	instructorRole := createMockRole(2, "instructor")

	mockRoleRepo.On("FindRoleByName", ctx, "instructor").Return(instructorRole, nil)
	mockUserRepo.On("CountByRoleID", ctx, instructorRole.ID).Return(int64(0), nil)
	mockUserRepo.On("FindByRoleIDWithPagination", ctx, instructorRole.ID, 0, limit).Return([]models.User{}, nil)

	response, err := userService.GetInstructorsWithPagination(ctx, page, limit)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, int64(0), response.Total)
	assert.Len(t, response.Data, 0)
	mockRoleRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestGetInstructorsWithPagination_PaginationCalculation(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := createMockUserService(mockUserRepo, mockRoleRepo)

	ctx := context.Background()
	page := 2
	limit := 5
	instructorRole := createMockRole(2, "instructor")
	instructors := []models.User{*createMockUser()}

	mockRoleRepo.On("FindRoleByName", ctx, "instructor").Return(instructorRole, nil)
	mockUserRepo.On("CountByRoleID", ctx, instructorRole.ID).Return(int64(12), nil)
	mockUserRepo.On("FindByRoleIDWithPagination", ctx, instructorRole.ID, 5, limit).Return(instructors, nil)

	response, err := userService.GetInstructorsWithPagination(ctx, page, limit)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, int64(12), response.Total)
	assert.Equal(t, page, response.Page)
	assert.Equal(t, limit, response.Limit)
	assert.Equal(t, 3, response.TotalPages) // 12/5 = 2.4, rounded up to 3
	mockRoleRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}
