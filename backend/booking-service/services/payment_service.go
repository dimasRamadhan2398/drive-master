package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"booking-service/models"
	"booking-service/models/dto"
	"booking-service/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// IPaymentService defines the payment service interface
type IPaymentService interface {
	CreatePayment(ctx context.Context, req dto.CreatePaymentRequest) (*dto.PaymentResponse, error)
	GetPayment(ctx context.Context, id uint) (*dto.PaymentResponse, error)
	GetPaymentByOrderID(ctx context.Context, orderID string) (*dto.PaymentResponse, error)
	GetPaymentDetail(ctx context.Context, orderID string) (*dto.PaymentDetailResponse, error)
	ListPayments(ctx context.Context, page, limit int) (*dto.PaymentListResponse, error)
	ListUserPayments(ctx context.Context, userID uuid.UUID, page, limit int) (*dto.PaymentListResponse, error)
	HandleCallback(ctx context.Context, callback dto.PaymentCallbackRequest) error
	CancelPayment(ctx context.Context, orderID string) error
	ExpirePendingPayments(ctx context.Context) error
}

// PaymentService handles payment operations
type PaymentService struct {
	paymentRepo     repositories.IPaymentRepository
	enrollmentRepo  repositories.IEnrollmentRepository
	transactionRepo repositories.ITransactionRepository
}

// NewPaymentService creates a new payment service
func NewPaymentService(
	paymentRepo repositories.IPaymentRepository,
	enrollmentRepo repositories.IEnrollmentRepository,
) IPaymentService {
	return &PaymentService{
		paymentRepo:    paymentRepo,
		enrollmentRepo: enrollmentRepo,
	}
}

// NewPaymentServiceWithTransaction creates a new payment service with transaction support
func NewPaymentServiceWithTransaction(
	paymentRepo repositories.IPaymentRepository,
	enrollmentRepo repositories.IEnrollmentRepository,
	transactionRepo repositories.ITransactionRepository,
) IPaymentService {
	return &PaymentService{
		paymentRepo:     paymentRepo,
		enrollmentRepo:  enrollmentRepo,
		transactionRepo: transactionRepo,
	}
}

// CreatePayment creates a new payment record
func (s *PaymentService) CreatePayment(ctx context.Context, req dto.CreatePaymentRequest) (*dto.PaymentResponse, error) {
	// Check if enrollment exists
	enrollment, err := s.enrollmentRepo.FindByID(ctx, req.EnrollmentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("enrollment not found")
		}
		return nil, err
	}

	// If transaction repo is available, validate amount against transaction
	var expectedAmount float64
	if s.transactionRepo != nil {
		tx, err := s.transactionRepo.FindByEnrollmentID(ctx, req.EnrollmentID)
		if err == nil && tx != nil {
			expectedAmount = tx.TotalAmount
			// Validate that the payment amount matches the transaction amount
			if req.Amount != expectedAmount {
				return nil, fmt.Errorf("payment amount (%.2f) does not match transaction amount (%.2f)", req.Amount, expectedAmount)
			}
		} else {
			// Fallback to enrollment total price if no transaction exists
			expectedAmount = enrollment.TotalPrice
		}
	} else {
		expectedAmount = enrollment.TotalPrice
	}

	// Generate unique order ID
	orderID := fmt.Sprintf("ORD-%s-%d", time.Now().Format("20060102150405"), uuid.New().ID())

	// Create payment record with pending status
	payment := &models.Payment{
		EnrollmentID:  req.EnrollmentID,
		UserID:        req.UserID,
		OrderID:       orderID,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		Status:        dto.PaymentStatusPending,
		ExpiresAt:     timePtr(time.Now().Add(24 * time.Hour)), // 24 hours expiry
	}

	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, err
	}

	// Generate mock payment URL (in production, this would call Midtrans)
	payment.PaymentURL = generatePaymentURL(orderID, req.PaymentMethod)
	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		fmt.Printf("Failed to update payment URL: %v\n", err)
	}

	resp := s.paymentRepo.ToResponse(payment)
	return &resp, nil
}

// GetPayment retrieves a payment by ID
func (s *PaymentService) GetPayment(ctx context.Context, id uint) (*dto.PaymentResponse, error) {
	payment, err := s.paymentRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}

	resp := s.paymentRepo.ToResponse(payment)
	return &resp, nil
}

// GetPaymentByOrderID retrieves a payment by order ID
func (s *PaymentService) GetPaymentByOrderID(ctx context.Context, orderID string) (*dto.PaymentResponse, error) {
	payment, err := s.paymentRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}

	resp := s.paymentRepo.ToResponse(payment)
	return &resp, nil
}

// GetPaymentDetail retrieves detailed payment information by order ID
func (s *PaymentService) GetPaymentDetail(ctx context.Context, orderID string) (*dto.PaymentDetailResponse, error) {
	payment, err := s.paymentRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}

	return &dto.PaymentDetailResponse{
		OrderID:       payment.OrderID,
		PaymentType:   payment.PaymentMethod,
		TransactionID: payment.TransactionID,
		GrossAmount:   payment.Amount,
		Status:        payment.Status,
		PaymentURL:    payment.PaymentURL,
		ExpiryTime:    payment.ExpiresAt,
	}, nil
}

// ListPayments retrieves all payments with pagination
func (s *PaymentService) ListPayments(ctx context.Context, page, limit int) (*dto.PaymentListResponse, error) {
	payments, err := s.paymentRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	total, err := s.paymentRepo.CountAll(ctx)
	if err != nil {
		return nil, err
	}

	resp := s.paymentRepo.ToListResponse(payments, total, page, limit)
	return &resp, nil
}

// ListUserPayments retrieves payments for a specific user
func (s *PaymentService) ListUserPayments(ctx context.Context, userID uuid.UUID, page, limit int) (*dto.PaymentListResponse, error) {
	payments, err := s.paymentRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	total := int64(len(payments))

	resp := s.paymentRepo.ToListResponse(payments, total, page, limit)
	return &resp, nil
}

// HandleCallback processes payment callback notifications
func (s *PaymentService) HandleCallback(ctx context.Context, callback dto.PaymentCallbackRequest) error {
	// Find payment by order ID
	payment, err := s.paymentRepo.FindByOrderID(ctx, callback.OrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("payment not found")
		}
		return err
	}

	// Map Midtrans status to our status
	var newStatus dto.PaymentStatus
	switch callback.TransactionStatus {
	case "capture", "settlement":
		newStatus = dto.PaymentStatusPaid
		now := time.Now()
		payment.PaidAt = &now
	case "cancel", "deny":
		newStatus = dto.PaymentStatusCancelled
	case "expire":
		newStatus = dto.PaymentStatusExpired
	case "pending":
		newStatus = dto.PaymentStatusPending
	case "failure":
		newStatus = dto.PaymentStatusFailed
	default:
		newStatus = dto.PaymentStatusPending
	}

	// Update payment status
	if err := s.paymentRepo.UpdateStatus(ctx, payment.ID, newStatus); err != nil {
		return err
	}

	// If payment is successful, update enrollment and transaction status
	if newStatus == dto.PaymentStatusPaid {
		if err := s.enrollmentRepo.MarkAsPaid(ctx, payment.EnrollmentID, time.Now(), payment.Amount); err != nil {
			// Log error but don't fail the callback
			fmt.Printf("Failed to mark enrollment as paid: %v\n", err)
		}

		// Also mark transaction as paid if transaction repo is available
		if s.transactionRepo != nil {
			tx, err := s.transactionRepo.FindByEnrollmentID(ctx, payment.EnrollmentID)
			if err == nil && tx != nil {
				if err := s.transactionRepo.MarkAsPaid(ctx, tx.ID); err != nil {
					fmt.Printf("Failed to mark transaction as paid: %v\n", err)
				}
			}
		}
	}

	// Store transaction ID if available
	if callback.TransactionStatus == "settlement" || callback.TransactionStatus == "capture" {
		transactionID := fmt.Sprintf("%s-%s", callback.OrderID, callback.PaymentType)
		if err := s.paymentRepo.UpdateStatusByOrderID(ctx, callback.OrderID, newStatus, transactionID); err != nil {
			fmt.Printf("Failed to update transaction ID: %v\n", err)
		}
	}

	return nil
}

// CancelPayment cancels a pending payment
func (s *PaymentService) CancelPayment(ctx context.Context, orderID string) error {
	payment, err := s.paymentRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("payment not found")
		}
		return err
	}

	if payment.Status != dto.PaymentStatusPending {
		return errors.New("only pending payments can be cancelled")
	}

	return s.paymentRepo.UpdateStatus(ctx, payment.ID, dto.PaymentStatusCancelled)
}

// ExpirePendingPayments expires payments that have passed their expiry time
func (s *PaymentService) ExpirePendingPayments(ctx context.Context) error {
	pendingPayments, err := s.paymentRepo.FindByStatus(ctx, dto.PaymentStatusPending)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, payment := range pendingPayments {
		// Check if payment has expired (24 hours)
		if payment.ExpiresAt != nil && payment.ExpiresAt.Before(now) {
			if err := s.paymentRepo.UpdateStatus(ctx, payment.ID, dto.PaymentStatusExpired); err != nil {
				fmt.Printf("Failed to expire payment %d: %v\n", payment.ID, err)
			}
		}
	}

	return nil
}

// Helper function to create time pointer
func timePtr(t time.Time) *time.Time {
	return &t
}

// generatePaymentURL generates a mock payment URL
// In production, this would call Midtrans Snap API
func generatePaymentURL(orderID string, method dto.PaymentMethod) string {
	return fmt.Sprintf("https://app.midtrans.com/mock/%s?method=%s", orderID, method)
}