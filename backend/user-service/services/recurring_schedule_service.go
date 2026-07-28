package services

import (
	"context"
	"errors"
	"user-service/models"
	"user-service/models/dto"
	"user-service/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IRecurringScheduleService interface {
	CreateRecurringSchedule(ctx context.Context, instructorID uuid.UUID, req dto.CreateRecurringScheduleRequest) (*dto.RecurringScheduleResponse, error)
	BulkCreateRecurringSchedules(ctx context.Context, instructorID uuid.UUID, req dto.BulkCreateRecurringScheduleRequest) ([]dto.RecurringScheduleResponse, error)
	GetRecurringSchedules(ctx context.Context, instructorID uuid.UUID) ([]dto.RecurringScheduleResponse, error)
	GetRecurringScheduleByID(ctx context.Context, id uuid.UUID) (*dto.RecurringScheduleResponse, error)
	UpdateRecurringSchedule(ctx context.Context, id uuid.UUID, req dto.UpdateRecurringScheduleRequest) (*dto.RecurringScheduleResponse, error)
	DeleteRecurringSchedule(ctx context.Context, id uuid.UUID) error
	DeleteAllRecurringSchedules(ctx context.Context, instructorID uuid.UUID) error
}

type RecurringScheduleService struct {
	recurringScheduleRepo repositories.IRecurringScheduleRepository
}

func NewRecurringScheduleService(
	recurringScheduleRepo repositories.IRecurringScheduleRepository,
) *RecurringScheduleService {
	return &RecurringScheduleService{
		recurringScheduleRepo: recurringScheduleRepo,
	}
}

// CreateRecurringSchedule creates a single recurring schedule
func (s *RecurringScheduleService) CreateRecurringSchedule(ctx context.Context, instructorID uuid.UUID, req dto.CreateRecurringScheduleRequest) (*dto.RecurringScheduleResponse, error) {
	// Validate day of week
	if req.DayOfWeek < 0 || req.DayOfWeek > 6 {
		return nil, errors.New("day_of_week must be between 0 (Sunday) and 6 (Saturday)")
	}

	// Validate time format
	if !isValidTimeFormat(req.StartTime) || !isValidTimeFormat(req.EndTime) {
		return nil, errors.New("invalid time format, expected HH:MM (e.g., 09:00)")
	}

	// Validate start time is before end time
	if req.StartTime >= req.EndTime {
		return nil, errors.New("start_time must be before end_time")
	}

	schedule := &models.InstructorRecurringSchedule{
		InstructorID: instructorID,
		DayOfWeek:    req.DayOfWeek,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		IsActive:     true,
	}

	if err := s.recurringScheduleRepo.Create(ctx, schedule); err != nil {
		return nil, err
	}

	return s.toResponse(schedule), nil
}

// BulkCreateRecurringSchedules creates multiple recurring schedules at once
// Example: Create schedules for Monday-Friday at 09:00-10:00 and 13:00-14:00
func (s *RecurringScheduleService) BulkCreateRecurringSchedules(ctx context.Context, instructorID uuid.UUID, req dto.BulkCreateRecurringScheduleRequest) ([]dto.RecurringScheduleResponse, error) {
	var createdSchedules []dto.RecurringScheduleResponse

	for _, slot := range req.Slots {
		// Validate day of week
		if slot.DayOfWeek < 0 || slot.DayOfWeek > 6 {
			return nil, errors.New("day_of_week must be between 0 (Sunday) and 6 (Saturday)")
		}

		// Validate time format
		if !isValidTimeFormat(slot.StartTime) || !isValidTimeFormat(slot.EndTime) {
			return nil, errors.New("invalid time format, expected HH:MM (e.g., 09:00)")
		}

		// Validate start time is before end time
		if slot.StartTime >= slot.EndTime {
			return nil, errors.New("start_time must be before end_time")
		}

		schedule := &models.InstructorRecurringSchedule{
			InstructorID: instructorID,
			DayOfWeek:    slot.DayOfWeek,
			StartTime:    slot.StartTime,
			EndTime:      slot.EndTime,
			IsActive:     true,
		}

		if err := s.recurringScheduleRepo.Create(ctx, schedule); err != nil {
			return nil, err
		}

		createdSchedules = append(createdSchedules, *s.toResponse(schedule))
	}

	return createdSchedules, nil
}

// GetRecurringSchedules retrieves all recurring schedules for an instructor
func (s *RecurringScheduleService) GetRecurringSchedules(ctx context.Context, instructorID uuid.UUID) ([]dto.RecurringScheduleResponse, error) {
	schedules, err := s.recurringScheduleRepo.GetByInstructorID(ctx, instructorID)
	if err != nil {
		return nil, err
	}

	var responses []dto.RecurringScheduleResponse
	for _, schedule := range schedules {
		responses = append(responses, *s.toResponse(&schedule))
	}

	return responses, nil
}

// GetRecurringScheduleByID retrieves a single recurring schedule by ID
func (s *RecurringScheduleService) GetRecurringScheduleByID(ctx context.Context, id uuid.UUID) (*dto.RecurringScheduleResponse, error) {
	schedule, err := s.recurringScheduleRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("recurring schedule not found")
		}
		return nil, err
	}

	return s.toResponse(schedule), nil
}

// UpdateRecurringSchedule updates a recurring schedule
func (s *RecurringScheduleService) UpdateRecurringSchedule(ctx context.Context, id uuid.UUID, req dto.UpdateRecurringScheduleRequest) (*dto.RecurringScheduleResponse, error) {
	schedule, err := s.recurringScheduleRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("recurring schedule not found")
		}
		return nil, err
	}

	if req.DayOfWeek != nil {
		if *req.DayOfWeek < 0 || *req.DayOfWeek > 6 {
			return nil, errors.New("day_of_week must be between 0 (Sunday) and 6 (Saturday)")
		}
		schedule.DayOfWeek = *req.DayOfWeek
	}

	if req.StartTime != nil {
		if !isValidTimeFormat(*req.StartTime) {
			return nil, errors.New("invalid start_time format, expected HH:MM")
		}
		schedule.StartTime = *req.StartTime
	}

	if req.EndTime != nil {
		if !isValidTimeFormat(*req.EndTime) {
			return nil, errors.New("invalid end_time format, expected HH:MM")
		}
		schedule.EndTime = *req.EndTime
	}

	if req.StartTime != nil && req.EndTime != nil {
		if *req.StartTime >= *req.EndTime {
			return nil, errors.New("start_time must be before end_time")
		}
	}

	if req.IsActive != nil {
		schedule.IsActive = *req.IsActive
	}

	if err := s.recurringScheduleRepo.Update(ctx, schedule); err != nil {
		return nil, err
	}

	return s.toResponse(schedule), nil
}

// DeleteRecurringSchedule deletes a single recurring schedule
func (s *RecurringScheduleService) DeleteRecurringSchedule(ctx context.Context, id uuid.UUID) error {
	return s.recurringScheduleRepo.Delete(ctx, id)
}

// DeleteAllRecurringSchedules deletes all recurring schedules for an instructor
func (s *RecurringScheduleService) DeleteAllRecurringSchedules(ctx context.Context, instructorID uuid.UUID) error {
	return s.recurringScheduleRepo.DeleteByInstructorID(ctx, instructorID)
}

// toResponse converts a model to a response DTO
func (s *RecurringScheduleService) toResponse(schedule *models.InstructorRecurringSchedule) *dto.RecurringScheduleResponse {
	return &dto.RecurringScheduleResponse{
		ID:           schedule.ID,
		InstructorID: schedule.InstructorID,
		DayOfWeek:    schedule.DayOfWeek,
		DayName:      models.DayOfWeekName(schedule.DayOfWeek),
		StartTime:    schedule.StartTime,
		EndTime:      schedule.EndTime,
		IsActive:     schedule.IsActive,
		CreatedAt:    schedule.CreatedAt,
		UpdatedAt:    schedule.UpdatedAt,
	}
}

// isValidTimeFormat validates HH:MM time format
func isValidTimeFormat(timeStr string) bool {
	if len(timeStr) != 5 {
		return false
	}
	if timeStr[2] != ':' {
		return false
	}

	// Check hour (00-23)
	hour := (int(timeStr[0]-'0') * 10) + int(timeStr[1]-'0')
	if hour < 0 || hour > 23 {
		return false
	}

	// Check minute (00-59)
	minute := (int(timeStr[3]-'0') * 10) + int(timeStr[4]-'0')
	if minute < 0 || minute > 59 {
		return false
	}

	return true
}
