package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"user-service/models"
)

type RegionClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewRegionClient(baseURL string) *RegionClient {
	return &RegionClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// RegionResponse represents the API response structure from core-service
type RegionResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// GetAllProvinces fetches all provinces from core-service
func (c *RegionClient) GetAllProvinces(ctx context.Context) ([]models.Province, error) {
	url := fmt.Sprintf("%s/api/v1/regions/provinces", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch provinces: status %d", resp.StatusCode)
	}

	var result RegionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Parse the data field into provinces
	data, err := json.Marshal(result.Data)
	if err != nil {
		return nil, err
	}

	var provinces []models.Province
	if err := json.Unmarshal(data, &provinces); err != nil {
		return nil, err
	}

	return provinces, nil
}

// GetRegenciesByProvince fetches regencies for a specific province
func (c *RegionClient) GetRegenciesByProvince(ctx context.Context, province string) ([]models.Regency, error) {
	url := fmt.Sprintf("%s/api/v1/regions/provinces/%s/regencies", c.baseURL, province)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch regencies: status %d", resp.StatusCode)
	}

	var result RegionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	data, err := json.Marshal(result.Data)
	if err != nil {
		return nil, err
	}

	var regencies []models.Regency
	if err := json.Unmarshal(data, &regencies); err != nil {
		return nil, err
	}

	return regencies, nil
}

// GetDistrictsByRegency fetches districts for a specific regency
func (c *RegionClient) GetDistrictsByRegency(ctx context.Context, province, regency string) ([]models.District, error) {
	url := fmt.Sprintf("%s/api/v1/regions/regencies/%s/districts?province=%s", c.baseURL, regency, province)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch districts: status %d", resp.StatusCode)
	}

	var result RegionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	data, err := json.Marshal(result.Data)
	if err != nil {
		return nil, err
	}

	var districts []models.District
	if err := json.Unmarshal(data, &districts); err != nil {
		return nil, err
	}

	return districts, nil
}