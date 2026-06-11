package services

import (
	"encoding/json"
	"fmt"
	"log"
	"payment-service/models"
	"payment-service/pkg/config"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

// PaymentType represents the type of payment method
type PaymentType string

const (
	PaymentTypeQRIS        PaymentType = "qris"
	PaymentTypeVA         PaymentType = "bank_transfer"
	PaymentTypeEWallet     PaymentType = "ewallet"
	PaymentTypeCreditCard  PaymentType = "credit_card"
)

// MidtransChargeResponse represents the response from Midtrans charge API
type MidtransChargeResponse struct {
	TransactionID   string                 `json:"transaction_id"`
	OrderID         string                 `json:"order_id"`
	PaymentType     string                 `json:"payment_type"`
	StatusCode      string                 `json:"status_code"`
	StatusMessage   string                 `json:"status_message"`
	TransactionTime  string                 `json:"transaction_time"`
	TransactionUUID string                 `json:"transaction_uuid"`
	MerchantID      string                 `json:"merchant_id"`
	FraudStatus     string                 `json:"fraud_status"`
	// QRIS specific
	QrCodeURL   string `json:"qr_code_url,omitempty"`
	QrCodeImage string `json:"qr_code_image,omitempty"`
	// VA specific
	VANumbers        []VANumber `json:"va_numbers,omitempty"`
	PermataVANumber  string     `json:"permata_va_number,omitempty"`
	BcaVANumber      string     `json:"bca_va_number,omitempty"`
	MandiriVANumber  string     `json:"mandiri_va_number,omitempty"`
	BniVANumber      string     `json:"bni_va_number,omitempty"`
	// E-wallet specific
	PaymentURL       string `json:"payment_url,omitempty"`
	AcquisitionDate  string `json:"acquisition_date,omitempty"`
	// Common
	GrossAmount      string `json:"gross_amount"`
}

// VANumber represents a Virtual Account number
type VANumber struct {
	Bank     string `json:"bank"`
	VANumber string `json:"va_number"`
}

// MidtransStatusResponse represents Midtrans transaction status
type MidtransStatusResponse struct {
	TransactionID     string `json:"transaction_id"`
	OrderID           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	TransactionStatus string `json:"transaction_status"`
	StatusCode        string `json:"status_code"`
	StatusMessage     string `json:"status_message"`
	ApprovalCode      string `json:"approval_code,omitempty"`
	Bank              string `json:"bank,omitempty"`
	MaskedCard        string `json:"masked_card,omitempty"`
	CardType          string `json:"card_type,omitempty"`
	PaymentCode       string `json:"payment_code,omitempty"`
	QrCodeURL         string `json:"qr_code_url,omitempty"`
	QrCodeImage       string `json:"qr_code_image,omitempty"`
	GrossAmount       string `json:"gross_amount"`
	SettlementTime    string `json:"settlement_time,omitempty"`
}

// IMidtransService defines the interface for Midtrans operations
type IMidtransService interface {
	// Charge operations
	ChargeQRIS(orderID string, amount float64, customerName, customerEmail string) (*MidtransChargeResponse, error)
	ChargeVA(orderID string, amount float64, customerName, customerEmail string, bank string) (*MidtransChargeResponse, error)
	ChargeEWallet(orderID string, amount float64, customerName, customerEmail string, walletType string) (*MidtransChargeResponse, error)

	// Status operations
	GetTransactionStatus(orderID string) (*MidtransStatusResponse, error)

	// Refund operations
	Refund(orderID string, amount float64, reason string) error

	// Cancel operations
	CancelTransaction(orderID string) error

	// Notification handling
	ParseNotification(payload map[string]interface{}) (*MidtransNotification, error)
	HandleNotification(notification *MidtransNotification) error

	// Utility
	IsEnabled() bool
	GetSnapURL() string
}

// MidtransNotification represents a parsed notification from Midtrans
type MidtransNotification struct {
	OrderID           string  `json:"order_id"`
	TransactionID     string  `json:"transaction_id"`
	TransactionStatus string  `json:"transaction_status"`
	TransactionType   string  `json:"transaction_type"`
	PaymentType       string  `json:"payment_type"`
	GrossAmount       float64 `json:"gross_amount"`
	ApprovalCode      string  `json:"approval_code"`
	Bank              string  `json:"bank"`
	FraudStatus       string  `json:"fraud_status"`
	StatusCode        string  `json:"status_code"`
}

// MidtransService handles all Midtrans payment operations
type MidtransService struct {
	snapClient    snap.Client
	coreAPIClient coreapi.Client
	cfg           *config.MidtransConfig
}

// NewMidtransService creates a new Midtrans service instance
func NewMidtransService(cfg *config.MidtransConfig) *MidtransService {
	env := midtrans.Sandbox
	if cfg.Environment == "production" {
		env = midtrans.Production
	}

	snapClient := snap.Client{}
	snapClient.New(cfg.ServerKey, env)

	coreAPIClient := coreapi.Client{}
	coreAPIClient.New(cfg.ServerKey, env)

	return &MidtransService{
		snapClient:    snapClient,
		coreAPIClient: coreAPIClient,
		cfg:           cfg,
	}
}

// IsEnabled checks if Midtrans is enabled
func (s *MidtransService) IsEnabled() bool {
	return s.cfg.Enabled
}

// GetSnapURL returns the Snap redirect URL
func (s *MidtransService) GetSnapURL() string {
	return s.cfg.SnapURL
}

// generateItemDetails creates item details for the charge request
func generateItemDetails(orderID string, amount float64) []snap.ItemDetail {
	return []snap.ItemDetail{
		{
			ID:    orderID,
			Name:  fmt.Sprintf("Payment for order %s", orderID),
			Price: int64(amount),
			Qty:   1,
		},
	}
}

// ChargeQRIS creates a QRIS charge request
func (s *MidtransService) ChargeQRIS(orderID string, amount float64, customerName, customerEmail string) (*MidtransChargeResponse, error) {
	req := &snap.CreateChargeReq{
		PaymentType: "qris",
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(amount),
		},
		CustomerDetails: &midtrans.CustomerDetails{
			FName: customerName,
			Email: customerEmail,
		},
	ItemDetails: generateItemDetails(orderID, amount),
	}

	resp, err := s.snapClient.ChargeTransaction(req)
	if err != nil {
		return nil, fmt.Errorf("failed to charge QRIS: %w", err)
	}

	return s.parseChargeResponse(resp)
}

// ChargeVA creates a Virtual Account charge request
func (s *MidtransService) ChargeVA(orderID string, amount float64, customerName, customerEmail string, bank string) (*MidtransChargeResponse, error) {
	paymentType := s.getVAPaymentType(bank)

	req := &snap.CreateChargeReq{
		PaymentType: paymentType,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(amount),
		},
		CustomerDetails: &midtrans.CustomerDetails{
			FName: customerName,
			Email: customerEmail,
		},
		PaymentType: paymentType,
	}

	// Set expiry time (default 24 hours)
	expiryTime := time.Now().Add(24 * time.Hour)
	req.CustomExpiry = &snap.CustomExpiry{
		ExpiryDuration: 24,
		Unit:           "hour",
	}

	resp, err := s.snapClient.ChargeTransaction(req)
	if err != nil {
		return nil, fmt.Errorf("failed to charge VA: %w", err)
	}

	return s.parseChargeResponse(resp)
}

// ChargeEWallet creates an e-wallet charge request
func (s *MidtransService) ChargeEWallet(orderID string, amount float64, customerName, customerEmail string, walletType string) (*MidtransChargeResponse, error) {
	paymentType := s.getEWalletPaymentType(walletType)

	req := &snap.CreateChargeReq{
		PaymentType: paymentType,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(amount),
		},
		CustomerDetails: &midtrans.CustomerDetails{
			FName: customerName,
			Email: customerEmail,
		},
		ItemDetails: generateItemDetails(orderID, amount),
	}

	// Set expiry for e-wallet (if applicable)
	expiryTime := time.Now().Add(24 * time.Hour)
	req.CustomExpiry = &snap.CustomExpiry{
		ExpiryDuration: 24,
		Unit:           "hour",
	}

	resp, err := s.snapClient.ChargeTransaction(req)
	if err != nil {
		return nil, fmt.Errorf("failed to charge e-wallet: %w", err)
	}

	return s.parseChargeResponse(resp)
}

// getVAPaymentType maps bank code to Midtrans VA payment type
func (s *MidtransService) getVAPaymentType(bank string) string {
	bankMap := map[string]string{
		"bca":     "bank_transfer",
		"mandiri": "echannel",
		"bni":     "bank_transfer",
		"permata": "bank_transfer",
		"bri":     "bank_transfer",
	}

	if pt, ok := bankMap[strings.ToLower(bank)]; ok {
		return pt
	}
	return "bank_transfer"
}

// getEWalletPaymentType maps wallet type to Midtrans payment type
func (s *MidtransService) getEWalletPaymentType(walletType string) string {
	walletMap := map[string]string{
		"gopay":    "gopay",
		"shopeepay": "shopeepay",
		"ovo":      "ovo",
		"dana":     "danazoo",
	}

	if pt, ok := walletMap[strings.ToLower(walletType)]; ok {
		return pt
	}
	return walletType
}

// GetTransactionStatus retrieves transaction status from Midtrans
func (s *MidtransService) GetTransactionStatus(orderID string) (*MidtransStatusResponse, error) {
	resp, err := s.coreAPIClient.GetTransactionStatus(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction status: %w", err)
	}

	return &MidtransStatusResponse{
		TransactionID:     resp.TransactionID,
		OrderID:           resp.OrderID,
		PaymentType:       resp.PaymentType,
		TransactionStatus: resp.TransactionStatus,
		StatusCode:        resp.StatusCode,
		StatusMessage:     resp.StatusMessage,
		ApprovalCode:      resp.ApprovalCode,
		Bank:              resp.Bank,
		MaskedCard:        resp.MaskedCard,
		CardType:          resp.CardType,
		PaymentCode:       resp.PaymentCode,
		QrCodeURL:         resp.QrCodeURL,
		QrCodeImage:       resp.QrCodeImage,
		GrossAmount:       resp.GrossAmount,
		SettlementTime:    resp.SettlementTime,
	}, nil
}

// Refund processes a refund request
func (s *MidtransService) Refund(orderID string, amount float64, reason string) error {
	req := &coreapi.RefundReq{
		RefundAmount: int64(amount),
		Reason:       reason,
	}

	_, err := s.coreAPIClient.RefundTransaction(orderID, req)
	if err != nil {
		return fmt.Errorf("failed to process refund: %w", err)
	}

	return nil
}

// CancelTransaction cancels a pending transaction
func (s *MidtransService) CancelTransaction(orderID string) error {
	err := s.coreAPIClient.CancelTransaction(orderID)
	if err != nil {
		return fmt.Errorf("failed to cancel transaction: %w", err)
	}
	return nil
}

// ParseNotification parses a Midtrans notification payload
func (s *MidtransService) ParseNotification(payload map[string]interface{}) (*MidtransNotification, error) {
	// Extract gross amount
	grossAmount := 0.0
	if ga, ok := payload["gross_amount"].(string); ok {
		fmt.Sscanf(ga, "%f", &grossAmount)
	}

	notification := &MidtransNotification{
		OrderID:           getStringValue(payload, "order_id"),
		TransactionID:     getStringValue(payload, "transaction_id"),
		TransactionStatus: getStringValue(payload, "transaction_status"),
		TransactionType:   getStringValue(payload, "transaction_type"),
		PaymentType:       getStringValue(payload, "payment_type"),
		GrossAmount:       grossAmount,
		ApprovalCode:      getStringValue(payload, "approval_code"),
		Bank:              getStringValue(payload, "bank"),
		FraudStatus:       getStringValue(payload, "fraud_status"),
		StatusCode:        getStringValue(payload, "status_code"),
	}

	return notification, nil
}

// HandleNotification processes a Midtrans notification
// This should be connected to the Payment service to update payment status
func (s *MidtransService) HandleNotification(notification *MidtransNotification) error {
	log.Printf("[Midtrans] Processing notification for OrderID: %s, Status: %s",
		notification.OrderID, notification.TransactionStatus)

	// The actual payment status update will be handled by the payment service
	// which will use this notification to update the payment and transaction records
	// This is typically done by publishing an event or directly updating through repository

	return nil
}

// parseChargeResponse converts Midtrans response to our response struct
func (s *MidtransService) parseChargeResponse(resp *snap.Response) (*MidtransChargeResponse, error) {
	if resp == nil {
		return nil, fmt.Errorf("nil response from Midtrans")
	}

	chargeResp := &MidtransChargeResponse{
		TransactionID:   resp.TransactionID,
		OrderID:         resp.OrderID,
		PaymentType:     resp.PaymentType,
		StatusCode:      resp.StatusCode,
		StatusMessage:   resp.StatusMessage,
		TransactionTime: resp.TransactionTime,
		FraudStatus:     resp.FraudStatus,
		QrCodeURL:       resp.QrCodeURL,
		QrCodeImage:     resp.QrCodeImage,
		GrossAmount:     resp.GrossAmount,
	}

	// Handle VA numbers
	if len(resp.VaNumbers) > 0 {
		chargeResp.VANumbers = make([]VANumber, len(resp.VaNumbers))
		for i, va := range resp.VaNumbers {
			chargeResp.VANumbers[i] = VANumber{
				Bank:     va.Bank,
				VANumber: va.VANumber,
			}
		}
	}

	// Handle Permata VA
	if resp.PermataVANumber != "" {
		chargeResp.PermitaVANumber = resp.PermataVANumber
	}

	return chargeResp, nil
}

// getStringValue safely extracts a string from a map
func getStringValue(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

// ToJSON converts charge response to JSON string
func (r *MidtransChargeResponse) ToJSON() string {
	b, err := json.Marshal(r)
	if err != nil {
		return "{}"
	}
	return string(b)
}

// MapMidtransStatusToPaymentStatus maps Midtrans transaction status to our payment status
func MapMidtransStatusToPaymentStatus(midtransStatus string) models.PaymentStatus {
	statusMap := map[string]models.PaymentStatus{
		"capture":   models.PaymentStatusSuccess,
		"settlement": models.PaymentStatusSuccess,
		"pending":   models.PaymentStatusPending,
		"deny":      models.PaymentStatusFailed,
		"cancel":    models.PaymentStatusCancelled,
		"expire":    models.PaymentStatusExpired,
		"refund":    models.PaymentStatusRefunded,
		"partial_refund": models.PaymentStatusRefunded,
	}

	if status, ok := statusMap[midtransStatus]; ok {
		return status
	}
	return models.PaymentStatusPending
}