package repositories

import (
	"payment-service/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IPaymentMethodRepository interface {
	GetAll() ([]models.PaymentMethod, error)
	GetByID(id uint) (*models.PaymentMethod, error)
	GetByCode(code string) (*models.PaymentMethod, error)
	GetActive() ([]models.PaymentMethod, error)
}

type PaymentMethodRepository struct {
	db *gorm.DB
}

func NewPaymentMethodRepository(db *gorm.DB) IPaymentMethodRepository {
	return &PaymentMethodRepository{db: db}
}

func (r *PaymentMethodRepository) GetAll() ([]models.PaymentMethod, error) {
	var methods []models.PaymentMethod
	err := r.db.Order("sort_order ASC").Find(&methods).Error
	return methods, err
}

func (r *PaymentMethodRepository) GetByID(id uint) (*models.PaymentMethod, error) {
	var method models.PaymentMethod
	err := r.db.First(&method, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &method, nil
}

func (r *PaymentMethodRepository) GetByCode(code string) (*models.PaymentMethod, error) {
	var method models.PaymentMethod
	err := r.db.First(&method, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &method, nil
}

func (r *PaymentMethodRepository) GetActive() ([]models.PaymentMethod, error) {
	var methods []models.PaymentMethod
	err := r.db.Where("is_active = ?", true).
		Order("sort_order ASC").
		Find(&methods).Error
	return methods, err
}

// IRefundRepository interface
type IRefundRepository interface {
	Create(refund *models.Refund) error
	GetByID(id uuid.UUID) (*models.Refund, error)
	GetByPaymentID(paymentID uuid.UUID) ([]models.Refund, error)
	Update(refund *models.Refund) error
}

type RefundRepository struct {
	db *gorm.DB
}

func NewRefundRepository(db *gorm.DB) IRefundRepository {
	return &RefundRepository{db: db}
}

func (r *RefundRepository) Create(refund *models.Refund) error {
	return r.db.Create(refund).Error
}

func (r *RefundRepository) GetByID(id uuid.UUID) (*models.Refund, error) {
	var refund models.Refund
	err := r.db.First(&refund, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &refund, nil
}

func (r *RefundRepository) GetByPaymentID(paymentID uuid.UUID) ([]models.Refund, error) {
	var refunds []models.Refund
	err := r.db.Where("payment_id = ?", paymentID).
		Order("created_at DESC").
		Find(&refunds).Error
	return refunds, err
}

func (r *RefundRepository) Update(refund *models.Refund) error {
	return r.db.Save(refund).Error
}