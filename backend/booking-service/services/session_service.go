package services

import (
	"context"
	"errors"

	"booking-service/models"
	"booking-service/models/dto"
	"booking-service/repositories"

	"gorm.io/gorm"
)

type SessionService struct {
	sessionRepo *repositories.SessionRepository
}

func NewSessionService(sessionRepo *repositories.SessionRepository) ISessionService {
	return &SessionService{
		sessionRepo: sessionRepo,
	}
}

func (s *SessionService) CreateSession(ctx context.Context, req dto.CreateSessionRequest) (*dto.SessionResponse, error) {
	session := &models.Session{
		UserID:        req.UserID,
		InstructorID:  req.InstructorID,
		DateOfSession: req.DateOfSession,
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
	session, err := s.sessionRepo.GetByID(ctx, id)
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
	sessions, total, err := s.sessionRepo.List(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	resp := s.sessionRepo.ToListResponse(sessions, total, page, limit)
	return &resp, nil
}

func (s *SessionService) ListUserSessions(ctx context.Context, userID uint, page, limit int) (*dto.SessionListResponse, error) {
	sessions, total, err := s.sessionRepo.GetByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	resp := s.sessionRepo.ToListResponse(sessions, total, page, limit)
	return &resp, nil
}

func (s *SessionService) ListInstructorSessions(ctx context.Context, instructorID uint, page, limit int) (*dto.SessionListResponse, error) {
	sessions, total, err := s.sessionRepo.GetByInstructorID(ctx, instructorID, page, limit)
	if err != nil {
		return nil, err
	}

	resp := s.sessionRepo.ToListResponse(sessions, total, page, limit)
	return &resp, nil
}