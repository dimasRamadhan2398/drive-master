package services

import (
	"context"
	"errors"

	"booking-service/models"
	"booking-service/models/dto"
	"booking-service/repositories"

	"gorm.io/gorm"
)

type EntitlementService struct {
	entitlementRepo *repositories.EntitlementRepository
}

func NewEntitlementService(entitlementRepo *repositories.EntitlementRepository) IEntitlementService {
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
		SessionsRemaining: req.SessionsRemaining,
		ExpiresAt:         req.ExpiresAt,
	}

	if err := s.entitlementRepo.Create(ctx, entitlement); err != nil {
		return nil, err
	}

	resp := s.entitlementRepo.ToResponse(entitlement)
	return &resp, nil
}

func (s *EntitlementService) GetEntitlement(ctx context.Context, id uint) (*dto.EntitlementResponse, error) {
	entitlement, err := s.entitlementRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("entitlement not found")
		}
		return nil, err
	}

	resp := s.entitlementRepo.ToResponse(entitlement)
	return &resp, nil
}

func (s *EntitlementService) UpdateEntitlement(ctx context.Context, id uint, req dto.UpdateEntitlementRequest) (*dto.EntitlementResponse, error) {
	entitlement, err := s.entitlementRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("entitlement not found")
		}
		return nil, err
	}

	if req.SessionsRemaining != nil {
		entitlement.SessionsRemaining = *req.SessionsRemaining
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

func (s *EntitlementService) DeleteEntitlement(ctx context.Context, id uint) error {
	return s.entitlementRepo.Delete(ctx, id)
}

func (s *EntitlementService) ListEntitlements(ctx context.Context, page, limit int) (*dto.EntitlementListResponse, error) {
	entitlements, total, err := s.entitlementRepo.List(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	resp := s.entitlementRepo.ToListResponse(entitlements, total, page, limit)
	return &resp, nil
}

func (s *EntitlementService) GetUserEntitlements(ctx context.Context, userID uint) ([]dto.EntitlementResponse, error) {
	entitlements, err := s.entitlementRepo.GetByUserID(ctx, userID)
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
	entitlements, err := s.entitlementRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.EntitlementResponse, len(entitlements))
	for i, e := range entitlements {
		responses[i] = s.entitlementRepo.ToResponse(&e)
	}
	return responses, nil
}