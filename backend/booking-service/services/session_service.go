package services

import (
	cClient "booking-service/clients/core"
	uClient "booking-service/clients/user"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"

	"booking-service/models"
	"booking-service/models/dto"
	"booking-service/pkg/kafka"
	"booking-service/repositories"

	"gorm.io/gorm"
)

type ISessionService interface {
	CreateSession(ctx context.Context, req dto.CreateSessionRequest) (*dto.SessionResponse, error)
	GetSession(ctx context.Context, id uint) (*dto.SessionResponse, error)
	ListSessions(ctx context.Context, page, limit int) (*dto.SessionListResponse, error)
	GetStats(ctx context.Context) (*dto.SessionStatsResponse, error)
	StartSession(ctx context.Context, id uint) (*dto.SessionResponse, error)
	CompleteSession(ctx context.Context, id uint) (*dto.SessionResponse, error)
	AdminCompleteSession(ctx context.Context, id uint) (*dto.SessionResponse, error)
	CancelSession(ctx context.Context, id uint) (*dto.SessionResponse, error)
	ListUserSessions(ctx context.Context, userID uuid.UUID, page, limit int) (*dto.SessionListResponse, error)
	ListInstructorSessions(ctx context.Context, instructorID uuid.UUID, page, limit int) (*dto.SessionListResponse, error)
	RateSession(ctx context.Context, id uint, rating float64, feedback string) (*dto.SessionResponse, error)
	AutoStartScheduledSessions(ctx context.Context) error
	AutoCompleteOngoingSessions(ctx context.Context) error
}

type SessionService struct {
	sessionRepo     repositories.ISessionRepository
	scheduleRepo    repositories.IScheduleRepository
	enrollmentRepo  repositories.IEnrollmentRepository
	eventPublisher  kafka.IEventPublisher
	userClient      uClient.IUserClient
	coreClient      cClient.ICoreClient
	db              *gorm.DB
}

func NewSessionService(
	sessionRepo repositories.ISessionRepository,
	scheduleRepo repositories.IScheduleRepository,
	userClient uClient.IUserClient,
) ISessionService {
	return &SessionService{
		sessionRepo:     sessionRepo,
		scheduleRepo:    scheduleRepo,
		userClient:      userClient,
	}
}

// NewSessionServiceWithEventPublisher creates a new session service with event publisher
func NewSessionServiceWithEventPublisher(
	sessionRepo repositories.ISessionRepository,
	scheduleRepo repositories.IScheduleRepository,
	eventPublisher kafka.IEventPublisher,
	userClient uClient.IUserClient,
) ISessionService {
	return &SessionService{
		sessionRepo:     sessionRepo,
		scheduleRepo:    scheduleRepo,
		eventPublisher:  eventPublisher,
		userClient:      userClient,
	}
}

// NewSessionServiceWithAllDeps creates a session service with all dependencies including enrollment repo and DB
func NewSessionServiceWithAllDeps(
	sessionRepo repositories.ISessionRepository,
	scheduleRepo repositories.IScheduleRepository,
	enrollmentRepo repositories.IEnrollmentRepository,
	eventPublisher kafka.IEventPublisher,
	userClient uClient.IUserClient,
	coreClient cClient.ICoreClient,
	db *gorm.DB,
) ISessionService {
	return &SessionService{
		sessionRepo:     sessionRepo,
		scheduleRepo:    scheduleRepo,
		enrollmentRepo:  enrollmentRepo,
		eventPublisher:  eventPublisher,
		userClient:      userClient,
		coreClient:      coreClient,
		db:              db,
	}
}

func (s *SessionService) CreateSession(ctx context.Context, req dto.CreateSessionRequest) (*dto.SessionResponse, error) {
	session := &models.Session{
		UserID:        req.UserID,
		InstructorID:  req.InstructorID,
		EntitlementID: req.EntitlementID,
		EnrollmentID:  req.EnrollmentID,
		ScheduleID:		req.ScheduleID,
		Date:          req.Date,
		Time:          req.Time,
		Duration:      req.Duration,
		CarID:         req.CarID,
		Area:          req.Area,
		Notes:         req.Notes,
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	resp := s.sessionRepo.ToResponse(session)
	return &resp, nil
}

func (s *SessionService) GetSession(ctx context.Context, id uint) (*dto.SessionResponse, error) {
	session, err := s.sessionRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("session not found")
		}
		return nil, err
	}

	resp := s.sessionRepo.ToResponse(session)
	return &resp, nil
}

func (s *SessionService) ListSessions(ctx context.Context, page, limit int) (*dto.SessionListResponse, error) {
	sessions, err := s.sessionRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	total, err := s.sessionRepo.CountAll(ctx)
	if err != nil {
		return nil, err
	}

	resp := s.sessionRepo.ToListResponse(sessions, total, page, limit)
	return &resp, nil
}

func (s *SessionService) ListUserSessions(ctx context.Context, userID uuid.UUID, page, limit int) (*dto.SessionListResponse, error) {
	sessions, err := s.sessionRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	total := int64(len(sessions))
	resp := s.sessionRepo.ToListResponse(sessions, total, page, limit)

	// Enrich with car name and instructor name if coreClient is available
	if s.coreClient != nil && len(resp.Data) > 0 {
		// Collect unique car IDs
		carIDs := make(map[uuid.UUID]struct{})
		for _, item := range resp.Data {
			if item.CarID != uuid.Nil {
				carIDs[item.CarID] = struct{}{}
			}
		}

		// Fetch car info concurrently
		carMap := make(map[uuid.UUID]string)
		var mu sync.Mutex
		var wg sync.WaitGroup
		for carID := range carIDs {
			wg.Add(1)
			go func(id uuid.UUID) {
				defer wg.Done()
				info, ferr := s.coreClient.GetCarByID(ctx, id)
				if ferr == nil && info != nil {
					mu.Lock()
					carMap[id] = info.Brand + " " + info.Model
					mu.Unlock()
				}
			}(carID)
		}
		wg.Wait()

		// Apply enriched names
		for i := range resp.Data {
			if name, ok := carMap[resp.Data[i].CarID]; ok {
				resp.Data[i].CarName = name
			}
		}
	}

	return &resp, nil
}

func (s *SessionService) ListInstructorSessions(ctx context.Context, instructorID uuid.UUID, page, limit int) (*dto.SessionListResponse, error) {
	sessions, err := s.sessionRepo.FindByInstructorID(ctx, instructorID)
	if err != nil {
		return nil, err
	}

	total := int64(len(sessions))

	resp := s.sessionRepo.ToListResponse(sessions, total, page, limit)
	return &resp, nil
}

func (s *SessionService) GetStats(ctx context.Context) (*dto.SessionStatsResponse, error) {
	stats, err := s.sessionRepo.GetStats(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.SessionStatsResponse{
		TotalSessions:     stats.Total,
		ActiveSessions:    stats.Active,
		CompletedSessions: stats.Completed,
		PendingSessions:   stats.Pending,
	}, nil
}

func (s *SessionService) StartSession(ctx context.Context, id uint) (*dto.SessionResponse, error) {
	// Check if session exists
	session, err := s.sessionRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("session not found")
		}
		return nil, err
	}

	// Validate session can be started
	if session.Status != "scheduled" {
		return nil, errors.New("session cannot be started: invalid status")
	}

	// Start the session
	startedAt := time.Now()
	if err := s.sessionRepo.StartSession(ctx, id, startedAt); err != nil {
		return nil, err
	}

	// Update associated schedule if any
	if session.ScheduleID != nil {
		if err := s.scheduleRepo.UpdateStatus(ctx, *session.ScheduleID, dto.ScheduleStatusInProgress); err != nil {
			// Log error but continue
		}
	}

	// Reload session
	session, err = s.sessionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := s.sessionRepo.ToResponse(session)
	return &resp, nil
}

func (s *SessionService) CompleteSession(ctx context.Context, id uint) (*dto.SessionResponse, error) {
	// Check if session exists
	session, err := s.sessionRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("session not found")
		}
		return nil, err
	}

	// Validate session can be completed
	if session.Status != "in_progress" {
		return nil, errors.New("session cannot be completed: must be in progress")
	}

	// Update associated schedule if any (do this before completing session)
	if session.ScheduleID != nil {
		if err := s.scheduleRepo.UpdateStatus(ctx, *session.ScheduleID, dto.ScheduleStatusCompleted); err != nil {
			// Log error but continue
		}
	}

	var allEntitlements []uClient.EntitlementInfo
	var getEntErr error
	if session.EntitlementID != uuid.Nil {
		allEntitlements, getEntErr = s.userClient.GetMemberEntitlements(ctx, session.UserID)
	}

	// If we have enrollment repo and db, use transaction for atomic updates
	if s.enrollmentRepo != nil && s.db != nil {
		err = s.db.Transaction(func(tx *gorm.DB) error {
			// Complete the session
			completedAt := time.Now()
			if err := s.sessionRepo.CompleteSessionTx(tx, id, completedAt); err != nil {
				return err
			}

			// Check if all sessions are used → update enrollment to completed
			if session.EnrollmentID != uuid.Nil && getEntErr == nil && len(allEntitlements) > 0 {
				enrollment, err := s.enrollmentRepo.FindByID(ctx, session.EnrollmentID)
				if err == nil && enrollment != nil {
					// Sum total sessions across all entitlements for this enrollment
					totalEntSessionLimit := 0
					hasEntitlementForEnrollment := false
					for _, ent := range allEntitlements {
						if ent.BookingID == session.EnrollmentID {
							totalEntSessionLimit += ent.TotalSessions
							hasEntitlementForEnrollment = true
						}
					}

					if hasEntitlementForEnrollment {
						// Retrieve all sessions for this enrollment in booking-service
						enrollmentSessions, err := s.sessionRepo.FindByEnrollmentID(ctx, session.EnrollmentID)
						if err == nil {
							completedCount := 0
							for _, s := range enrollmentSessions {
								// Include completed sessions and the current session we are completing
								if s.Status == "completed" || s.ID == id {
									completedCount++
								}
							}

							// If all sessions are completed, mark enrollment as completed
							if completedCount >= totalEntSessionLimit {
								if enrollment.Status != models.EnrollmentStatusCompleted {
									if err := s.enrollmentRepo.UpdateStatusTx(tx, enrollment.ID, models.EnrollmentStatusCompleted); err != nil {
										return err
									}
								}
							} else if enrollment.Status == models.EnrollmentStatusPendingPayment || enrollment.Status == models.EnrollmentStatusPaid {
								// Mark as in_progress when session is completed if it wasn't already
								if err := s.enrollmentRepo.UpdateStatusTx(tx, enrollment.ID, models.EnrollmentStatusInProgress); err != nil {
									return err
								}
							}
						}
					}
				}
			}

			return nil
		})

		if err != nil {
			return nil, err
		}
	} else {
		// Fallback: non-transactional update (backward compatibility)
		completedAt := time.Now()
		if err := s.sessionRepo.CompleteSession(ctx, id, completedAt); err != nil {
			return nil, err
		}
	}

	// Synchronously call user-service via HTTP to decrement remaining sessions & increment usedSessions
	if session.EntitlementID != uuid.Nil && session.UserID != uuid.Nil && s.userClient != nil {
		if err := s.userClient.UseSession(ctx, session.UserID, session.EntitlementID); err != nil {
			log.Printf("Warning: failed to decrement entitlement in user-service via direct HTTP: %v", err)
		}
	}

	// Publish Kafka event for session completion (for external services like core-service)
	if s.eventPublisher != nil {
		var packageID uuid.UUID
		if session.EnrollmentID != uuid.Nil {
			enrollment, err := s.enrollmentRepo.FindByID(ctx, session.EnrollmentID)
			if err == nil && enrollment != nil {
				packageID = enrollment.PackageID
			}
		}

		// Get updated entitlement for remaining sessions count
		var sessionsRemaining int
		if session.EntitlementID != uuid.Nil {
			entitlement, err := s.userClient.GetEntitlement(ctx, session.UserID, session.EntitlementID)
			if err == nil && entitlement != nil {
				sessionsRemaining = entitlement.RemainingSessions - 1
				if sessionsRemaining < 0 {
					sessionsRemaining = 0
				}
			}
		}

		_ = s.eventPublisher.PublishSessionCompletedWithEnrollment(
			ctx,
			id,
			session.EntitlementID,
			session.EnrollmentID,
			session.UserID,
			packageID,
			1,
			sessionsRemaining,
		)
	}

	// Reload session
	session, err = s.sessionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := s.sessionRepo.ToResponse(session)
	return &resp, nil
}

// AdminCompleteSession force-completes a session on behalf of the admin.
// Unlike CompleteSession it does NOT require the session to be in_progress,
// sets is_ended_by_admin=true and end_time so the auto-scheduler never reverts it.
func (s *SessionService) AdminCompleteSession(ctx context.Context, id uint) (*dto.SessionResponse, error) {
	// Load session
	session, err := s.sessionRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("session not found")
		}
		return nil, err
	}

	if session.Status == "cancelled" {
		return nil, errors.New("cannot complete a cancelled session")
	}
	if session.Status == "completed" && session.IsEndedByAdmin {
		// Already admin-completed, nothing to do – return current state
		resp := s.sessionRepo.ToResponse(session)
		return &resp, nil
	}

	// Update associated schedule to completed
	if session.ScheduleID != nil {
		_ = s.scheduleRepo.UpdateStatus(ctx, *session.ScheduleID, dto.ScheduleStatusCompleted)
	}

	endTime := time.Now()
	if err := s.sessionRepo.ForceCompleteByAdmin(ctx, id, endTime); err != nil {
		return nil, err
	}

	// Synchronously call user-service via HTTP to decrement remaining sessions & increment usedSessions
	if session.EntitlementID != uuid.Nil && session.UserID != uuid.Nil && s.userClient != nil {
		if err := s.userClient.UseSession(ctx, session.UserID, session.EntitlementID); err != nil {
			log.Printf("Warning: failed to decrement entitlement in user-service via direct HTTP: %v", err)
		}
	}

	// Publish Kafka event so external services get notified
	if s.eventPublisher != nil {
		var packageID uuid.UUID
		if session.EnrollmentID != uuid.Nil && s.enrollmentRepo != nil {
			enrollment, ferr := s.enrollmentRepo.FindByID(ctx, session.EnrollmentID)
			if ferr == nil && enrollment != nil {
				packageID = enrollment.PackageID
			}
		}
		_ = s.eventPublisher.PublishSessionCompletedWithEnrollment(
			ctx,
			id,
			session.EntitlementID,
			session.EnrollmentID,
			session.UserID,
			packageID,
			1,
			0,
		)
	}

	// Reload and return
	session, err = s.sessionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := s.sessionRepo.ToResponse(session)
	return &resp, nil
}

// RateSession rates a completed driving session and updates the instructor's average score
func (s *SessionService) RateSession(ctx context.Context, id uint, rating float64, feedback string) (*dto.SessionResponse, error) {
	session, err := s.sessionRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("session not found")
		}
		return nil, err
	}

	if session.Status != "completed" {
		return nil, errors.New("cannot rate a session that is not completed")
	}

	if err := s.sessionRepo.RateSession(ctx, id, rating, feedback); err != nil {
		return nil, err
	}

	// Update instructor rating via userClient
	if s.userClient != nil && session.InstructorID != uuid.Nil {
		if err := s.userClient.RateInstructor(ctx, session.InstructorID, rating); err != nil {
			log.Printf("Warning: failed to update instructor rating in user-service: %v", err)
		}
	}

	session, err = s.sessionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := s.sessionRepo.ToResponse(session)
	return &resp, nil
}

func (s *SessionService) CancelSession(ctx context.Context, id uint) (*dto.SessionResponse, error) {
	// Check if session exists
	session, err := s.sessionRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("session not found")
		}
		return nil, err
	}

	// Validate session can be cancelled (not already completed or cancelled)
	if session.Status == "completed" {
		return nil, errors.New("session cannot be cancelled: already completed")
	}
	if session.Status == "cancelled" {
		return nil, errors.New("session cannot be cancelled: already cancelled")
	}

	// Cancel the session
	if err := s.sessionRepo.CancelSession(ctx, id); err != nil {
		return nil, err
	}

	// Update associated schedule if any
	if session.ScheduleID != nil {
		if err := s.scheduleRepo.ReleaseSlot(ctx, *session.ScheduleID); err != nil {
			// Log error
		}
	}



	// Reload session
	session, err = s.sessionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := s.sessionRepo.ToResponse(session)
	return &resp, nil
}

func parseDurationMinutes(duration int) time.Duration {
	if duration <= 0 {
		return 60 * time.Minute
	}
	if duration <= 10 {
		return time.Duration(duration * 60) * time.Minute
	}
	return time.Duration(duration) * time.Minute
}

func (s *SessionService) AutoStartScheduledSessions(ctx context.Context) error {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.Local
	}
	now := time.Now().In(loc)

	// 1. Process scheduled sessions
	sessions, err := s.sessionRepo.FindByStatus(ctx, "scheduled")
	if err == nil {
		for _, session := range sessions {
			dateStr := fmt.Sprintf("%s %s", session.Date.Format("2006-01-02"), session.Time)
			startTime, err := time.ParseInLocation("2006-01-02 15:04", dateStr, loc)
			if err != nil {
				startTime, _ = time.Parse("2006-01-02 15:04", dateStr)
			}

			dur := parseDurationMinutes(session.Duration)
			endTime := startTime.Add(dur)

			if now.After(startTime) || now.Equal(startTime) {
				if now.Before(endTime) {
					log.Printf("[SessionMonitor] Session ID %d of user %s reached start time %s. Auto-starting...",
						session.ID, session.UserID, startTime.Format("15:04"))

					_, startErr := s.StartSession(ctx, session.ID)
					if startErr != nil {
						log.Printf("[SessionMonitor] Failed to auto-start session ID %d: %v", session.ID, startErr)
					} else {
						log.Printf("[SessionMonitor] Successfully auto-started session ID %d", session.ID)
					}
				} else {
					log.Printf("[SessionMonitor] Session ID %d past end time %s. Starting and completing...",
						session.ID, endTime.Format("15:04"))
					_, _ = s.StartSession(ctx, session.ID)
					_, _ = s.CompleteSession(ctx, session.ID)
				}
			}
		}
	}

	// 2. Process all schedules for status updates based on time and student booking
	schedules, err := s.scheduleRepo.FindAll(ctx)
	if err == nil {
		for _, sched := range schedules {
			dateStr := fmt.Sprintf("%s %s", sched.Date.Format("2006-01-02"), sched.Time)
			startTime, err := time.ParseInLocation("2006-01-02 15:04", dateStr, loc)
			if err != nil {
				startTime, _ = time.Parse("2006-01-02 15:04", dateStr)
			}
			dur := parseDurationMinutes(sched.Duration)
			endTime := startTime.Add(dur)

			hasStudent := sched.UserID != nil

			// Case 1: Current time is WITHIN the slot window (e.g., 11:00-12:00 at 11:01 AM)
			if (now.After(startTime) || now.Equal(startTime)) && now.Before(endTime) {
				if hasStudent {
					// Student booked: set to in-progress
					sess, sessErr := s.sessionRepo.FindByScheduleID(ctx, sched.ID)
					if sessErr == nil && sess != nil {
						// Never revert a session that was explicitly ended by an admin
						if sess.IsEndedByAdmin {
							continue
						}
						if sess.Status == "scheduled" {
							_, _ = s.StartSession(ctx, sess.ID)
						} else if sess.Status == "completed" {
							// Only revert to in_progress if NOT admin-ended
							if !sess.IsEndedByAdmin {
								_ = s.sessionRepo.UpdateStatus(ctx, sess.ID, "in_progress")
							}
						}
					}
					if sched.Status != dto.ScheduleStatusInProgress {
						_ = s.scheduleRepo.UpdateStatus(ctx, sched.ID, dto.ScheduleStatusInProgress)
						log.Printf("[ScheduleMonitor] Slot ID %d (has student) started -> in-progress", sched.ID)
					}
				} else {
					// No student booked: mark as blocked (passed/expired unbooked slot)
					if sched.Status == dto.ScheduleStatusAvailable {
						_ = s.scheduleRepo.UpdateStatus(ctx, sched.ID, dto.ScheduleStatusBlocked)
						log.Printf("[ScheduleMonitor] Slot ID %d (no student) started -> blocked (passed)", sched.ID)
					}
				}
			} else if now.After(endTime) || now.Equal(endTime) {
				// Case 2: Current time is AFTER end time (e.g. 12:00 PM or later)
				if hasStudent {
					// Student booked: complete the session & schedule
					sess, sessErr := s.sessionRepo.FindByScheduleID(ctx, sched.ID)
					if sessErr == nil && sess != nil {
						if sess.Status != "completed" {
							_, _ = s.CompleteSession(ctx, sess.ID)
						}
					}
					if sched.Status != dto.ScheduleStatusCompleted {
						_ = s.scheduleRepo.UpdateStatus(ctx, sched.ID, dto.ScheduleStatusCompleted)
						log.Printf("[ScheduleMonitor] Slot ID %d (has student) reached end time -> completed", sched.ID)
					}
				} else {
					// No student booked: mark as blocked (passed/expired unbooked slot)
					if sched.Status == dto.ScheduleStatusAvailable {
						_ = s.scheduleRepo.UpdateStatus(ctx, sched.ID, dto.ScheduleStatusBlocked)
						log.Printf("[ScheduleMonitor] Slot ID %d (no student) reached end time -> blocked (passed)", sched.ID)
					}
				}
			}
		}
	}

	return nil
}

func (s *SessionService) AutoCompleteOngoingSessions(ctx context.Context) error {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.Local
	}
	now := time.Now().In(loc)

	// 1. Process in-progress sessions
	sessions, err := s.sessionRepo.FindByStatus(ctx, "in_progress")
	if err == nil {
		for _, session := range sessions {
			dateStr := fmt.Sprintf("%s %s", session.Date.Format("2006-01-02"), session.Time)
			startTime, err := time.ParseInLocation("2006-01-02 15:04", dateStr, loc)
			if err != nil {
				startTime, _ = time.Parse("2006-01-02 15:04", dateStr)
			}

			dur := parseDurationMinutes(session.Duration)
			endTime := startTime.Add(dur)

			if now.After(endTime) || now.Equal(endTime) {
				log.Printf("[SessionMonitor] Session ID %d of user %s reached end time %s. Completing session...",
					session.ID, session.UserID, endTime.Format("15:04"))

				_, err := s.CompleteSession(ctx, session.ID)
				if err != nil {
					log.Printf("[SessionMonitor] Failed to auto-complete session ID %d: %v", session.ID, err)
				} else {
					log.Printf("[SessionMonitor] Successfully completed session ID %d", session.ID)
				}
			}
		}
	}

	// 2. Process in-progress schedules directly
	schedules, err := s.scheduleRepo.FindAll(ctx)
	if err == nil {
		for _, sched := range schedules {
			if sched.Status == dto.ScheduleStatusInProgress {
				dateStr := fmt.Sprintf("%s %s", sched.Date.Format("2006-01-02"), sched.Time)
				startTime, err := time.ParseInLocation("2006-01-02 15:04", dateStr, loc)
				if err != nil {
					startTime, _ = time.Parse("2006-01-02 15:04", dateStr)
				}
				dur := parseDurationMinutes(sched.Duration)
				endTime := startTime.Add(dur)

				if now.After(endTime) || now.Equal(endTime) {
					sess, sessErr := s.sessionRepo.FindByScheduleID(ctx, sched.ID)
					if sessErr == nil && sess != nil && sess.Status != "completed" {
						_, _ = s.CompleteSession(ctx, sess.ID)
					}
					if err := s.scheduleRepo.UpdateStatus(ctx, sched.ID, dto.ScheduleStatusCompleted); err != nil {
						log.Printf("[ScheduleMonitor] Failed to update in-progress schedule ID %d to completed: %v", sched.ID, err)
					} else {
						log.Printf("[ScheduleMonitor] Auto-updated schedule ID %d to completed", sched.ID)
					}
				}
			}
		}
	}

	return nil
}