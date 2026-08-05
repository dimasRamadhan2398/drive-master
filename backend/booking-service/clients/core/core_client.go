package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"booking-service/models/dto"

	"github.com/google/uuid"
)

// ICoreClient defines the interface for core-service operations
type ICoreClient interface {
	GetCars(ctx context.Context, page, limit int) (*dto.PagedData[CarResponse], error)
	GetCarByID(ctx context.Context, carID uuid.UUID) (*CarInfo, error)
	GetSalesOverview(ctx context.Context, startDate, endDate string) (*SalesOverviewResponse, error)
	IncrementPackageCount(ctx context.Context, packageID uuid.UUID) error
	GetPackageByID(ctx context.Context, packageID uuid.UUID) (*PackageResponse, error)
	GetAddOnByID(ctx context.Context, addOnID uuid.UUID) (*AddOnResponse, error)
	CreateSale(ctx context.Context, req CreateSaleRequest) error
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
func (c *CoreClient) GetCarByID(ctx context.Context, carID uuid.UUID) (*CarInfo, error) {
	url := fmt.Sprintf("%s/api/v1/cars/%s", c.baseURL, carID.String())

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

// SalesOverviewResponse is the response DTO for sales overview
type SalesOverviewResponse struct {
	TotalRevenue     float64 `json:"totalRevenue"`
	TotalSales       int64   `json:"totalSales"`
	TotalRefunds     float64 `json:"totalRefunds"`
	NetRevenue       float64 `json:"netRevenue"`
	AvgOrderValue    float64 `json:"avgOrderValue"`
	CompletedSales   int64   `json:"completedSales"`
	PendingSales     int64   `json:"pendingSales"`
	CanceledSales    int64   `json:"canceledSales"`
	RefundedSales    int64   `json:"refundedSales"`
	GrowthRate       float64 `json:"growthRate"`
}

// GetSalesOverview retrieves sales overview from core-service
func (c *CoreClient) GetSalesOverview(ctx context.Context, startDate, endDate string) (*SalesOverviewResponse, error) {
	url := fmt.Sprintf("%s/api/v1/admin/sales/analytics/overview?startDate=%s&endDate=%s", c.baseURL, startDate, endDate)

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

	var overview SalesOverviewResponse
	if err := json.Unmarshal(apiResp.Data, &overview); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sales overview: %w", err)
	}

	return &overview, nil
}

// IncrementPackageCount increments the enrollment count for a package in core-service
func (c *CoreClient) IncrementPackageCount(ctx context.Context, packageID uuid.UUID) error {
	url := fmt.Sprintf("%s/api/v1/packages/%s/increment-count", c.baseURL, packageID.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call core-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("core-service returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetPackageByID retrieves a package by ID from core-service
func (c *CoreClient) GetPackageByID(ctx context.Context, packageID uuid.UUID) (*PackageResponse, error) {
	url := fmt.Sprintf("%s/api/v1/packages/%s", c.baseURL, packageID.String())

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

	var pkg PackageResponse
	if err := json.Unmarshal(apiResp.Data, &pkg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal package: %w", err)
	}

	return &pkg, nil
}

// GetAddOnByID retrieves an add-on by ID from core-service
func (c *CoreClient) GetAddOnByID(ctx context.Context, addOnID uuid.UUID) (*AddOnResponse, error) {
	url := fmt.Sprintf("%s/api/v1/addons/%s", c.baseURL, addOnID.String())

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

	var addon AddOnResponse
	if err := json.Unmarshal(apiResp.Data, &addon); err != nil {
		return nil, fmt.Errorf("failed to unmarshal add-on: %w", err)
	}

	return &addon, nil
}

// CreateSaleRequest is the payload for creating a sale in core-service
type CreateSaleRequest struct {
	UserID        string              `json:"userId"`
	PackageID     string              `json:"packageId,omitempty"`
	Items         []CreateSaleItem    `json:"items"`
	PaymentMethod string              `json:"paymentMethod"`
	Source        string              `json:"source"`
	Notes         string              `json:"notes"`
}

// CreateSaleItem represents a single item in a sale
type CreateSaleItem struct {
	PackageID   string  `json:"packageId"`
	PackageName string  `json:"packageName"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unitPrice"`
	Discount    float64 `json:"discount"`
}

// CreateSale creates a sale record in core-service
func (c *CoreClient) CreateSale(ctx context.Context, req CreateSaleRequest) error {
	url := fmt.Sprintf("%s/api/v1/admin/sales", c.baseURL)

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal sale request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to call core-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("core-service returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
