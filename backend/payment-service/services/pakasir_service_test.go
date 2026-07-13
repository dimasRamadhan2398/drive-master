package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"payment-service/models"
	"payment-service/pkg/config"
	"payment-service/repositories"

	"github.com/stretchr/testify/assert"
)

type mockPaymentRepo struct {
	repositories.IPaymentRepository
	payment *models.Payment
	err     error
}

func (m *mockPaymentRepo) GetByOrderID(orderID string) (*models.Payment, error) {
	return m.payment, m.err
}

func TestPakasirService_CreateCheckout(t *testing.T) {
	cfg := &config.PakasirConfig{
		Slug:    "drivemaster",
		APIKey:  "test_key",
		BaseURL: "https://app.pakasir.com",
	}

	repo := &mockPaymentRepo{}
	service := NewPakasirService(cfg, repo)

	resp, err := service.CreateCheckout("INV-12345", 150000, "Premium Package", "John Doe", "john@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "INV-12345", resp.OrderID)
	assert.True(t, strings.HasPrefix(resp.RedirectURL, "https://app.pakasir.com/pay/drivemaster/150000"))
	assert.Contains(t, resp.RedirectURL, "order_id=INV-12345")
}

func TestPakasirService_GetTransactionStatus(t *testing.T) {
	// Setup mock httptest server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/transactiondetail", r.URL.Path)
		assert.Equal(t, "drivemaster", r.URL.Query().Get("project"))
		assert.Equal(t, "100000", r.URL.Query().Get("amount"))
		assert.Equal(t, "INV-999", r.URL.Query().Get("order_id"))
		assert.Equal(t, "test_key", r.URL.Query().Get("api_key"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"transaction": {
				"amount": 100000,
				"order_id": "INV-999",
				"project": "drivemaster",
				"status": "completed",
				"payment_method": "qris",
				"completed_at": "2026-07-09T13:00:00Z"
			}
		}`))
	}))
	defer server.Close()

	cfg := &config.PakasirConfig{
		Slug:    "drivemaster",
		APIKey:  "test_key",
		BaseURL: server.URL,
	}

	mockPayment := &models.Payment{
		OrderID: "INV-999",
		Amount:  100000,
	}
	repo := &mockPaymentRepo{payment: mockPayment}
	service := NewPakasirService(cfg, repo)

	status, err := service.GetTransactionStatus("INV-999")
	assert.NoError(t, err)
	assert.Equal(t, models.PaymentStatusSuccess, status)
}

func TestPakasirService_VerifyNotification(t *testing.T) {
	// Setup mock httptest server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"transaction": {
				"amount": 100000,
				"order_id": "INV-999",
				"project": "drivemaster",
				"status": "completed",
				"payment_method": "qris",
				"completed_at": "2026-07-09T13:00:00Z"
			}
		}`))
	}))
	defer server.Close()

	cfg := &config.PakasirConfig{
		Slug:    "drivemaster",
		APIKey:  "test_key",
		BaseURL: server.URL,
	}

	mockPayment := &models.Payment{
		OrderID: "INV-999",
		Amount:  100000,
	}
	repo := &mockPaymentRepo{payment: mockPayment}
	service := NewPakasirService(cfg, repo)

	webhookBody, _ := json.Marshal(map[string]interface{}{
		"amount":         100000,
		"order_id":       "INV-999",
		"project":        "drivemaster",
		"status":         "completed",
		"payment_method": "qris",
		"completed_at":   "2026-07-09T13:00:00Z",
	})

	payload, err := service.VerifyNotification(nil, webhookBody)
	assert.NoError(t, err)
	assert.NotNil(t, payload)
	assert.Equal(t, "INV-999", payload.OrderID)
	assert.Equal(t, models.PaymentStatusSuccess, payload.TransactionStatus)
	assert.Equal(t, "qris", payload.PaymentType)
	assert.Equal(t, 100000.0, payload.GrossAmount)
}
