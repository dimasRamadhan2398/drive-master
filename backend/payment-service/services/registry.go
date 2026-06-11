package services

import (
	"payment-service/pkg/config"
	"payment-service/repositories"
)

type Registry struct {
	repos     *repositories.Registry
	midtrans  *MidtransService
	txService ITransactionService
}

func NewServiceRegistry(repos *repositories.Registry, cfg *config.MidtransConfig) *Registry {
	return &Registry{
		repos:     repos,
		midtrans:  NewMidtransService(cfg),
		txService: NewTransactionService(repos.GetTransaction()),
	}
}

// GetMidtrans returns the Midtrans service
func (r *Registry) GetMidtrans() IMidtransService {
	return r.midtrans
}

// GetTransaction returns the transaction service
func (r *Registry) GetTransaction() ITransactionService {
	return r.txService
}

type IServiceRegistry interface {
	GetMidtrans() IMidtransService
	GetTransaction() ITransactionService
}