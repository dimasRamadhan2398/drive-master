package clients

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// ITransactionClient defines the interface for transaction operations
type ITransactionClient interface {
	GetRecentTransactions(ctx context.Context, page, limit int) (*TransactionsResponse, error)
	GetTransactionByID(ctx context.Context, id string) (*TransactionResponse, error)
}

// TransactionResponse represents a transaction from payment service
type TransactionResponse struct {
	ID            uuid.UUID `json:"id"`
	PaymentID     uuid.UUID `json:"paymentId"`
	Type          string    `json:"type"`
	Status        string    `json:"status"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	Gateway       string    `json:"gateway"`
	GatewayTxnID  string    `json:"gatewayTxnId"`
	ErrorCode     string    `json:"errorCode,omitempty"`
	ErrorMessage  string    `json:"errorMessage,omitempty"`
	ProcessedAt   *time.Time `json:"processedAt,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
}

// TransactionsResponse represents the response for listing transactions
type TransactionsResponse struct {
	Data  []TransactionResponse `json:"data"`
	Total int64                `json:"total"`
	Page  int                  `json:"page"`
	Limit int                  `json:"limit"`
}

// APIResponse is the standard response format from payment service
type APIResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// TransactionClient implements ITransactionClient
type TransactionClient struct {
	baseURL       string
	serviceName   string
	signatureKey  string
	httpClient    *http.Client
}

// NewTransactionClient creates a new transaction client with signature-based auth
func NewTransactionClient(baseURL, serviceName, signatureKey string) ITransactionClient {
	return &TransactionClient{
		baseURL:      baseURL,
		serviceName:  serviceName,
		signatureKey: signatureKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// generateAPIKey generates the API key signature
func generateAPIKey(serviceName, signatureKey, requestAt string) string {
	data := fmt.Sprintf("%s:%s:%s", serviceName, signatureKey, requestAt)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// GetRecentTransactions retrieves recent transactions from payment service
func (c *TransactionClient) GetRecentTransactions(ctx context.Context, page, limit int) (*TransactionsResponse, error) {
	url := fmt.Sprintf("%s/transactions/all?page=%d&limit=%d", c.baseURL, page, limit)
	return c.doGet(ctx, url)
}

// GetTransactionByID retrieves a single transaction by ID
func (c *TransactionClient) GetTransactionByID(ctx context.Context, id string) (*TransactionResponse, error) {
	url := fmt.Sprintf("%s/transactions/%s", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add signature-based auth headers
	requestAt := time.Now().UTC().Format(time.RFC3339)
	apiKey := generateAPIKey(c.serviceName, c.signatureKey, requestAt)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-service-name", c.serviceName)
	req.Header.Set("x-request-at", requestAt)
	req.Header.Set("x-api-key", apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API error: %s", apiResp.Message)
	}

	var tx TransactionResponse
	if err := json.Unmarshal(apiResp.Data, &tx); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}

	return &tx, nil
}

// doGet performs a GET request with signature-based auth headers
func (c *TransactionClient) doGet(ctx context.Context, url string) (*TransactionsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add signature-based auth headers
	requestAt := time.Now().UTC().Format(time.RFC3339)
	apiKey := generateAPIKey(c.serviceName, c.signatureKey, requestAt)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-service-name", c.serviceName)
	req.Header.Set("x-request-at", requestAt)
	req.Header.Set("x-api-key", apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API error: %s", apiResp.Message)
	}

	// Parse the nested data structure
	var data struct {
		Data  []TransactionResponse `json:"data"`
		Total int64                `json:"total"`
		Page  int                  `json:"page"`
		Limit int                  `json:"limit"`
	}
	if err := json.Unmarshal(apiResp.Data, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transactions: %w", err)
	}

	return &TransactionsResponse{
		Data:  data.Data,
		Total: data.Total,
		Page:  data.Page,
		Limit: data.Limit,
	}, nil
}

// doRequest performs the HTTP request with headers (for future use)
func doRequest(ctx context.Context, method, url string, body interface{}, headers map[string]string) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		reqBody = io.NopCloser(bytes.NewReader(jsonData))
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	return client.Do(req)
}
