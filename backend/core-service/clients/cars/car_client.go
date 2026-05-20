package cars

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/drive-master/clients/core-client/dto"
)

// ICarClient defines the interface for car operations
type ICarClient interface {
	GetAllCars(ctx context.Context) ([]dto.CarResponse, error)
	GetCarByID(ctx context.Context, id string) (*dto.CarResponse, error)
	CreateCar(ctx context.Context, req *dto.CreateCarRequest) (*dto.CarResponse, error)
	UpdateCar(ctx context.Context, id string, req *dto.UpdateCarRequest) (*dto.CarResponse, error)
	DeleteCar(ctx context.Context, id string) error
}

// CarClient implements ICarClient
type CarClient struct {
	baseURL    string
	httpClient *http.Client
	headers    map[string]string
}

// NewCarClient creates a new car client
func NewCarClient(baseURL string, headers map[string]string) ICarClient {
	return &CarClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * 1e9, // 30 seconds in nanoseconds
		},
		headers: headers,
	}
}

// NewCarClientWithHeaders creates a car client with custom headers (for auth)
func NewCarClientWithHeaders(baseURL string, headers map[string]string) ICarClient {
	return &CarClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * 1e9,
		},
		headers: headers,
	}
}

// APIResponse is the standard response format
type APIResponse struct {
	Success bool            `json:"success"`
	Message string         `json:"message,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// GetAllCars retrieves all cars
func (c *CarClient) GetAllCars(ctx context.Context) ([]dto.CarResponse, error) {
	resp, err := c.doRequest(ctx, "GET", "/cars", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API error: %s", apiResp.Message)
	}

	var cars []dto.CarResponse
	if err := json.Unmarshal(apiResp.Data, &cars); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cars: %w", err)
	}

	return cars, nil
}

// GetCarByID retrieves a car by ID
func (c *CarClient) GetCarByID(ctx context.Context, id string) (*dto.CarResponse, error) {
	resp, err := c.doRequest(ctx, "GET", "/cars/"+id, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API error: %s", apiResp.Message)
	}

	var car dto.CarResponse
	if err := json.Unmarshal(apiResp.Data, &car); err != nil {
		return nil, fmt.Errorf("failed to unmarshal car: %w", err)
	}

	return &car, nil
}

// CreateCar creates a new car
func (c *CarClient) CreateCar(ctx context.Context, req *dto.CreateCarRequest) (*dto.CarResponse, error) {
	resp, err := c.doRequest(ctx, "POST", "/cars", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API error: %s", apiResp.Message)
	}

	var car dto.CarResponse
	if err := json.Unmarshal(apiResp.Data, &car); err != nil {
		return nil, fmt.Errorf("failed to unmarshal car: %w", err)
	}

	return &car, nil
}

// UpdateCar updates an existing car
func (c *CarClient) UpdateCar(ctx context.Context, id string, req *dto.UpdateCarRequest) (*dto.CarResponse, error) {
	resp, err := c.doRequest(ctx, "PUT", "/cars/"+id, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API error: %s", apiResp.Message)
	}

	var car dto.CarResponse
	if err := json.Unmarshal(apiResp.Data, &car); err != nil {
		return nil, fmt.Errorf("failed to unmarshal car: %w", err)
	}

	return &car, nil
}

// DeleteCar deletes a car
func (c *CarClient) DeleteCar(ctx context.Context, id string) error {
	resp, err := c.doRequest(ctx, "DELETE", "/cars/"+id, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		return fmt.Errorf("API error: %s", apiResp.Message)
	}

	return nil
}

// doRequest performs the HTTP request with headers
func (c *CarClient) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	url := c.baseURL + path

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

	// Add custom headers
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	return c.httpClient.Do(req)
}