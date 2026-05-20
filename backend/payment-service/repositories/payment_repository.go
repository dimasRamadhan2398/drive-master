package repositories

import (
	"payment-service/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IPaymentRepository interface {
	Create(payment *models.Payment) error
	GetByID(id uuid.UUID) (*models.Payment, error)
	GetByOrderID(orderID string) (*models.Payment, error)
	GetByBookingID(bookingID uuid.UUID) (*models.Payment, error)
	Update(payment *models.Payment) error
	List(filter ListPaymentsFilter, page, pageSize int) ([]models.Payment, int64, error)
	GetByUserID(userID uuid.UUID, page, pageSize int) ([]models.Payment, int64, error)
	GetExpiredPendingPayments() ([]models.Payment, error)
}

type ListPaymentsFilter struct {
	UserID    *uuid.UUID
	BookingID *uuid.UUID
	Status    models.PaymentStatus
}

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) IPaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *PaymentRepository) GetByID(id uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("PaymentMethod").Preload("Transactions").
		First(&payment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) GetByOrderID(orderID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("PaymentMethod").Preload("Transactions").
		First(&payment, "order_id = ?", orderID).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) GetByBookingID(bookingID uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("PaymentMethod").Preload("Transactions").
		First(&payment, "booking_id = ?", bookingID).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) Update(payment *models.Payment) error {
	return r.db.Save(payment).Error
}

func (r *PaymentRepository) List(filter ListPaymentsFilter, page, pageSize int) ([]models.Payment, int64, error) {
	var payments []models.Payment
	var total int64

	query := r.db.Model(&models.Payment{})

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.BookingID != nil {
		query = query.Where("booking_id = ?", *filter.BookingID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Preload("PaymentMethod").Preload("Transactions").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&payments).Error

	return payments, total, err
}

func (r *PaymentRepository) GetByUserID(userID uuid.UUID, page, pageSize int) ([]models.Payment, int64, error) {
	var payments []models.Payment
	var total int64

	offset := (page - 1) * pageSize

	r.db.Model(&models.Payment{}).Where("user_id = ?", userID).Count(&total)

	err := r.db.Preload("PaymentMethod").Preload("Transactions").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&payments).Error

	return payments, total, err
}

func (r *PaymentRepository) GetExpiredPendingPayments() ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.
		Where("status = ? AND expiry_time < NOW()", models.PaymentStatusPending).
		Find(&payments).Error
	return payments, err
}