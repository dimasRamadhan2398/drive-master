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

type AddOnInfo struct {
	ID       uuid.UUID
	Name     string
	Price    float64
	Sessions int
}

type ITransactionService interface {
	CreateTransaction(ctx context.Context, enrollmentID, userID uuid.UUID, packageID uuid.UUID, packageName string, packagePrice float64, addOns []AddOnInfo) (*models.Transaction, error)
	GetTransactionByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error)
	GetTransactionByEnrollmentID(ctx context.Context, enrollmentID uuid.UUID) (*models.Transaction, error)
	GetUserTransactions(ctx context.Context, userID uuid.UUID, page, limit int) (*dto.TransactionListResponse, error)
	MarkAsPaid(ctx context.Context, id uuid.UUID) error
	ExistsForEnrollment(ctx context.Context, enrollmentID uuid.UUID) (bool, error)
}

type TransactionService struct {
	transactionRepo repositories.ITransactionRepository
}

func NewTransactionService(transactionRepo repositories.ITransactionRepository) ITransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
	}
}

// CreateTransaction creates a new transaction with line items for an enrollment
func (s *TransactionService) CreateTransaction(
	ctx context.Context,
	enrollmentID, userID uuid.UUID,
	packageID uuid.UUID, packageName string, packagePrice float64,
	addOns []AddOnInfo,
) (*models.Transaction, error) {
	// Check if transaction already exists for this enrollment
	exists, err := s.transactionRepo.ExistsByEnrollmentID(ctx, enrollmentID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("transaction already exists for this enrollment")
	}

	// Calculate totals
	var addOnsTotal float64
	items := make([]models.TransactionItem, 0, 1+len(addOns))

	// Add package as first line item
	items = append(items, models.TransactionItem{
		ItemType:  models.TransactionItemTypePackage,
		ItemID:    packageID,
		ItemName:  packageName,
		Quantity:  1,
		UnitPrice: packagePrice,
		Subtotal:  packagePrice,
	})

	// Add add-ons as line items
	for _, addon := range addOns {
		subtotal := addon.Price
		addOnsTotal += subtotal
		items = append(items, models.TransactionItem{
			ItemType:  models.TransactionItemTypeAddOn,
			ItemID:    addon.ID,
			ItemName:  addon.Name,
			Quantity:  1,
			UnitPrice: addon.Price,
			Subtotal:  subtotal,
			Sessions:  addon.Sessions,
		})
	}

	// Create transaction
	transaction := &models.Transaction{
		EnrollmentID: enrollmentID,
		UserID:      userID,
		BasePrice:   packagePrice,
		AddOnsTotal: addOnsTotal,
		TotalAmount: packagePrice + addOnsTotal,
		Status:      models.TransactionStatusPending,
		Items:       items,
	}

	if err := s.transactionRepo.CreateWithItems(ctx, transaction, items); err != nil {
		return nil, err
	}

	return transaction, nil
}

// GetTransactionByID retrieves a transaction by ID
func (s *TransactionService) GetTransactionByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	return s.transactionRepo.FindByID(ctx, id)
}

// GetTransactionByEnrollmentID retrieves a transaction by enrollment ID
func (s *TransactionService) GetTransactionByEnrollmentID(ctx context.Context, enrollmentID uuid.UUID) (*models.Transaction, error) {
	return s.transactionRepo.FindByEnrollmentID(ctx, enrollmentID)
}

// GetUserTransactions retrieves transactions for a user with pagination
func (s *TransactionService) GetUserTransactions(ctx context.Context, userID uuid.UUID, page, limit int) (*dto.TransactionListResponse, error) {
	transactions, total, err := s.transactionRepo.FindByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}
	resp := s.transactionRepo.ToListResponse(transactions, total, page, limit)
	return &resp, nil
}

// MarkAsPaid marks a transaction as paid
func (s *TransactionService) MarkAsPaid(ctx context.Context, id uuid.UUID) error {
	// Verify transaction exists
	_, err := s.transactionRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("transaction not found")
		}
		return err
	}

	// Update status
	if err := s.transactionRepo.MarkAsPaid(ctx, id); err != nil {
		return err
	}

	return nil
}

// ExistsForEnrollment checks if a transaction exists for an enrollment
func (s *TransactionService) ExistsForEnrollment(ctx context.Context, enrollmentID uuid.UUID) (bool, error) {
	return s.transactionRepo.ExistsByEnrollmentID(ctx, enrollmentID)
}
