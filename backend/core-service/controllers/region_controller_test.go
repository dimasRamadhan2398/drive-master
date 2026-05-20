package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"core-service/models"
	"core-service/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRegionService mocks the region service for testing
type MockRegionService struct {
	mock.Mock
}

func (m *MockRegionService) GetAllProvinces(ctx context.Context) ([]models.Province, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Province), args.Error(1)
}

func (m *MockRegionService) GetRegenciesByProvince(ctx context.Context, province string) ([]models.Regency, error) {
	args := m.Called(ctx, province)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Regency), args.Error(1)
}

func (m *MockRegionService) GetDistrictsByRegency(ctx context.Context, province, regency string) ([]models.District, error) {
	args := m.Called(ctx, province, regency)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.District), args.Error(1)
}

// setupRegionRouter creates a test router for region endpoints
func setupRegionRouter(controller IRegionController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	v1 := router.Group("/api/v1")
	regions := v1.Group("/regions")
	{
		regions.GET("/provinces", controller.GetAllProvinces)
		regions.GET("/provinces/:province/regencies", controller.GetRegenciesByProvince)
		regions.GET("/regencies/:regency/districts", controller.GetDistrictsByRegency)
	}

	return router
}

// ========== GetAllProvinces Tests ==========

func TestGetAllProvinces_Success(t *testing.T) {
	mockService := new(MockRegionService)
	controller := NewRegionController(mockService)
	router := setupRegionRouter(controller)

	provinces := []models.Province{
		{ID: 11, Name: "ACEH"},
		{ID: 31, Name: "DKI JAKARTA"},
		{ID: 36, Name: "BANTEN"},
	}

	mockService.On("GetAllProvinces", mock.Anything).Return(provinces, nil)

	req, _ := http.NewRequest("GET", "/api/v1/regions/provinces", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "Provinces fetched successfully", resp.Message)

	data, ok := resp.Data.([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 3, len(data))

	mockService.AssertExpectations(t)
}

func TestGetAllProvinces_ServiceError(t *testing.T) {
	mockService := new(MockRegionService)
	controller := NewRegionController(mockService)
	router := setupRegionRouter(controller)

	mockService.On("GetAllProvinces", mock.Anything).Return(nil, assert.AnError)

	req, _ := http.NewRequest("GET", "/api/v1/regions/provinces", nil)
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

func TestGetAllProvinces_EmptyResult(t *testing.T) {
	mockService := new(MockRegionService)
	controller := NewRegionController(mockService)
	router := setupRegionRouter(controller)

	provinces := []models.Province{}

	mockService.On("GetAllProvinces", mock.Anything).Return(provinces, nil)

	req, _ := http.NewRequest("GET", "/api/v1/regions/provinces", nil)
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

// ========== GetRegenciesByProvince Tests ==========

func TestGetRegenciesByProvince_Success(t *testing.T) {
	mockService := new(MockRegionService)
	controller := NewRegionController(mockService)
	router := setupRegionRouter(controller)

	regencies := []models.Regency{
		{ID: 1101, ProvinceID: 11, Name: "KABUPATEN SIMEULUE", Type: "kabupaten"},
		{ID: 1102, ProvinceID: 11, Name: "KABUPATEN ACEH SINGKIL", Type: "kabupaten"},
	}

	mockService.On("GetRegenciesByProvince", mock.Anything, "11").Return(regencies, nil)

	req, _ := http.NewRequest("GET", "/api/v1/regions/provinces/11/regencies", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "Regencies fetched successfully", resp.Message)

	data, ok := resp.Data.([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 2, len(data))

	mockService.AssertExpectations(t)
}

func TestGetRegenciesByProvince_ServiceError(t *testing.T) {
	mockService := new(MockRegionService)
	controller := NewRegionController(mockService)
	router := setupRegionRouter(controller)

	mockService.On("GetRegenciesByProvince", mock.Anything, "11").Return(nil, assert.AnError)

	req, _ := http.NewRequest("GET", "/api/v1/regions/provinces/11/regencies", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)

	mockService.AssertExpectations(t)
}

func TestGetRegenciesByProvince_EmptyResult(t *testing.T) {
	mockService := new(MockRegionService)
	controller := NewRegionController(mockService)
	router := setupRegionRouter(controller)

	regencies := []models.Regency{}

	mockService.On("GetRegenciesByProvince", mock.Anything, "99").Return(regencies, nil)

	req, _ := http.NewRequest("GET", "/api/v1/regions/provinces/99/regencies", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)

	mockService.AssertExpectations(t)
}

// ========== GetDistrictsByRegency Tests ==========

func TestGetDistrictsByRegency_Success(t *testing.T) {
	mockService := new(MockRegionService)
	controller := NewRegionController(mockService)
	router := setupRegionRouter(controller)

	districts := []models.District{
		{ID: 110101, Name: "TAMIN", RegencyID: 1101},
		{ID: 110102, Name: "SIMEULUE BARAT", RegencyID: 1101},
	}

	mockService.On("GetDistrictsByRegency", mock.Anything, "11", "1101").Return(districts, nil)

	req, _ := http.NewRequest("GET", "/api/v1/regions/regencies/1101/districts?province=11", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "Districts fetched successfully", resp.Message)

	data, ok := resp.Data.([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 2, len(data))

	mockService.AssertExpectations(t)
}

func TestGetDistrictsByRegency_MissingProvinceQueryParam(t *testing.T) {
	mockService := new(MockRegionService)
	controller := NewRegionController(mockService)
	router := setupRegionRouter(controller)

	req, _ := http.NewRequest("GET", "/api/v1/regions/regencies/1101/districts", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.NotNil(t, resp.Error)
	assert.Equal(t, "Province ID is required as query parameter", resp.Error.Message)

	mockService.AssertNotCalled(t, "GetDistrictsByRegency")
}

func TestGetDistrictsByRegency_ServiceError(t *testing.T) {
	mockService := new(MockRegionService)
	controller := NewRegionController(mockService)
	router := setupRegionRouter(controller)

	mockService.On("GetDistrictsByRegency", mock.Anything, "11", "1101").Return(nil, assert.AnError)

	req, _ := http.NewRequest("GET", "/api/v1/regions/regencies/1101/districts?province=11", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)

	mockService.AssertExpectations(t)
}

func TestGetDistrictsByRegency_EmptyResult(t *testing.T) {
	mockService := new(MockRegionService)
	controller := NewRegionController(mockService)
	router := setupRegionRouter(controller)

	districts := []models.District{}

	mockService.On("GetDistrictsByRegency", mock.Anything, "11", "9999").Return(districts, nil)

	req, _ := http.NewRequest("GET", "/api/v1/regions/regencies/9999/districts?province=11", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)

	mockService.AssertExpectations(t)
}

// ========== Edge Cases ==========

func TestGetRegenciesByProvince_EmptyProvinceID(t *testing.T) {
	mockService := new(MockRegionService)
	controller := NewRegionController(mockService)
	router := setupRegionRouter(controller)

	req, _ := http.NewRequest("GET", "/api/v1/regions/provinces//regencies", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetDistrictsByRegency_EmptyRegencyID(t *testing.T) {
	mockService := new(MockRegionService)
	controller := NewRegionController(mockService)
	router := setupRegionRouter(controller)

	req, _ := http.NewRequest("GET", "/api/v1/regions/regencies//districts?province=11", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}