package services

import (
	"context"
	"fmt"
	"time"

	"user-service/models"
	"user-service/models/dto"
	"user-service/repositories"

	"github.com/google/uuid"
)

type IEntitlementService interface {
	CreateEntitlement(ctx context.Context, memberID uuid.UUID, input dto.CreateEntitlementInput) (*dto.EntitlementResponse, error)
	UpdateEntitlement(ctx context.Context, memberID, entitlementID uuid.UUID, input dto.UpdateEntitlementInput) (*dto.EntitlementResponse, error)
	DeleteEntitlement(ctx context.Context, memberID, entitlementID uuid.UUID) error
	GetEntitlement(ctx context.Context, memberID, entitlementID uuid.UUID) (*dto.EntitlementResponse, error)
	ListEntitlements(ctx context.Context, memberID uuid.UUID, page, limit int) (*dto.EntitlementListResponse, error)
	UseSession(ctx context.Context, memberID, entitlementID uuid.UUID, input dto.UseSessionInput) (*dto.EntitlementResponse, error)
	SyncEntitlementFromBooking(ctx context.Context, memberID, bookingID uuid.UUID, packageID uuid.UUID, packageName string, totalSessions int) (*dto.EntitlementResponse, error)
}

type EntitlementService struct {
	entitlementRepo repositories.IEntitlementRepository
}

func NewEntitlementService(entitlementRepo repositories.IEntitlementRepository) IEntitlementService {
	return &EntitlementService{
		entitlementRepo: entitlementRepo,
	}
}

func (s *EntitlementService) CreateEntitlement(ctx context.Context, memberID uuid.UUID, input dto.CreateEntitlementInput) (*dto.EntitlementResponse, error) {
	// Check if entitlement already exists for this booking
	existing, err := s.entitlementRepo.FindByBookingID(ctx, input.BookingID)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("entitlement already exists for booking %s", input.BookingID)
	}

	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}

	entitlement := &models.Entitlement{
		ID:            uuid.New(),
		MemberID:      memberID,
		BookingID:     input.BookingID,
		PackageID:     input.PackageID,
		PackageName:   input.PackageName,
		TotalSessions: input.TotalSessions,
		Remaining:     input.TotalSessions,
		UsedSessions:  0,
		StartDate:     startDate,
		Status:        models.EntitlementStatusActive,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if input.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", input.EndDate)
		if err == nil {
			entitlement.EndDate = &endDate
		}
	}

	if err := s.entitlementRepo.Create(ctx, entitlement); err != nil {
		return nil, err
	}

	return toEntitlementResponse(entitlement), nil
}

func (s *EntitlementService) UpdateEntitlement(ctx context.Context, memberID, entitlementID uuid.UUID, input dto.UpdateEntitlementInput) (*dto.EntitlementResponse, error) {
	entitlement, err := s.entitlementRepo.FindByMemberAndID(ctx, memberID, entitlementID)
	if err != nil {
		return nil, err
	}

	if input.Remaining >= 0 {
		entitlement.Remaining = input.Remaining
	}
	if input.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", input.EndDate)
		if err == nil {
			entitlement.EndDate = &endDate
		}
	}
	if input.Status != "" {
		entitlement.Status = models.EntitlementStatus(input.Status)
	}

	entitlement.UpdatedAt = time.Now()

	if err := s.entitlementRepo.Update(ctx, entitlement); err != nil {
		return nil, err
	}

	return toEntitlementResponse(entitlement), nil
}

func (s *EntitlementService) DeleteEntitlement(ctx context.Context, memberID, entitlementID uuid.UUID) error {
	_, err := s.entitlementRepo.FindByMemberAndID(ctx, memberID, entitlementID)
	if err != nil {
		return err
	}
	return s.entitlementRepo.Delete(ctx, entitlementID)
}

func (s *EntitlementService) GetEntitlement(ctx context.Context, memberID, entitlementID uuid.UUID) (*dto.EntitlementResponse, error) {
	entitlement, err := s.entitlementRepo.FindByMemberAndID(ctx, memberID, entitlementID)
	if err != nil {
		return nil, err
	}
	return toEntitlementResponse(entitlement), nil
}

func (s *EntitlementService) ListEntitlements(ctx context.Context, memberID uuid.UUID, page, limit int) (*dto.EntitlementListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	entitlements, total, err := s.entitlementRepo.FindByMemberID(ctx, memberID, page, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.EntitlementResponse, len(entitlements))
	for i, ent := range entitlements {
		responses[i] = *toEntitlementResponse(&ent)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &dto.EntitlementListResponse{
		Data:       responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *EntitlementService) UseSession(ctx context.Context, memberID, entitlementID uuid.UUID, input dto.UseSessionInput) (*dto.EntitlementResponse, error) {
	entitlement, err := s.entitlementRepo.FindByMemberAndID(ctx, memberID, entitlementID)
	if err != nil {
		return nil, err
	}

	if entitlement.Status != models.EntitlementStatusActive {
		return nil, fmt.Errorf("entitlement is not active")
	}

	if entitlement.Remaining <= 0 {
		return nil, fmt.Errorf("no remaining sessions")
	}

	// Atomic decrement
	if err := s.entitlementRepo.DecrementRemaining(ctx, entitlementID); err != nil {
		return nil, err
	}

	// Refresh the entitlement to get updated values
	entitlement, err = s.entitlementRepo.FindByMemberAndID(ctx, memberID, entitlementID)
	if err != nil {
		return nil, err
	}

	// Update status if no remaining sessions
	if entitlement.Remaining == 0 {
		entitlement.Status = models.EntitlementStatusUsed
		entitlement.UpdatedAt = time.Now()
		s.entitlementRepo.Update(ctx, entitlement)
	}

	return toEntitlementResponse(entitlement), nil
}

func (s *EntitlementService) SyncEntitlementFromBooking(ctx context.Context, memberID, bookingID uuid.UUID, packageID uuid.UUID, packageName string, totalSessions int) (*dto.EntitlementResponse, error) {
	// Check if already exists
	existing, err := s.entitlementRepo.FindByBookingID(ctx, bookingID)
	if err == nil && existing != nil {
		return toEntitlementResponse(existing), nil
	}

	input := dto.CreateEntitlementInput{
		BookingID:     bookingID,
		PackageID:     packageID,
		PackageName:   packageName,
		TotalSessions: totalSessions,
		StartDate:     time.Now().Format("2006-01-02"),
	}

	return s.CreateEntitlement(ctx, memberID, input)
}

func toEntitlementResponse(ent *models.Entitlement) *dto.EntitlementResponse {
	resp := &dto.EntitlementResponse{
		ID:            ent.ID,
		MemberID:      ent.MemberID,
		BookingID:     ent.BookingID,
		PackageID:     ent.PackageID,
		PackageName:   ent.PackageName,
		TotalSessions: ent.TotalSessions,
		Remaining:     ent.Remaining,
		UsedSessions:  ent.UsedSessions,
		StartDate:     ent.StartDate.Format("2006-01-02"),
		Status:        string(ent.Status),
		CreatedAt:     ent.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     ent.UpdatedAt.Format(time.RFC3339),
	}

	if ent.EndDate != nil {
		endStr := ent.EndDate.Format("2006-01-02")
		resp.EndDate = &endStr
	}

	return resp
}