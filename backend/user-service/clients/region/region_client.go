package region

import (
	"context"
	"fmt"

	"user-service/clients"
)

// Client wraps the configuration and provides region operations
type Client struct {
	config clients.IClientConfig
}

// NewClient creates a new region client with the given config
func NewClient(config clients.IClientConfig) *Client {
	return &Client{
		config: config,
	}
}

// GetAllProvinces fetches all provinces from core-service
func (c *Client) GetAllProvinces(ctx context.Context) ([]Province, error) {
	url := fmt.Sprintf("%s/regions/provinces", c.config.BaseURL())

	var result APIResponse
	_, _, errs := c.config.Client().Get(url).EndStruct(&result)
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to fetch provinces: %v", errs[0])
	}

	if !result.Success {
		return nil, fmt.Errorf("API error: %s", result.Message)
	}

	return result.Data.Provinces, nil
}

// GetRegenciesByProvince fetches regencies for a specific province
func (c *Client) GetRegenciesByProvince(ctx context.Context, provinceID string) ([]Regency, error) {
	url := fmt.Sprintf("%s/regions/provinces/%s/regencies", c.config.BaseURL(), provinceID)

	var result APIResponse
	_, _, errs := c.config.Client().Get(url).EndStruct(&result)
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to fetch regencies: %v", errs[0])
	}

	if !result.Success {
		return nil, fmt.Errorf("API error: %s", result.Message)
	}

	return result.Data.Regencies, nil
}

// GetDistrictsByRegency fetches districts for a specific regency
func (c *Client) GetDistrictsByRegency(ctx context.Context, provinceID, regencyID string) ([]District, error) {
	url := fmt.Sprintf("%s/regions/regencies/%s/districts?province=%s", c.config.BaseURL(), regencyID, provinceID)

	var result APIResponse
	_, _, errs := c.config.Client().Get(url).EndStruct(&result)
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to fetch districts: %v", errs[0])
	}

	if !result.Success {
		return nil, fmt.Errorf("API error: %s", result.Message)
	}

	return result.Data.Districts, nil
}