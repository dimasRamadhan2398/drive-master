package services

import (
	"time"
	"user-service/models"
	"user-service/models/dto"
	"user-service/repositories"

	"github.com/google/uuid"
)

type IInstructorService interface {
	GetInstructorProfile(userID uuid.UUID) (*dto.InstructorProfileResponse, error)
	CreateInstructorProfile(userID uuid.UUID) (*dto.InstructorProfileResponse, error)
	UpdateInstructorProfile(profile *models.InstructorProfile) error
	DeleteInstructorProfile(userID uuid.UUID) error
}

type InstructorService struct {
	instructorRepo repositories.IInstructorRepository
}

func NewInstructorService(
	instructorRepo repositories.IInstructorRepository,
) *InstructorService {
	return &InstructorService{
		instructorRepo: instructorRepo,
	}
}

// GetInstructorProfile retrieves an instructor profile by user ID
func (s *InstructorService) GetInstructorProfile(userID uuid.UUID) (*dto.InstructorProfileResponse, error) {
	result, err := s.instructorRepo.FindInstructorProfileByUserID(userID);
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
	}

	return response, nil
}

// CreateInstructorProfile creates a new instructor profile for a user
func (s *InstructorService) CreateInstructorProfile(userID uuid.UUID) (*dto.InstructorProfileResponse, error) {
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

	if err := s.instructorRepo.CreateInstructorProfile(profile); err != nil {
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
		Bio:               profile.Bio,
	}, nil
}

// DeleteInstructorProfile deletes an instructor profile by user ID
func (s *InstructorService) DeleteInstructorProfile(userID uuid.UUID) error {
	return s.instructorRepo.DeleteInstructorProfile(userID)
}

// UpdateInstructorProfile updates an instructor profile
func (s *InstructorService) UpdateInstructorProfile(profile *models.InstructorProfile) error {
	return s.instructorRepo.UpdateInstructorProfile(profile)
}