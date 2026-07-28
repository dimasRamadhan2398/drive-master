package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"booking-service/models"
	"booking-service/models/dto"
	"booking-service/pkg/base"
	apperrors "booking-service/pkg/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ITransactionRepository interface {
	Create(ctx context.Context, tx *models.Transaction) error
	CreateWithItems(ctx context.Context, tx *models.Transaction, items []models.TransactionItem) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error)
	FindByEnrollmentID(ctx context.Context, enrollmentID uuid.UUID) (*models.Transaction, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]models.Transaction, int64, error)
	Update(ctx context.Context, tx *models.Transaction) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.TransactionStatus) error
	MarkAsPaid(ctx context.Context, id uuid.UUID) error
	ExistsByEnrollmentID(ctx context.Context, enrollmentID uuid.UUID) (bool, error)
	ToResponse(tx *models.Transaction) dto.TransactionResponse
	ToListResponse(txs []models.Transaction, total int64, page, limit int) dto.TransactionListResponse
}

type TransactionRepository struct {
	*base.BaseRepository
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) ITransactionRepository {
	return &TransactionRepository{
		BaseRepository: base.NewBaseRepository(db),
		db:              db,
	}
}

// Create creates a new transaction
func (r *TransactionRepository) Create(ctx context.Context, tx *models.Transaction) error {
	return r.BaseRepository.Create(tx)
}

// CreateWithItems creates a transaction with its items in a transaction
func (r *TransactionRepository) CreateWithItems(ctx context.Context, tx *models.Transaction, items []models.TransactionItem) error {
	return r.db.Transaction(func(dbTx *gorm.DB) error {
		// Create transaction
		if err := dbTx.Create(tx).Error; err != nil {
			return err
		}

		// Create items with transaction ID
		for i := range items {
			items[i].TransactionID = tx.ID
			if items[i].ID == uuid.Nil {
				items[i].ID = uuid.New()
			}
			if err := dbTx.Create(&items[i]).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// FindByID finds a transaction by ID with items
func (r *TransactionRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	var tx models.Transaction
	if err := r.BaseRepository.FindByIDWithPreload(&tx, id, "Items"); err != nil {
		return nil, err
	}
	return &tx, nil
}

// FindByEnrollmentID finds a transaction by enrollment ID with items
func (r *TransactionRepository) FindByEnrollmentID(ctx context.Context, enrollmentID uuid.UUID) (*models.Transaction, error) {
	var tx models.Transaction
	if err := r.db.Preload("Items").Where("enrollment_id = ?", enrollmentID).First(&tx).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", apperrors.ErrDatabase, err)
	}
	return &tx, nil
}

// FindByUserID finds transactions by user ID with pagination
func (r *TransactionRepository) FindByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]models.Transaction, int64, error) {
	var txs []models.Transaction

	// Count total
	total, err := r.BaseRepository.Count(&models.Transaction{}, base.NewQueryOptions().WithWhere(map[string]any{"user_id": userID}))
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	opts := base.NewQueryOptions().
		WithPagination(offset, limit).
		WithPreloads("Items").
		WithOrder("created_at DESC").
		WithWhere(map[string]any{"user_id": userID})

	if err := r.BaseRepository.FindMany(&models.Transaction{}, &txs, opts); err != nil {
		return nil, 0, err
	}

	return txs, total, nil
}

// Update updates a transaction
func (r *TransactionRepository) Update(ctx context.Context, tx *models.Transaction) error {
	return r.BaseRepository.Update(tx)
}

// UpdateStatus updates the transaction status
func (r *TransactionRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.TransactionStatus) error {
	return r.BaseRepository.Exec(
		"UPDATE transactions SET status = ?, updated_at = ? WHERE id = ?",
		status, time.Now(), id,
	)
}

// MarkAsPaid marks a transaction as paid
func (r *TransactionRepository) MarkAsPaid(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.BaseRepository.Exec(
		"UPDATE transactions SET status = ?, paid_at = ?, updated_at = ? WHERE id = ?",
		models.TransactionStatusCompleted, now, now, id,
	)
}

// ExistsByEnrollmentID checks if a transaction exists for an enrollment
func (r *TransactionRepository) ExistsByEnrollmentID(ctx context.Context, enrollmentID uuid.UUID) (bool, error) {
	return r.BaseRepository.Exists(&models.Transaction{}, "enrollment_id = ?", enrollmentID)
}

// ToResponse converts a Transaction to TransactionResponse DTO
func (r *TransactionRepository) ToResponse(tx *models.Transaction) dto.TransactionResponse {
	items := make([]dto.TransactionItemResponse, len(tx.Items))
	for i, item := range tx.Items {
		items[i] = dto.TransactionItemResponse{
			ID:          item.ID,
			ItemType:    string(item.ItemType),
			ItemID:      item.ItemID,
			ItemName:    item.ItemName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Subtotal:    item.Subtotal,
			Sessions:    item.Sessions,
		}
	}

	resp := dto.TransactionResponse{
		ID:           tx.ID,
		EnrollmentID: tx.EnrollmentID,
		UserID:       tx.UserID,
		BasePrice:    tx.BasePrice,
		AddOnsTotal:  tx.AddOnsTotal,
		TotalAmount:  tx.TotalAmount,
		Status:       string(tx.Status),
		PaidAt:       tx.PaidAt,
		Items:        items,
		CreatedAt:    tx.CreatedAt,
		UpdatedAt:    tx.UpdatedAt,
	}
	return resp
}

// ToListResponse converts a slice of Transactions to TransactionListResponse DTO
func (r *TransactionRepository) ToListResponse(txs []models.Transaction, total int64, page, limit int) dto.TransactionListResponse {
	items := make([]dto.TransactionResponse, len(txs))
	for i, tx := range txs {
		items[i] = r.ToResponse(&tx)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return dto.TransactionListResponse{
		Data:       items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}

// Helper function to find by ID with preloads
func (r *TransactionRepository) FindByIDWithPreload(entity any, id any, preloads ...string) error {
	db := r.db
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	if err := db.First(entity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return apperrors.ErrNotFound
		}
		return err
	}
	return nil
}
