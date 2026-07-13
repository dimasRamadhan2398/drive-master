package services

import (
	"payment-service/models"
)

// CheckoutResponse represents the response when initiating a checkout/payment link
type CheckoutResponse struct {
	Token       string `json:"token,omitempty"`
	RedirectURL string `json:"redirectUrl"`
	OrderID     string `json:"orderId"`
}

// NotificationPayload represents parsed notification/webhook details normalized across gateways
type NotificationPayload struct {
	OrderID           string
	TransactionID     string
	TransactionStatus models.PaymentStatus
	PaymentType       string
	GrossAmount       float64
	RawPayload        map[string]interface{}
}

// IPaymentGatewayService defines the unified interface for payment gateway integrations
type IPaymentGatewayService interface {
	CreateCheckout(orderID string, amount float64, packageName, customerName, customerEmail string) (*CheckoutResponse, error)
	GetTransactionStatus(orderID string) (models.PaymentStatus, error)
	VerifyNotification(headers map[string]string, rawBody []byte) (*NotificationPayload, error)
	GetName() string
}
