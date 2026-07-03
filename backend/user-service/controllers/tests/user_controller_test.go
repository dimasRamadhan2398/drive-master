package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"user-service/controllers"
	"user-service/models"
	"user-service/models/dto"
	apperrors "user-service/pkg/errors"
	"user-service/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, input dto.CreateUserRequest) (*models.User, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUserByIDWithProfiles(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetAllUsersWithProfiles(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetInstructorsWithPagination(ctx context.Context, page, limit int) (*dto.InstructorListResponse, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.InstructorListResponse), args.Error(1)
}

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(ctx context.Context, req *dto.LoginInput) (*dto.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LoginResponse), args.Error(1)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, refreshToken string) (*dto.LoginResponse, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LoginResponse), args.Error(1)
}

func (m *MockAuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RegisterResponse), args.Error(1)
}

func (m *MockAuthService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	args := m.Called(ctx, userID, oldPassword, newPassword)
	return args.Error(0)
}

func (m *MockAuthService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) GenerateAndSendOTP(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockAuthService) VerifyOTP(ctx context.Context, email, otp string) error {
	args := m.Called(ctx, email, otp)
	return args.Error(0)
}

func (m *MockAuthService) ResendOTP(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

type MockMemberService struct {
	mock.Mock
}

func (m *MockMemberService) GetMemberProfile(ctx context.Context, userID uuid.UUID) (*dto.MemberProfileResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MemberProfileResponse), args.Error(1)
}

func (m *MockMemberService) GetRawMemberProfile(ctx context.Context, userID uuid.UUID) (*models.MemberProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MemberProfile), args.Error(1)
}

func (m *MockMemberService) CreateMemberProfile(ctx context.Context, userID uuid.UUID) (*models.MemberProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MemberProfile), args.Error(1)
}

func (m *MockMemberService) UpdateMemberProfile(ctx context.Context, profile *models.MemberProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockMemberService) DeleteMemberProfile(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

type MockInstructorService struct {
	mock.Mock
}

func (m *MockInstructorService) GetInstructorProfile(ctx context.Context, userID uuid.UUID) (*dto.InstructorProfileResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.InstructorProfileResponse), args.Error(1)
}

func (m *MockInstructorService) CreateInstructorProfile(ctx context.Context, userID uuid.UUID) (*dto.InstructorProfileResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.InstructorProfileResponse), args.Error(1)
}

func (m *MockInstructorService) UpdateInstructorProfile(ctx context.Context, profile *models.InstructorProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockInstructorService) DeleteInstructorProfile(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

type MockRoleService struct {
	mock.Mock
}

func (m *MockRoleService) GetRole(ctx context.Context, id uint) (*models.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleService) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleService) FindAllRoles(ctx context.Context) ([]models.Role, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Role), args.Error(1)
}

func (m *MockRoleService) GetAllRoles() ([]models.Role, error) {
	args := m.Called()
	return args.Get(0).([]models.Role), args.Error(1)
}

func (m *MockRoleService) UpdateUserRole(ctx context.Context, userID uuid.UUID, roleID uint) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendEmail(ctx context.Context, input dto.SendEmailRequest) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *MockEmailService) SendWelcomeEmail(ctx context.Context, toEmail, username string) error {
	args := m.Called(ctx, toEmail, username)
	return args.Error(0)
}

func (m *MockEmailService) SendPasswordResetEmail(ctx context.Context, toEmail, resetToken string) error {
	args := m.Called(ctx, toEmail, resetToken)
	return args.Error(0)
}

func (m *MockEmailService) SendOTPEmail(ctx context.Context, toEmail, otp string) error {
	args := m.Called(ctx, toEmail, otp)
	return args.Error(0)
}

func (m *MockEmailService) SendBookingConfirmationEmail(ctx context.Context, toEmail, studentName, instructorName, dateTime, lessonType string) error {
	args := m.Called(ctx, toEmail, studentName, instructorName, dateTime, lessonType)
	return args.Error(0)
}

func (m *MockEmailService) SendLessonReminderEmail(ctx context.Context, toEmail, studentName, instructorName, dateTime, lessonType string) error {
	args := m.Called(ctx, toEmail, studentName, instructorName, dateTime, lessonType)
	return args.Error(0)
}

func (m *MockEmailService) SendLessonCancellationEmail(ctx context.Context, toEmail, studentName, instructorName, dateTime, reason string) error {
	args := m.Called(ctx, toEmail, studentName, instructorName, dateTime, reason)
	return args.Error(0)
}

type MockMediaService struct {
	mock.Mock
}

func (m *MockMediaService) UploadMedia(ctx context.Context, input services.UploadMediaInput) (*services.MediaUploadResponse, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.MediaUploadResponse), args.Error(1)
}

func (m *MockMediaService) UploadBase64Media(ctx context.Context, base64Data, fileName, folder string) (*services.MediaUploadResponse, error) {
	args := m.Called(ctx, base64Data, fileName, folder)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.MediaUploadResponse), args.Error(1)
}

func (m *MockMediaService) UploadUserProfileImage(ctx context.Context, fileData []byte, userID uint, fileName string) (*services.MediaUploadResponse, error) {
	args := m.Called(ctx, fileData, userID, fileName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.MediaUploadResponse), args.Error(1)
}

func (m *MockMediaService) DeleteMedia(ctx context.Context, fileID string) error {
	args := m.Called(ctx, fileID)
	return args.Error(0)
}

func (m *MockMediaService) GetMediaMetadata(ctx context.Context, fileID string) (*services.MediaUploadResponse, error) {
	args := m.Called(ctx, fileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.MediaUploadResponse), args.Error(1)
}

func (m *MockMediaService) BulkDeleteMedia(ctx context.Context, fileIDs []string) error {
	args := m.Called(ctx, fileIDs)
	return args.Error(0)
}

// ==================== Test Setup ====================

func setupUserControllerRouter(userController controllers.IUserController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	users := router.Group("/users")
	users.POST("", userController.CreateUser)
	users.GET("", userController.GetAllUsers)
	users.GET("/:id", userController.GetUserByID)
	users.PUT("/:id", userController.UpdateUser)
	users.DELETE("/:id", userController.DeleteUser)

	return router
}

func createUserController() (controllers.IUserController, *MockUserService) {
	mockUserService := new(MockUserService)
	mockAuthService := new(MockAuthService)
	mockMemberService := new(MockMemberService)
	mockInstructorService := new(MockInstructorService)
	mockRoleService := new(MockRoleService)
	mockEmailService := new(MockEmailService)
	mockMediaService := new(MockMediaService)

	controller := controllers.NewUserController(
		mockUserService,
		mockAuthService,
		mockMemberService,
		mockInstructorService,
		mockRoleService,
		mockEmailService,
		mockMediaService,
	)

	return controller, mockUserService
}

// ==================== CreateUser Tests ====================

func TestCreateUser_Success(t *testing.T) {
	controller, mockUserService := createUserController()
	router := setupUserControllerRouter(controller)

	createUserJSON := `{
		"firstName": "John",
		"lastName": "Doe",
		"username": "johndoe",
		"password": "password123",
		"emailAddress": "johndoe@example.com",
		"phoneNumber": "081234567890",
		"dateOfBirth": "1999-08-22T00:00:00Z",
		"image": "https://example.com/image.jpg",
		"roleId": 1
	}`

	mockUserService.On("CreateUser", mock.Anything, mock.AnythingOfType("dto.CreateUserRequest")).Return(&models.User{
		ID:           uuid.New(),
		Username:     "johndoe",
		EmailAddress: "johndoe@example.com",
		FirstName:    "John",
		LastName:     "Doe",
		PhoneNumber:  "081234567890",
		RoleID:       1,
		IsActive:     true,
		IsVerified:   false,
	}, nil)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(createUserJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	assert.Equal(t, "User created successfully", response["message"])

	mockUserService.AssertExpectations(t)
}

func TestCreateUser_InvalidJSON(t *testing.T) {
	controller, _ := createUserController()
	router := setupUserControllerRouter(controller)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUser_MissingRequiredFields(t *testing.T) {
	controller, _ := createUserController()
	router := setupUserControllerRouter(controller)

	// Missing required fields - only name provided
	createUserJSON := `{"firstName": "John"}`

	req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(createUserJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUser_InvalidEmail(t *testing.T) {
	controller, _ := createUserController()
	router := setupUserControllerRouter(controller)

	createUserJSON := `{
		"firstName": "John",
		"lastName": "Doe",
		"username": "johndoe",
		"password": "password123",
		"emailAddress": "invalid-email",
		"phoneNumber": "081234567890",
		"dateOfBirth": "1999-08-22T00:00:00Z",
		"image": "https://example.com/image.jpg",
		"roleId": 1
	}`

	req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(createUserJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUser_EmailAlreadyExists(t *testing.T) {
	controller, mockUserService := createUserController()
	router := setupUserControllerRouter(controller)

	createUserJSON := `{
		"firstName": "John",
		"lastName": "Doe",
		"username": "johndoe",
		"password": "password123",
		"emailAddress": "johndoe@example.com",
		"phoneNumber": "081234567890",
		"dateOfBirth": "1999-08-22T00:00:00Z",
		"image": "https://example.com/image.jpg",
		"roleId": 1
	}`

	mockUserService.On("CreateUser", mock.Anything, mock.AnythingOfType("dto.CreateUserRequest")).Return(nil, apperrors.ErrEmailExist)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(createUserJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	mockUserService.AssertExpectations(t)
}

func TestCreateUser_UsernameAlreadyExists(t *testing.T) {
	controller, mockUserService := createUserController()
	router := setupUserControllerRouter(controller)

	createUserJSON := `{
		"firstName": "John",
		"lastName": "Doe",
		"username": "johndoe",
		"password": "password123",
		"emailAddress": "johndoe@example.com",
		"phoneNumber": "081234567890",
		"dateOfBirth": "1999-08-22T00:00:00Z",
		"image": "https://example.com/image.jpg",
		"roleId": 1
	}`

	mockUserService.On("CreateUser", mock.Anything, mock.AnythingOfType("dto.CreateUserRequest")).Return(nil, apperrors.ErrUsernameExist)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(createUserJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	mockUserService.AssertExpectations(t)
}

// ==================== GetAllUsers Tests ====================

func TestGetAllUsers_Success(t *testing.T) {
	controller, mockUserService := createUserController()
	router := setupUserControllerRouter(controller)

	users := []models.User{
		{
			ID:           uuid.New(),
			Username:     "johndoe",
			EmailAddress: "johndoe@example.com",
			FirstName:    "John",
			LastName:     "Doe",
			PhoneNumber:  "081234567890",
			RoleID:       1,
			IsActive:     true,
			IsVerified:   true,
		},
		{
			ID:           uuid.New(),
			Username:     "janedoe",
			EmailAddress: "janedoe@example.com",
			FirstName:    "Jane",
			LastName:     "Doe",
			PhoneNumber:  "081234567891",
			RoleID:       2,
			IsActive:     true,
			IsVerified:   true,
		},
	}

	mockUserService.On("GetAllUsersWithProfiles", mock.Anything).Return(users, nil)

	req, _ := http.NewRequest("GET", "/users", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	assert.Equal(t, "Users retrieved successfully", response["message"])

	mockUserService.AssertExpectations(t)
}

func TestGetAllUsers_EmptyList(t *testing.T) {
	controller, mockUserService := createUserController()
	router := setupUserControllerRouter(controller)

	mockUserService.On("GetAllUsersWithProfiles", mock.Anything).Return([]models.User{}, nil)

	req, _ := http.NewRequest("GET", "/users", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])

	mockUserService.AssertExpectations(t)
}

func TestGetAllUsers_InternalError(t *testing.T) {
	controller, mockUserService := createUserController()
	router := setupUserControllerRouter(controller)

	mockUserService.On("GetAllUsersWithProfiles", mock.Anything).Return(nil, apperrors.ErrInternalServer)

	req, _ := http.NewRequest("GET", "/users", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockUserService.AssertExpectations(t)
}

// ==================== GetUserByID Tests ====================

func TestGetUserByID_Success(t *testing.T) {
	controller, mockUserService := createUserController()
	router := setupUserControllerRouter(controller)

	userID := uuid.New()

	mockUserService.On("GetUserByIDWithProfiles", mock.Anything, userID).Return(&models.User{
		ID:           userID,
		Username:     "johndoe",
		EmailAddress: "johndoe@example.com",
		FirstName:    "John",
		LastName:     "Doe",
		PhoneNumber:  "081234567890",
		RoleID:       1,
		IsActive:     true,
		IsVerified:   true,
	}, nil)

	req, _ := http.NewRequest("GET", "/users/"+userID.String(), nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	assert.Equal(t, "User retrieved successfully", response["message"])

	mockUserService.AssertExpectations(t)
}

func TestGetUserByID_InvalidUUID(t *testing.T) {
	controller, _ := createUserController()
	router := setupUserControllerRouter(controller)

	req, _ := http.NewRequest("GET", "/users/invalid-uuid", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUserByID_NotFound(t *testing.T) {
	controller, mockUserService := createUserController()
	router := setupUserControllerRouter(controller)

	userID := uuid.New()

	mockUserService.On("GetUserByIDWithProfiles", mock.Anything, userID).Return(nil, apperrors.ErrUserNotFound)

	req, _ := http.NewRequest("GET", "/users/"+userID.String(), nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	mockUserService.AssertExpectations(t)
}

// ==================== UpdateUser Tests ====================

func TestUpdateUser_Success(t *testing.T) {
	controller, mockUserService := createUserController()
	router := setupUserControllerRouter(controller)

	userID := uuid.New()

	existingUser := &models.User{
		ID:           userID,
		Username:     "johndoe",
		EmailAddress: "johndoe@example.com",
		FirstName:    "John",
		LastName:     "Doe",
		PhoneNumber:  "081234567890",
		RoleID:       1,
		IsActive:     true,
		IsVerified:   true,
	}

	updateUserJSON := `{
		"firstName": "Johnny",
		"lastName": "Doe",
		"username": "johnnydoe"
	}`

	mockUserService.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	mockUserService.On("UpdateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)

	req, _ := http.NewRequest("PUT", "/users/"+userID.String(), bytes.NewBufferString(updateUserJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	assert.Equal(t, "User updated successfully", response["message"])

	mockUserService.AssertExpectations(t)
}

func TestUpdateUser_InvalidUUID(t *testing.T) {
	controller, _ := createUserController()
	router := setupUserControllerRouter(controller)

	updateUserJSON := `{"firstName": "Johnny"}`

	req, _ := http.NewRequest("PUT", "/users/invalid-uuid", bytes.NewBufferString(updateUserJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateUser_InvalidJSON(t *testing.T) {
	controller, _ := createUserController()
	router := setupUserControllerRouter(controller)

	userID := uuid.New()

	req, _ := http.NewRequest("PUT", "/users/"+userID.String(), bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Invalid JSON causes validation error with 400 response (early return)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateUser_UserNotFound(t *testing.T) {
	controller, mockUserService := createUserController()
	router := setupUserControllerRouter(controller)

	userID := uuid.New()

	// Both firstName and lastName are required in UpdateUserRequest
	updateUserJSON := `{"firstName": "Johnny", "lastName": "Doe"}`

	mockUserService.On("GetUserByID", mock.Anything, userID).Return(nil, apperrors.ErrUserNotFound)

	req, _ := http.NewRequest("PUT", "/users/"+userID.String(), bytes.NewBufferString(updateUserJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	mockUserService.AssertExpectations(t)
}

func TestUpdateUser_ValidationError(t *testing.T) {
	controller, mockUserService := createUserController()
	router := setupUserControllerRouter(controller)

	userID := uuid.New()

	existingUser := &models.User{
		ID:           userID,
		Username:     "johndoe",
		EmailAddress: "johndoe@example.com",
	}

	// Invalid email format - both firstName and lastName required
	updateUserJSON := `{
		"firstName": "Jo",
		"lastName": "Do",
		"emailAddress": "invalid-email"
	}`

	mockUserService.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil).Maybe()

	req, _ := http.NewRequest("PUT", "/users/"+userID.String(), bytes.NewBufferString(updateUserJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockUserService.AssertExpectations(t)
}

// ==================== DeleteUser Tests ====================

func TestDeleteUser_Success(t *testing.T) {
	controller, mockUserService := createUserController()
	router := setupUserControllerRouter(controller)

	userID := uuid.New()

	existingUser := &models.User{
		ID:           userID,
		Username:     "johndoe",
		EmailAddress: "johndoe@example.com",
	}

	mockUserService.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	mockUserService.On("DeleteUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)

	req, _ := http.NewRequest("DELETE", "/users/"+userID.String(), nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	assert.Equal(t, "User deleted successfully", response["message"])

	mockUserService.AssertExpectations(t)
}

func TestDeleteUser_InvalidUUID(t *testing.T) {
	controller, _ := createUserController()
	router := setupUserControllerRouter(controller)

	req, _ := http.NewRequest("DELETE", "/users/invalid-uuid", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteUser_UserNotFound(t *testing.T) {
	controller, mockUserService := createUserController()
	router := setupUserControllerRouter(controller)

	userID := uuid.New()

	mockUserService.On("GetUserByID", mock.Anything, userID).Return(nil, apperrors.ErrUserNotFound)

	req, _ := http.NewRequest("DELETE", "/users/"+userID.String(), nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	mockUserService.AssertExpectations(t)
}

func TestDeleteUser_InternalError(t *testing.T) {
	controller, mockUserService := createUserController()
	router := setupUserControllerRouter(controller)

	userID := uuid.New()

	existingUser := &models.User{
		ID:           userID,
		Username:     "johndoe",
		EmailAddress: "johndoe@example.com",
	}

	mockUserService.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	mockUserService.On("DeleteUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(apperrors.ErrInternalServer)

	req, _ := http.NewRequest("DELETE", "/users/"+userID.String(), nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockUserService.AssertExpectations(t)
}
