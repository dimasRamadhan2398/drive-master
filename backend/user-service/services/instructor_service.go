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

	// Work Experience
	CreateWorkExperience(instructorID uuid.UUID, input dto.CreateWorkExperienceRequest) (*models.WorkExperience, error)
	UpdateWorkExperience(workExp *models.WorkExperience) error
	DeleteWorkExperience(id uint) error
	GetWorkExperiences(instructorID uuid.UUID) ([]models.WorkExperience, error)

	// Coverage Areas
	AddCoverageArea(instructorID uuid.UUID, areaID uint) error
	RemoveCoverageArea(instructorID uuid.UUID, areaID uint) error
}

type InstructorService struct {
	instructorRepo repositories.IInstructorRepository
	workExpRepo    repositories.IWorkExperienceRepository
}

func NewInstructorService(
	instructorRepo repositories.IInstructorRepository,
	workExpRepo repositories.IWorkExperienceRepository,
) *InstructorService {
	return &InstructorService{
		instructorRepo: instructorRepo,
		workExpRepo:    workExpRepo,
	}
}

// GetInstructorProfile retrieves an instructor profile by user ID
func (s *InstructorService) GetInstructorProfile(userID uuid.UUID) (*dto.InstructorProfileResponse, error) {
	result, err := s.instructorRepo.FindInstructorProfileByUserID(userID);
	if err != nil {
		return nil, err
	}

	workExperiences := []dto.WorkExperienceResponse{}

	for _, workExp := range result.WorkExperiences {
		workExperiences = append(workExperiences, dto.WorkExperienceResponse{
			ID: workExp.ID,
			CompanyName: workExp.CompanyName,
			Role: workExp.Role,
			StartDate: workExp.StartDate,
			EndDate: workExp.EndDate,
			Description: workExp.Description,
		})
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
		WorkExperiences: workExperiences,
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

// CreateWorkExperience creates a new work experience for an instructor
func (s *InstructorService) CreateWorkExperience(instructorID uuid.UUID, input dto.CreateWorkExperienceRequest) (*models.WorkExperience, error) {
	workExp := &models.WorkExperience{
		InstructorID: instructorID,
		CompanyName:  input.CompanyName,
		Role:         input.Role,
		StartDate:    input.StartDate,
		EndDate:      input.EndDate,
		Description:  input.Description,
	}

	if err := s.workExpRepo.Create(workExp); err != nil {
		return nil, err
	}

	return workExp, nil
}

// UpdateWorkExperience updates an existing work experience
func (s *InstructorService) UpdateWorkExperience(workExp *models.WorkExperience) error {
	return s.workExpRepo.Update(workExp)
}

// DeleteWorkExperience deletes a work experience by ID
func (s *InstructorService) DeleteWorkExperience(id uint) error {
	return s.workExpRepo.Delete(id)
}

// GetWorkExperiences retrieves all work experiences for an instructor
func (s *InstructorService) GetWorkExperiences(instructorID uuid.UUID) ([]models.WorkExperience, error) {
	return s.workExpRepo.FindByInstructorID(instructorID)
}

// AddCoverageArea adds a coverage area to an instructor (placeholder implementation)
func (s *InstructorService) AddCoverageArea(instructorID uuid.UUID, areaID uint) error {
	
	return nil
}

// RemoveCoverageArea removes a coverage area from an instructor (placeholder implementation)
func (s *InstructorService) RemoveCoverageArea(instructorID uuid.UUID, areaID uint) error {
	// TODO: Implement coverage area functionality
	return nil
}