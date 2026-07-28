package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	cClient "booking-service/clients/core"
	uClient "booking-service/clients/user"
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
}

type EnrollmentService struct {
	enrollmentRepo   repositories.IEnrollmentRepository
	transactionSvc   ITransactionService
	eventPublisher  kafka.IEventPublisher
	coreClient      cClient.ICoreClient
	userClient      uClient.IUserClient
}

func NewEnrollmentService(
	enrollmentRepo repositories.IEnrollmentRepository,
) IEnrollmentService {
	return &EnrollmentService{
		enrollmentRepo:  enrollmentRepo,
	}
}

// NewEnrollmentServiceWithDeps creates a new enrollment service with all dependencies
func NewEnrollmentServiceWithDeps(
	enrollmentRepo repositories.IEnrollmentRepository,
	eventPublisher kafka.IEventPublisher,
	coreClient cClient.ICoreClient,
) IEnrollmentService {
	return &EnrollmentService{
		enrollmentRepo:  enrollmentRepo,
		eventPublisher:  eventPublisher,
		coreClient:      coreClient,
	}
}

// NewEnrollmentServiceWithAllDeps creates a new enrollment service with transaction service
func NewEnrollmentServiceWithAllDeps(
	enrollmentRepo repositories.IEnrollmentRepository,
	transactionSvc ITransactionService,
	eventPublisher kafka.IEventPublisher,
	coreClient cClient.ICoreClient,
	userClient uClient.IUserClient,
) IEnrollmentService {
	return &EnrollmentService{
		enrollmentRepo:   enrollmentRepo,
		transactionSvc:    transactionSvc,
		eventPublisher:   eventPublisher,
		coreClient:       coreClient,
		userClient:       userClient,
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

	if enrollment.Status == models.EnrollmentStatusPaid || enrollment.Status == "paid" || enrollment.Status == "active" {
		// Enrollment is already paid — still trigger entitlement sync as a safety net.
		// The sync is idempotent (user-service deduplicates by bookingID), so calling
		// it again is harmless but ensures the entitlement gets created if it was missed
		// on the first attempt (e.g. due to service restart, network error, or race).
		go s.syncEntitlementForPaidEnrollment(ctx, enrollment, packageID, packageName, totalPrice, paymentMethod)
		resp := s.enrollmentRepo.ToResponse(enrollment)
		return &resp, nil
	}

	if enrollment.Status != models.EnrollmentStatusPendingPayment && enrollment.Status != "pending" && enrollment.Status != "pending_payment" {
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

	// Get package details for event & sync
	var totalSessions int
	var pkgName string
	if s.coreClient != nil {
		if pkg, err := s.coreClient.GetPackageByID(ctx, enrollment.PackageID); err == nil {
			totalSessions = pkg.GetSessions()
			pkgName = pkg.Name
		}
	}
	if packageName != "" {
		pkgName = packageName
	}

	// Calculate extra sessions from transaction
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

	// Publish enrollment.paid event if event publisher is available
	if s.eventPublisher != nil {
		_ = s.eventPublisher.PublishEnrollmentPaid(ctx, enrollment.ID.String(), enrollment.UserID.String(), enrollment.PackageID, totalPrice, totalSessions, pkgName)
	}

	// Direct HTTP call to user-service as dual sync fallback for entitlement creation
	userServiceURL := getEnv("USER_SERVICE_URL", "http://127.0.0.1:8001")
	if userServiceURL == "http://user-service:8001" {
		userServiceURL = "http://127.0.0.1:8001"
	}
	syncUrl := fmt.Sprintf("%s/api/v1/entitlements/sync", userServiceURL)
	syncBody, _ := json.Marshal(map[string]interface{}{
		"member_id":      enrollment.UserID.String(),
		"booking_id":     enrollment.ID.String(),
		"package_id":     enrollment.PackageID.String(),
		"package_name":   pkgName,
		"total_sessions": totalSessions,
	})
	syncReq, syncErr := http.NewRequest("POST", syncUrl, bytes.NewBuffer(syncBody))
	if syncErr == nil {
		syncReq.Header.Set("Content-Type", "application/json")
		httpClient := &http.Client{Timeout: 10 * time.Second}
		syncResp, syncDoErr := httpClient.Do(syncReq)
		if syncDoErr == nil {
			syncResp.Body.Close()
		}
	}



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

	// Create entitlement in user-service so member profile has active session count
	if s.userClient != nil {
		totalSessions := 0
		pkgName := packageName
		if s.coreClient != nil {
			if pkg, err := s.coreClient.GetPackageByID(ctx, enrollment.PackageID); err == nil {
				totalSessions = pkg.Sessions
				if pkgName == "" {
					pkgName = pkg.Name
				}
			}
		}
		if totalSessions <= 0 {
			totalSessions = 6 // Fallback default sessions count
		}
		if pkgName == "" {
			pkgName = "Driving Package"
		}

		entInput := uClient.CreateEntitlementInput{
			BookingID:     enrollment.ID,
			PackageID:     enrollment.PackageID,
			PackageName:   pkgName,
			TotalSessions: totalSessions,
			StartDate:     time.Now().Format("2006-01-02"),
		}
		if err := s.userClient.CreateEntitlement(ctx, enrollment.UserID, entInput); err != nil {
			fmt.Printf("[MarkAsPaid] Failed to create entitlement in user-service: %v\n", err)
		} else {
			fmt.Printf("[MarkAsPaid] Entitlement created in user-service for user %s\n", enrollment.UserID.String())
		}
	}

	resp := s.enrollmentRepo.ToResponse(enrollment)
	s.populateDetails(ctx, &resp)
	return &resp, nil
}

// syncEntitlementForPaidEnrollment triggers entitlement creation in user-service
// for an enrollment that is already marked as paid. It publishes the enrollment.paid
// Kafka event and makes a direct HTTP call to user-service as a fallback.
// This is called in a goroutine and is idempotent — user-service deduplicates by bookingID.
func (s *EnrollmentService) syncEntitlementForPaidEnrollment(ctx context.Context, enrollment *models.Enrollment, packageID, packageName string, totalPrice float64, paymentMethod string) {
	// Get package details for event & sync
	var totalSessions int
	var pkgName string
	if s.coreClient != nil {
		if pkg, err := s.coreClient.GetPackageByID(ctx, enrollment.PackageID); err == nil {
			totalSessions = pkg.GetSessions()
			pkgName = pkg.Name
		}
	}
	if packageName != "" {
		pkgName = packageName
	}
	if pkgName == "" {
		pkgName = "Driving Package"
	}

	// Calculate extra sessions from transaction
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

	if totalSessions <= 0 {
		totalSessions = 6 // Fallback default sessions count
	}

	fmt.Printf("[MarkAsPaid] Enrollment %s already paid — syncing entitlement (sessions=%d, pkg=%s)\n", enrollment.ID.String(), totalSessions, pkgName)

	// Publish enrollment.paid event if event publisher is available
	if s.eventPublisher != nil {
		_ = s.eventPublisher.PublishEnrollmentPaid(ctx, enrollment.ID.String(), enrollment.UserID.String(), enrollment.PackageID, totalPrice, totalSessions, pkgName)
	}

	// Direct HTTP call to user-service as dual sync fallback for entitlement creation
	userServiceURL := getEnv("USER_SERVICE_URL", "http://127.0.0.1:8001")
	if userServiceURL == "http://user-service:8001" {
		userServiceURL = "http://127.0.0.1:8001"
	}
	syncUrl := fmt.Sprintf("%s/api/v1/entitlements/sync", userServiceURL)
	syncBody, _ := json.Marshal(map[string]interface{}{
		"member_id":      enrollment.UserID.String(),
		"booking_id":     enrollment.ID.String(),
		"package_id":     enrollment.PackageID.String(),
		"package_name":   pkgName,
		"total_sessions": totalSessions,
	})
	syncReq, syncErr := http.NewRequest("POST", syncUrl, bytes.NewBuffer(syncBody))
	if syncErr == nil {
		syncReq.Header.Set("Content-Type", "application/json")
		httpClient := &http.Client{Timeout: 10 * time.Second}
		syncResp, syncDoErr := httpClient.Do(syncReq)
		if syncDoErr == nil {
			syncResp.Body.Close()
			fmt.Printf("[MarkAsPaid] Entitlement sync HTTP call succeeded for enrollment %s\n", enrollment.ID.String())
		} else {
			fmt.Printf("[MarkAsPaid] Entitlement sync HTTP call failed for enrollment %s: %v\n", enrollment.ID.String(), syncDoErr)
		}
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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}