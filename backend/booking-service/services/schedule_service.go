package services

import (
	"context"
	"errors"
	"time"

	"booking-service/models/dto"
	"booking-service/repositories"

	"github.com/google/uuid"
)

type ScheduleService struct {
	scheduleRepo        IScheduleRepository
	enrollmentRepo      repositories.IEnrollmentRepository
	availabilityService IAvailabilityService
}

type IScheduleRepository interface {
	Create(ctx context.Context, schedule *dto.Schedule) error
	FindByID(ctx context.Context, id uint) (*dto.Schedule, error)
	Update(ctx context.Context, schedule *dto.Schedule) error
	Delete(ctx context.Context, schedule *dto.Schedule) error
	FindAll(ctx context.Context) ([]dto.Schedule, error)
	FindByDateAndInstructor(ctx context.Context, date time.Time, instructorID uuid.UUID) ([]dto.Schedule, error)
	FindByDateAndTime(ctx context.Context, date time.Time, time string, instructorID uuid.UUID, carID uint) (*dto.Schedule, error)
	FindAvailableByDateRange(ctx context.Context, startDate, endDate time.Time) ([]dto.Schedule, error)
	UpdateStatus(ctx context.Context, id uint, status dto.ScheduleStatus) error
	BookSlot(ctx context.Context, id uint, userID, enrollmentID uint) error
	ReleaseSlot(ctx context.Context, id uint) error
	ToResponse(schedule *dto.Schedule) dto.ScheduleResponse
	ToListResponse(schedules []dto.Schedule, total int64, page, limit int) dto.ScheduleListResponse
	ExistsForInstructorAndDateTime(ctx context.Context, instructorID uuid.UUID, date time.Time, timeStr string) (bool, error)
	CountAll(ctx context.Context) (int64, error)
}

type IScheduleService interface {
	CreateSchedule(ctx context.Context, req dto.CreateScheduleRequest) (*dto.ScheduleResponse, error)
	GetSchedule(ctx context.Context, id uint) (*dto.ScheduleResponse, error)
	UpdateSchedule(ctx context.Context, id uint, req dto.UpdateScheduleRequest) (*dto.ScheduleResponse, error)
	DeleteSchedule(ctx context.Context, id uint) error
	ListSchedules(ctx context.Context, page, limit int) (*dto.ScheduleListResponse, error)
	ListSchedulesFiltered(ctx context.Context, params dto.ScheduleFilterParams) (*dto.ScheduleListResponse, error)
	GetAvailableSchedules(ctx context.Context, startDate, endDate string) (*dto.ScheduleListResponse, error)
	BookSlot(ctx context.Context, slotID uint, req dto.BookSlotRequest) (*dto.ScheduleResponse, error)
	CancelBooking(ctx context.Context, slotID uint) error
}

func NewScheduleService(
	scheduleRepo IScheduleRepository,
	enrollmentRepo repositories.IEnrollmentRepository,
	availabilityService IAvailabilityService,
) IScheduleService {
	return &ScheduleService{
		scheduleRepo:        scheduleRepo,
		enrollmentRepo:      enrollmentRepo,
		availabilityService: availabilityService,
	}
}

func (s *ScheduleService) CreateSchedule(ctx context.Context, req dto.CreateScheduleRequest) (*dto.ScheduleResponse, error) {
	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format, expected YYYY-MM-DD")
	}

	duration := req.Duration
	if duration == 0 {
		duration = 60 // Default 60 minutes
	}

	// Check instructor availability against recurring schedules
	if err := s.availabilityService.CheckAvailability(ctx, req.InstructorID, parsedDate, req.Time, duration); err != nil {
		return nil, err
	}

	schedule := &dto.Schedule{
		Date:         parsedDate,
		Time:         req.Time,
		Duration:     duration,
		InstructorID: req.InstructorID,
		CarID:        req.CarID,
		Status:       dto.ScheduleStatusAvailable,
		Notes:        req.Notes,
	}

	if err := s.scheduleRepo.Create(ctx, schedule); err != nil {
		return nil, err
	}

	resp := s.scheduleRepo.ToResponse(schedule)
	return &resp, nil
}

func (s *ScheduleService) GetSchedule(ctx context.Context, id uint) (*dto.ScheduleResponse, error) {
	schedule, err := s.scheduleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := s.scheduleRepo.ToResponse(schedule)
	return &resp, nil
}

func (s *ScheduleService) UpdateSchedule(ctx context.Context, id uint, req dto.UpdateScheduleRequest) (*dto.ScheduleResponse, error) {
	schedule, err := s.scheduleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Date != nil {
		parsedDate, err := time.Parse("2006-01-02", *req.Date)
		if err != nil {
			return nil, errors.New("invalid date format, expected YYYY-MM-DD")
		}
		schedule.Date = parsedDate
	}
	if req.Time != nil {
		schedule.Time = *req.Time
	}
	if req.Duration != nil {
		schedule.Duration = *req.Duration
	}
	if req.InstructorID != nil {
		schedule.InstructorID = *req.InstructorID
	}
	if req.CarID != nil {
		schedule.CarID = *req.CarID
	}
	if req.Notes != nil {
		schedule.Notes = *req.Notes
	}
	if req.Status != nil {
		schedule.Status = dto.ScheduleStatus(*req.Status)
	}

	if err := s.scheduleRepo.Update(ctx, schedule); err != nil {
		return nil, err
	}

	resp := s.scheduleRepo.ToResponse(schedule)
	return &resp, nil
}

func (s *ScheduleService) DeleteSchedule(ctx context.Context, id uint) error {
	schedule, err := s.scheduleRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if schedule.Status != dto.ScheduleStatusAvailable {
		return errors.New("cannot delete a non-available schedule")
	}

	return s.scheduleRepo.Delete(ctx, schedule)
}

func (s *ScheduleService) ListSchedules(ctx context.Context, page, limit int) (*dto.ScheduleListResponse, error) {
	schedules, err := s.scheduleRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	total, err := s.scheduleRepo.CountAll(ctx)
	if err != nil {
		return nil, err
	}

	resp := s.scheduleRepo.ToListResponse(schedules, total, page, limit)
	return &resp, nil
}

func (s *ScheduleService) GetAvailableSchedules(ctx context.Context, startDate, endDate string) (*dto.ScheduleListResponse, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, errors.New("invalid start date format, expected YYYY-MM-DD")
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, errors.New("invalid end date format, expected YYYY-MM-DD")
	}

	schedules, err := s.scheduleRepo.FindAvailableByDateRange(ctx, start, end)
	if err != nil {
		return nil, err
	}

	// Use page=1 and limit=len(schedules) for the response
	resp := s.scheduleRepo.ToListResponse(schedules, int64(len(schedules)), 1, len(schedules))
	return &resp, nil
}

func (s *ScheduleService) ListSchedulesFiltered(ctx context.Context, params dto.ScheduleFilterParams) (*dto.ScheduleListResponse, error) {
	var schedules []dto.Schedule
	var err error

	// Get all schedules and filter manually (in a real app, you'd use DB queries)
	allSchedules, err := s.scheduleRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// Filter schedules based on params
	for _, sched := range allSchedules {
		match := true

		if params.Date != "" {
			parsedDate, _ := time.Parse("2006-01-02", params.Date)
			if !sched.Date.Equal(parsedDate) {
				match = false
			}
		}

		if params.StartDate != "" && params.EndDate != "" {
			startDate, _ := time.Parse("2006-01-02", params.StartDate)
			endDate, _ := time.Parse("2006-01-02", params.EndDate)
			if sched.Date.Before(startDate) || sched.Date.After(endDate) {
				match = false
			}
		}

		if params.InstructorID != "" && sched.InstructorID.String() != params.InstructorID {
			match = false
		}

		if params.CarID != 0 && sched.CarID != params.CarID {
			match = false
		}

		if params.Status != "" && string(sched.Status) != params.Status {
			match = false
		}

		if match {
			schedules = append(schedules, sched)
		}
	}

	total := int64(len(schedules))
	resp := s.scheduleRepo.ToListResponse(schedules, total, params.Page, params.Limit)
	return &resp, nil
}

func (s *ScheduleService) BookSlot(ctx context.Context, slotID uint, req dto.BookSlotRequest) (*dto.ScheduleResponse, error) {
	// Check if schedule exists and is available
	schedule, err := s.scheduleRepo.FindByID(ctx, slotID)
	if err != nil {
		return nil, errors.New("schedule slot not found")
	}

	if schedule.Status != dto.ScheduleStatusAvailable {
		return nil, errors.New("schedule slot is not available")
	}

	// Re-validate instructor availability (schedule was created with availability check,
	// but we check again to handle cases where schedule might have been created externally)
	if err := s.availabilityService.CheckAvailability(ctx, schedule.InstructorID, schedule.Date, schedule.Time, schedule.Duration); err != nil {
		return nil, err
	}

	// Check if enrollment exists
	_, err = s.enrollmentRepo.FindByID(ctx, req.EntitlementID)
	if err != nil {
		return nil, errors.New("enrollment not found")
	}

	// Book the slot
	if err := s.scheduleRepo.BookSlot(ctx, slotID, req.UserID, req.EntitlementID); err != nil {
		return nil, err
	}

	// Reload schedule
	schedule, err = s.scheduleRepo.FindByID(ctx, slotID)
	if err != nil {
		return nil, err
	}

	resp := s.scheduleRepo.ToResponse(schedule)
	return &resp, nil
}

func (s *ScheduleService) CancelBooking(ctx context.Context, slotID uint) error {
	schedule, err := s.scheduleRepo.FindByID(ctx, slotID)
	if err != nil {
		return errors.New("schedule slot not found")
	}

	if schedule.Status != dto.ScheduleStatusBooked {
		return errors.New("schedule slot is not booked")
	}

	return s.scheduleRepo.ReleaseSlot(ctx, slotID)
}