package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ClientConfig holds configuration for the core client
type ClientConfig struct {
	BaseURL    string        // e.g., "http://localhost:8001/api/v1"
	Timeout    time.Duration // e.g., 30 * time.Second
	Retries    int           // number of retries on failure
	RetryDelay time.Duration // delay between retries
}

// DefaultClientConfig returns a default configuration
func DefaultClientConfig(baseURL string) ClientConfig {
	return ClientConfig{
		BaseURL:    baseURL,
		Timeout:    30 * time.Second,
		Retries:    3,
		RetryDelay: 1 * time.Second,
	}
}

// HTTPClient is the base HTTP client for making requests
type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
	retries    int
	retryDelay time.Duration
}

// NewHTTPClient creates a new HTTP client
func NewHTTPClient(cfg ClientConfig) *HTTPClient {
	return &HTTPClient{
		baseURL: cfg.BaseURL,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		retries:    cfg.Retries,
		retryDelay: cfg.RetryDelay,
	}
}

// APIResponse represents a standard API response from core-service
type APIResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   *APIError       `json:"error,omitempty"`
}

// APIError represents an error from the API
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// doRequest performs an HTTP request with retry logic
func (c *HTTPClient) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := fmt.Sprintf("%s%s", c.baseURL, path)
	var lastErr error

	for i := 0; i <= c.retries; i++ {
		if i > 0 {
			time.Sleep(c.retryDelay)
		}

		req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
		if err != nil {
			lastErr = err
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		// Don't retry on success or client errors (4xx except 429)
		if resp.StatusCode < 500 && resp.StatusCode != http.StatusTooManyRequests {
			return resp, nil
		}

		// For server errors or rate limiting, retry
		resp.Body.Close()
		lastErr = fmt.Errorf("server error: %d", resp.StatusCode)
	}

	return nil, lastErr
}

// decodeResponse decodes the API response
func decodeResponse[T any](resp *http.Response) (*T, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	if apiResp.Error != nil {
		return nil, fmt.Errorf("%s: %s", apiResp.Error.Code, apiResp.Error.Message)
	}

	if apiResp.Data == nil {
		return nil, fmt.Errorf("no data in response")
	}

	var result T
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return &result, nil
}

// decodeResponseRaw decodes the response directly without APIResponse wrapper
func decodeResponseRaw[T any](resp *http.Response) (*T, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}