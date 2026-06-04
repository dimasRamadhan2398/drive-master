package services

import (
	"context"
	"fmt"
	"time"
	"user-service/models"
	"user-service/models/dto"
	"user-service/pkg/base"
	"user-service/pkg/config"
	"user-service/pkg/errors"
	"user-service/pkg/logger"
	"user-service/pkg/redis"
	"user-service/repositories"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IInstructorService interface {
	GetInstructorProfile(ctx context.Context, userID uuid.UUID) (*dto.InstructorProfileResponse, error)
	CreateInstructorProfile(ctx context.Context, userID uuid.UUID) (*dto.InstructorProfileResponse, error)
	CreateInstructorProfileWithInput(ctx context.Context, userID uuid.UUID, req dto.InstructorProfileRequest) (*dto.InstructorProfileResponse, error)
	UpdateInstructorProfile(ctx context.Context, profile *models.InstructorProfile) error
	DeleteInstructorProfile(ctx context.Context, userID uuid.UUID) error
	CreateInstructorWithUser(ctx context.Context, req dto.CreateInstructorWithUserRequest) (*dto.CreateInstructorWithUserResponse, error)
}

type InstructorService struct {
	*base.BaseService
	instructorRepo repositories.IInstructorRepository
	userRepo       repositories.IUserRepository
	roleRepo       repositories.IRoleRepository
	emailService   IMailtrapEmailService
	redisCli       *redis.Client
}

func NewInstructorService(
	instructorRepo repositories.IInstructorRepository,
	userRepo repositories.IUserRepository,
	roleRepo repositories.IRoleRepository,
	emailService IMailtrapEmailService,
	redisCli *redis.Client,
) *InstructorService {
	return &InstructorService{
		instructorRepo: instructorRepo,
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		emailService:   emailService,
		redisCli:       redisCli,
	}
}

// GetInstructorProfile retrieves an instructor profile by user ID
func (s *InstructorService) GetInstructorProfile(ctx context.Context, userID uuid.UUID) (*dto.InstructorProfileResponse, error) {
	result, err := s.instructorRepo.FindInstructorProfileByUserID(ctx, userID);
	if err != nil {
		return nil, err
	}

	response := &dto.InstructorProfileResponse{
		UserID:            result.UserID,
		IsActive:          result.IsActive,
		NumberOfStudents:  result.NumberOfStudents,
		YearsOfExperience: result.YearsOfExperience,
		SessionsCompleted: result.SessionsCompleted,
		AverageRating:     result.AverageRating,
		Bio: 				result.Bio,
		LicenseNumber: result.LicenseNumber,
		LicenseExpiry: result.LicenseExpiry,
		PhotoURL: result.PhotoURL,
		BNSPCertificateNumber: result.BNSPCertificateNumber,
		Description: result.Description,

	}

	return response, nil
}

// CreateInstructorProfileWithInput creates a new instructor profile with provided data
func (s *InstructorService) CreateInstructorProfileWithInput(ctx context.Context, userID uuid.UUID, req dto.InstructorProfileRequest) (*dto.InstructorProfileResponse, error) {
	profile := &models.InstructorProfile{
		UserID:            userID,
		IsActive:          true,
		NumberOfStudents:   0,
		YearsOfExperience:  0,
		SessionsCompleted: 0,
		AverageRating:     0,
		PhotoURL:          "",
		CreatedAt:         time.Now(),
	}

	// Apply values from request if provided
	if req.LicenseNumber != nil {
		profile.LicenseNumber = *req.LicenseNumber
	}
	if req.LicenseExpiry.Unix() > 0 {
		profile.LicenseExpiry = req.LicenseExpiry
	}
	if req.BNSPCertificateNumber != nil {
		profile.BNSPCertificateNumber = *req.BNSPCertificateNumber
	}
	if req.Description != nil {
		profile.Bio = *req.Description
	}
	if req.YearsOfExperience != nil {
		profile.YearsOfExperience = *req.YearsOfExperience
	}
	if req.Specialization != nil {
		profile.Specialization = *req.Specialization
	}

	if err := s.instructorRepo.CreateInstructorProfile(ctx, profile); err != nil {
		return nil, err
	}

	return &dto.InstructorProfileResponse{
		UserID:               profile.UserID,
		BNSPCertificateNumber: profile.BNSPCertificateNumber,
		NumberOfStudents:     profile.NumberOfStudents,
		SessionsCompleted:    profile.SessionsCompleted,
		AverageRating:        profile.AverageRating,
		Description:         profile.Bio,
		LicenseNumber:        profile.LicenseNumber,
		YearsOfExperience:    profile.YearsOfExperience,
		LicenseExpiry:        profile.LicenseExpiry,
		IsActive:            profile.IsActive,
		PhotoURL:            profile.PhotoURL,
	}, nil
}

// CreateInstructorProfile creates a new instructor profile with default empty values
func (s *InstructorService) CreateInstructorProfile(ctx context.Context, userID uuid.UUID) (*dto.InstructorProfileResponse, error) {
	profile := &models.InstructorProfile{
		UserID:            userID,
		LicenseNumber:     "",
		LicenseExpiry:     time.Now(),
		BNSPCertificateNumber: "",
		Bio:               "",
		IsActive:          true,
		NumberOfStudents:  0,
		YearsOfExperience: 0,
		SessionsCompleted: 0,
		AverageRating:     0,
		PhotoURL:          "",
		CreatedAt:         time.Now(),
	}

	if err := s.instructorRepo.CreateInstructorProfile(ctx, profile); err != nil {
		return nil, err
	}
	return &dto.InstructorProfileResponse{
		UserID:            profile.UserID,
		BNSPCertificateNumber: profile.BNSPCertificateNumber,
		NumberOfStudents:  profile.NumberOfStudents,
		SessionsCompleted: profile.SessionsCompleted,
		AverageRating:     profile.AverageRating,
		Description:      profile.Bio,
		LicenseNumber:     profile.LicenseNumber,
		YearsOfExperience: profile.YearsOfExperience,
		LicenseExpiry:     profile.LicenseExpiry,
		IsActive:          profile.IsActive,
		PhotoURL:          profile.PhotoURL,
	
		Specialization: profile.Specialization,

	}, nil
}

// DeleteInstructorProfile deletes an instructor profile by user ID
func (s *InstructorService) DeleteInstructorProfile(ctx context.Context, userID uuid.UUID) error {
	return s.instructorRepo.DeleteInstructorProfile(ctx, userID)
}

// UpdateInstructorProfile updates an instructor profile
func (s *InstructorService) UpdateInstructorProfile(ctx context.Context, profile *models.InstructorProfile) error {
	return s.instructorRepo.UpdateInstructorProfile(ctx, profile)
}

// CreateInstructorWithUser creates both a user and an instructor profile in a single transaction
func (s *InstructorService) CreateInstructorWithUser(ctx context.Context, req dto.CreateInstructorWithUserRequest) (*dto.CreateInstructorWithUserResponse, error) {
	// Validate email uniqueness
	if s.isEmailExist(ctx, req.Email) {
		return nil, errors.ErrEmailExist
	}

	// Validate username uniqueness
	if s.isUsernameExist(ctx, req.Username) {
		return nil, errors.ErrUsernameExist
	}

	// Get instructor role ID
	instructorRoleID, err := s.getInstructorRoleID(ctx)
	if err != nil {
		return nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user and instructor profile in transaction
	db := s.userRepo.(*repositories.UserRepository).GetDB()

	var createdUser *models.User
	var createdProfile *models.InstructorProfile

	err = db.Transaction(func(tx *gorm.DB) error {
		// Create user
		t, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			t = time.Now()
		}
		userModel := models.User{
			FirstName:    req.FirstName,
			LastName:     req.LastName,
			Username:     req.Username,
			EmailAddress: req.Email,
			PhoneNumber:  req.PhoneNumber,
			PasswordHash: string(hashedPassword),
			DateOfBirth:  t,
			RoleID:       instructorRoleID,
			IsActive:     true,
		}

		if err := tx.Create(&userModel).Error; err != nil {
			return err
		}
		createdUser = &userModel

		// Create instructor profile
		profile := &models.InstructorProfile{
			UserID:            userModel.ID,
			IsActive:          true,
			NumberOfStudents:  0,
			YearsOfExperience: 0,
			SessionsCompleted: 0,
			AverageRating:     0,
			PhotoURL:          "",
			CreatedAt:         time.Now(),
		}

		// Apply optional fields
		if req.LicenseNumber != nil {
			profile.LicenseNumber = *req.LicenseNumber
		}
		if req.LicenseExpiry != nil {
			parsedTime, err := time.Parse("02/01/2006", *req.LicenseExpiry)
			if err == nil {
				profile.LicenseExpiry = parsedTime
			}
		}
		if req.BNSPCertificateNumber != nil {
			profile.BNSPCertificateNumber = *req.BNSPCertificateNumber
		}
		if req.Description != nil {
			profile.Bio = *req.Description
		}
		if req.YearsOfExperience != nil {
			profile.YearsOfExperience = *req.YearsOfExperience
		}
		if req.Specialization != nil {
			profile.Specialization = *req.Specialization
		}

		if err := tx.Create(profile).Error; err != nil {
			return err
		}
		createdProfile = profile

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Send OTP email asynchronously
	go func() {
		if err := s.sendOTPEmail(context.Background(), createdUser.EmailAddress); err != nil {
			s.LogError("Failed to send OTP after instructor registration", logger.LogField("error", err))
		}
	}()

	// Generate access token (unused but kept for future extension)
	cfg := config.Get()
	accessExpirationTime := time.Now().Add(time.Duration(cfg.JWT.ExpiryHour) * time.Hour).Unix()
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Unix(accessExpirationTime, 0)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	_, err = token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	// Generate refresh token
	refreshToken := generateRefreshToken()
	refreshKey := fmt.Sprintf("refresh_token:%s", refreshToken)
	expiryDuration := time.Duration(cfg.JWT.RefreshTokenExpiryDays) * 24 * time.Hour
	if err := s.redisCli.Client.Set(ctx, refreshKey, createdUser.ID.String(), expiryDuration).Err(); err != nil {
		return nil, errors.ErrInternalServer
	}

	// Build response
	return &dto.CreateInstructorWithUserResponse{
		User: dto.CreateUserResponse{
			UserID:      createdUser.ID,
			Email:       createdUser.EmailAddress,
			Username:    createdUser.Username,
			FirstName:   createdUser.FirstName,
			LastName:    createdUser.LastName,
			PhoneNumber: createdUser.PhoneNumber,
			RoleID:      createdUser.RoleID,
		},
		Profile: &dto.InstructorProfileResponse{
			UserID:               createdProfile.UserID,
			BNSPCertificateNumber: createdProfile.BNSPCertificateNumber,
			NumberOfStudents:      createdProfile.NumberOfStudents,
			SessionsCompleted:     createdProfile.SessionsCompleted,
			AverageRating:         createdProfile.AverageRating,
			Description:          createdProfile.Bio,
			LicenseNumber:         createdProfile.LicenseNumber,
			YearsOfExperience:     createdProfile.YearsOfExperience,
			LicenseExpiry:         createdProfile.LicenseExpiry,
			IsActive:             createdProfile.IsActive,
			PhotoURL:             createdProfile.PhotoURL,
			Specialization:       createdProfile.Specialization,
		},
	}, nil
}

// getInstructorRoleID retrieves the role ID for instructor
func (s *InstructorService) getInstructorRoleID(ctx context.Context) (uint, error) {
	roles, err := s.roleRepo.FindAllRoles(ctx)
	if err != nil {
		return 0, err
	}

	for _, role := range roles {
		if role.Name == "instructor" {
			return role.ID, nil
		}
	}

	return 0, errors.ErrInternalServer
}

// isEmailExist checks if email is already registered
func (s *InstructorService) isEmailExist(ctx context.Context, email string) bool {
	_, err := s.userRepo.FindByEmail(ctx, email)
	return err == nil
}

// isUsernameExist checks if username is already taken
func (s *InstructorService) isUsernameExist(ctx context.Context, username string) bool {
	_, err := s.userRepo.FindByUsername(ctx, username)
	return err == nil
}

// sendOTPEmail sends an OTP email to the user
func (s *InstructorService) sendOTPEmail(ctx context.Context, email string) error {
	return s.emailService.SendOTPEmail(ctx, email, "")
}