package services

import (
	"context"
	"user-service/models"
	"user-service/models/dto"
	"user-service/repositories"

	"github.com/google/uuid"
)

type IWorkExperienceService interface {
	CreateWorkExperience(ctx context.Context, instructorID uuid.UUID, input dto.CreateWorkExperienceRequest) (*models.WorkExperience, error)
	UpdateWorkExperience(ctx context.Context, workExp *models.WorkExperience) error
	DeleteWorkExperience(ctx context.Context, id uint) error
	GetWorkExperiences(ctx context.Context, instructorID uuid.UUID) ([]models.WorkExperience, error)
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
func (s *WorkExperienceService) CreateWorkExperience(ctx context.Context, instructorID uuid.UUID, input dto.CreateWorkExperienceRequest) (*models.WorkExperience, error) {
	profile, err := s.instructorRepo.FindInstructorProfileByUserID(ctx, instructorID)
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

	if err := s.workExpRepo.Create(ctx, workExp); err != nil {
		return nil, err
	}

	return workExp, nil
}

// UpdateWorkExperience updates an existing work experience
func (s *WorkExperienceService) UpdateWorkExperience(ctx context.Context, workExp *models.WorkExperience) error {
	return s.workExpRepo.Update(ctx, workExp)
}

// DeleteWorkExperience deletes a work experience by ID
func (s *WorkExperienceService) DeleteWorkExperience(ctx context.Context, id uint) error {
	return s.workExpRepo.Delete(ctx, id)
}

// GetWorkExperiences retrieves all work experiences for an instructor
func (s *WorkExperienceService) GetWorkExperiences(ctx context.Context, instructorID uuid.UUID) ([]models.WorkExperience, error) {
	return s.workExpRepo.FindByInstructorID(ctx, instructorID)
}
