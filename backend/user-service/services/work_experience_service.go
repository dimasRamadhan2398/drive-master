package services

import (
	"user-service/models"
	"user-service/models/dto"
	"user-service/repositories"

	"github.com/google/uuid"
)

type IWorkExperienceService interface {
	CreateWorkExperience(instructorID uuid.UUID, input dto.CreateWorkExperienceRequest) (*models.WorkExperience, error)
	UpdateWorkExperience(workExp *models.WorkExperience) error
	DeleteWorkExperience(id uint) error
	GetWorkExperiences(instructorID uuid.UUID) ([]models.WorkExperience, error)
}

type WorkExperienceService struct {
	workExpRepo    repositories.IWorkExperienceRepository
	instructorRepo repositories.IInstructorRepository
}

func NewWorkExperienceService(
	workExpRepo repositories.IWorkExperienceRepository,
	instructorRepo repositories.IInstructorRepository,
) IWorkExperienceService {
	return &WorkExperienceService{
		workExpRepo:    workExpRepo,
		instructorRepo: instructorRepo,
	}
}

// CreateWorkExperience creates a new work experience for an instructor
func (s *WorkExperienceService) CreateWorkExperience(instructorID uuid.UUID, input dto.CreateWorkExperienceRequest) (*models.WorkExperience, error) {
	profile, err := s.instructorRepo.FindInstructorProfileByUserID(instructorID)
	if err != nil {
		return nil, err
	}

	workExp := &models.WorkExperience{
		InstructorID: profile.UserID,
		CompanyName:  input.CompanyName,
		Role:         input.Role,
		StartDate:    input.StartDate,
		EndDate:      input.EndDate,
		Description:  input.Description,
		IsVerified:   false, // always system-set
	}

	if err := s.workExpRepo.Create(workExp); err != nil {
		return nil, err
	}

	return workExp, nil
}

// UpdateWorkExperience updates an existing work experience
func (s *WorkExperienceService) UpdateWorkExperience(workExp *models.WorkExperience) error {
	return s.workExpRepo.Update(workExp)
}

// DeleteWorkExperience deletes a work experience by ID
func (s *WorkExperienceService) DeleteWorkExperience(id uint) error {
	return s.workExpRepo.Delete(id)
}

// GetWorkExperiences retrieves all work experiences for an instructor
func (s *WorkExperienceService) GetWorkExperiences(instructorID uuid.UUID) ([]models.WorkExperience, error) {
	return s.workExpRepo.FindByInstructorID(instructorID)
}
