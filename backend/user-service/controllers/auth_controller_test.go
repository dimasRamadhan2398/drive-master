package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"user-service/models"
	"user-service/models/dto"
	apperrors "user-service/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Services

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

type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendEmail(ctx context.Context, input dto.SendEmailRequest) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *MockEmailService) SendWelcomeEmail(ctx context.Context, email, username string) error {
	args := m.Called(ctx, email, username)
	return args.Error(0)
}

func (m *MockEmailService) SendPasswordResetEmail(ctx context.Context, email, resetToken string) error {
	args := m.Called(ctx, email, resetToken)
	return args.Error(0)
}

func (m *MockEmailService) SendOTPEmail(ctx context.Context, email, otp string) error {
	args := m.Called(ctx, email, otp)
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

// Test Setup

func setupRouter(authController *AuthController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	auth := router.Group("/auth")
	auth.POST("/register", authController.Register)
	auth.POST("/login", authController.Login)
	auth.POST("/verify-otp", authController.VerifyOTP)
	auth.POST("/resend-otp", authController.ResendOTP)

	return router
}

func TestRegister_Success(t *testing.T) {
	mockAuthService := new(MockAuthService)
	mockUserService := new(MockUserService)
	mockEmailService := new(MockEmailService)
	mockRoleService := new(MockRoleService)

	controller := &AuthController{
		authService:  mockAuthService,
		userService:  mockUserService,
		emailService: mockEmailService,
		roleService:  mockRoleService,
	}

	router := setupRouter(controller)

	// Use raw JSON string with the correct date format (DD/MM/YYYY)
	registerJSON := `{
		"name": "John Doe",
		"username": "johndoe",
		"password": "password123",
		"confirmPassword": "password123",
		"email": "johndoe@example.com",
		"phoneNumber": "081234567890",
		"dateOfBirth": "1999-08-22",
		"roleId": 1
	}`

	mockAuthService.On("Register", mock.Anything, mock.AnythingOfType("*dto.RegisterRequest")).Return(&dto.RegisterResponse{
		User: dto.CreateUserResponse{
			UserID:      uuid.New(),
			Email:       "johndoe@example.com",
			Username:    "johndoe",
			PhoneNumber: "081234567890",
			RoleID:      1,
		},
	}, nil)

	// Mock for the goroutine call to SendWelcomeEmail
	mockEmailService.On("SendWelcomeEmail", mock.Anything, "johndoe@example.com", "johndoe").Return(nil).Maybe()

	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBufferString(registerJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	assert.Equal(t, "User registered successfully. Please check your email for OTP verification", response["message"])

	mockAuthService.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
}

func TestRegister_InvalidJSON(t *testing.T) {
	mockAuthService := new(MockAuthService)
	mockUserService := new(MockUserService)
	mockEmailService := new(MockEmailService)
	mockRoleService := new(MockRoleService)

	controller := &AuthController{
		authService:  mockAuthService,
		userService:  mockUserService,
		emailService: mockEmailService,
		roleService:  mockRoleService,
	}

	router := setupRouter(controller)

	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegister_MissingRequiredFields(t *testing.T) {
	mockAuthService := new(MockAuthService)
	mockUserService := new(MockUserService)
	mockEmailService := new(MockEmailService)
	mockRoleService := new(MockRoleService)

	controller := &AuthController{
		authService:  mockAuthService,
		userService:  mockUserService,
		emailService: mockEmailService,
		roleService:  mockRoleService,
	}

	router := setupRouter(controller)

	// Missing required fields - only name provided
	registerJSON := `{"name": "John Doe"}`

	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBufferString(registerJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Gin binding errors (missing required fields) return 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegister_InvalidEmail(t *testing.T) {
	mockAuthService := new(MockAuthService)
	mockUserService := new(MockUserService)
	mockEmailService := new(MockEmailService)
	mockRoleService := new(MockRoleService)

	controller := &AuthController{
		authService:  mockAuthService,
		userService:  mockUserService,
		emailService: mockEmailService,
		roleService:  mockRoleService,
	}

	router := setupRouter(controller)

	// Invalid email format - missing @ and domain
	registerJSON := `{
		"name": "John Doe",
		"username": "johndoe",
		"password": "password123",
		"confirmPassword": "password123",
		"email": "invalid-email",
		"phoneNumber": "081234567890",
		"dateOfBirth": "1999-08-25",
		"roleId": 1
	}`

	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBufferString(registerJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Gin binding errors (invalid email format) return 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	mockAuthService := new(MockAuthService)
	mockUserService := new(MockUserService)
	mockEmailService := new(MockEmailService)
	mockRoleService := new(MockRoleService)

	controller := &AuthController{
		authService:  mockAuthService,
		userService:  mockUserService,
		emailService: mockEmailService,
		roleService:  mockRoleService,
	}

	router := setupRouter(controller)

	registerJSON := `{
		"name": "John Doe",
		"username": "johndoe",
		"password": "password123",
		"confirmPassword": "password123",
		"email": "johndoe@example.com",
		"phoneNumber": "081234567890",
		"dateOfBirth": "25/08/1999",
		"roleId": 1
	}`

	mockAuthService.On("Register", mock.Anything, mock.AnythingOfType("*dto.RegisterRequest")).Return(nil, apperrors.ErrEmailExist)

	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBufferString(registerJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	mockAuthService.AssertExpectations(t)
}

func TestRegister_UsernameAlreadyExists(t *testing.T) {
	mockAuthService := new(MockAuthService)
	mockUserService := new(MockUserService)
	mockEmailService := new(MockEmailService)
	mockRoleService := new(MockRoleService)

	controller := &AuthController{
		authService:  mockAuthService,
		userService:  mockUserService,
		emailService: mockEmailService,
		roleService:  mockRoleService,
	}

	router := setupRouter(controller)

	registerJSON := `{
		"name": "John Doe",
		"username": "johndoe",
		"password": "password123",
		"confirmPassword": "password123",
		"email": "johndoe@example.com",
		"phoneNumber": "081234567890",
		"dateOfBirth": "25/08/1999",
		"roleId": 1
	}`

	mockAuthService.On("Register", mock.Anything, mock.AnythingOfType("*dto.RegisterRequest")).Return(nil, apperrors.ErrUsernameExist)

	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBufferString(registerJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	mockAuthService.AssertExpectations(t)
}

func TestRegister_PasswordMismatch(t *testing.T) {
	mockAuthService := new(MockAuthService)
	mockUserService := new(MockUserService)
	mockEmailService := new(MockEmailService)
	mockRoleService := new(MockRoleService)

	controller := &AuthController{
		authService:  mockAuthService,
		userService:  mockUserService,
		emailService: mockEmailService,
		roleService:  mockRoleService,
	}

	router := setupRouter(controller)

	registerJSON := `{
		"name": "John Doe",
		"username": "johndoe",
		"password": "password123",
		"confirmPassword": "password456",
		"email": "johndoe@example.com",
		"phoneNumber": "081234567890",
		"dateOfBirth": "1999-08-25",
		"roleId": 1
	}`

	mockAuthService.On("Register", mock.Anything, mock.AnythingOfType("*dto.RegisterRequest")).Return(nil, apperrors.ErrPasswordDoesNotMatch)

	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBufferString(registerJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// ErrPasswordDoesNotMatch returns 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	mockAuthService.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	mockAuthService := new(MockAuthService)
	mockUserService := new(MockUserService)
	mockEmailService := new(MockEmailService)
	mockRoleService := new(MockRoleService)

	controller := &AuthController{
		authService:  mockAuthService,
		userService:  mockUserService,
		emailService: mockEmailService,
		roleService:  mockRoleService,
	}

	router := setupRouter(controller)

	loginJSON := `{
		"email": "johndoe@example.com",
		"password": "password123"
	}`

	mockAuthService.On("Login", mock.Anything, mock.AnythingOfType("*dto.LoginInput")).Return(&dto.LoginResponse{
		AccessToken: "jwt-token",
		ExpiresIn:   1234567890,
		User: dto.GetUserResponse{
			UserID:      uuid.New(),
			Email:       "johndoe@example.com",
			Username:    "johndoe",
			PhoneNumber: "081234567890",
			RoleID:      1,
		},
	}, nil)

	mockRoleService.On("GetRole", mock.Anything, uint(1)).Return(&models.Role{
		ID:   1,
		Name: "user",
	}, nil)

	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBufferString(loginJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	assert.Equal(t, "Login successful", response["message"])

	mockAuthService.AssertExpectations(t)
	mockRoleService.AssertExpectations(t)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	mockAuthService := new(MockAuthService)
	mockUserService := new(MockUserService)
	mockEmailService := new(MockEmailService)
	mockRoleService := new(MockRoleService)

	controller := &AuthController{
		authService:  mockAuthService,
		userService:  mockUserService,
		emailService: mockEmailService,
		roleService:  mockRoleService,
	}

	router := setupRouter(controller)

	loginJSON := `{
		"email": "johndoe@example.com",
		"password": "wrongpassword"
	}`

	mockAuthService.On("Login", mock.Anything, mock.AnythingOfType("*dto.LoginInput")).Return(nil, apperrors.ErrInvalidCredentials)

	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBufferString(loginJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	mockAuthService.AssertExpectations(t)
}

func TestLogin_InvalidJSON(t *testing.T) {
	mockAuthService := new(MockAuthService)
	mockUserService := new(MockUserService)
	mockEmailService := new(MockEmailService)
	mockRoleService := new(MockRoleService)

	controller := &AuthController{
		authService:  mockAuthService,
		userService:  mockUserService,
		emailService: mockEmailService,
		roleService:  mockRoleService,
	}

	router := setupRouter(controller)

	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestVerifyOTP_Success(t *testing.T) {
	mockAuthService := new(MockAuthService)
	mockUserService := new(MockUserService)
	mockEmailService := new(MockEmailService)
	mockRoleService := new(MockRoleService)

	controller := &AuthController{
		authService:  mockAuthService,
		userService:  mockUserService,
		emailService: mockEmailService,
		roleService:  mockRoleService,
	}

	router := setupRouter(controller)

	verifyJSON := `{
		"email": "johndoe@example.com",
		"otp": "123456"
	}`

	mockAuthService.On("VerifyOTP", mock.Anything, "johndoe@example.com", "123456").Return(nil)

	req, _ := http.NewRequest("POST", "/auth/verify-otp", bytes.NewBufferString(verifyJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	assert.Equal(t, "Email verified successfully", response["message"])

	mockAuthService.AssertExpectations(t)
}

func TestVerifyOTP_InvalidOTP(t *testing.T) {
	mockAuthService := new(MockAuthService)
	mockUserService := new(MockUserService)
	mockEmailService := new(MockEmailService)
	mockRoleService := new(MockRoleService)

	controller := &AuthController{
		authService:  mockAuthService,
		userService:  mockUserService,
		emailService: mockEmailService,
		roleService:  mockRoleService,
	}

	router := setupRouter(controller)

	verifyJSON := `{
		"email": "johndoe@example.com",
		"otp": "000000"
	}`

	mockAuthService.On("VerifyOTP", mock.Anything, "johndoe@example.com", "000000").Return(apperrors.ErrInvalidOTP)

	req, _ := http.NewRequest("POST", "/auth/verify-otp", bytes.NewBufferString(verifyJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	mockAuthService.AssertExpectations(t)
}

func TestResendOTP_Success(t *testing.T) {
	mockAuthService := new(MockAuthService)
	mockUserService := new(MockUserService)
	mockEmailService := new(MockEmailService)
	mockRoleService := new(MockRoleService)

	controller := &AuthController{
		authService:  mockAuthService,
		userService:  mockUserService,
		emailService: mockEmailService,
		roleService:  mockRoleService,
	}

	router := setupRouter(controller)

	resendJSON := `{
		"email": "johndoe@example.com"
	}`

	mockAuthService.On("ResendOTP", mock.Anything, "johndoe@example.com").Return(nil)

	req, _ := http.NewRequest("POST", "/auth/resend-otp", bytes.NewBufferString(resendJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	assert.Equal(t, "OTP has been sent to your email", response["message"])

	mockAuthService.AssertExpectations(t)
}
