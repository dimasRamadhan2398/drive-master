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
	ValidateCredentials(ctx context.Context, emailOrUsername, password string) (*dto.LoginResponse, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	HashPassword(password string) (string, error)
}

type AuthService struct {
	*base.BaseService
	userRepo repositories.IUserRepository
	redisCli *redis.Client
}

func NewAuthService(userRepo repositories.IUserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if s.isUsernameExist(ctx, req.Username) {
		return nil, apperrors.ErrUsernameExist
	}

	if s.isEmailExist(ctx, req.Email) {
		return nil, apperrors.ErrEmailExist
	}

	if req.Password != req.ConfirmPassword {
		return nil, apperrors.ErrPasswordDoesNotMatch
	}

	registerReq := &dto.RegisterRequest{
		Name:        req.Name,
		Username:    req.Username,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    string(hashedPassword),
		RoleID:      req.RoleID,
	}

	user, err := s.userRepo.Create(ctx, registerReq)
	if err != nil {
		return nil, err
	}

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

// ValidateCredentials validates user credentials and returns the user
func (s *AuthService) ValidateCredentials(ctx context.Context, emailOrUsername, password string) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, emailOrUsername)
	if err != nil {
		// If not found by email, try by username
		user, err = s.userRepo.FindByUsername(ctx, emailOrUsername)
		if err != nil {
			return nil, err
		}
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		User: dto.GetUserResponse{
			UserID:      user.ID,
			Email:       user.Email,
			Username:    user.Username,
			PhoneNumber: user.PhoneNumber,
			Image:       user.Image,
			DateOfBirth: user.DateOfBirth,
			Address:     user.Address,
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

func (s *AuthService) generateDeviceOTP(ctx context.Context, userID uuid.UUID, email string) (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)

	if err != nil {
		return "", errors.ErrGenerateOTP
	}

	otpCode := fmt.Sprintf("%06d", n.Int64())

	// Store OTP in Redis for 15 minutes
	otpKey := fmt.Sprintf("device:otp:%s", userID.String())
	if err := s.redisCli.Client.Set(ctx, otpKey, otpCode, 15*time.Minute).Err(); err != nil {
		s.LogError("Failed to store OTP in redis", logger.LogField("error", err))
		return "", errors.ErrInternalServer
	}

	return otpCode, nil
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
