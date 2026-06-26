package services

import (
	"context"
	"errors"
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
}

type SessionService struct {
	sessionRepo     repositories.ISessionRepository
	scheduleRepo    repositories.IScheduleRepository
	entitlementRepo repositories.IEntitlementRepository
	enrollmentRepo  repositories.IEnrollmentRepository
	eventPublisher  kafka.IEventPublisher
	db              *gorm.DB
}

func NewSessionService(
	sessionRepo repositories.ISessionRepository,
	scheduleRepo repositories.IScheduleRepository,
	entitlementRepo repositories.IEntitlementRepository,
) ISessionService {
	return &SessionService{
		sessionRepo:     sessionRepo,
		scheduleRepo:    scheduleRepo,
		entitlementRepo: entitlementRepo,
	}
}

// NewSessionServiceWithEventPublisher creates a new session service with event publisher
func NewSessionServiceWithEventPublisher(
	sessionRepo repositories.ISessionRepository,
	scheduleRepo repositories.IScheduleRepository,
	entitlementRepo repositories.IEntitlementRepository,
	eventPublisher kafka.IEventPublisher,
) ISessionService {
	return &SessionService{
		sessionRepo:     sessionRepo,
		scheduleRepo:    scheduleRepo,
		entitlementRepo: entitlementRepo,
		eventPublisher:  eventPublisher,
	}
}

// NewSessionServiceWithAllDeps creates a session service with all dependencies including enrollment repo and DB
func NewSessionServiceWithAllDeps(
	sessionRepo repositories.ISessionRepository,
	scheduleRepo repositories.IScheduleRepository,
	entitlementRepo repositories.IEntitlementRepository,
	enrollmentRepo repositories.IEnrollmentRepository,
	eventPublisher kafka.IEventPublisher,
	db *gorm.DB,
) ISessionService {
	return &SessionService{
		sessionRepo:     sessionRepo,
		scheduleRepo:    scheduleRepo,
		entitlementRepo: entitlementRepo,
		enrollmentRepo:  enrollmentRepo,
		eventPublisher:  eventPublisher,
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
	if session.Status != "in-progress" {
		return nil, errors.New("session cannot be completed: must be in progress")
	}

	// Update associated schedule if any (do this before completing session)
	if session.ScheduleID != nil {
		if err := s.scheduleRepo.UpdateStatus(ctx, *session.ScheduleID, dto.ScheduleStatusCompleted); err != nil {
			// Log error but continue
		}
	}

	// If we have enrollment repo and db, use transaction for atomic updates
	if s.enrollmentRepo != nil && s.db != nil {
		err = s.db.Transaction(func(tx *gorm.DB) error {
			// Complete the session
			completedAt := time.Now()
			if err := s.sessionRepo.CompleteSessionTx(tx, id, completedAt); err != nil {
				return err
			}

			// Increment used sessions in entitlement
			if session.EntitlementID != uuid.Nil {
				entitlement, err := s.entitlementRepo.FindByID(ctx, session.EntitlementID)
				if err != nil || entitlement == nil {
					return err
				}

				newUsedSessions := entitlement.UsedSessions + 1
				if err := s.entitlementRepo.UpdateUsedSessionsTx(tx, entitlement.ID, newUsedSessions); err != nil {
					return err
				}

				// Check if all sessions are used → update enrollment to completed
				if session.EnrollmentID != uuid.Nil {
					enrollment, err := s.enrollmentRepo.FindByID(ctx, session.EnrollmentID)
					if err == nil && enrollment != nil {
						// Check if all entitlements for this enrollment are exhausted
						allEntitlements, err := s.entitlementRepo.FindByEnrollmentID(ctx, session.EnrollmentID)
						if err == nil && len(allEntitlements) > 0 {
							allExhausted := true
							for _, ent := range allEntitlements {
								if ent.ID != entitlement.ID {
									// Check remaining for other entitlements
									// For simplicity, we only mark completed if this entitlement is the last one
									// and all others have been used up
								}
								if ent.UsedSessions < ent.TotalSessions {
									allExhausted = false
								}
							}

							// If this was the last entitlement to be exhausted, mark enrollment as completed
							if allExhausted || (entitlement.UsedSessions >= entitlement.TotalSessions) {
								if enrollment.Status != models.EnrollmentStatusCompleted {
									if err := s.enrollmentRepo.UpdateStatusTx(tx, enrollment.ID, models.EnrollmentStatusCompleted); err != nil {
										return err
									}
								}
							} else if enrollment.Status == models.EnrollmentStatusPendingPayment {
								// If enrollment was pending, mark as in_progress when first session is completed
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

		// Increment used sessions in entitlement
		if session.EntitlementID != uuid.Nil {
			entitlement, err := s.entitlementRepo.FindByID(ctx, session.EntitlementID)
			if err == nil && entitlement != nil {
				newUsedSessions := entitlement.UsedSessions + 1
				if err := s.entitlementRepo.UpdateUsedSessions(ctx, entitlement.ID, newUsedSessions); err != nil {
					// Log error but continue
				}
			}
		}
	}

	// Publish Kafka event for session completion (for external services like core-service)
	if s.eventPublisher != nil {
		var packageID uuid.UUID
		if session.EnrollmentID != uuid.Nil {
			enrollments, _ := s.entitlementRepo.FindByEnrollmentID(ctx, session.EnrollmentID)
			if len(enrollments) > 0 {
				packageID, _ = uuid.Parse(enrollments[0].SourceID)
			}
		}

		// Get updated entitlement for remaining sessions count
		var sessionsRemaining int
		if session.EntitlementID != uuid.Nil {
			entitlement, err := s.entitlementRepo.FindByID(ctx, session.EntitlementID)
			if err == nil && entitlement != nil {
				sessionsRemaining = entitlement.TotalSessions - entitlement.UsedSessions
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

	// Decrement used sessions in entitlement (session cancelled, so it doesn't count)
	entitlement, err := s.entitlementRepo.FindByID(ctx, session.EntitlementID)
	if err == nil && entitlement != nil {
		if entitlement.UsedSessions > 0 {
			if err := s.entitlementRepo.UpdateUsedSessions(ctx, entitlement.ID, entitlement.UsedSessions-1); err != nil {
				// Log error
			}
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