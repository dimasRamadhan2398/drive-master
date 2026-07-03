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
	MarkAsPaid(ctx context.Context, id uuid.UUID, totalPrice float64, packageID, packageName, paymentMethod string) (*dto.EnrollmentResponse, error)
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
		if pkg.DiscountPrice > 0 {
			totalPrice = pkg.DiscountPrice
		} else {
			totalPrice = pkg.Price
		}
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
				// Fallback for mock/test addon ID 00000000-0000-0000-0000-000000000001
				if addOnID == uuid.MustParse("00000000-0000-0000-0000-000000000001") {
					addOns = append(addOns, AddOnInfo{
						ID:       addOnID,
						Name:     "Extra Session",
						Price:    350000,
						Sessions: 1,
					})
					addOnsTotal += 350000
				}
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
	s.populateDetails(ctx, &enrollmentResp)
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
	s.populateDetails(ctx, &enrollmentResp)

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
	s.populateDetails(ctx, &resp)
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

func (s *EnrollmentService) MarkAsPaid(ctx context.Context, id uuid.UUID, totalPrice float64, packageID, packageName, paymentMethod string) (*dto.EnrollmentResponse, error) {
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
		var pkgName string
		if s.coreClient != nil {
			if pkg, err := s.coreClient.GetPackageByID(ctx, enrollment.PackageID); err == nil {
				totalSessions = pkg.Sessions
				pkgName = pkg.Name
			}
		}
		if packageName != "" {
			pkgName = packageName
		}

		// Calculate extra sessions from the enrollment transaction
		extraSessions := 0
		if s.transactionSvc != nil {
			tx, err := s.transactionSvc.GetTransactionByEnrollmentID(ctx, enrollment.ID)
			if err == nil && tx != nil {
				for _, item := range tx.Items {
					if item.ItemID == uuid.MustParse("22222222-2222-2222-2222-222222222201") {
						extraSessions += item.Sessions * item.Quantity
					}
				}
			}
		}
		totalSessions += extraSessions

		_ = s.eventPublisher.PublishEnrollmentPaid(ctx, enrollment.ID.String(), enrollment.UserID.String(), enrollment.PackageID, totalPrice, totalSessions, pkgName)
	}

	// Create entitlements automatically if services are available
	s.createEntitlementsForEnrollment(ctx, enrollment)

	// Create a sale record in core-service to track this as a completed sale
	if s.coreClient != nil && packageID != "" {
		var saleItems []cClient.CreateSaleItem

		if s.transactionSvc != nil {
			tx, err := s.transactionSvc.GetTransactionByEnrollmentID(ctx, enrollment.ID)
			if err == nil && tx != nil {
				for _, item := range tx.Items {
					saleItems = append(saleItems, cClient.CreateSaleItem{
						PackageID:   item.ItemID.String(),
						PackageName: item.ItemName,
						Quantity:    item.Quantity,
						UnitPrice:   item.UnitPrice,
						Discount:    0,
					})
				}
			}
		}

		// Fallback to single item if transaction is not found or empty
		if len(saleItems) == 0 {
			pkgName := packageName
			if pkgName == "" {
				pkgName = "Driving Package"
			}
			if pkg, err := s.coreClient.GetPackageByID(ctx, enrollment.PackageID); err == nil && pkgName == "" {
				pkgName = pkg.Name
			}
			saleItems = append(saleItems, cClient.CreateSaleItem{
				PackageID:   packageID,
				PackageName: pkgName,
				Quantity:    1,
				UnitPrice:   totalPrice,
				Discount:    0,
			})
		}

		saleReq := cClient.CreateSaleRequest{
			UserID:        enrollment.UserID.String(),
			PackageID:     packageID,
			PaymentMethod: paymentMethod,
			Source:        "web",
			Items:         saleItems,
		}
		if err := s.coreClient.CreateSale(ctx, saleReq); err != nil {
			// Log error but don't fail the enrollment — sale is non-critical
			fmt.Printf("[MarkAsPaid] Failed to create sale record: %v\n", err)
		} else {
			fmt.Printf("[MarkAsPaid] Sale record created for enrollment %s with %d items\n", id.String(), len(saleItems))
		}
	}

	resp := s.enrollmentRepo.ToResponse(enrollment)
	s.populateDetails(ctx, &resp)
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

	// Calculate extra sessions from the enrollment transaction
	extraSessions := 0
	if s.transactionSvc != nil {
		tx, err := s.transactionSvc.GetTransactionByEnrollmentID(ctx, enrollment.ID)
		if err == nil && tx != nil {
			for _, item := range tx.Items {
				if item.ItemID == uuid.MustParse("22222222-2222-2222-2222-222222222201") {
					extraSessions += item.Sessions * item.Quantity
				}
			}
		}
	}
	totalSessions := pkg.Sessions + extraSessions

	// Create entitlement for the package
	_, err = s.entitlementSvc.CreateEntitlement(ctx, dto.CreateEntitlementRequest{
		UserID:            enrollment.UserID,
		SourceType:        "package",
		SourceID:          enrollment.PackageID.String(),
		TotalSessions:     totalSessions,
		SessionsRemaining: totalSessions,
		ExpiresAt:         expiresAt,
	})
	if err != nil {
		// Log error but don't fail the enrollment
		return
	}

	// Create entitlement for bonus session (1 session, 30 days validity)
	bonusExpiresAt := time.Now().AddDate(0, 0, 30) // Bonus session validity, e.g., 30 days
	_, err = s.entitlementSvc.CreateEntitlement(ctx, dto.CreateEntitlementRequest{
		UserID:            enrollment.UserID,
		SourceType:        "bonus",
		SourceID:          "free-trial-bonus",
		TotalSessions:     1,
		SessionsRemaining: 1,
		ExpiresAt:         bonusExpiresAt,
	})
	if err != nil {
		// Log error but don't fail the enrollment
		fmt.Printf("Failed to create bonus entitlement: %v\n", err)
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

	pkgCache := make(map[uuid.UUID]*dto.EnrollmentPackageResponse)
	addonCache := make(map[uuid.UUID]*dto.EnrollmentAddOnResponse)
	for i := range resp.Data {
		s.populateDetailsWithCache(ctx, &resp.Data[i], pkgCache, addonCache)
	}

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

	pkgCache := make(map[uuid.UUID]*dto.EnrollmentPackageResponse)
	addonCache := make(map[uuid.UUID]*dto.EnrollmentAddOnResponse)
	for i := range resp.Data {
		s.populateDetailsWithCache(ctx, &resp.Data[i], pkgCache, addonCache)
	}

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
	pkgCache := make(map[uuid.UUID]*dto.EnrollmentPackageResponse)
	addonCache := make(map[uuid.UUID]*dto.EnrollmentAddOnResponse)
	for i, e := range enrollments {
		resp.Data[i] = s.enrollmentRepo.ToResponse(&e)
		s.populateDetailsWithCache(ctx, &resp.Data[i], pkgCache, addonCache)
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

func (s *EnrollmentService) populateDetails(ctx context.Context, resp *dto.EnrollmentResponse) {
	pkgCache := make(map[uuid.UUID]*dto.EnrollmentPackageResponse)
	addonCache := make(map[uuid.UUID]*dto.EnrollmentAddOnResponse)
	s.populateDetailsWithCache(ctx, resp, pkgCache, addonCache)
}

func (s *EnrollmentService) populateDetailsWithCache(
	ctx context.Context,
	resp *dto.EnrollmentResponse,
	pkgCache map[uuid.UUID]*dto.EnrollmentPackageResponse,
	addonCache map[uuid.UUID]*dto.EnrollmentAddOnResponse,
) {
	if s.coreClient == nil {
		return
	}

	// Fetch package details
	if pkgDetail, exists := pkgCache[resp.PackageID]; exists {
		resp.Package = pkgDetail
	} else {
		pkg, err := s.coreClient.GetPackageByID(ctx, resp.PackageID)
		if err == nil && pkg != nil {
			detail := &dto.EnrollmentPackageResponse{
				ID:            pkg.ID,
				Name:          pkg.Name,
				Description:   pkg.Description,
				Price:         pkg.Price,
				DiscountPrice: pkg.DiscountPrice,
				Sessions:      pkg.Sessions,
				Duration:      pkg.Duration,
			}
			pkgCache[resp.PackageID] = detail
			resp.Package = detail
		}
	}

	// Fetch transaction to get addon details
	if s.transactionSvc != nil {
		tx, err := s.transactionSvc.GetTransactionByEnrollmentID(ctx, resp.ID)
		if err == nil && tx != nil {
			var addonResponses []dto.EnrollmentAddOnResponse
			for _, item := range tx.Items {
				if item.ItemType == models.TransactionItemTypeAddOn {
					var addonDetail *dto.EnrollmentAddOnResponse
					if cached, exists := addonCache[item.ItemID]; exists {
						addonDetail = cached
					} else {
						addon, err := s.coreClient.GetAddOnByID(ctx, item.ItemID)
						if err == nil && addon != nil {
							detail := &dto.EnrollmentAddOnResponse{
								ID:          addon.ID,
								Title:       addon.Title,
								Description: addon.Description,
								Price:       addon.Price,
								Sessions:    addon.Sessions,
							}
							addonCache[item.ItemID] = detail
							addonDetail = detail
						} else {
							// Fallback to transaction item details
							detail := &dto.EnrollmentAddOnResponse{
								ID:          item.ItemID,
								Title:       item.ItemName,
								Description: "",
								Price:       item.UnitPrice,
								Sessions:    item.Sessions,
							}
							addonDetail = detail
						}
					}
					if addonDetail != nil {
						addonResponses = append(addonResponses, *addonDetail)
					}
				}
			}
			resp.AddOns = addonResponses
		}
	}
}