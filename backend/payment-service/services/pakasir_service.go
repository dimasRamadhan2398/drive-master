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

// CreateCheckout generates a Pakasir hosted payment URL
func (s *PakasirService) CreateCheckout(orderID string, amount float64, packageName, customerName, customerEmail string) (*CheckoutResponse, error) {
	// Format return redirect URL
	returnURL := fmt.Sprintf("https://drivemaster.id/auth/payment-status?orderId=%s", orderID)

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
			PaymentMethod string  `json:"payment_method"`
			CompletedAt   string  `json:"completed_at"`
		} `json:"transaction"`
	}

	var res PakasirDetailResponse
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return models.PaymentStatusPending, fmt.Errorf("failed to parse Pakasir status response: %w", err)
	}

	switch strings.ToLower(res.Transaction.Status) {
	case "completed":
		return models.PaymentStatusSuccess, nil
	case "pending":
		return models.PaymentStatusPending, nil
	default:
		return models.PaymentStatusFailed, nil
	}
}

// VerifyNotification validates incoming Pakasir webhook and verifies it with the Transaction Detail API
func (s *PakasirService) VerifyNotification(headers map[string]string, rawBody []byte) (*NotificationPayload, error) {
	type PakasirWebhookPayload struct {
		Amount        float64 `json:"amount"`
		OrderID       string  `json:"order_id"`
		Project       string  `json:"project"`
		Status        string  `json:"status"`
		PaymentMethod string  `json:"payment_method"`
		CompletedAt   string  `json:"completed_at"`
	}

	var payload PakasirWebhookPayload
	if err := json.Unmarshal(rawBody, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse Pakasir webhook body: %w", err)
	}

	// 1. Verify project name matches our slug
	if payload.Project != s.cfg.Slug {
		return nil, fmt.Errorf("webhook project '%s' does not match configured project '%s'", payload.Project, s.cfg.Slug)
	}

	// 2. Query Pakasir detail API to double-check status safely (highly recommended by Pakasir docs)
	status, err := s.GetTransactionStatus(payload.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify transaction status via Pakasir API: %w", err)
	}

	return &NotificationPayload{
		OrderID:           payload.OrderID,
		TransactionID:     payload.OrderID, // Pakasir doesn't provide a separate transaction ID in webhook
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
