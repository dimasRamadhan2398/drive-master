package services

import (
	"context"
	"errors"

	"booking-service/models"
	"booking-service/models/dto"
	"booking-service/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EntitlementService struct {
	entitlementRepo repositories.IEntitlementRepository
}

func NewEntitlementService(entitlementRepo repositories.IEntitlementRepository) IEntitlementService {
	return &EntitlementService{
		entitlementRepo: entitlementRepo,
	}
}

func (s *EntitlementService) CreateEntitlement(ctx context.Context, req dto.CreateEntitlementRequest) (*dto.EntitlementResponse, error) {
	entitlement := &models.UserEntitlement{
		UserID:            req.UserID,
		SourceType:        req.SourceType,
		SourceID:          req.SourceID,
		TotalSessions:     req.TotalSessions,
		UsedSessions: 		req.SessionsRemaining,
		ExpiresAt:         req.ExpiresAt,
	}

	if err := s.entitlementRepo.Create(ctx, entitlement); err != nil {
		return nil, err
	}

	resp := s.entitlementRepo.ToResponse(entitlement)
	return &resp, nil
}

func (s *EntitlementService) GetEntitlement(ctx context.Context, id uuid.UUID) (*dto.EntitlementResponse, error) {
	entitlement, err := s.entitlementRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("entitlement not found")
		}
		return nil, err
	}

	resp := s.entitlementRepo.ToResponse(entitlement)
	return &resp, nil
}

func (s *EntitlementService) UpdateEntitlement(ctx context.Context, id uuid.UUID, req dto.UpdateEntitlementRequest) (*dto.EntitlementResponse, error) {
	entitlement, err := s.entitlementRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("entitlement not found")
		}
		return nil, err
	}

	if req.SessionsRemaining != nil {
		entitlement.UsedSessions = *req.SessionsRemaining
	}
	if req.ExpiresAt != nil {
		entitlement.ExpiresAt = *req.ExpiresAt
	}

	if err := s.entitlementRepo.Update(ctx, entitlement); err != nil {
		return nil, err
	}

	resp := s.entitlementRepo.ToResponse(entitlement)
	return &resp, nil
}

func (s *EntitlementService) DeleteEntitlement(ctx context.Context, id uuid.UUID) error {
	entitlement, err := s.entitlementRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return s.entitlementRepo.Delete(ctx, entitlement)
}

func (s *EntitlementService) ListEntitlements(ctx context.Context, page, limit int) (*dto.EntitlementListResponse, error) {
	entitlements, err := s.entitlementRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	total, err := s.entitlementRepo.CountAll(ctx)
	if err != nil {
		return nil, err
	}

	resp := s.entitlementRepo.ToListResponse(entitlements, total, page, limit)
	return &resp, nil
}

func (s *EntitlementService) GetUserEntitlements(ctx context.Context, userID uint) ([]dto.EntitlementResponse, error) {
	entitlements, err := s.entitlementRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.EntitlementResponse, len(entitlements))
	for i, e := range entitlements {
		responses[i] = s.entitlementRepo.ToResponse(&e)
	}
	return responses, nil
}

func (s *EntitlementService) GetActiveEntitlements(ctx context.Context, userID uint) ([]dto.EntitlementResponse, error) {
	entitlements, err := s.entitlementRepo.FindActiveByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.EntitlementResponse, len(entitlements))
	for i, e := range entitlements {
		responses[i] = s.entitlementRepo.ToResponse(&e)
	}
	return responses, nil
}