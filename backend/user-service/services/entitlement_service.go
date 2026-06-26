package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"user-service/models"
	"user-service/models/dto"
	"user-service/repositories"
	"user-service/services/listeners"

	"github.com/google/uuid"
)

type IEntitlementService interface {
	CreateEntitlement(ctx context.Context, memberID uuid.UUID, input dto.CreateEntitlementInput) (*dto.EntitlementResponse, error)
	UpdateEntitlement(ctx context.Context, memberID, entitlementID uuid.UUID, input dto.UpdateEntitlementInput) (*dto.EntitlementResponse, error)
	DeleteEntitlement(ctx context.Context, memberID, entitlementID uuid.UUID) error
	GetEntitlement(ctx context.Context, memberID, entitlementID uuid.UUID) (*dto.EntitlementResponse, error)
	GetEntitlementByID(ctx context.Context, entitlementID uuid.UUID) (*dto.EntitlementResponse, error)
	ListEntitlements(ctx context.Context, memberID uuid.UUID, page, limit int) (*dto.EntitlementListResponse, error)
	UseSession(ctx context.Context, memberID, entitlementID uuid.UUID, input dto.UseSessionInput) (*dto.EntitlementResponse, error)
	SyncEntitlementFromBooking(ctx context.Context, memberID, bookingID uuid.UUID, packageID uuid.UUID, packageName string, totalSessions int) (*dto.EntitlementResponse, error)
	CalculateTotalAvailableSessions(ctx context.Context, memberID uuid.UUID) (int, error)
	FindSessionsStatsByMemberIDs(ctx context.Context, memberIDs []uuid.UUID) (map[uuid.UUID][]models.Entitlement, error)
}

type EntitlementService struct {
	entitlementRepo repositories.IEntitlementRepository
	memberRepo      repositories.IMemberRepository
	certService     ICertificationService
	eventPublisher  listeners.EventPublisherInterface
	listener        listeners.IEntitlementCompletedListener
}

func NewEntitlementService(
	entitlementRepo repositories.IEntitlementRepository,
	memberRepo repositories.IMemberRepository,
	certService ICertificationService,
	eventPublisher listeners.EventPublisherInterface,
	listener listeners.IEntitlementCompletedListener,
) IEntitlementService {
	return &EntitlementService{
		entitlementRepo: entitlementRepo,
		memberRepo:      memberRepo,
		certService:     certService,
		eventPublisher:  eventPublisher,
		listener:        listener,
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

	// Update member profile total available sessions
	if err := s.updateMemberProfileTotalSessions(ctx, memberID); err != nil {
		// Log error but don't fail the request
		fmt.Printf("failed to update member profile total sessions: %v\n", err)
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

	// Update member profile total available sessions
	if err := s.updateMemberProfileTotalSessions(ctx, memberID); err != nil {
		fmt.Printf("failed to update member profile total sessions: %v\n", err)
	}

	return toEntitlementResponse(entitlement), nil
}

func (s *EntitlementService) DeleteEntitlement(ctx context.Context, memberID, entitlementID uuid.UUID) error {
	_, err := s.entitlementRepo.FindByMemberAndID(ctx, memberID, entitlementID)
	if err != nil {
		return err
	}
	if err := s.entitlementRepo.Delete(ctx, entitlementID); err != nil {
		return err
	}

	// Update member profile total available sessions
	if err := s.updateMemberProfileTotalSessions(ctx, memberID); err != nil {
		fmt.Printf("failed to update member profile total sessions: %v\n", err)
	}

	return nil
}

func (s *EntitlementService) GetEntitlement(ctx context.Context, memberID, entitlementID uuid.UUID) (*dto.EntitlementResponse, error) {
	entitlement, err := s.entitlementRepo.FindByMemberAndID(ctx, memberID, entitlementID)
	if err != nil {
		return nil, err
	}
	return toEntitlementResponse(entitlement), nil
}

func (s *EntitlementService) GetEntitlementByID(ctx context.Context, entitlementID uuid.UUID) (*dto.EntitlementResponse, error) {
	entitlement, err := s.entitlementRepo.FindByID(ctx, entitlementID)
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

	return &dto.EntitlementListResponse{
		Data:       responses,
		Pagination: dto.NewPaginationMeta(total, page, limit),
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

	if err := entitlement.ValidateSessionCount(); err != nil {
		// log critical error — data is inconsistent
		log.Printf("[CRITICAL] entitlement %s failed integrity check: %v",
			entitlementID, err)
		return nil, err
	}

	// Update status if no remaining sessions
	if entitlement.Remaining == 0 {
		entitlement.Status = models.EntitlementStatusUsed
		entitlement.UpdatedAt = time.Now()
		if err := s.entitlementRepo.Update(ctx, entitlement); err != nil {
			return nil, err
		}

		// Trigger completion listener: issue certification and publish event
		if s.listener != nil {
			if err := s.listener.OnEntitlementCompleted(ctx, entitlement); err != nil {
				// Log but don't fail the request - certification can be issued later
				log.Printf("[WARN] failed to trigger entitlement completion listener: %v", err)
			}
		}
	}

	// Update member profile total available sessions
	if err := s.updateMemberProfileTotalSessions(ctx, memberID); err != nil {
		fmt.Printf("failed to update member profile total sessions: %v\n", err)
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

// CalculateTotalAvailableSessions sums up remaining sessions from all active entitlements for a member
func (s *EntitlementService) CalculateTotalAvailableSessions(ctx context.Context, memberID uuid.UUID) (int, error) {
	// Get all active entitlements for the member
	entitlements, _, err := s.entitlementRepo.FindByMemberID(ctx, memberID, 1, 1000)
	if err != nil {
		return 0, err
	}

	total := 0
	for _, ent := range entitlements {
		if ent.Status == models.EntitlementStatusActive {
			total += ent.Remaining
		}
	}

	return total, nil
}

// FindSessionsStatsByMemberIDs returns active entitlements for multiple members
func (s *EntitlementService) FindSessionsStatsByMemberIDs(ctx context.Context, memberIDs []uuid.UUID) (map[uuid.UUID][]models.Entitlement, error) {
	return s.entitlementRepo.FindActiveByMemberIDs(ctx, memberIDs)
}

// updateMemberProfileTotalSessions updates the member profile with the calculated total available sessions
func (s *EntitlementService) updateMemberProfileTotalSessions(ctx context.Context, memberID uuid.UUID) error {
	total, err := s.CalculateTotalAvailableSessions(ctx, memberID)
	if err != nil {
		return err
	}
	return s.memberRepo.UpdateTotalAvailableSessions(ctx, memberID, total)
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
