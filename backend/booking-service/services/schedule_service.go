package services

import (
	"context"
	"errors"
	"time"

	"booking-service/models"
	"booking-service/models/dto"
	"booking-service/repositories"

	"gorm.io/gorm"
)

type ScheduleService struct {
	scheduleRepo        *repositories.ScheduleRepository
	enrollmentRepo      *repositories.EnrollmentRepository
	availabilityService IAvailabilityService
}

func NewScheduleService(
	scheduleRepo *repositories.ScheduleRepository,
	enrollmentRepo *repositories.EnrollmentRepository,
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

	schedule := &models.Schedule{
		Date:         parsedDate,
		Time:         req.Time,
		Duration:     duration,
		InstructorID: req.InstructorID,
		CarID:        req.CarID,
		Status:       models.ScheduleStatusAvailable,
		Notes:        req.Notes,
	}

	if err := s.scheduleRepo.Create(ctx, schedule); err != nil {
		return nil, err
	}

	resp := s.scheduleRepo.ToResponse(schedule)
	return &resp, nil
}

func (s *ScheduleService) GetSchedule(ctx context.Context, id uint) (*dto.ScheduleResponse, error) {
	schedule, err := s.scheduleRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("schedule not found")
		}
		return nil, err
	}

	resp := s.scheduleRepo.ToResponse(schedule)
	return &resp, nil
}

func (s *ScheduleService) UpdateSchedule(ctx context.Context, id uint, req dto.UpdateScheduleRequest) (*dto.ScheduleResponse, error) {
	schedule, err := s.scheduleRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("schedule not found")
		}
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
		schedule.Status = models.ScheduleStatus(*req.Status)
	}

	if err := s.scheduleRepo.Update(ctx, schedule); err != nil {
		return nil, err
	}

	resp := s.scheduleRepo.ToResponse(schedule)
	return &resp, nil
}

func (s *ScheduleService) DeleteSchedule(ctx context.Context, id uint) error {
	schedule, err := s.scheduleRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("schedule not found")
		}
		return err
	}

	if schedule.Status != models.ScheduleStatusAvailable {
		return errors.New("cannot delete a non-available schedule")
	}

	return s.scheduleRepo.Delete(ctx, id)
}

func (s *ScheduleService) ListSchedules(ctx context.Context, page, limit int) (*dto.ScheduleListResponse, error) {
	schedules, total, err := s.scheduleRepo.List(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	resp := s.scheduleRepo.ToListResponse(schedules, total, page, limit)
	return &resp, nil
}

func (s *ScheduleService) ListSchedulesFiltered(ctx context.Context, params dto.ScheduleFilterParams) (*dto.ScheduleListResponse, error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Page <= 0 {
		params.Page = 1
	}

	schedules, total, err := s.scheduleRepo.ListFiltered(ctx, params)
	if err != nil {
		return nil, err
	}

	resp := s.scheduleRepo.ToListResponse(schedules, total, params.Page, params.Limit)
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

	schedules, err := s.scheduleRepo.GetAvailableByDateRange(ctx, start, end)
	if err != nil {
		return nil, err
	}

	// Use page=1 and limit=len(schedules) for the response
	resp := s.scheduleRepo.ToListResponse(schedules, int64(len(schedules)), 1, len(schedules))
	return &resp, nil
}

func (s *ScheduleService) BookSlot(ctx context.Context, slotID uint, req dto.BookSlotRequest) (*dto.ScheduleResponse, error) {
	// Check if schedule exists and is available
	schedule, err := s.scheduleRepo.GetByID(ctx, slotID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("schedule slot not found")
		}
		return nil, err
	}

	if schedule.Status != models.ScheduleStatusAvailable {
		return nil, errors.New("schedule slot is not available")
	}

	// Re-validate instructor availability (schedule was created with availability check,
	// but we check again to handle cases where schedule might have been created externally)
	if err := s.availabilityService.CheckAvailability(ctx, schedule.InstructorID, schedule.Date, schedule.Time, schedule.Duration); err != nil {
		return nil, err
	}

	// Check if enrollment exists
	_, err = s.enrollmentRepo.GetByID(ctx, req.EntitlementID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("enrollment not found")
		}
		return nil, err
	}

	// Book the slot
	if err := s.scheduleRepo.BookSlot(ctx, slotID, req.UserID, req.EntitlementID); err != nil {
		return nil, err
	}

	// Reload schedule
	schedule, err = s.scheduleRepo.GetByID(ctx, slotID)
	if err != nil {
		return nil, err
	}

	resp := s.scheduleRepo.ToResponse(schedule)
	return &resp, nil
}

func (s *ScheduleService) CancelBooking(ctx context.Context, slotID uint) error {
	schedule, err := s.scheduleRepo.GetByID(ctx, slotID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("schedule slot not found")
		}
		return err
	}

	if schedule.Status != models.ScheduleStatusBooked {
		return errors.New("schedule slot is not booked")
	}

	return s.scheduleRepo.ReleaseSlot(ctx, slotID)
}