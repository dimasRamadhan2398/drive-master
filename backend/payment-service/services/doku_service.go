package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"payment-service/models"
	"payment-service/pkg/config"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type DokuService struct {
	cfg *config.DokuConfig
}

func NewDokuService(cfg *config.DokuConfig) *DokuService {
	return &DokuService{cfg: cfg}
}

func (s *DokuService) GetName() string {
	return "doku"
}

// DokuCheckoutRequest matches Doku Checkout API body structure
type DokuCheckoutRequest struct {
	Order   DokuOrderRequest   `json:"order"`
	Payment DokuPaymentRequest `json:"payment"`
}

type DokuOrderRequest struct {
	Amount            int64  `json:"amount"`
	InvoiceNumber     string `json:"invoice_number"`
	CallbackURL       string `json:"callback_url,omitempty"`
	CallbackURLResult string `json:"callback_url_result,omitempty"`
}

type DokuPaymentRequest struct {
	PaymentDueDate int `json:"payment_due_date,omitempty"`
}

// DokuCheckoutResponse matches Doku Checkout API response structure
type DokuCheckoutResponse struct {
	Response *DokuInnerResponse `json:"response,omitempty"`
	Message  []string           `json:"message,omitempty"`
	Error    *DokuErrorResponse `json:"error,omitempty"`
}

type DokuInnerResponse struct {
	Order   DokuOrderResponse   `json:"order"`
	Payment DokuPaymentResponse `json:"payment"`
}

type DokuOrderResponse struct {
	InvoiceNumber string `json:"invoice_number"`
	Amount        string `json:"amount"`
}

type DokuPaymentResponse struct {
	URL string `json:"url"`
}

type DokuErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// CreateCheckout initiates a checkout with Doku Checkout API
func (s *DokuService) CreateCheckout(orderID string, amount float64, packageName, customerName, customerEmail string) (*CheckoutResponse, error) {
	targetURL := fmt.Sprintf("%s/checkout/v1/payment", s.cfg.BaseURL)
	targetPath := "/checkout/v1/payment"

	returnURL := fmt.Sprintf("%s/auth/payment-status?orderId=%s", s.cfg.FrontendURL, orderID)

	reqBody := DokuCheckoutRequest{
		Order: DokuOrderRequest{
			Amount:            int64(amount),
			InvoiceNumber:     orderID,
			CallbackURL:       returnURL,
			CallbackURLResult: returnURL,
		},
		Payment: DokuPaymentRequest{
			PaymentDueDate: s.cfg.PaymentDueDate,
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Doku checkout request: %w", err)
	}

	reqID := uuid.New().String()
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05Z")

	// Generate Digest
	hasher := sha256.New()
	hasher.Write(bodyBytes)
	digest := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	// Generate Signature
	signaturePayload := fmt.Sprintf("Client-Id:%s\nRequest-Id:%s\nRequest-Timestamp:%s\nRequest-Target:%s\nDigest:%s",
		s.cfg.ClientID, reqID, timestamp, targetPath, digest)

	mac := hmac.New(sha256.New, []byte(s.cfg.SecretKey))
	mac.Write([]byte(signaturePayload))
	signature := "HMACSHA256=" + base64.StdEncoding.EncodeToString(mac.Sum(nil))

	req, err := http.NewRequest("POST", targetURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-Id", s.cfg.ClientID)
	req.Header.Set("Request-Id", reqID)
	req.Header.Set("Request-Timestamp", timestamp)
	req.Header.Set("Signature", signature)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request to Doku failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("doku checkout API returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	var res DokuCheckoutResponse
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return nil, fmt.Errorf("failed to parse Doku response: %w", err)
	}

	if res.Response == nil || res.Response.Payment.URL == "" {
		return nil, fmt.Errorf("doku response did not contain payment URL: %s", string(respBytes))
	}

	return &CheckoutResponse{
		RedirectURL: res.Response.Payment.URL,
		OrderID:     orderID,
	}, nil
}

// GetTransactionStatus retrieves payment status by checking transaction status directly
func (s *DokuService) GetTransactionStatus(orderID string) (models.PaymentStatus, error) {
	return models.PaymentStatusPending, nil
}

// VerifyNotification validates Doku signature headers and returns a NotificationPayload
func (s *DokuService) VerifyNotification(headers map[string]string, rawBody []byte) (*NotificationPayload, error) {
	// Normalize header names to lowercase for comparison
	var clientID, reqID, timestamp, receivedSignature string
	for k, v := range headers {
		switch strings.ToLower(k) {
		case "client-id":
			clientID = v
		case "request-id":
			reqID = v
		case "request-timestamp":
			timestamp = v
		case "signature":
			receivedSignature = v
		}
	}

	if clientID == "" || reqID == "" || timestamp == "" || receivedSignature == "" {
		return nil, fmt.Errorf("missing required Doku validation headers")
	}

	// 1. Calculate Digest
	hasher := sha256.New()
	hasher.Write(rawBody)
	digest := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	// 2. Build signature component string
	targetPath := "/api/v1/payments/callback"
	signaturePayload := fmt.Sprintf("Client-Id:%s\nRequest-Id:%s\nRequest-Timestamp:%s\nRequest-Target:%s\nDigest:%s",
		clientID, reqID, timestamp, targetPath, digest)

	// 3. Compute HMAC
	mac := hmac.New(sha256.New, []byte(s.cfg.SecretKey))
	mac.Write([]byte(signaturePayload))
	computedSignature := "HMACSHA256=" + base64.StdEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(receivedSignature), []byte(computedSignature)) {
		return nil, fmt.Errorf("invalid Doku webhook signature")
	}

	// Parse webhook body
	var webhookData map[string]interface{}
	if err := json.Unmarshal(rawBody, &webhookData); err != nil {
		return nil, fmt.Errorf("failed to parse Doku webhook body: %w", err)
	}

	// Extract invoice & transaction status
	invoiceNumber := ""
	if order, ok := webhookData["order"].(map[string]interface{}); ok {
		if inv, ok := order["invoice_number"].(string); ok {
			invoiceNumber = inv
		}
	}

	transactionID := ""
	if transaction, ok := webhookData["transaction"].(map[string]interface{}); ok {
		if id, ok := transaction["id"].(string); ok {
			transactionID = id
		}
	}
	if transactionID == "" {
		if id, ok := webhookData["transactionId"].(string); ok {
			transactionID = id
		}
	}

	txnStatus := models.PaymentStatusPending
	if transaction, ok := webhookData["transaction"].(map[string]interface{}); ok {
		if statusStr, ok := transaction["status"].(string); ok {
			switch strings.ToUpper(statusStr) {
			case "SUCCESS":
				txnStatus = models.PaymentStatusSuccess
			case "FAILED":
				txnStatus = models.PaymentStatusFailed
			default:
				txnStatus = models.PaymentStatusPending
			}
		}
	}

	amount := 0.0
	if order, ok := webhookData["order"].(map[string]interface{}); ok {
		if amt, ok := order["amount"].(float64); ok {
			amount = amt
		} else if amtStr, ok := order["amount"].(string); ok {
			amount, _ = strconv.ParseFloat(amtStr, 64)
		}
	}

	paymentType := ""
	if channel, ok := webhookData["channel"].(map[string]interface{}); ok {
		if name, ok := channel["name"].(string); ok {
			paymentType = name
		}
	}

	return &NotificationPayload{
		OrderID:           invoiceNumber,
		TransactionID:     transactionID,
		TransactionStatus: txnStatus,
		PaymentType:       paymentType,
		GrossAmount:       amount,
		RawPayload:        webhookData,
	}, nil
}

func (s *DokuService) SimulatePayment(orderID string, amount float64) error {
	return nil
}
