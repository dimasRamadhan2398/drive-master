package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
	"user-service/models/dto"
	"user-service/pkg/base"
	"user-service/pkg/config"
	"user-service/pkg/errors"
	"user-service/repositories"

	"user-service/pkg/logger"
	"user-service/pkg/redis"

	apperrors "user-service/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Login(ctx context.Context, req *dto.LoginInput) (*dto.LoginResponse, error)
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	HashPassword(password string) (string, error)

	// OTP methods
	GenerateAndSendOTP(ctx context.Context, email string) error
	VerifyOTP(ctx context.Context, email, otp string) error
	ResendOTP(ctx context.Context, email string) error
}

type AuthService struct {
	*base.BaseService
	userRepo     repositories.IUserRepository
	redisCli    *redis.Client
	emailService IMailtrapEmailService
}

func NewAuthService(userRepo repositories.IUserRepository, redisCli *redis.Client, emailService IMailtrapEmailService) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		redisCli:     redisCli,
		emailService: emailService,
	}
}

type Claims struct {
	User *dto.GetUserResponse
	jwt.RegisteredClaims
}

// Login authenticates a user with email and password
func (s *AuthService) Login(ctx context.Context, req *dto.LoginInput) (*dto.LoginResponse, error) {
	// Find user by email or username
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		// If not found by email, try by username
		user, err = s.userRepo.FindByUsername(ctx, req.Email)
		if err != nil {
			return nil, apperrors.ErrInvalidCredentials
		}
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	// Check if user is verified (optional - uncomment if required)
	// if !user.IsVerified {
	// 	return nil, apperrors.ErrEmailNotVerified
	// }

	// Set expiration time
	expirationTime := time.Now().Add(time.Duration(config.AppCfg.JWT.ExpiryHour) * time.Minute).Unix()

	// Build user response
	userResp := dto.GetUserResponse{
		UserID:      user.ID,
		Email:       user.Email,
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
		Image:       user.Image,
		DateOfBirth: user.DateOfBirth,
		Address:     user.Address,
		RoleID:      user.RoleID,
	}

	// Create JWT claims
	claims := &Claims{
		User: &userResp,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expirationTime, 0)),
		},
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.AppCfg.JWT.Secret))
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}

	return &dto.LoginResponse{
		User:         userResp,
		AccessToken:  tokenString,
		ExpiresIn:    expirationTime,
	}, nil
}

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	if s.isUsernameExist(ctx, req.Username) {
		return nil, apperrors.ErrUsernameExist
	}

	if s.isEmailExist(ctx, req.Email) {
		return nil, apperrors.ErrEmailExist
	}

	if req.Password != req.ConfirmPassword {
		return nil, apperrors.ErrPasswordDoesNotMatch
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	registerReq := &dto.RegisterRequest{
		Name:         req.Name,
		Username:     req.Username,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		DateOfBirth:  req.DateOfBirth,
		Password:     string(hashedPassword),
		RoleID:       req.RoleID,
	}

	user, err := s.userRepo.Create(ctx, registerReq)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := s.GenerateAndSendOTP(context.Background(), user.Email); err != nil {
			s.LogError("Failed to send OTP after registration", logger.LogField("error", err))
		}
	}()

	return &dto.RegisterResponse{
		User: dto.CreateUserResponse{
			UserID:      user.ID,
			Email:       user.Email,
			Username:    user.Username,
			PhoneNumber: user.PhoneNumber,
			RoleID:      user.RoleID,
		},
	}, nil
}

// ChangePassword changes a user's password
func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	return s.userRepo.Update(ctx, user)
}

// HashPassword hashes a password and returns the hash
func (s *AuthService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// GenerateAndSendOTP generates a 6-digit OTP, stores it in Redis, and sends it via email
func (s *AuthService) GenerateAndSendOTP(ctx context.Context, email string) error {
	// Verify email exists
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return apperrors.ErrUserNotFound
	}

	// Check if already verified
	if user.IsVerified {
		return apperrors.ErrAlreadyVerified
	}

	// Generate 6-digit OTP
	otp, err := s.generateOTP()
	if err != nil {
		return err
	}

	// Store OTP in Redis for 15 minutes with user ID key
	otpKey := fmt.Sprintf("email:otp:%s", user.ID.String())
	if err := s.redisCli.Client.Set(ctx, otpKey, otp, 15*time.Minute).Err(); err != nil {
		s.LogError("Failed to store OTP in redis", logger.LogField("error", err))
		return errors.ErrInternalServer
	}

	// Send OTP via email
	if err := s.emailService.SendOTPEmail(ctx, email, otp); err != nil {
		s.LogError("Failed to send OTP email", logger.LogField("error", err))
		return errors.ErrInternalServer
	}

	return nil
}

// VerifyOTP verifies the OTP and marks the user as verified
func (s *AuthService) VerifyOTP(ctx context.Context, email, otp string) error {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return apperrors.ErrUserNotFound
	}

	// Check if already verified
	if user.IsVerified {
		return apperrors.ErrAlreadyVerified
	}

	// Get stored OTP from Redis
	otpKey := fmt.Sprintf("email:otp:%s", user.ID.String())
	storedOtp, err := s.redisCli.Client.Get(ctx, otpKey).Result()
	if err != nil {
		return apperrors.ErrOTPExpired
	}

	// Verify OTP matches
	if storedOtp != otp {
		return apperrors.ErrInvalidOTP
	}

	// Mark user as verified
	user.IsVerified = true
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Delete OTP from Redis after successful verification
	s.redisCli.Client.Del(ctx, otpKey)

	return nil
}

// ResendOTP generates a new OTP and sends it via email
func (s *AuthService) ResendOTP(ctx context.Context, email string) error {
	// Verify email exists
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return apperrors.ErrUserNotFound
	}

	// Check if already verified
	if user.IsVerified {
		return apperrors.ErrAlreadyVerified
	}

	// Delete existing OTP if any
	otpKey := fmt.Sprintf("email:otp:%s", user.ID.String())
	s.redisCli.Client.Del(ctx, otpKey)

	// Generate new OTP
	otp, err := s.generateOTP()
	if err != nil {
		return err
	}

	// Store new OTP in Redis for 15 minutes
	if err := s.redisCli.Client.Set(ctx, otpKey, otp, 15*time.Minute).Err(); err != nil {
		s.LogError("Failed to store OTP in redis", logger.LogField("error", err))
		return errors.ErrInternalServer
	}

	// Send new OTP via email
	if err := s.emailService.SendOTPEmail(ctx, email, otp); err != nil {
		s.LogError("Failed to send OTP email", logger.LogField("error", err))
		return errors.ErrInternalServer
	}

	return nil
}

// generateOTP generates a cryptographically secure 6-digit OTP
func (s *AuthService) generateOTP() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", errors.ErrGenerateOTP
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func (s *AuthService) isUsernameExist(ctx context.Context, username string) bool {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return false
	}

	if user != nil {
		return true
	}

	return false
}

func (s *AuthService) isEmailExist(ctx context.Context, email string) bool {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return false
	}

	if user != nil {
		return true
	}

	return false
}
