package regions

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/drive-master/clients/core-client/dto"
)

// IRegionClient defines the interface for region operations
type IRegionClient interface {
	GetAllProvinces(ctx context.Context) ([]dto.ProvinceResponse, error)
	GetRegenciesByProvince(ctx context.Context, provinceID string) ([]dto.RegencyResponse, error)
	GetDistrictsByRegency(ctx context.Context, provinceID, regencyID string) ([]dto.DistrictResponse, error)
}

// RegionClient implements IRegionClient
type RegionClient struct {
	baseURL    string
	httpClient *http.Client
	headers    map[string]string
}

// NewRegionClient creates a new region client
func NewRegionClient(baseURL string, headers map[string]string) IRegionClient {
	return &RegionClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * 1e9, // 30 seconds in nanoseconds
		},
		headers: headers,
	}
}

// APIResponse is the standard response format
type APIResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// GetAllProvinces retrieves all provinces
func (c *RegionClient) GetAllProvinces(ctx context.Context) ([]dto.ProvinceResponse, error) {
	resp, err := c.doRequest(ctx, "GET", "/regions/provinces", nil)
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

	var provinces []dto.ProvinceResponse
	if err := json.Unmarshal(apiResp.Data, &provinces); err != nil {
		return nil, fmt.Errorf("failed to unmarshal provinces: %w", err)
	}

	return provinces, nil
}

// GetRegenciesByProvince retrieves all regencies for a specific province
func (c *RegionClient) GetRegenciesByProvince(ctx context.Context, provinceID string) ([]dto.RegencyResponse, error) {
	resp, err := c.doRequest(ctx, "GET", "/regions/provinces/"+provinceID+"/regencies", nil)
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

	var regencies []dto.RegencyResponse
	if err := json.Unmarshal(apiResp.Data, &regencies); err != nil {
		return nil, fmt.Errorf("failed to unmarshal regencies: %w", err)
	}

	return regencies, nil
}

// GetDistrictsByRegency retrieves all districts for a specific regency
func (c *RegionClient) GetDistrictsByRegency(ctx context.Context, provinceID, regencyID string) ([]dto.DistrictResponse, error) {
	resp, err := c.doRequest(ctx, "GET", "/regions/regencies/"+regencyID+"/districts?province="+provinceID, nil)
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

	var districts []dto.DistrictResponse
	if err := json.Unmarshal(apiResp.Data, &districts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal districts: %w", err)
	}

	return districts, nil
}

// doRequest performs the HTTP request with headers
func (c *RegionClient) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
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