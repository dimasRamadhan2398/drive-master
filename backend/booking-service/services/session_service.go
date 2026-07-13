package services

import (
	uClient "booking-service/clients/user"
	"context"
	"errors"
	"fmt"
	"log"
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
	CancelSession(ctx context.Context, id uint) (*dto.SessionResponse, error)
	ListUserSessions(ctx context.Context, userID uuid.UUID, page, limit int) (*dto.SessionListResponse, error)
	ListInstructorSessions(ctx context.Context, instructorID uuid.UUID, page, limit int) (*dto.SessionListResponse, error)
	AutoCompleteOngoingSessions(ctx context.Context) error
}

type SessionService struct {
	sessionRepo     repositories.ISessionRepository
	scheduleRepo    repositories.IScheduleRepository
	enrollmentRepo  repositories.IEnrollmentRepository
	eventPublisher  kafka.IEventPublisher
	userClient      uClient.IUserClient
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
	db *gorm.DB,
) ISessionService {
	return &SessionService{
		sessionRepo:     sessionRepo,
		scheduleRepo:    scheduleRepo,
		enrollmentRepo:  enrollmentRepo,
		eventPublisher:  eventPublisher,
		userClient:      userClient,
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
				sessionsRemaining = entitlement.RemainingSessions
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

func (s *SessionService) AutoCompleteOngoingSessions(ctx context.Context) error {
	// Find all sessions that are in progress
	sessions, err := s.sessionRepo.FindByStatus(ctx, "in_progress")
	if err != nil {
		return err
	}

	now := time.Now()
	for _, session := range sessions {
		// Parse session start date and time
		dateStr := fmt.Sprintf("%s %s", session.Date.Format("2006-01-02"), session.Time)
		// Assume server/database time zone matches local time zone of lessons
		startTime, err := time.ParseInLocation("2006-01-02 15:04", dateStr, time.Local)
		if err != nil {
			startTime, _ = time.Parse("2006-01-02 15:04", dateStr)
		}

		endTime := startTime.Add(time.Duration(session.Duration) * time.Minute)

		// If session reaches end time, complete it
		if now.After(endTime) {
			log.Printf("[SessionMonitor] Session ID %d of user %s reached end time %s (duration %d mins). Completing session...",
				session.ID, session.UserID, endTime.Format("15:04"), session.Duration)

			_, err := s.CompleteSession(ctx, session.ID)
			if err != nil {
				log.Printf("[SessionMonitor] Failed to auto-complete session ID %d: %v", session.ID, err)
			} else {
				log.Printf("[SessionMonitor] Successfully completed session ID %d", session.ID)
			}
		}
	}
	return nil
}