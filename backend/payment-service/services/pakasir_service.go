package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"payment-service/models"
	"payment-service/pkg/config"
	"payment-service/repositories"
)

type PakasirService struct {
	cfg         *config.PakasirConfig
	paymentRepo repositories.IPaymentRepository
}

func NewPakasirService(cfg *config.PakasirConfig, paymentRepo repositories.IPaymentRepository) *PakasirService {
	return &PakasirService{
		cfg:         cfg,
		paymentRepo: paymentRepo,
	}
}

func (s *PakasirService) GetName() string {
	return "pakasir"
}

// IsSandbox returns true when the environment is set to sandbox
func (s *PakasirService) IsSandbox() bool {
	return strings.ToLower(s.cfg.Environment) == "sandbox"
}

// ─────────────────────────────────────────────────────────────────────────────
// B. URL-based checkout (Section B in docs)
// ─────────────────────────────────────────────────────────────────────────────

// CreateCheckout generates a Pakasir hosted payment URL
func (s *PakasirService) CreateCheckout(orderID string, amount float64, packageName, customerName, customerEmail string) (*CheckoutResponse, error) {
	// Use configured frontend URL, fall back to production if not set
	frontendURL := s.cfg.FrontendURL
	if frontendURL == "" {
		frontendURL = "https://drivemaster.id"
	}
	returnURL := fmt.Sprintf("%s/auth/payment-status?orderId=%s", frontendURL, orderID)

	// URL-based integration: https://app.pakasir.com/pay/{slug}/{amount}?order_id={order_id}&redirect={return_url}
	redirectURL := fmt.Sprintf("%s/pay/%s/%d?order_id=%s&redirect=%s",
		s.cfg.BaseURL,
		s.cfg.Slug,
		int(amount),
		orderID,
		url.QueryEscape(returnURL),
	)

	return &CheckoutResponse{
		Token:       orderID,
		RedirectURL: redirectURL,
		OrderID:     orderID,
	}, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// C. API-based integration (Sections C.2 – C.5 in docs)
// ─────────────────────────────────────────────────────────────────────────────

// PakasirPaymentMethod lists the available payment methods supported by Pakasir API.
// Docs Section C.3.
type PakasirPaymentMethod string

const (
	PakasirMethodQRIS         PakasirPaymentMethod = "qris"
	PakasirMethodBNIVA        PakasirPaymentMethod = "bni_va"
	PakasirMethodBRIVA        PakasirPaymentMethod = "bri_va"
	PakasirMethodCIMBNiagaVA  PakasirPaymentMethod = "cimb_niaga_va"
	PakasirMethodSampoernaVA  PakasirPaymentMethod = "sampoerna_va"
	PakasirMethodBNCVA        PakasirPaymentMethod = "bnc_va"
	PakasirMethodMaybankVA    PakasirPaymentMethod = "maybank_va"
	PakasirMethodPermataVA    PakasirPaymentMethod = "permata_va"
	PakasirMethodATMBersamaVA PakasirPaymentMethod = "atm_bersama_va"
	PakasirMethodArthaGrahaVA PakasirPaymentMethod = "artha_graha_va"
)

// PakasirTransactionResponse is the response body from /api/transactioncreate/{method}
// Docs Section C.2.
type PakasirTransactionResponse struct {
	Payment struct {
		Project       string  `json:"project"`
		OrderID       string  `json:"order_id"`
		Amount        float64 `json:"amount"`
		Fee           float64 `json:"fee"`
		TotalPayment  float64 `json:"total_payment"`
		PaymentMethod string  `json:"payment_method"`
		// PaymentNumber holds the QRIS string or VA number
		PaymentNumber string `json:"payment_number"`
		ExpiredAt     string `json:"expired_at"`
	} `json:"payment"`
}

// CreateTransaction calls the Pakasir API to create a transaction and returns
// the payment number (QRIS string / VA number), total amount, and expiry.
// Docs Section C.2: POST /api/transactioncreate/{method}
func (s *PakasirService) CreateTransaction(orderID string, amount float64, method PakasirPaymentMethod) (*PakasirTransactionResponse, error) {
	targetURL := fmt.Sprintf("%s/api/transactioncreate/%s", s.cfg.BaseURL, method)

	body := map[string]interface{}{
		"project":  s.cfg.Slug,
		"order_id": orderID,
		"amount":   int(amount),
		"api_key":  s.cfg.APIKey,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Pakasir create-transaction body: %w", err)
	}

	req, err := http.NewRequest("POST", targetURL, strings.NewReader(string(bodyBytes)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request to Pakasir create-transaction API failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pakasir create-transaction API returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	var res PakasirTransactionResponse
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return nil, fmt.Errorf("failed to parse Pakasir create-transaction response: %w", err)
	}

	return &res, nil
}

// SimulatePayment triggers a simulated payment on a sandbox transaction so the webhook is fired.
// Only works when the project is in Sandbox mode.
// Docs Section C.4: POST /api/paymentsimulation
func (s *PakasirService) SimulatePayment(orderID string, amount float64) error {
	if !s.IsSandbox() {
		return fmt.Errorf("SimulatePayment is only available in sandbox mode")
	}

	targetURL := fmt.Sprintf("%s/api/paymentsimulation", s.cfg.BaseURL)

	body := map[string]interface{}{
		"project":  s.cfg.Slug,
		"order_id": orderID,
		"amount":   int(amount),
		"api_key":  s.cfg.APIKey,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal Pakasir simulate-payment body: %w", err)
	}

	req, err := http.NewRequest("POST", targetURL, strings.NewReader(string(bodyBytes)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http request to Pakasir payment-simulation API failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("pakasir payment-simulation API returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	return nil
}

// CancelTransaction cancels an existing Pakasir transaction.
// Docs Section C.5: POST /api/transactioncancel
func (s *PakasirService) CancelTransaction(orderID string, amount float64) error {
	targetURL := fmt.Sprintf("%s/api/transactioncancel", s.cfg.BaseURL)

	body := map[string]interface{}{
		"project":  s.cfg.Slug,
		"order_id": orderID,
		"amount":   int(amount),
		"api_key":  s.cfg.APIKey,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal Pakasir cancel-transaction body: %w", err)
	}

	req, err := http.NewRequest("POST", targetURL, strings.NewReader(string(bodyBytes)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http request to Pakasir cancel-transaction API failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("pakasir cancel-transaction API returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// E. Transaction Detail API (Section E in docs)
// ─────────────────────────────────────────────────────────────────────────────

// GetTransactionStatus queries the Pakasir Transaction Detail API
func (s *PakasirService) GetTransactionStatus(orderID string) (models.PaymentStatus, error) {
	// 1. Fetch payment from database to retrieve amount
	payment, err := s.paymentRepo.GetByOrderID(orderID)
	if err != nil {
		return models.PaymentStatusPending, fmt.Errorf("failed to fetch payment for order %s: %w", orderID, err)
	}

	amountInt := int(payment.Amount)

	// API: GET https://app.pakasir.com/api/transactiondetail?project={slug}&amount={amount}&order_id={order_id}&api_key={api_key}
	targetURL := fmt.Sprintf("%s/api/transactiondetail?project=%s&amount=%d&order_id=%s&api_key=%s",
		s.cfg.BaseURL,
		s.cfg.Slug,
		amountInt,
		orderID,
		s.cfg.APIKey,
	)

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return models.PaymentStatusPending, err
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return models.PaymentStatusPending, fmt.Errorf("http request to Pakasir status API failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.PaymentStatusPending, err
	}

	if resp.StatusCode != http.StatusOK {
		return models.PaymentStatusPending, fmt.Errorf("pakasir status API returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	type PakasirDetailResponse struct {
		Transaction struct {
			Amount        float64 `json:"amount"`
			OrderID       string  `json:"order_id"`
			Project       string  `json:"project"`
			Status        string  `json:"status"`
			PaymentStatus string  `json:"payment_status"`
			State         string  `json:"state"`
			PaymentMethod string  `json:"payment_method"`
			CompletedAt   string  `json:"completed_at"`
		} `json:"transaction"`
		Status        string  `json:"status"`
		PaymentStatus string  `json:"payment_status"`
		State         string  `json:"state"`
		Amount        float64 `json:"amount"`
		OrderID       string  `json:"order_id"`
	}

	var res PakasirDetailResponse
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return models.PaymentStatusPending, fmt.Errorf("failed to parse Pakasir status response: %w", err)
	}

	statusStr := res.Transaction.Status
	if statusStr == "" {
		statusStr = res.Transaction.PaymentStatus
	}
	if statusStr == "" {
		statusStr = res.Transaction.State
	}
	if statusStr == "" {
		statusStr = res.Status
	}
	if statusStr == "" {
		statusStr = res.PaymentStatus
	}
	if statusStr == "" {
		statusStr = res.State
	}

	statusLower := strings.ToLower(strings.TrimSpace(statusStr))
	switch statusLower {
	case "completed", "success", "paid", "done", "settlement", "1":
		return models.PaymentStatusSuccess, nil
	case "pending", "process", "processing", "waiting":
		return models.PaymentStatusPending, nil
	case "failed", "expire", "expired", "cancel", "cancelled", "deny", "denied":
		return models.PaymentStatusFailed, nil
	default:
		if statusLower == "" {
			return models.PaymentStatusPending, nil
		}
		return models.PaymentStatusFailed, nil
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// D. Webhook handler (Section D in docs)
// ─────────────────────────────────────────────────────────────────────────────

// VerifyNotification validates incoming Pakasir webhook and verifies it with the Transaction Detail API
func (s *PakasirService) VerifyNotification(headers map[string]string, rawBody []byte) (*NotificationPayload, error) {
	type PakasirWebhookPayload struct {
		Amount        float64 `json:"amount"`
		OrderID       string  `json:"order_id"`
		Project       string  `json:"project"`
		Status        string  `json:"status"`
		PaymentStatus string  `json:"payment_status"`
		PaymentMethod string  `json:"payment_method"`
		CompletedAt   string  `json:"completed_at"`
	}

	var payload PakasirWebhookPayload
	if err := json.Unmarshal(rawBody, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse Pakasir webhook body: %w", err)
	}

	if payload.Project != "" && s.cfg.Slug != "" && payload.Project != s.cfg.Slug {
		return nil, fmt.Errorf("webhook project '%s' does not match configured project '%s'", payload.Project, s.cfg.Slug)
	}

	statusStr := payload.Status
	if statusStr == "" {
		statusStr = payload.PaymentStatus
	}

	statusLower := strings.ToLower(strings.TrimSpace(statusStr))
	var status models.PaymentStatus
	switch statusLower {
	case "completed", "success", "paid", "done", "settlement", "1":
		status = models.PaymentStatusSuccess
	case "pending", "process", "processing", "waiting":
		status = models.PaymentStatusPending
	default:
		var err error
		status, err = s.GetTransactionStatus(payload.OrderID)
		if err != nil && statusLower == "" {
			return nil, fmt.Errorf("failed to verify transaction status via Pakasir API: %w", err)
		}
		if statusLower != "" && status == models.PaymentStatusPending {
			status = models.PaymentStatusFailed
		}
	}

	return &NotificationPayload{
		OrderID:           payload.OrderID,
		TransactionID:     payload.OrderID,
		TransactionStatus: status,
		PaymentType:       payload.PaymentMethod,
		GrossAmount:       payload.Amount,
		RawPayload: map[string]interface{}{
			"amount":         payload.Amount,
			"order_id":       payload.OrderID,
			"project":        payload.Project,
			"status":         payload.Status,
			"payment_method": payload.PaymentMethod,
			"completed_at":   payload.CompletedAt,
		},
	}, nil
}
