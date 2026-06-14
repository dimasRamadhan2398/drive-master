package repositories

import (
	"context"
	"time"

	"booking-service/models"
	"booking-service/models/dto"
	"booking-service/pkg/base"

	"gorm.io/gorm"
)

type IPaymentRepository interface {
	Create(ctx context.Context, payment *models.Payment) error
	CreateTx(tx *gorm.DB, payment *models.Payment) error
	FindByID(ctx context.Context, id uint) (*models.Payment, error)
	FindByOrderID(ctx context.Context, orderID string) (*models.Payment, error)
	Update(ctx context.Context, payment *models.Payment) error
	Delete(ctx context.Context, payment *models.Payment) error
	FindAll(ctx context.Context) ([]models.Payment, error)
	FindByEnrollmentID(ctx context.Context, enrollmentID uint) ([]models.Payment, error)
	FindByUserID(ctx context.Context, userID uint) ([]models.Payment, error)
	FindByStatus(ctx context.Context, status dto.PaymentStatus) ([]models.Payment, error)
	UpdateStatus(ctx context.Context, id uint, status dto.PaymentStatus) error
	UpdateStatusByOrderID(ctx context.Context, orderID string, status dto.PaymentStatus, transactionID string) error
	CountAll(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, status dto.PaymentStatus) (int64, error)
	Exists(ctx context.Context, condition any, args ...any) (bool, error)
	ToResponse(payment *models.Payment) dto.PaymentResponse
	ToListResponse(payments []models.Payment, total int64, page, limit int) dto.PaymentListResponse
}

type PaymentRepository struct {
	*base.BaseRepository
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) IPaymentRepository {
	return &PaymentRepository{BaseRepository: base.NewBaseRepository(db), db: db}
}

func (r *PaymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	return r.BaseRepository.Create(payment)
}

func (r *PaymentRepository) CreateTx(tx *gorm.DB, payment *models.Payment) error {
	return r.BaseRepository.CreateTx(tx, payment)
}

func (r *PaymentRepository) FindByID(ctx context.Context, id uint) (*models.Payment, error) {
	var payment models.Payment
	if err := r.BaseRepository.FindByIDWithPreload(&payment, id); err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) FindByOrderID(ctx context.Context, orderID string) (*models.Payment, error) {
	var payment models.Payment
	opts := base.NewQueryOptions().WithWhere(map[string]any{"order_id": orderID})
	if err := r.BaseRepository.FindOne(&models.Payment{}, &payment, opts); err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) Update(ctx context.Context, payment *models.Payment) error {
	return r.BaseRepository.Update(payment)
}

func (r *PaymentRepository) Delete(ctx context.Context, payment *models.Payment) error {
	return r.BaseRepository.Delete(payment)
}

func (r *PaymentRepository) FindAll(ctx context.Context) ([]models.Payment, error) {
	var payments []models.Payment
	opts := base.NewQueryOptions().WithOrder("created_at DESC")
	if err := r.BaseRepository.FindMany(&models.Payment{}, &payments, opts); err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *PaymentRepository) FindByEnrollmentID(ctx context.Context, enrollmentID uint) ([]models.Payment, error) {
	var payments []models.Payment
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"enrollment_id": enrollmentID}).
		WithOrder("created_at DESC")
	if err := r.BaseRepository.FindMany(&models.Payment{}, &payments, opts); err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *PaymentRepository) FindByUserID(ctx context.Context, userID uint) ([]models.Payment, error) {
	var payments []models.Payment
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"user_id": userID}).
		WithOrder("created_at DESC")
	if err := r.BaseRepository.FindMany(&models.Payment{}, &payments, opts); err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *PaymentRepository) FindByStatus(ctx context.Context, status dto.PaymentStatus) ([]models.Payment, error) {
	var payments []models.Payment
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"status": status}).
		WithOrder("created_at DESC")
	if err := r.BaseRepository.FindMany(&models.Payment{}, &payments, opts); err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *PaymentRepository) UpdateStatus(ctx context.Context, id uint, status dto.PaymentStatus) error {
	return r.BaseRepository.Exec(
		"UPDATE payments SET status = ?, updated_at = ? WHERE id = ?",
		status, time.Now(), id,
	)
}

func (r *PaymentRepository) UpdateStatusByOrderID(ctx context.Context, orderID string, status dto.PaymentStatus, transactionID string) error {
	return r.BaseRepository.Exec(
		"UPDATE payments SET status = ?, transaction_id = ?, updated_at = ? WHERE order_id = ?",
		status, transactionID, time.Now(), orderID,
	)
}

func (r *PaymentRepository) CountAll(ctx context.Context) (int64, error) {
	return r.BaseRepository.Count(&models.Payment{}, base.NewQueryOptions())
}

func (r *PaymentRepository) CountByStatus(ctx context.Context, status dto.PaymentStatus) (int64, error) {
	opts := base.NewQueryOptions().WithWhere(map[string]any{"status": status})
	return r.BaseRepository.Count(&models.Payment{}, opts)
}

func (r *PaymentRepository) Exists(ctx context.Context, condition any, args ...any) (bool, error) {
	return r.BaseRepository.Exists(&models.Payment{}, condition, args...)
}

// ToResponse converts a Payment model to PaymentResponse DTO
func (r *PaymentRepository) ToResponse(payment *models.Payment) dto.PaymentResponse {
	return dto.PaymentResponse{
		ID:            payment.ID,
		EnrollmentID:  payment.EnrollmentID,
		UserID:        payment.UserID,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount,
		PaymentMethod: payment.PaymentMethod,
		Status:        payment.Status,
		PaymentURL:    payment.PaymentURL,
		PaidAt:        payment.PaidAt,
		ExpiresAt:     payment.ExpiresAt,
		CreatedAt:     payment.CreatedAt,
		UpdatedAt:     payment.UpdatedAt,
	}
}

// ToListResponse converts a slice of Payments to PaymentListResponse DTO
func (r *PaymentRepository) ToListResponse(payments []models.Payment, total int64, page, limit int) dto.PaymentListResponse {
	items := make([]dto.PaymentResponse, len(payments))
	for i, p := range payments {
		items[i] = r.ToResponse(&p)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return dto.PaymentListResponse{
		Data:       items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}