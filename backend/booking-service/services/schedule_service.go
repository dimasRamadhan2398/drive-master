package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	cClient "booking-service/clients/core"
	uClient "booking-service/clients/user"
	"booking-service/models/dto"
	"booking-service/pkg/utils"
	"booking-service/repositories"

	"github.com/google/uuid"
)

type ScheduleService struct {
	scheduleRepo        repositories.IScheduleRepository
	enrollmentRepo      repositories.IEnrollmentRepository
	sessionRepo         repositories.ISessionRepository
	entitlementRepo     repositories.IEntitlementRepository
	sessionService      ISessionService
	availabilityService IAvailabilityService
	userClient          uClient.IUserClient
	coreClient          cClient.ICoreClient
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
	StartSession(ctx context.Context, slotID uint) (*dto.ScheduleResponse, error)
	CompleteSession(ctx context.Context, slotID uint) (*dto.ScheduleResponse, error)
	GetStats(ctx context.Context) (*dto.ScheduleStatsResponse, error)
}

func NewScheduleService(
	scheduleRepo repositories.IScheduleRepository,
	enrollmentRepo repositories.IEnrollmentRepository,
	sessionRepo repositories.ISessionRepository,
	entitlementRepo repositories.IEntitlementRepository,
	sessionService ISessionService,
	availabilityService IAvailabilityService,
	userClient uClient.IUserClient,
	coreClient cClient.ICoreClient,
) IScheduleService {
	return &ScheduleService{
		scheduleRepo:        scheduleRepo,
		enrollmentRepo:      enrollmentRepo,
		sessionRepo:         sessionRepo,
		entitlementRepo:     entitlementRepo,
		sessionService:      sessionService,
		availabilityService: availabilityService,
		userClient:          userClient,
		coreClient:          coreClient,
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

	return s.enrichSchedule(ctx, schedule)
}

func (s *ScheduleService) GetSchedule(ctx context.Context, id uint) (*dto.ScheduleResponse, error) {
	schedule, err := s.scheduleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.enrichSchedule(ctx, schedule)
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

	return s.enrichSchedule(ctx, schedule)
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

	total := int64(len(schedules))
	if limit <= 0 {
		limit = 10
	}

	start := (page - 1) * limit
	if start > int(total) {
		start = int(total)
	}
	end := start + limit
	if end > int(total) {
		end = int(total)
	}

	paginatedSchedules := schedules[start:end]

	// Enrich schedules with instructor, car, and user names
	enrichedSchedules, err := s.enrichSchedules(ctx, paginatedSchedules)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &dto.ScheduleListResponse{
		Data: enrichedSchedules,
		Pagination: dto.PaginationMeta{
			Page:       page,
			Total:      total,
			Limit:      limit,
			TotalPages: totalPages,
		},
	}, nil
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

	enriched, err := s.enrichSchedules(ctx, schedules)
	if err != nil {
		return nil, err
	}

	return &dto.ScheduleListResponse{
		Data: enriched,
		Pagination: dto.PaginationMeta{
			Page:       1,
			Total:      int64(len(schedules)),
			Limit:      len(schedules),
			TotalPages: 1,
		},
	}, nil
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

		if params.CarID != "" && sched.CarID.String() != params.CarID {
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
	limit := params.Limit
	if limit <= 0 {
		limit = 10
	}

	start := (params.Page - 1) * limit
	if start > int(total) {
		start = int(total)
	}
	end := start + limit
	if end > int(total) {
		end = int(total)
	}

	paginatedSchedules := schedules[start:end]

	enriched, err := s.enrichSchedules(ctx, paginatedSchedules)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &dto.ScheduleListResponse{
		Data: enriched,
		Pagination: dto.PaginationMeta{
			Page:       params.Page,
			Total:      total,
			Limit:      limit,
			TotalPages: totalPages,
		},
	}, nil
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

	// Check if entitlement exists and is active
	entitlement, err := s.entitlementRepo.FindByID(ctx, req.EntitlementID)
	if err != nil {
		return nil, errors.New("entitlement not found")
	}

	if entitlement.UsedSessions >= entitlement.TotalSessions {
		return nil, errors.New("entitlement has no sessions remaining")
	}

	if entitlement.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("entitlement has expired")
	}

	// Book the slot
	if err := s.scheduleRepo.BookSlot(ctx, slotID, req.UserID, entitlement.EnrollmentID); err != nil {
		return nil, err
	}

	// Create a driving session record using SessionService
	_, err = s.sessionService.CreateSession(ctx, dto.CreateSessionRequest{
		EnrollmentID:  entitlement.EnrollmentID,
		EntitlementID: entitlement.ID,
		UserID:        entitlement.UserID,
		InstructorID:  schedule.InstructorID,
		CarID:         schedule.CarID,
		ScheduleID:    &slotID,
		Date:          schedule.Date,
		Time:          schedule.Time,
		Duration:      schedule.Duration,
		Notes:         req.Notes,
	})
	if err != nil {
		log.Printf("Failed to create driving session: %v", err)
	}

	// Increment used sessions
	if err := s.entitlementRepo.UpdateUsedSessions(ctx, entitlement.ID, entitlement.UsedSessions+1); err != nil {
		log.Printf("Failed to increment used sessions: %v", err)
	}

	// Reload schedule
	schedule, err = s.scheduleRepo.FindByID(ctx, slotID)
	if err != nil {
		return nil, err
	}

	return s.enrichSchedule(ctx, schedule)
}

func (s *ScheduleService) CancelBooking(ctx context.Context, slotID uint) error {
	schedule, err := s.scheduleRepo.FindByID(ctx, slotID)
	if err != nil {
		return errors.New("schedule slot not found")
	}

	if schedule.Status != dto.ScheduleStatusBooked {
		return errors.New("schedule slot is not booked")
	}

	// Find the session and cancel it via SessionService
	session, err := s.sessionRepo.FindByScheduleID(ctx, slotID)
	if err == nil && session != nil {
		if _, err := s.sessionService.CancelSession(ctx, session.ID); err != nil {
			log.Printf("Warning: failed to cancel session: %v", err)
		}
	} else {
		// Fallback for direct schedule release if session not found
		return s.scheduleRepo.ReleaseSlot(ctx, slotID)
	}

	return nil
}

func (s *ScheduleService) StartSession(ctx context.Context, slotID uint) (*dto.ScheduleResponse, error) {
	session, err := s.sessionRepo.FindByScheduleID(ctx, slotID)
	if err != nil || session == nil {
		return nil, errors.New("booked session not found for this schedule")
	}

	if _, err := s.sessionService.StartSession(ctx, session.ID); err != nil {
		return nil, err
	}

	// Reload schedule
	schedule, err := s.scheduleRepo.FindByID(ctx, slotID)
	if err != nil {
		return nil, err
	}

	return s.enrichSchedule(ctx, schedule)
}

func (s *ScheduleService) CompleteSession(ctx context.Context, slotID uint) (*dto.ScheduleResponse, error) {
	session, err := s.sessionRepo.FindByScheduleID(ctx, slotID)
	if err != nil || session == nil {
		return nil, errors.New("booked session not found for this schedule")
	}

	if _, err := s.sessionService.CompleteSession(ctx, session.ID); err != nil {
		return nil, err
	}

	// Reload schedule
	schedule, err := s.scheduleRepo.FindByID(ctx, slotID)
	if err != nil {
		return nil, err
	}

	return s.enrichSchedule(ctx, schedule)
}

func (s *ScheduleService) GetStats(ctx context.Context) (*dto.ScheduleStatsResponse, error) {
	return s.scheduleRepo.GetStats(ctx)
}

// enrichSchedules fetches instructor, car, and student names concurrently
// and assembles the final enriched response list.
func (s *ScheduleService) enrichSchedules(ctx context.Context, schedules []dto.Schedule) ([]dto.ScheduleResponse, error) {
	// --- Collect unique IDs to minimize external calls ---
	instructorIDSet := make(map[string]struct{})
	userIDSet       := make(map[string]struct{})
	carIDSet        := make(map[uuid.UUID]struct{})

	for _, sched := range schedules {
		instructorIDSet[sched.InstructorID.String()] = struct{}{}
		carIDSet[sched.CarID] = struct{}{}
		if sched.UserID != nil {
			userIDSet[sched.UserID.String()] = struct{}{}
		}
	}

	allUserIDs := make(map[string]struct{})
	for id := range instructorIDSet { allUserIDs[id] = struct{}{} }
	for id := range userIDSet       { allUserIDs[id] = struct{}{} }
	carIDs := utils.SliceFromSet(carIDSet)

	// --- Fan-out: call user-service and core-service concurrently ---
	var (
		mu       sync.Mutex
		wg       sync.WaitGroup
		fetchErr error
		userMap  = make(map[string]uClient.UserInfo)
		carMap   = make(map[uuid.UUID]cClient.CarInfo)
	)

	for id := range allUserIDs {
		wg.Add(1)
		go func(userID string) {
			defer wg.Done()
			parsedID, err := uuid.Parse(userID)
			if err != nil {
				mu.Lock()
				fetchErr = fmt.Errorf("failed to parse user ID %s: %w", userID, err)
				mu.Unlock()
				return
			}
			info, err := s.userClient.GetUserByID(ctx, parsedID)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				fetchErr = fmt.Errorf("failed to fetch user %s: %w", userID, err)
				return
			}
			userMap[userID] = *info
		}(id)
	}

	for _, id := range carIDs {
		wg.Add(1)
		go func(carID uuid.UUID) {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()

			// Check for zero UUID or invalid car ID - use default fallback
			if carID == uuid.Nil {
				carMap[carID] = cClient.CarInfo{
					ID:    "898170a5-08db-4467-b33e-049660a4231c",
					Brand: "BYD",
					Model: "Atto 1",
				}
				return
			}

			info, err := s.coreClient.GetCarByID(ctx, carID)
			if err != nil {
				// If car not found (404) or other error, use default fallback
				log.Printf("Warning: failed to fetch car %s, using default fallback: %v", carID, err)
				carMap[carID] = cClient.CarInfo{
					ID:    "898170a5-08db-4467-b33e-049660a4231c",
					Brand: "BYD",
					Model: "Atto 1",
				}
				return
			}
			carMap[carID] = *info
		}(id)
	}

	wg.Wait()

	if fetchErr != nil {
		return nil, fetchErr
	}

	// --- Assemble enriched responses ---
	result := make([]dto.ScheduleResponse, 0, len(schedules))
	for _, sched := range schedules {
		resp := s.scheduleRepo.ToResponse(&sched)

		if u, ok := userMap[sched.InstructorID.String()]; ok {
			resp.InstructorName = u.FirstName + " " + u.LastName
		}
		if c, ok := carMap[sched.CarID]; ok {
			resp.CarName = c.Brand + " " + c.Model
		}
		if sched.UserID != nil {
			key := sched.UserID.String()
			if u, ok := userMap[key]; ok {
				name := u.FirstName + " " + u.LastName
				resp.UserName = &name
			}
		}

		result = append(result, resp)
	}

	return result, nil
}

// enrichSchedule enriches a single schedule response.
func (s *ScheduleService) enrichSchedule(ctx context.Context, schedule *dto.Schedule) (*dto.ScheduleResponse, error) {
	enriched, err := s.enrichSchedules(ctx, []dto.Schedule{*schedule})
	if err != nil {
		return nil, err
	}
	if len(enriched) == 0 {
		return nil, fmt.Errorf("failed to enrich schedule")
	}
	return &enriched[0], nil
}