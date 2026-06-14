package services

import (
	"context"
	"errors"

	"booking-service/models"
	"booking-service/models/dto"
	"booking-service/repositories"

	"gorm.io/gorm"
)

type ISessionService interface {
	CreateSession(ctx context.Context, req dto.CreateSessionRequest) (*dto.SessionResponse, error)
	GetSession(ctx context.Context, id uint) (*dto.SessionResponse, error)
	ListSessions(ctx context.Context, page, limit int) (*dto.SessionListResponse, error)
	GetStats(ctx context.Context) (*dto.SessionStatsResponse, error)
}

type SessionService struct {
	sessionRepo repositories.ISessionRepository
}

func NewSessionService(sessionRepo repositories.ISessionRepository) ISessionService {
	return &SessionService{
		sessionRepo: sessionRepo,
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

func (s *SessionService) ListUserSessions(ctx context.Context, userID uint, page, limit int) (*dto.SessionListResponse, error) {
	sessions, err := s.sessionRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	total := int64(len(sessions))

	resp := s.sessionRepo.ToListResponse(sessions, total, page, limit)
	return &resp, nil
}

func (s *SessionService) ListInstructorSessions(ctx context.Context, instructorID uint, page, limit int) (*dto.SessionListResponse, error) {
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

	return&dto.SessionStatsResponse{
		TotalSessions:     stats.Total,
		ActiveSessions:    stats.Active,
		CompletedSessions: stats.Completed,
		PendingSessions:   stats.Pending,
	}, nil
}