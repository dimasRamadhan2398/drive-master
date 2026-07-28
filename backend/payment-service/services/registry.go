package services

import (
	"strings"
	"payment-service/pkg/config"
	"payment-service/repositories"
)

type Registry struct {
	repos          *repositories.Registry
	paymentGateway IPaymentGatewayService
	txService      ITransactionService
}

func NewServiceRegistry(repos *repositories.Registry, cfg *config.Config) *Registry {
	var gateway IPaymentGatewayService
	if strings.ToLower(cfg.App.PaymentGateway) == "doku" {
		gateway = NewDokuService(&cfg.Doku)
	} else if strings.ToLower(cfg.App.PaymentGateway) == "pakasir" {
		gateway = NewPakasirService(&cfg.Pakasir, repos.GetPayment())
	} else {
		gateway = NewMidtransService(&cfg.Midtrans)
	}

	return &Registry{
		repos:          repos,
		paymentGateway: gateway,
		txService:      NewTransactionService(repos.GetTransaction()),
	}
}

// GetPaymentGateway returns the active payment gateway service
func (r *Registry) GetPaymentGateway() IPaymentGatewayService {
	return r.paymentGateway
}

// GetTransaction returns the transaction service
func (r *Registry) GetTransaction() ITransactionService {
	return r.txService
}

type IServiceRegistry interface {
	GetPaymentGateway() IPaymentGatewayService
	GetTransaction() ITransactionService
}