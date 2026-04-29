package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
	"user-service/models"
	"user-service/pkg/base"
	"user-service/pkg/errors"
	"user-service/repositories"

	"user-service/pkg/logger"
	"user-service/pkg/redis"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Login(email string, password string) (*models.User, error)
	ValidateNewUser(user *models.User) error
	ValidateCredentials(emailOrUsername, password string) (*models.User, error)
	ChangePassword(userID uuid.UUID, oldPassword, newPassword string) error
	HashPassword(password string) (string, error)
}

type AuthService struct {
	*base.BaseService
	userRepo repositories.IUserRepository
	redisCli    *redis.Client
}

func NewAuthService(userRepo repositories.IUserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// Login authenticates a user with email and password
func (s *AuthService) Login(email string, password string) (*models.User, error) {
	return s.ValidateCredentials(email, password)
}

// ValidateNewUser validates a new user
func (s *AuthService) ValidateNewUser(user *models.User) error {
	// Add validation logic here if needed
	return nil
}

// ValidateCredentials validates user credentials and returns the user
func (s *AuthService) ValidateCredentials(emailOrUsername, password string) (*models.User, error) {
	var user *models.User
	var err error

	// Try to find by email first
	user, err = s.userRepo.FindByEmail(emailOrUsername)
	if err != nil {
		// If not found by email, try by username
		user, err = s.userRepo.FindByUsername(emailOrUsername)
		if err != nil {
			return nil, err
		}
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword changes a user's password
func (s *AuthService) ChangePassword(userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
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
	return s.userRepo.Update(user)
}

// HashPassword hashes a password and returns the hash
func (s *AuthService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// AuthenticateWithHashedPassword validates credentials where the user already has a hashed password
func (s *AuthService) AuthenticateWithHashedPassword(emailOrUsername, hashedPassword string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(emailOrUsername)
	if err != nil {
		user, err = s.userRepo.FindByUsername(emailOrUsername)
		if err != nil {
			return nil, err
		}
	}

	if user.PasswordHash != hashedPassword {
		return nil, err
	}

	return user, nil
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