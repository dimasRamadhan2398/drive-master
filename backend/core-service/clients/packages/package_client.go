package packages

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/drive-master/clients/core-client/dto"
)

// IPackageClient defines the interface for package operations
type IPackageClient interface {
	GetAllPackages(ctx context.Context) ([]dto.PackageResponse, error)
	GetPackageByID(ctx context.Context, id string) (*dto.PackageResponse, error)
	CreatePackage(ctx context.Context, req *dto.CreatePackageRequest) (*dto.PackageResponse, error)
	UpdatePackage(ctx context.Context, id string, req *dto.UpdatePackageRequest) (*dto.PackageResponse, error)
	DeletePackage(ctx context.Context, id string) error
}

// PackageClient implements IPackageClient
type PackageClient struct {
	baseURL    string
	httpClient *http.Client
	headers    map[string]string
}

// NewPackageClient creates a new package client
func NewPackageClient(baseURL string, headers map[string]string) IPackageClient {
	return &PackageClient{
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

// GetAllPackages retrieves all packages
func (c *PackageClient) GetAllPackages(ctx context.Context) ([]dto.PackageResponse, error) {
	resp, err := c.doRequest(ctx, "GET", "/packages", nil)
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

	var packages []dto.PackageResponse
	if err := json.Unmarshal(apiResp.Data, &packages); err != nil {
		return nil, fmt.Errorf("failed to unmarshal packages: %w", err)
	}

	return packages, nil
}

// GetPackageByID retrieves a package by ID
func (c *PackageClient) GetPackageByID(ctx context.Context, id string) (*dto.PackageResponse, error) {
	resp, err := c.doRequest(ctx, "GET", "/packages/"+id, nil)
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

	var pkg dto.PackageResponse
	if err := json.Unmarshal(apiResp.Data, &pkg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal package: %w", err)
	}

	return &pkg, nil
}

// CreatePackage creates a new package
func (c *PackageClient) CreatePackage(ctx context.Context, req *dto.CreatePackageRequest) (*dto.PackageResponse, error) {
	resp, err := c.doRequest(ctx, "POST", "/packages", req)
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

	var pkg dto.PackageResponse
	if err := json.Unmarshal(apiResp.Data, &pkg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal package: %w", err)
	}

	return &pkg, nil
}

// UpdatePackage updates an existing package
func (c *PackageClient) UpdatePackage(ctx context.Context, id string, req *dto.UpdatePackageRequest) (*dto.PackageResponse, error) {
	resp, err := c.doRequest(ctx, "PUT", "/packages/"+id, req)
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

	var pkg dto.PackageResponse
	if err := json.Unmarshal(apiResp.Data, &pkg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal package: %w", err)
	}

	return &pkg, nil
}

// DeletePackage deletes a package
func (c *PackageClient) DeletePackage(ctx context.Context, id string) error {
	resp, err := c.doRequest(ctx, "DELETE", "/packages/"+id, nil)
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
func (c *PackageClient) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
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