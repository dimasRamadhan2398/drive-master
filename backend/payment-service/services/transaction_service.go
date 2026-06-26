package services

import (
	"encoding/json"
	"payment-service/models"
	"payment-service/models/dto"
	"payment-service/repositories"
	"time"

	"github.com/google/uuid"
)

type ITransactionService interface {
	// Core operations
	CreateTransaction(paymentID uuid.UUID, txType models.TransactionType, amount float64, currency string, gateway, gatewayTxnID string, response interface{}) (*models.Transaction, error)
	GetTransaction(id uuid.UUID) (*models.Transaction, error)
	GetTransactionsByPaymentID(paymentID uuid.UUID) ([]models.Transaction, error)
	ListTransactions(page, pageSize int) ([]models.Transaction, int64, error)

	// Transaction status updates
	UpdateTransactionStatus(id uuid.UUID, status models.TransactionStatus, errorCode, errorMsg string) error
	MarkTransactionSuccess(id uuid.UUID, gatewayTxnID string, response interface{}) error
	MarkTransactionFailed(id uuid.UUID, errorCode, errorMsg string) error
}

type TransactionService struct {
	repo repositories.ITransactionRepository
}

func NewTransactionService(repo repositories.ITransactionRepository) ITransactionService {
	return &TransactionService{repo: repo}
}

// CreateTransaction creates a new transaction record
func (s *TransactionService) CreateTransaction(
	paymentID uuid.UUID,
	txType models.TransactionType,
	amount float64,
	currency string,
	gateway string,
	gatewayTxnID string,
	response interface{},
) (*models.Transaction, error) {
	gatewayResponse := "{}"
	if response != nil {
		respBytes, err := json.Marshal(response)
		if err == nil {
			gatewayResponse = string(respBytes)
		}
	}

	tx := &models.Transaction{
		ID:              uuid.New(),
		PaymentID:       paymentID,
		Type:            txType,
		Status:          models.TransactionStatusPending,
		Amount:          amount,
		Currency:        currency,
		Gateway:         gateway,
		GatewayTxnID:    gatewayTxnID,
		GatewayResponse: gatewayResponse,
		CreatedAt:      time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.repo.Create(tx); err != nil {
		return nil, err
	}

	return tx, nil
}

// GetTransaction retrieves a transaction by ID
func (s *TransactionService) GetTransaction(id uuid.UUID) (*models.Transaction, error) {
	return s.repo.GetByID(id)
}

// GetTransactionsByPaymentID retrieves all transactions for a payment
func (s *TransactionService) GetTransactionsByPaymentID(paymentID uuid.UUID) ([]models.Transaction, error) {
	return s.repo.GetByPaymentID(paymentID)
}

// ListTransactions retrieves transactions with pagination
func (s *TransactionService) ListTransactions(page, pageSize int) ([]models.Transaction, int64, error) {
	offset := (page - 1) * pageSize
	transactions, err := s.repo.ListWithPagination(offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountAll()
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// UpdateTransactionStatus updates the status of a transaction
func (s *TransactionService) UpdateTransactionStatus(id uuid.UUID, status models.TransactionStatus, errorCode, errorMsg string) error {
	tx, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	tx.Status = status
	tx.ErrorCode = errorCode
	tx.ErrorMessage = errorMsg
	tx.UpdatedAt = time.Now()

	if status == models.TransactionStatusSuccess || status == models.TransactionStatusFailed {
		now := time.Now()
		tx.ProcessedAt = &now
	}

	return s.repo.Update(tx)
}

// MarkTransactionSuccess marks a transaction as successful
func (s *TransactionService) MarkTransactionSuccess(id uuid.UUID, gatewayTxnID string, response interface{}) error {
	tx, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	tx.Status = models.TransactionStatusSuccess
	tx.GatewayTxnID = gatewayTxnID
	tx.ProcessedAt = func() *time.Time { now := time.Now(); return &now }()
	tx.UpdatedAt = time.Now()

	if response != nil {
		if respBytes, err := json.Marshal(response); err == nil {
			tx.GatewayResponse = string(respBytes)
		}
	}

	return s.repo.Update(tx)
}

// MarkTransactionFailed marks a transaction as failed
func (s *TransactionService) MarkTransactionFailed(id uuid.UUID, errorCode, errorMsg string) error {
	tx, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	tx.Status = models.TransactionStatusFailed
	tx.ErrorCode = errorCode
	tx.ErrorMessage = errorMsg
	tx.ProcessedAt = func() *time.Time { now := time.Now(); return &now }()
	tx.UpdatedAt = time.Now()

	return s.repo.Update(tx)
}

// ToDTO converts a Transaction to TransactionDTO
func (s *TransactionService) ToDTO(tx *models.Transaction) dto.TransactionDTO {
	return dto.TransactionDTO{
		ID:           tx.ID,
		PaymentID:    tx.PaymentID,
		Type:         tx.Type,
		Status:       tx.Status,
		Amount:       tx.Amount,
		Currency:     tx.Currency,
		Gateway:      tx.Gateway,
		GatewayTxnID: tx.GatewayTxnID,
		ErrorCode:    tx.ErrorCode,
		ErrorMessage: tx.ErrorMessage,
		ProcessedAt:  tx.ProcessedAt,
		CreatedAt:    tx.CreatedAt,
	}
}