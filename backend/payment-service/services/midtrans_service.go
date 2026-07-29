package services

import (
	"encoding/json"
	"fmt"
	"log"
	"payment-service/models"
	"payment-service/pkg/config"
	"strings"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

// PaymentType represents the type of payment method
type PaymentType string

const (
	PaymentTypeQRIS       PaymentType = "qris"
	PaymentTypeVA         PaymentType = "bank_transfer"
	PaymentTypeEWallet    PaymentType = "ewallet"
	PaymentTypeCreditCard PaymentType = "credit_card"
)

type SnapResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirectUrl"`
}

// MidtransChargeResponse represents the response from Midtrans charge API
type MidtransChargeResponse struct {
	TransactionID   string `json:"transaction_id"`
	OrderID         string `json:"order_id"`
	PaymentType     string `json:"payment_type"`
	StatusCode      string `json:"status_code"`
	StatusMessage   string `json:"status_message"`
	TransactionTime string `json:"transaction_time"`
	TransactionUUID string `json:"transaction_uuid"`
	MerchantID      string `json:"merchant_id"`
	FraudStatus     string `json:"fraud_status"`
	// QRIS specific
	QrCodeURL   string `json:"qr_code_url,omitempty"`
	QrCodeImage string `json:"qr_code_image,omitempty"`
	// VA specific
	VANumbers       []VANumber `json:"va_numbers,omitempty"`
	PermataVANumber string     `json:"permata_va_number,omitempty"`
	BcaVANumber     string     `json:"bca_va_number,omitempty"`
	MandiriVANumber string     `json:"mandiri_va_number,omitempty"`
	BniVANumber     string     `json:"bni_va_number,omitempty"`
	BillKey         string     `json:"bill_key,omitempty"`    // Mandiri
	BillerCode      string     `json:"biller_code,omitempty"` // Mandiri
	// E-wallet specific
	Actions []map[string]string `json:"actions,omitempty"` // GoPay, ShopeePay deeplink/QR
	// Common
	GrossAmount string `json:"gross_amount"`
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
	GrossAmount       string `json:"gross_amount"`
	SettlementTime    string `json:"settlement_time,omitempty"`
}

// IMidtransService defines the interface for Midtrans operations
type IMidtransService interface {
	CreateSnapTransaction(orderID string, amount float64, packageName, customerName, customerEmail string) (*SnapResponse, error)
	ChargeQRIS(orderID string, amount float64, packageName, customerName, customerEmail string) (*MidtransChargeResponse, error)
	ChargeVA(orderID string, amount float64, packageName, customerName, customerEmail string, bank string) (*MidtransChargeResponse, error)
	ChargeEWallet(orderID string, amount float64, packageName, customerName, customerEmail string, walletType string) (*MidtransChargeResponse, error)

	// Status operations
	GetMidtransTransactionStatus(orderID string) (*MidtransStatusResponse, error)

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

// CreateSnapTransaction uses Snap API — returns a token and redirect URL
// CreateSnapTransaction implements [IMidtransService].
func (s *MidtransService) CreateSnapTransaction(orderID string, amount float64, packageName, customerName string, customerEmail string) (*SnapResponse, error) {
	items := generateItemDetails(orderID, amount, packageName)

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(amount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: customerName,
			Email: customerEmail,
		},
		Items: &items,
		Callbacks: &snap.Callbacks{
			Finish: fmt.Sprintf("%s/auth/payment-status?status=success&orderId=%s", s.cfg.FrontendURL, orderID),
		},
	}

	resp, midErr := s.snapClient.CreateTransaction(req)
	if midErr != nil {
		return nil, fmt.Errorf("failed to create SNAP transaction: %w", midErr)
	}

	return &SnapResponse{
		Token:       resp.Token,
		RedirectURL: resp.RedirectURL,
	}, nil
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

// generateItemDetails builds item details for a charge request.
// FIX: return type is []midtrans.ItemDetails
func generateItemDetails(orderID string, amount float64, packageName string) []midtrans.ItemDetails {
	name := fmt.Sprintf("Paket %s", packageName)
	if len(name) > 50 {
		name = name[:50]
	}
	return []midtrans.ItemDetails{
		{
			ID:           orderID,
			Name:         name,
			Price:        int64(amount),
			Qty:          1,
			Brand:        "Drive Master",
			Category:     "Pembelian Paket Mengemudi",
			MerchantName: "PT. Drive Master Indonesia",
		},
	}
}

// ChargeQRIS creates a QRIS charge request
func (s *MidtransService) ChargeQRIS(orderID string, amount float64, packageName string, customerName, customerEmail string) (*MidtransChargeResponse, error) {
	items := generateItemDetails(orderID, amount, packageName)

	req := &coreapi.ChargeReq{
		PaymentType: coreapi.PaymentTypeQris,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(amount),
		},
		CustomerDetails: &midtrans.CustomerDetails{
			FName: customerName,
			Email: customerEmail,
		},
		Items: &items,
	}

	resp, midErr := s.coreAPIClient.ChargeTransaction(req)
	if midErr != nil {
		return nil, fmt.Errorf("failed to charge QRIS: %s", midErr.GetMessage())
	}

	return parseCoreAPIResponse(resp), nil
}

// req := &snap.Request{
// 	TransactionDetails: midtrans.TransactionDetails{
// 		OrderID:  orderID,
// 		GrossAmt: int64(amount),
// 	},
// 	CustomerDetail: &midtrans.CustomerDetails{
// 		FName: customerName,
// 		Email: customerEmail,
// 	},
// 	Items: generateItemDetails(orderID, amount),
// }

// ChargeVA charges via Core API — returns a Virtual Account number directly.
func (s *MidtransService) ChargeVA(orderID string, amount float64, packageName, customerName, customerEmail, bank string) (*MidtransChargeResponse, error) {
	items := generateItemDetails(orderID, amount, packageName)

	// FIX: use coreapi.ChargeReq with BankTransfer details
	req := &coreapi.ChargeReq{
		PaymentType: coreapi.PaymentTypeBankTransfer,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(amount),
		},
		CustomerDetails: &midtrans.CustomerDetails{
			FName: customerName,
			Email: customerEmail,
		},
		Items: &items,
		BankTransfer: &coreapi.BankTransferDetails{
			// FIX: use midtrans.Bank constants, not plain strings
			Bank: midtrans.Bank(strings.ToLower(bank)),
		},
	}

	resp, midErr := s.coreAPIClient.ChargeTransaction(req)
	if midErr != nil {
		return nil, fmt.Errorf("failed to charge VA: %s", midErr.GetMessage())
	}

	return parseCoreAPIResponse(resp), nil
}

// ChargeEWallet charges via Core API — returns deeplink/QR actions for GoPay, ShopeePay, etc.
func (s *MidtransService) ChargeEWallet(orderID string, amount float64, packageName, customerName, customerEmail, walletType string) (*MidtransChargeResponse, error) {
	items := generateItemDetails(orderID, amount, packageName)

	paymentType := getEWalletPaymentType(walletType)

	req := &coreapi.ChargeReq{
		PaymentType: paymentType,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(amount),
		},
		CustomerDetails: &midtrans.CustomerDetails{
			FName: customerName,
			Email: customerEmail,
		},
		Items: &items,
	}

	// GoPay needs a callback URL
	if paymentType == coreapi.PaymentTypeGopay {
		req.Gopay = &coreapi.GopayDetails{
			EnableCallback: true,
			CallbackUrl:    s.cfg.GopayCallbackURL, // add this to your config
		}
	}

	resp, midErr := s.coreAPIClient.ChargeTransaction(req)
	if midErr != nil {
		return nil, fmt.Errorf("failed to charge e-wallet: %s", midErr.GetMessage())
	}

	return parseCoreAPIResponse(resp), nil
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
		"gopay":     "gopay",
		"shopeepay": "shopeepay",
		"ovo":       "ovo",
		"dana":      "danazoo",
	}

	if pt, ok := walletMap[strings.ToLower(walletType)]; ok {
		return pt
	}
	return walletType
}

// GetMidtransTransactionStatus retrieves transaction status from Midtrans
func (s *MidtransService) GetMidtransTransactionStatus(orderID string) (*MidtransStatusResponse, error) {
	resp, err := s.coreAPIClient.CheckTransaction(orderID)
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
		GrossAmount:       resp.GrossAmount,
		SettlementTime:    resp.SettlementTime,
	}, nil
}

// Refund processes a refund request
func (s *MidtransService) Refund(orderID string, amount float64, reason string) error {
	req := &coreapi.RefundReq{
		RefundKey: "refund_" + orderID,
		Amount:    int64(amount),
		Reason:    reason,
	}

	_, err := s.coreAPIClient.RefundTransaction(orderID, req)
	if err != nil {
		return fmt.Errorf("failed to process refund: %w", err)
	}

	return nil
}

// CancelTransaction cancels a pending transaction
func (s *MidtransService) CancelTransaction(orderID string) error {
	_, err := s.coreAPIClient.CancelTransaction(orderID)
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

// // parseChargeResponse converts Midtrans response to our response struct
// func (s *MidtransService) parseChargeResponse(resp *snap.Response) (*MidtransChargeResponse, error) {
// 	if resp == nil {
// 		return nil, fmt.Errorf("nil response from Midtrans")
// 	}

// 	chargeResp := &MidtransChargeResponse{
// 		TransactionID:   resp.TransactionID,
// 		OrderID:         resp.OrderID,
// 		PaymentType:     resp.PaymentType,
// 		StatusCode:      resp.StatusCode,
// 		StatusMessage:   resp.StatusMessage,
// 		TransactionTime: resp.TransactionTime,
// 		FraudStatus:     resp.FraudStatus,
// 		QrCodeURL:       resp.QrCodeURL,
// 		QrCodeImage:     resp.QrCodeImage,
// 		GrossAmount:     resp.GrossAmount,
// 	}

// 	// Handle VA numbers
// 	if len(resp.VaNumbers) > 0 {
// 		chargeResp.VANumbers = make([]VANumber, len(resp.VaNumbers))
// 		for i, va := range resp.VaNumbers {
// 			chargeResp.VANumbers[i] = VANumber{
// 				Bank:     va.Bank,
// 				VANumber: va.VANumber,
// 			}
// 		}
// 	}

// 	// Handle Permata VA
// 	if resp.PermataVANumber != "" {
// 		chargeResp.PermitaVANumber = resp.PermataVANumber
// 	}

// 	return chargeResp, nil
// }

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
		"capture":        models.PaymentStatusSuccess,
		"settlement":     models.PaymentStatusSuccess,
		"pending":        models.PaymentStatusPending,
		"deny":           models.PaymentStatusFailed,
		"cancel":         models.PaymentStatusCancelled,
		"expire":         models.PaymentStatusExpired,
		"refund":         models.PaymentStatusRefunded,
		"partial_refund": models.PaymentStatusRefunded,
	}

	if status, ok := statusMap[midtransStatus]; ok {
		return status
	}
	return models.PaymentStatusPending
}

func getEWalletPaymentType(walletType string) coreapi.CoreapiPaymentType {
	walletMap := map[string]coreapi.CoreapiPaymentType{
		"gopay":     coreapi.PaymentTypeGopay,
		"shopeepay": coreapi.PaymentTypeShopeepay,
	}
	if pt, ok := walletMap[strings.ToLower(walletType)]; ok {
		return pt
	}
	return coreapi.PaymentTypeGopay
}

// parseCoreAPIResponse maps coreapi.Response to our internal struct.
func parseCoreAPIResponse(resp *coreapi.ChargeResponse) *MidtransChargeResponse {
	result := &MidtransChargeResponse{
		TransactionID:   resp.TransactionID,
		OrderID:         resp.OrderID,
		PaymentType:     resp.PaymentType,
		StatusCode:      resp.StatusCode,
		StatusMessage:   resp.StatusMessage,
		TransactionTime: resp.TransactionTime,
		FraudStatus:     resp.FraudStatus,
		QrCodeURL:       resp.QRString,
		GrossAmount:     resp.GrossAmount,
		// FIX: PermataVANumber, not PermitaVANumber
		PermataVANumber: resp.PermataVaNumber,
		BillKey:         resp.BillKey,
		BillerCode:      resp.BillerCode,
	}

	if len(resp.VaNumbers) > 0 {
		result.VANumbers = make([]VANumber, len(resp.VaNumbers))
		for i, va := range resp.VaNumbers {
			result.VANumbers[i] = VANumber{
				Bank:     va.Bank,
				VANumber: va.VANumber,
			}
		}
	}

	return result
}

// GetName returns the gateway name
func (s *MidtransService) GetName() string {
	return "midtrans"
}

// CreateCheckout wraps CreateSnapTransaction to conform to IPaymentGatewayService
func (s *MidtransService) CreateCheckout(orderID string, amount float64, packageName, customerName, customerEmail string) (*CheckoutResponse, error) {
	snapResp, err := s.CreateSnapTransaction(orderID, amount, packageName, customerName, customerEmail)
	if err != nil {
		return nil, err
	}
	return &CheckoutResponse{
		Token:       snapResp.Token,
		RedirectURL: snapResp.RedirectURL,
		OrderID:     orderID,
	}, nil
}

// GetTransactionStatus conforms to IPaymentGatewayService and returns internal models.PaymentStatus
func (s *MidtransService) GetTransactionStatus(orderID string) (models.PaymentStatus, error) {
	resp, err := s.GetMidtransTransactionStatus(orderID)
	if err != nil {
		return models.PaymentStatusPending, err
	}
	return MapMidtransStatusToPaymentStatus(resp.TransactionStatus), nil
}

// VerifyNotification parses notification payload and maps to generic format
func (s *MidtransService) VerifyNotification(headers map[string]string, rawBody []byte) (*NotificationPayload, error) {
	var payload map[string]interface{}
	if err := json.Unmarshal(rawBody, &payload); err != nil {
		return nil, err
	}

	notif, err := s.ParseNotification(payload)
	if err != nil {
		return nil, err
	}

	return &NotificationPayload{
		OrderID:           notif.OrderID,
		TransactionID:     notif.TransactionID,
		TransactionStatus: MapMidtransStatusToPaymentStatus(notif.TransactionStatus),
		PaymentType:       notif.PaymentType,
		GrossAmount:       notif.GrossAmount,
		RawPayload:        payload,
	}, nil
}

func (s *MidtransService) SimulatePayment(orderID string, amount float64) error {
	return nil
}
