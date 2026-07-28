package clients

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/drive-master/clients/core-client/cars"
	"github.com/drive-master/clients/core-client/dto"
	"github.com/drive-master/clients/core-client/packages"
	"github.com/drive-master/clients/core-client/regions"
)

// Client is the main entry point for the core-service client
type Client struct {
	BaseURL      string
	ServiceName  string
	SignatureKey string
	Cars      cars.ICarClient
	Packages packages.IPackageClient
	Regions  regions.IRegionClient
}

// Config holds the configuration for the client
type Config struct {
	BaseURL      string // e.g., "http://localhost:8001/api/v1"
	ServiceName  string // e.g., "user-service"
	SignatureKey string // The signature key for API authentication
	Timeout      int    // Timeout in seconds (default: 30)
}

// NewClient creates a new core-service client
func NewClient(cfg Config) *Client {
	if cfg.Timeout == 0 {
		cfg.Timeout = 30
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:8001/api/v1"
	}

	headers := generateAuthHeaders(cfg.ServiceName, cfg.SignatureKey)

	return &Client{
		BaseURL:      baseURL,
		ServiceName:  cfg.ServiceName,
		SignatureKey: cfg.SignatureKey,
		Cars:      cars.NewCarClient(baseURL, headers),
		Packages: packages.NewPackageClient(baseURL, headers),
		Regions:   regions.NewRegionClient(baseURL, headers),
	}
}

// generateAuthHeaders generates the authentication headers
func generateAuthHeaders(serviceName, signatureKey string) map[string]string {
	requestAt := time.Now().UTC().Format(time.RFC3339)
	apiKey := generateAPIKey(serviceName, signatureKey, requestAt)

	return map[string]string{
		"x-service-name": serviceName,
		"x-request-at":  requestAt,
		"x-api-key":     apiKey,
	}
}

// generateAPIKey generates the API key signature
func generateAPIKey(serviceName, signatureKey, requestAt string) string {
	data := fmt.Sprintf("%s:%s:%s", serviceName, signatureKey, requestAt)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// NewClientWithBearerToken creates a client with Bearer token authentication
func NewClientWithBearerToken(baseURL, bearerToken string) *Client {
	if baseURL == "" {
		baseURL = "http://localhost:8001/api/v1"
	}

	headers := map[string]string{
		"Authorization": "Bearer " + bearerToken,
	}

	return &Client{
		BaseURL: baseURL,
		Cars:    cars.NewCarClient(baseURL, headers),
		Packages: packages.NewPackageClient(baseURL, headers),
		Regions:  regions.NewRegionClient(baseURL, headers),
	}
}

// CarService provides a simple interface for car operations
type CarService struct {
	client *Client
}

// NewCarService creates a new car service
func NewCarService(client *Client) *CarService {
	return &CarService{client: client}
}

// GetAll returns all cars
func (s *CarService) GetAll(ctx interface{}) ([]dto.CarResponse, error) {
	return s.client.Cars.GetAllCars(nil)
}

// GetByID returns a car by ID
func (s *CarService) GetByID(id string) (*dto.CarResponse, error) {
	return s.client.Cars.GetCarByID(nil, id)
}

// PackageService provides a simple interface for package operations
type PackageService struct {
	client *Client
}

// NewPackageService creates a new package service
func NewPackageService(client *Client) *PackageService {
	return &PackageService{client: client}
}

// GetAll returns all packages
func (s *PackageService) GetAll() ([]dto.PackageResponse, error) {
	return s.client.Packages.GetAllPackages(nil)
}

// GetByID returns a package by ID
func (s *PackageService) GetByID(id string) (*dto.PackageResponse, error) {
	return s.client.Packages.GetPackageByID(nil, id)
}

// RegionService provides a simple interface for region operations
type RegionService struct {
	client *Client
}

// NewRegionService creates a new region service
func NewRegionService(client *Client) *RegionService {
	return &RegionService{client: client}
}

// GetProvinces returns all provinces
func (s *RegionService) GetProvinces() ([]dto.ProvinceResponse, error) {
	return s.client.Regions.GetAllProvinces(nil)
}

// GetRegencies returns all regencies for a province
func (s *RegionService) GetRegencies(provinceID string) ([]dto.RegencyResponse, error) {
	return s.client.Regions.GetRegenciesByProvince(nil, provinceID)
}

// GetDistricts returns all districts for a regency
func (s *RegionService) GetDistricts(provinceID, regencyID string) ([]dto.DistrictResponse, error) {
	return s.client.Regions.GetDistrictsByRegency(nil, provinceID, regencyID)
}