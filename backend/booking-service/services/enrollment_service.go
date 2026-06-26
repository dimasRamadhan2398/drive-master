package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	cClient "booking-service/clients/core"
	"booking-service/models"
	"booking-service/models/dto"
	"booking-service/pkg/kafka"
	"booking-service/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type IEnrollmentService interface {
	CreateEnrollment(ctx context.Context, req dto.CreateEnrollmentRequest) (*dto.EnrollmentWithTransactionResponse, error)
	GetEnrollment(ctx context.Context, id uuid.UUID) (*dto.EnrollmentWithTransactionResponse, error)
	UpdateEnrollment(ctx context.Context, id uuid.UUID, req dto.UpdateEnrollmentRequest) (*dto.EnrollmentResponse, error)
	CancelEnrollment(ctx context.Context, id uuid.UUID) error
	MarkAsPaid(ctx context.Context, id uuid.UUID, totalPrice float64) (*dto.EnrollmentResponse, error)
	ListEnrollments(ctx context.Context, page, limit int) (*dto.EnrollmentListResponse, error)
	ListUserEnrollments(ctx context.Context, userID uuid.UUID, page, limit int) (*dto.EnrollmentListResponse, error)
	ListEnrollmentsByStatus(ctx context.Context, status string, page, limit int) (*dto.EnrollmentListResponse, error)
	CreateEntitlementFromEnrollment(ctx context.Context, enrollmentID uuid.UUID, sourceType, sourceID string, totalSessions int, expiresAt time.Time) (*dto.EntitlementResponse, error)
}

type EnrollmentService struct {
	enrollmentRepo   repositories.IEnrollmentRepository
	entitlementRepo repositories.IEntitlementRepository
	transactionSvc   ITransactionService
	eventPublisher  kafka.IEventPublisher
	entitlementSvc  IEntitlementService
	coreClient      cClient.ICoreClient
}

func NewEnrollmentService(
	enrollmentRepo repositories.IEnrollmentRepository,
	entitlementRepo repositories.IEntitlementRepository,
) IEnrollmentService {
	return &EnrollmentService{
		enrollmentRepo:  enrollmentRepo,
		entitlementRepo: entitlementRepo,
	}
}

// NewEnrollmentServiceWithDeps creates a new enrollment service with all dependencies
func NewEnrollmentServiceWithDeps(
	enrollmentRepo repositories.IEnrollmentRepository,
	entitlementRepo repositories.IEntitlementRepository,
	eventPublisher kafka.IEventPublisher,
	entitlementSvc IEntitlementService,
	coreClient cClient.ICoreClient,
) IEnrollmentService {
	return &EnrollmentService{
		enrollmentRepo:  enrollmentRepo,
		entitlementRepo: entitlementRepo,
		eventPublisher:  eventPublisher,
		entitlementSvc:  entitlementSvc,
		coreClient:      coreClient,
	}
}

// NewEnrollmentServiceWithAllDeps creates a new enrollment service with transaction service
func NewEnrollmentServiceWithAllDeps(
	enrollmentRepo repositories.IEnrollmentRepository,
	entitlementRepo repositories.IEntitlementRepository,
	transactionSvc ITransactionService,
	eventPublisher kafka.IEventPublisher,
	entitlementSvc IEntitlementService,
	coreClient cClient.ICoreClient,
) IEnrollmentService {
	return &EnrollmentService{
		enrollmentRepo:   enrollmentRepo,
		entitlementRepo:  entitlementRepo,
		transactionSvc:    transactionSvc,
		eventPublisher:   eventPublisher,
		entitlementSvc:   entitlementSvc,
		coreClient:       coreClient,
	}
}

func (s *EnrollmentService) CreateEnrollment(ctx context.Context, req dto.CreateEnrollmentRequest) (*dto.EnrollmentWithTransactionResponse, error) {
	// Fetch package details from core-service to get price and validity
	var packageName string
	var totalPrice float64
	var expiresAt time.Time

	if s.coreClient != nil {
		pkg, err := s.coreClient.GetPackageByID(ctx, req.PackageID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch package details: %w", err)
		}
		packageName = pkg.Name
		totalPrice = pkg.Price
		// Calculate expiration date based on package validity
		expiresAt = time.Now().AddDate(0, 0, pkg.ValidityDays)
		if expiresAt.Before(time.Now()) {
			expiresAt = time.Now().AddDate(1, 0, 0) // Default 1 year if validity is invalid
		}
	} else {
		// Fallback defaults if core-client is not available
		packageName = "Package"
		totalPrice = 0
		expiresAt = time.Now().AddDate(1, 0, 0) // Default 1 year validity
	}

	// Fetch add-on details if provided
	var addOns []AddOnInfo
	var addOnsTotal float64
	for _, addOnID := range req.AddOns {
		if s.coreClient != nil {
			addOn, err := s.coreClient.GetAddOnByID(ctx, addOnID)
			if err != nil {
				// Skip add-ons that can't be found
				continue
			}
			addOns = append(addOns, AddOnInfo{
				ID:       addOn.ID,
				Name:     addOn.Title,
				Price:    addOn.Price,
				Sessions: addOn.Sessions,
			})
			addOnsTotal += addOn.Price
		}
	}

	// Calculate final total price
	totalPrice += addOnsTotal

	// Create enrollment with pending_payment status
	enrollment := &models.Enrollment{
		UserID:     req.UserID,
		PackageID:  req.PackageID,
		Status:     models.EnrollmentStatusPendingPayment,
		TotalPrice: totalPrice,
		ExpiresAt:  expiresAt,
	}

	if err := s.enrollmentRepo.Create(ctx, enrollment); err != nil {
		return nil, err
	}

	// Create transaction if transaction service is available
	var transaction *dto.TransactionResponse
	if s.transactionSvc != nil {
		tx, err := s.transactionSvc.CreateTransaction(ctx, enrollment.ID, req.UserID, req.PackageID, packageName, totalPrice-addOnsTotal, addOns)
		if err != nil {
			// Log error but don't fail enrollment creation
			// In production, you might want to handle this differently
		} else {
			txResp, err := s.transactionSvc.GetTransactionByID(ctx, tx.ID)
			if err == nil && txResp != nil {
				transaction = ptrToTransactionResponse(nil, txResp)
			}
		}
	}

	enrollmentResp := s.enrollmentRepo.ToResponse(enrollment)
	return &dto.EnrollmentWithTransactionResponse{
		Enrollment:  enrollmentResp,
		Transaction: transaction,
	}, nil
}

// Helper function to convert transaction model to response
func ptrToTransactionResponse(_ *repositories.TransactionRepository, tx *models.Transaction) *dto.TransactionResponse {
	items := make([]dto.TransactionItemResponse, len(tx.Items))
	for i, item := range tx.Items {
		items[i] = dto.TransactionItemResponse{
			ID:        item.ID,
			ItemType:  string(item.ItemType),
			ItemID:    item.ItemID,
			ItemName:  item.ItemName,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			Subtotal:  item.Subtotal,
			Sessions:  item.Sessions,
		}
	}
	return &dto.TransactionResponse{
		ID:           tx.ID,
		EnrollmentID: tx.EnrollmentID,
		UserID:      tx.UserID,
		BasePrice:   tx.BasePrice,
		AddOnsTotal: tx.AddOnsTotal,
		TotalAmount: tx.TotalAmount,
		Status:      string(tx.Status),
		PaidAt:      tx.PaidAt,
		Items:       items,
		CreatedAt:   tx.CreatedAt,
		UpdatedAt:   tx.UpdatedAt,
	}
}

func (s *EnrollmentService) GetEnrollment(ctx context.Context, id uuid.UUID) (*dto.EnrollmentWithTransactionResponse, error) {
	enrollment, err := s.enrollmentRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("enrollment not found")
		}
		return nil, err
	}

	enrollmentResp := s.enrollmentRepo.ToResponse(enrollment)

	// Get transaction if transaction service is available
	var transaction *dto.TransactionResponse
	if s.transactionSvc != nil {
		tx, err := s.transactionSvc.GetTransactionByEnrollmentID(ctx, id)
		if err == nil && tx != nil {
			transaction = ptrToTransactionResponse(nil, tx)
		}
	}

	return &dto.EnrollmentWithTransactionResponse{
		Enrollment:  enrollmentResp,
		Transaction: transaction,
	}, nil
}

func (s *EnrollmentService) UpdateEnrollment(ctx context.Context, id uuid.UUID, req dto.UpdateEnrollmentRequest) (*dto.EnrollmentResponse, error) {
	enrollment, err := s.enrollmentRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("enrollment not found")
		}
		return nil, err
	}

	if req.ExpiresAt != nil {
		enrollment.ExpiresAt = *req.ExpiresAt
	}
	if req.Status != nil {
		enrollment.Status = models.EnrollmentStatus(*req.Status)
	}

	if err := s.enrollmentRepo.Update(ctx, enrollment); err != nil {
		return nil, err
	}

	resp := s.enrollmentRepo.ToResponse(enrollment)
	return &resp, nil
}

func (s *EnrollmentService) CancelEnrollment(ctx context.Context, id uuid.UUID) error {
	enrollment, err := s.enrollmentRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("enrollment not found")
		}
		return err
	}

	if enrollment.Status == models.EnrollmentStatusCancelled {
		return errors.New("enrollment already cancelled")
	}

	if enrollment.Status == models.EnrollmentStatusCompleted {
		return errors.New("cannot cancel a completed enrollment")
	}

	return s.enrollmentRepo.UpdateStatus(ctx, id, models.EnrollmentStatusCancelled)
}

func (s *EnrollmentService) MarkAsPaid(ctx context.Context, id uuid.UUID, totalPrice float64) (*dto.EnrollmentResponse, error) {
	enrollment, err := s.enrollmentRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("enrollment not found")
		}
		return nil, err
	}

	if enrollment.Status != models.EnrollmentStatusPendingPayment {
		return nil, errors.New("enrollment is not in pending_payment status")
	}

	now := time.Now()
	if err := s.enrollmentRepo.MarkAsPaid(ctx, id, now, totalPrice); err != nil {
		return nil, err
	}

	// Reload enrollment
	enrollment, err = s.enrollmentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Publish enrollment.paid event if event publisher is available
	if s.eventPublisher != nil {
		// Get package details for the event
		var totalSessions int
		var packageName string
		if s.coreClient != nil {
			if pkg, err := s.coreClient.GetPackageByID(ctx, enrollment.PackageID); err == nil {
				totalSessions = pkg.Sessions
				packageName = pkg.Name
			}
		}

		_ = s.eventPublisher.PublishEnrollmentPaid(ctx, enrollment.ID.String(), enrollment.UserID.String(), enrollment.PackageID, totalPrice, totalSessions, packageName)
	}

	// Create entitlements automatically if services are available
	s.createEntitlementsForEnrollment(ctx, enrollment)

	resp := s.enrollmentRepo.ToResponse(enrollment)
	return &resp, nil
}

// createEntitlementsForEnrollment creates user entitlements based on the package
func (s *EnrollmentService) createEntitlementsForEnrollment(ctx context.Context, enrollment *models.Enrollment) {
	// Skip if we don't have the required services
	if s.entitlementSvc == nil || s.coreClient == nil {
		return
	}

	// Get package details from core-service
	pkg, err := s.coreClient.GetPackageByID(ctx, enrollment.PackageID)
	if err != nil {
		// Log error but don't fail the enrollment
		return
	}

	// Calculate expiration date based on package validity
	expiresAt := enrollment.CreatedAt.AddDate(0, 0, pkg.ValidityDays)
	if expiresAt.Before(time.Now()) {
		expiresAt = time.Now().AddDate(1, 0, 0) // Default 1 year if validity is invalid
	}

	// Create entitlement for the package
	_, err = s.entitlementSvc.CreateEntitlement(ctx, dto.CreateEntitlementRequest{
		UserID:            enrollment.UserID,
		SourceType:        "package",
		SourceID:          enrollment.PackageID.String(),
		TotalSessions:     pkg.Sessions,
		SessionsRemaining: pkg.Sessions,
		ExpiresAt:         expiresAt,
	})
	if err != nil {
		// Log error but don't fail the enrollment
		return
	}
}

func (s *EnrollmentService) ListEnrollments(ctx context.Context, page, limit int) (*dto.EnrollmentListResponse, error) {
	enrollments, err := s.enrollmentRepo.FindAllPaginated(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	total, err := s.enrollmentRepo.CountAll(ctx)
	if err != nil {
		return nil, err
	}

	resp := s.enrollmentRepo.ToListResponse(enrollments, total, page, limit)
	return &resp, nil
}

func (s *EnrollmentService) ListUserEnrollments(ctx context.Context, userID uuid.UUID, page, limit int) (*dto.EnrollmentListResponse, error) {
	enrollments, err := s.enrollmentRepo.FindByUserIDPaginated(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	total, err := s.enrollmentRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	resp := s.enrollmentRepo.ToListResponse(enrollments, total, page, limit)
	return &resp, nil
}

func (s *EnrollmentService) ListEnrollmentsByStatus(ctx context.Context, status string, page, limit int) (*dto.EnrollmentListResponse, error) {
	enrollments, total, err := s.enrollmentRepo.FindByStatus(ctx, models.EnrollmentStatus(status))
	if err != nil {
		return nil, err
	}

	resp := &dto.EnrollmentListResponse{
		Data:       make([]dto.EnrollmentResponse, len(enrollments)),
		Pagination: dto.NewPaginationMeta(total, page, limit),
	}
	for i, e := range enrollments {
		resp.Data[i] = s.enrollmentRepo.ToResponse(&e)
	}
	return resp, nil
}

func (s *EnrollmentService) CreateEntitlementFromEnrollment(ctx context.Context, enrollmentID uuid.UUID, sourceType, sourceID string, totalSessions int, expiresAt time.Time) (*dto.EntitlementResponse, error) {
	enrollment, err := s.enrollmentRepo.FindByID(ctx, enrollmentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("enrollment not found")
		}
		return nil, err
	}

	entitlement := &models.UserEntitlement{
		EnrollmentID:  enrollment.ID,
		UserID:        enrollment.UserID,
		SourceType:    sourceType,
		SourceID:      sourceID,
		TotalSessions: totalSessions,
		UsedSessions:  0,
		ExpiresAt:     expiresAt,
	}

	if err := s.entitlementRepo.Create(ctx, entitlement); err != nil {
		return nil, err
	}

	resp := s.entitlementRepo.ToResponse(entitlement)
	return &resp, nil
}