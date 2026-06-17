package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"booking-service/models/dto"
)

// ICoreClient defines the interface for core-service operations
type ICoreClient interface {
	GetCars(ctx context.Context, page, limit int) (*dto.PagedData[CarResponse], error)
	GetCarByID(ctx context.Context, carID uint) (*CarInfo, error)
}

// CoreClient implements ICoreClient
type CoreClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewCoreClient creates a new core service client
func NewCoreClient(baseURL string) ICoreClient {
	return &CoreClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetCars retrieves a paginated list of cars from core-service
func (c *CoreClient) GetCars(ctx context.Context, page, limit int) (*dto.PagedData[CarResponse], error) {
	url := fmt.Sprintf("%s/api/v1/cars?page=%d&limit=%d", c.baseURL, page, limit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call core-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("core-service returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API error: %s", apiResp.Message)
	}

	var cars dto.PagedData[CarResponse]
	if err := json.Unmarshal(apiResp.Data, &cars); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cars: %w", err)
	}

	return &cars, nil
}

// APIResponse is the standard response format from core-service
type APIResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// GetCarByID retrieves a single car by ID from core-service
func (c *CoreClient) GetCarByID(ctx context.Context, carID uint) (*CarInfo, error) {
	url := fmt.Sprintf("%s/api/v1/cars/%d", c.baseURL, carID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call core-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("core-service returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API error: %s", apiResp.Message)
	}

	var car CarInfo
	if err := json.Unmarshal(apiResp.Data, &car); err != nil {
		return nil, fmt.Errorf("failed to unmarshal car: %w", err)
	}

	return &car, nil
}
