package repositories

import (
	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

func NewRepositoryRegistry(db *gorm.DB) *Registry {
	return &Registry{db: db}
}

// GetPayment returns payment repository
func (r *Registry) GetPayment() IPaymentRepository {
	return NewPaymentRepository(r.db)
}

// GetTransaction returns transaction repository
func (r *Registry) GetTransaction() ITransactionRepository {
	return NewTransactionRepository(r.db)
}

// GetPaymentMethod returns payment method repository
func (r *Registry) GetPaymentMethod() IPaymentMethodRepository {
	return NewPaymentMethodRepository(r.db)
}

// GetRefund returns refund repository
func (r *Registry) GetRefund() IRefundRepository {
	return NewRefundRepository(r.db)
}

type IRepositoryRegistry interface {
	GetPayment() IPaymentRepository
	GetTransaction() ITransactionRepository
	GetPaymentMethod() IPaymentMethodRepository
	GetRefund() IRefundRepository
}