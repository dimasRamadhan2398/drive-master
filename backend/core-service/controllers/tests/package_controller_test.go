package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"core-service/models"
	"core-service/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPackageService mocks the package service for testing
type MockPackageService struct {
	mock.Mock
}

func (m *MockPackageService) CreatePackage(ctx context.Context, pkg *models.Package) error {
	args := m.Called(ctx, pkg)
	return args.Error(0)
}

func (m *MockPackageService) GetPackageByID(ctx context.Context, id uuid.UUID) (*models.Package, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Package), args.Error(1)
}

func (m *MockPackageService) GetPackageByIDWithBenefits(ctx context.Context, id uuid.UUID) (*models.Package, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Package), args.Error(1)
}

func (m *MockPackageService) GetAllPackages(ctx context.Context) ([]models.Package, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Package), args.Error(1)
}

func (m *MockPackageService) GetPackagesByType(ctx context.Context, packageType models.PackageType) ([]models.Package, error) {
	args := m.Called(ctx, packageType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Package), args.Error(1)
}

func (m *MockPackageService) GetPackagesByStatus(ctx context.Context, status models.PackageStatus) ([]models.Package, error) {
	args := m.Called(ctx, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Package), args.Error(1)
}

func (m *MockPackageService) UpdatePackage(ctx context.Context, pkg *models.Package) error {
	args := m.Called(ctx, pkg)
	return args.Error(0)
}

func (m *MockPackageService) DeletePackage(ctx context.Context, pkg *models.Package) error {
	args := m.Called(ctx, pkg)
	return args.Error(0)
}

func (m *MockPackageService) CountPackages(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// setupPackageRouter creates a test router for package endpoints
func setupPackageRouter(controller IPackageController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/packages", controller.GetAllPackages)
		v1.GET("/packages/:id", controller.GetPackageByID)
		v1.POST("/packages", controller.CreatePackage)
		v1.PUT("/packages/:id", controller.UpdatePackage)
		v1.DELETE("/packages/:id", controller.DeletePackage)
	}

	return router
}

// ========== GetAllPackages Tests ==========

func TestGetAllPackages_Success(t *testing.T) {
	mockService := new(MockPackageService)
	controller := NewPackageController(mockService)
	router := setupPackageRouter(controller)

	packages := []models.Package{
		{
			ID:              uuid.New(),
			Name:            "Basic Package",
			Description:     "A basic driving package",
			PackageType:     models.PackageTypeBronze,
			Price:           500000,
			DiscountPrice:   450000,
			DurationMinutes: 60,
			TotalSessions:   4,
			Status:          models.PackageStatusActive,
		},
		{
			ID:              uuid.New(),
			Name:            "Premium Package",
			Description:     "A premium driving package",
			PackageType:     models.PackageTypeGold,
			Price:           1000000,
			DiscountPrice:   900000,
			DurationMinutes: 90,
			TotalSessions:   8,
			Status:          models.PackageStatusActive,
		},
	}

	mockService.On("GetAllPackages", mock.Anything).Return(packages, nil)

	req, _ := http.NewRequest("GET", "/api/v1/packages", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "Packages fetched successfully", resp.Message)

	data, ok := resp.Data.([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 2, len(data))

	mockService.AssertExpectations(t)
}

func TestGetAllPackages_ServiceError(t *testing.T) {
	mockService := new(MockPackageService)
	controller := NewPackageController(mockService)
	router := setupPackageRouter(controller)

	mockService.On("GetAllPackages", mock.Anything).Return(nil, assert.AnError)

	req, _ := http.NewRequest("GET", "/api/v1/packages", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.NotNil(t, resp.Error)

	mockService.AssertExpectations(t)
}

func TestGetAllPackages_EmptyResult(t *testing.T) {
	mockService := new(MockPackageService)
	controller := NewPackageController(mockService)
	router := setupPackageRouter(controller)

	packages := []models.Package{}

	mockService.On("GetAllPackages", mock.Anything).Return(packages, nil)

	req, _ := http.NewRequest("GET", "/api/v1/packages", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)

	data, ok := resp.Data.([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 0, len(data))

	mockService.AssertExpectations(t)
}

// ========== GetPackageByID Tests ==========

func TestGetPackageByID_Success(t *testing.T) {
	mockService := new(MockPackageService)
	controller := NewPackageController(mockService)
	router := setupPackageRouter(controller)

	pkgID := uuid.New()
	pkg := &models.Package{
		ID:              pkgID,
		Name:            "Basic Package",
		Description:     "A basic driving package",
		PackageType:     models.PackageTypeBronze,
		Price:           500000,
		DiscountPrice:   450000,
		DurationMinutes: 60,
		TotalSessions:   4,
		Status:          models.PackageStatusActive,
	}

	mockService.On("GetPackageByID", mock.Anything, pkgID).Return(pkg, nil)

	req, _ := http.NewRequest("GET", "/api/v1/packages/"+pkgID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "Package fetched successfully", resp.Message)
	assert.NotNil(t, resp.Data)

	mockService.AssertExpectations(t)
}

func TestGetPackageByID_NotFound(t *testing.T) {
	mockService := new(MockPackageService)
	controller := NewPackageController(mockService)
	router := setupPackageRouter(controller)

	pkgID := uuid.New()

	mockService.On("GetPackageByID", mock.Anything, pkgID).Return(nil, nil)

	req, _ := http.NewRequest("GET", "/api/v1/packages/"+pkgID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.NotNil(t, resp.Error)

	mockService.AssertExpectations(t)
}

func TestGetPackageByID_InvalidIDFormat(t *testing.T) {
	mockService := new(MockPackageService)
	controller := NewPackageController(mockService)
	router := setupPackageRouter(controller)

	req, _ := http.NewRequest("GET", "/api/v1/packages/invalid-uuid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.NotNil(t, resp.Error)
	assert.Equal(t, "Invalid package ID format", resp.Error.Message)

	mockService.AssertNotCalled(t, "GetPackageByID")
}

func TestGetPackageByID_ServiceError(t *testing.T) {
	mockService := new(MockPackageService)
	controller := NewPackageController(mockService)
	router := setupPackageRouter(controller)

	pkgID := uuid.New()

	mockService.On("GetPackageByID", mock.Anything, pkgID).Return(nil, assert.AnError)

	req, _ := http.NewRequest("GET", "/api/v1/packages/"+pkgID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)

	mockService.AssertExpectations(t)
}

// ========== CreatePackage Tests ==========

func TestCreatePackage_Success(t *testing.T) {
	mockService := new(MockPackageService)
	controller := NewPackageController(mockService)
	router := setupPackageRouter(controller)

	mockService.On("CreatePackage", mock.Anything, mock.AnythingOfType("*models.Package")).Return(nil)

	createPackageJSON := `{
		"name": "Basic Package",
		"description": "A basic driving package",
		"packageType": "bronze",
		"price": 500000,
		"discountPrice": 450000,
		"durationMinutes": 60,
		"totalSessions": 4
	}`

	req, _ := http.NewRequest("POST", "/api/v1/packages", bytes.NewBufferString(createPackageJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "Package created successfully", resp.Message)

	mockService.AssertExpectations(t)
}

func TestCreatePackage_WithBenefits(t *testing.T) {
	mockService := new(MockPackageService)
	controller := NewPackageController(mockService)
	router := setupPackageRouter(controller)

	mockService.On("CreatePackage", mock.Anything, mock.AnythingOfType("*models.Package")).Return(nil)

	createPackageJSON := `{
		"name": "Premium Package",
		"description": "A premium driving package",
		"packageType": "gold",
		"price": 1000000,
		"durationMinutes": 90,
		"totalSessions": 8,
		"benefits": [
			{
				"title": "8 Sessions",
				"description": "8 driving sessions with instructor",
				"icon": "car",
				"sortOrder": 1
			},
			{
				"title": "Free Pickup",
				"description": "Free pickup service",
				"icon": "location",
				"sortOrder": 2
			}
		]
	}`

	req, _ := http.NewRequest("POST", "/api/v1/packages", bytes.NewBufferString(createPackageJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "Package created successfully", resp.Message)

	mockService.AssertExpectations(t)
}

func TestCreatePackage_InvalidJSON(t *testing.T) {
	mockService := new(MockPackageService)
	controller := NewPackageController(mockService)
	router := setupPackageRouter(controller)

	req, _ := http.NewRequest("POST", "/api/v1/packages", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)

	mockService.AssertNotCalled(t, "CreatePackage")
}

func TestCreatePackage_MissingRequiredFields(t *testing.T) {
	mockService := new(MockPackageService)
	controller := NewPackageController(mockService)
	router := setupPackageRouter(controller)

	// Missing name, packageType, price
	createPackageJSON := `{
		"description": "A basic driving package"
	}`

	req, _ := http.NewRequest("POST", "/api/v1/packages", bytes.NewBufferString(createPackageJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)

	mockService.AssertNotCalled(t, "CreatePackage")
}

func TestCreatePackage_ServiceError(t *testing.T) {
	mockService := new(MockPackageService)
	controller := NewPackageController(mockService)
	router := setupPackageRouter(controller)

	mockService.On("CreatePackage", mock.Anything, mock.AnythingOfType("*models.Package")).Return(assert.AnError)

	createPackageJSON := `{
		"name": "Basic Package",
		"packageType": "bronze",
		"price": 500000
	}`

	req, _ := http.NewRequest("POST", "/api/v1/packages", bytes.NewBufferString(createPackageJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)

	mockService.AssertExpectations(t)
}

// ========== UpdatePackage Tests ==========

func TestUpdatePackage_Success(t *testing.T) {
	mockService := new(MockPackageService)
	controller := NewPackageController(mockService)
	router := setupPackageRouter(controller)

	pkgID := uuid.New()
	existingPkg := &models.Package{
		ID:              pkgID,
		Name:            "Basic Package",
		Description:     "A basic driving package",
		PackageType:     models.PackageTypeBronze,
		Price:           500000,
		DiscountPrice:   450000,
		DurationMinutes: 60,
		TotalSessions:   4,
		Status:          models.PackageStatusActive,
	}

	mockService.On("GetPackageByID", mock.Anything, pkgID).Return(existingPkg, nil).Once()
	mockService.On("UpdatePackage", mock.Anything, mock.AnythingOfType("*models.Package")).Return(nil)

	updatePackageJSON := `{
		"description": "Updated description",
		"price": 550000,
		"discountPrice": 500000
	}`

	req, _ := http.NewRequest("PUT", "/api/v1/packages/"+pkgID.String(), bytes.NewBufferString(updatePackageJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "Package updated successfully", resp.Message)

	mockService.AssertExpectations(t)
}

func TestUpdatePackage_NotFound(t *testing.T) {
	mockService := new(MockPackageService)
	controller := NewPackageController(mockService)
	router := setupPackageRouter(controller)

	pkgID := uuid.New()

	mockService.On("GetPackageByID", mock.Anything, pkgID).Return(nil, nil)

	updatePackageJSON := `{
		"price": 550000
	}`

	req, _ := http.NewRequest("PUT", "/api/v1/packages/"+pkgID.String(), bytes.NewBufferString(updatePackageJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)

	mockService.AssertExpectations(t)
}

func TestUpdatePackage_InvalidIDFormat(t *testing.T) {
	mockService := new(MockPackageService)
	controller := NewPackageController(mockService)
	router := setupPackageRouter(controller)

	updatePackageJSON := `{
		"price": 550000
	}`

	req, _ := http.NewRequest("PUT", "/api/v1/packages/invalid-uuid", bytes.NewBufferString(updatePackageJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "Invalid package ID format", resp.Error.Message)

	mockService.AssertNotCalled(t, "GetPackageByID")
}

func TestUpdatePackage_InvalidJSON(t *testing.T) {
	mockService := new(MockPackageService)
	controller := NewPackageController(mockService)
	router := setupPackageRouter(controller)

	pkgID := uuid.New()
	existingPkg := &models.Package{
		ID:              pkgID,
		Name:            "Basic Package",
		Description:     "A basic driving package",
		PackageType:     models.PackageTypeBronze,
		Price:           500000,
		DiscountPrice:   450000,
		DurationMinutes: 60,
		TotalSessions:   4,
		Status:          models.PackageStatusActive,
	}

	mockService.On("GetPackageByID", mock.Anything, pkgID).Return(existingPkg, nil)

	req, _ := http.NewRequest("PUT", "/api/v1/packages/"+pkgID.String(), bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)

	// GetPackageByID was called but UpdatePackage was not due to invalid JSON
	mockService.AssertCalled(t, "GetPackageByID", mock.Anything, pkgID)
}

func TestUpdatePackage_ServiceError(t *testing.T) {
	mockService := new(MockPackageService)
	controller := NewPackageController(mockService)
	router := setupPackageRouter(controller)

	pkgID := uuid.New()
	existingPkg := &models.Package{
		ID:              pkgID,
		Name:            "Basic Package",
		Description:     "A basic driving package",
		PackageType:     models.PackageTypeBronze,
		Price:           500000,
		DiscountPrice:   450000,
		DurationMinutes: 60,
		TotalSessions:   4,
		Status:          models.PackageStatusActive,
	}

	mockService.On("GetPackageByID", mock.Anything, pkgID).Return(existingPkg, nil)
	mockService.On("UpdatePackage", mock.Anything, mock.AnythingOfType("*models.Package")).Return(assert.AnError)

	updatePackageJSON := `{
		"price": 550000
	}`

	req, _ := http.NewRequest("PUT", "/api/v1/packages/"+pkgID.String(), bytes.NewBufferString(updatePackageJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)

	mockService.AssertExpectations(t)
}

// ========== DeletePackage Tests ==========

func TestDeletePackage_Success(t *testing.T) {
	mockService := new(MockPackageService)
	controller := NewPackageController(mockService)
	router := setupPackageRouter(controller)

	pkgID := uuid.New()
	existingPkg := &models.Package{
		ID:              pkgID,
		Name:            "Basic Package",
		Description:     "A basic driving package",
		PackageType:     models.PackageTypeBronze,
		Price:           500000,
		DiscountPrice:   450000,
		DurationMinutes: 60,
		TotalSessions:   4,
		Status:          models.PackageStatusActive,
	}

	mockService.On("GetPackageByID", mock.Anything, pkgID).Return(existingPkg, nil)
	mockService.On("DeletePackage", mock.Anything, existingPkg).Return(nil)

	req, _ := http.NewRequest("DELETE", "/api/v1/packages/"+pkgID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "Package deleted successfully", resp.Message)

	mockService.AssertExpectations(t)
}

func TestDeletePackage_NotFound(t *testing.T) {
	mockService := new(MockPackageService)
	controller := NewPackageController(mockService)
	router := setupPackageRouter(controller)

	pkgID := uuid.New()

	mockService.On("GetPackageByID", mock.Anything, pkgID).Return(nil, nil)

	req, _ := http.NewRequest("DELETE", "/api/v1/packages/"+pkgID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)

	mockService.AssertExpectations(t)
}

func TestDeletePackage_InvalidIDFormat(t *testing.T) {
	mockService := new(MockPackageService)
	controller := NewPackageController(mockService)
	router := setupPackageRouter(controller)

	req, _ := http.NewRequest("DELETE", "/api/v1/packages/invalid-uuid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "Invalid package ID format", resp.Error.Message)

	mockService.AssertNotCalled(t, "GetPackageByID")
}

func TestDeletePackage_ServiceError(t *testing.T) {
	mockService := new(MockPackageService)
	controller := NewPackageController(mockService)
	router := setupPackageRouter(controller)

	pkgID := uuid.New()
	existingPkg := &models.Package{
		ID:              pkgID,
		Name:            "Basic Package",
		Description:     "A basic driving package",
		PackageType:     models.PackageTypeBronze,
		Price:           500000,
		DiscountPrice:   450000,
		DurationMinutes: 60,
		TotalSessions:   4,
		Status:          models.PackageStatusActive,
	}

	mockService.On("GetPackageByID", mock.Anything, pkgID).Return(existingPkg, nil)
	mockService.On("DeletePackage", mock.Anything, existingPkg).Return(assert.AnError)

	req, _ := http.NewRequest("DELETE", "/api/v1/packages/"+pkgID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)

	mockService.AssertExpectations(t)
}