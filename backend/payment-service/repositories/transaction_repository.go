package repositories

import (
	"payment-service/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ITransactionRepository interface {
	Create(tx *models.Transaction) error
	GetByID(id uuid.UUID) (*models.Transaction, error)
	GetByPaymentID(paymentID uuid.UUID) ([]models.Transaction, error)
	Update(tx *models.Transaction) error
	ListWithPagination(offset, limit int) ([]models.Transaction, error)
	CountAll() (int64, error)
}

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) ITransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(tx *models.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *TransactionRepository) GetByID(id uuid.UUID) (*models.Transaction, error) {
	var tx models.Transaction
	err := r.db.First(&tx, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *TransactionRepository) GetByPaymentID(paymentID uuid.UUID) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Where("payment_id = ?", paymentID).
		Order("created_at DESC").
		Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) Update(tx *models.Transaction) error {
	return r.db.Save(tx).Error
}

func (r *TransactionRepository) ListWithPagination(offset, limit int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) CountAll() (int64, error) {
	var count int64
	err := r.db.Model(&models.Transaction{}).Count(&count).Error
	return count, err
}