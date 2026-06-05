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

// MockCarService mocks the car service for testing
type MockCarService struct {
	mock.Mock
}

func (m *MockCarService) CreateCar(ctx context.Context, car *models.Car) error {
	args := m.Called(ctx, car)
	return args.Error(0)
}

func (m *MockCarService) GetCarByID(ctx context.Context, id uuid.UUID) (*models.Car, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Car), args.Error(1)
}

func (m *MockCarService) GetAllCars(ctx context.Context) ([]models.Car, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Car), args.Error(1)
}

func (m *MockCarService) GetCarsByStatus(ctx context.Context, status models.CarStatus) ([]models.Car, error) {
	args := m.Called(ctx, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Car), args.Error(1)
}

func (m *MockCarService) GetCarsByTransmission(ctx context.Context, transmission models.TransmissionType) ([]models.Car, error) {
	args := m.Called(ctx, transmission)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Car), args.Error(1)
}

func (m *MockCarService) UpdateCar(ctx context.Context, car *models.Car) error {
	args := m.Called(ctx, car)
	return args.Error(0)
}

func (m *MockCarService) DeleteCar(ctx context.Context, car *models.Car) error {
	args := m.Called(ctx, car)
	return args.Error(0)
}

func (m *MockCarService) CountCars(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// setupCarRouter creates a test router for car endpoints
func setupCarRouter(controller ICarController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/cars", controller.GetAllCars)
		v1.GET("/cars/:id", controller.GetCarByID)
		v1.POST("/cars", controller.CreateCar)
		v1.PUT("/cars/:id", controller.UpdateCar)
		v1.DELETE("/cars/:id", controller.DeleteCar)
	}

	return router
}

// ========== GetAllCars Tests ==========

func TestGetAllCars_Success(t *testing.T) {
	mockService := new(MockCarService)
	controller := NewCarController(mockService)
	router := setupCarRouter(controller)

	cars := []models.Car{
		{
			ID:           uuid.New(),
			Brand:        "Toyota",
			Model:        "Vios",
			Year:         2023,
			LicensePlate: "B 1234 XYZ",
			Color:        "Black",
			Transmission: models.TransmissionAutomatic,
			Status:       models.CarStatusAvailable,
		},
		{
			ID:           uuid.New(),
			Brand:        "Honda",
			Model:        "Civic",
			Year:         2022,
			LicensePlate: "B 5678 ABC",
			Color:        "White",
			Transmission: models.TransmissionManual,
			Status:       models.CarStatusInUse,
		},
	}

	mockService.On("GetAllCars", mock.Anything).Return(cars, nil)

	req, _ := http.NewRequest("GET", "/api/v1/cars", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "Cars fetched successfully", resp.Message)

	data, ok := resp.Data.([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 2, len(data))

	mockService.AssertExpectations(t)
}

func TestGetAllCars_ServiceError(t *testing.T) {
	mockService := new(MockCarService)
	controller := NewCarController(mockService)
	router := setupCarRouter(controller)

	mockService.On("GetAllCars", mock.Anything).Return(nil, assert.AnError)

	req, _ := http.NewRequest("GET", "/api/v1/cars", nil)
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

func TestGetAllCars_EmptyResult(t *testing.T) {
	mockService := new(MockCarService)
	controller := NewCarController(mockService)
	router := setupCarRouter(controller)

	cars := []models.Car{}

	mockService.On("GetAllCars", mock.Anything).Return(cars, nil)

	req, _ := http.NewRequest("GET", "/api/v1/cars", nil)
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

// ========== GetCarByID Tests ==========

func TestGetCarByID_Success(t *testing.T) {
	mockService := new(MockCarService)
	controller := NewCarController(mockService)
	router := setupCarRouter(controller)

	carID := uuid.New()
	car := &models.Car{
		ID:           carID,
		Brand:        "Toyota",
		Model:        "Vios",
		Year:         2023,
		LicensePlate: "B 1234 XYZ",
		Color:        "Black",
		Transmission: models.TransmissionAutomatic,
		Status:       models.CarStatusAvailable,
	}

	mockService.On("GetCarByID", mock.Anything, carID).Return(car, nil)

	req, _ := http.NewRequest("GET", "/api/v1/cars/"+carID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "Car fetched successfully", resp.Message)
	assert.NotNil(t, resp.Data)

	mockService.AssertExpectations(t)
}

func TestGetCarByID_NotFound(t *testing.T) {
	mockService := new(MockCarService)
	controller := NewCarController(mockService)
	router := setupCarRouter(controller)

	carID := uuid.New()

	mockService.On("GetCarByID", mock.Anything, carID).Return(nil, nil)

	req, _ := http.NewRequest("GET", "/api/v1/cars/"+carID.String(), nil)
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

func TestGetCarByID_InvalidIDFormat(t *testing.T) {
	mockService := new(MockCarService)
	controller := NewCarController(mockService)
	router := setupCarRouter(controller)

	req, _ := http.NewRequest("GET", "/api/v1/cars/invalid-uuid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.NotNil(t, resp.Error)
	assert.Equal(t, "Invalid car ID format", resp.Error.Message)

	mockService.AssertNotCalled(t, "GetCarByID")
}

func TestGetCarByID_ServiceError(t *testing.T) {
	mockService := new(MockCarService)
	controller := NewCarController(mockService)
	router := setupCarRouter(controller)

	carID := uuid.New()

	mockService.On("GetCarByID", mock.Anything, carID).Return(nil, assert.AnError)

	req, _ := http.NewRequest("GET", "/api/v1/cars/"+carID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)

	mockService.AssertExpectations(t)
}

// ========== CreateCar Tests ==========

func TestCreateCar_Success(t *testing.T) {
	mockService := new(MockCarService)
	controller := NewCarController(mockService)
	router := setupCarRouter(controller)

	mockService.On("CreateCar", mock.Anything, mock.AnythingOfType("*models.Car")).Return(nil)

	createCarJSON := `{
		"brand": "Toyota",
		"model": "Vios",
		"year": 2023,
		"licensePlate": "B 1234 XYZ",
		"color": "Black",
		"transmission": "automatic"
	}`

	req, _ := http.NewRequest("POST", "/api/v1/cars", bytes.NewBufferString(createCarJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "Car created successfully", resp.Message)

	mockService.AssertExpectations(t)
}

func TestCreateCar_InvalidJSON(t *testing.T) {
	mockService := new(MockCarService)
	controller := NewCarController(mockService)
	router := setupCarRouter(controller)

	req, _ := http.NewRequest("POST", "/api/v1/cars", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)

	mockService.AssertNotCalled(t, "CreateCar")
}

func TestCreateCar_MissingRequiredFields(t *testing.T) {
	mockService := new(MockCarService)
	controller := NewCarController(mockService)
	router := setupCarRouter(controller)

	// Missing brand, model, year, licensePlate
	createCarJSON := `{
		"color": "Black"
	}`

	req, _ := http.NewRequest("POST", "/api/v1/cars", bytes.NewBufferString(createCarJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)

	mockService.AssertNotCalled(t, "CreateCar")
}

func TestCreateCar_ServiceError(t *testing.T) {
	mockService := new(MockCarService)
	controller := NewCarController(mockService)
	router := setupCarRouter(controller)

	mockService.On("CreateCar", mock.Anything, mock.AnythingOfType("*models.Car")).Return(assert.AnError)

	createCarJSON := `{
		"brand": "Toyota",
		"model": "Vios",
		"year": 2023,
		"licensePlate": "B 1234 XYZ"
	}`

	req, _ := http.NewRequest("POST", "/api/v1/cars", bytes.NewBufferString(createCarJSON))
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

// ========== UpdateCar Tests ==========

func TestUpdateCar_Success(t *testing.T) {
	mockService := new(MockCarService)
	controller := NewCarController(mockService)
	router := setupCarRouter(controller)

	carID := uuid.New()
	existingCar := &models.Car{
		ID:           carID,
		Brand:        "Toyota",
		Model:        "Vios",
		Year:         2023,
		LicensePlate: "B 1234 XYZ",
		Color:        "Black",
		Transmission: models.TransmissionAutomatic,
		Status:       models.CarStatusAvailable,
	}

	mockService.On("GetCarByID", mock.Anything, carID).Return(existingCar, nil).Once()
	mockService.On("UpdateCar", mock.Anything, mock.AnythingOfType("*models.Car")).Return(nil)

	updateCarJSON := `{
		"color": "White"
	}`

	req, _ := http.NewRequest("PUT", "/api/v1/cars/"+carID.String(), bytes.NewBufferString(updateCarJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "Car updated successfully", resp.Message)

	mockService.AssertExpectations(t)
}

func TestUpdateCar_NotFound(t *testing.T) {
	mockService := new(MockCarService)
	controller := NewCarController(mockService)
	router := setupCarRouter(controller)

	carID := uuid.New()

	mockService.On("GetCarByID", mock.Anything, carID).Return(nil, nil)

	updateCarJSON := `{
		"color": "White"
	}`

	req, _ := http.NewRequest("PUT", "/api/v1/cars/"+carID.String(), bytes.NewBufferString(updateCarJSON))
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

func TestUpdateCar_InvalidIDFormat(t *testing.T) {
	mockService := new(MockCarService)
	controller := NewCarController(mockService)
	router := setupCarRouter(controller)

	updateCarJSON := `{
		"color": "White"
	}`

	req, _ := http.NewRequest("PUT", "/api/v1/cars/invalid-uuid", bytes.NewBufferString(updateCarJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "Invalid car ID format", resp.Error.Message)

	mockService.AssertNotCalled(t, "GetCarByID")
}

func TestUpdateCar_InvalidJSON(t *testing.T) {
	mockService := new(MockCarService)
	controller := NewCarController(mockService)
	router := setupCarRouter(controller)

	carID := uuid.New()
	existingCar := &models.Car{
		ID:           carID,
		Brand:        "Toyota",
		Model:        "Vios",
		Year:         2023,
		LicensePlate: "B 1234 XYZ",
		Color:        "Black",
		Transmission: models.TransmissionAutomatic,
		Status:       models.CarStatusAvailable,
	}

	mockService.On("GetCarByID", mock.Anything, carID).Return(existingCar, nil)

	req, _ := http.NewRequest("PUT", "/api/v1/cars/"+carID.String(), bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)

	// GetCarByID was called but UpdateCar was not due to invalid JSON
	mockService.AssertCalled(t, "GetCarByID", mock.Anything, carID)
}

func TestUpdateCar_ServiceError(t *testing.T) {
	mockService := new(MockCarService)
	controller := NewCarController(mockService)
	router := setupCarRouter(controller)

	carID := uuid.New()
	existingCar := &models.Car{
		ID:           carID,
		Brand:        "Toyota",
		Model:        "Vios",
		Year:         2023,
		LicensePlate: "B 1234 XYZ",
		Color:        "Black",
		Transmission: models.TransmissionAutomatic,
		Status:       models.CarStatusAvailable,
	}

	mockService.On("GetCarByID", mock.Anything, carID).Return(existingCar, nil)
	mockService.On("UpdateCar", mock.Anything, mock.AnythingOfType("*models.Car")).Return(assert.AnError)

	updateCarJSON := `{
		"color": "White"
	}`

	req, _ := http.NewRequest("PUT", "/api/v1/cars/"+carID.String(), bytes.NewBufferString(updateCarJSON))
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

// ========== DeleteCar Tests ==========

func TestDeleteCar_Success(t *testing.T) {
	mockService := new(MockCarService)
	controller := NewCarController(mockService)
	router := setupCarRouter(controller)

	carID := uuid.New()
	existingCar := &models.Car{
		ID:           carID,
		Brand:        "Toyota",
		Model:        "Vios",
		Year:         2023,
		LicensePlate: "B 1234 XYZ",
		Color:        "Black",
		Transmission: models.TransmissionAutomatic,
		Status:       models.CarStatusAvailable,
	}

	mockService.On("GetCarByID", mock.Anything, carID).Return(existingCar, nil)
	mockService.On("DeleteCar", mock.Anything, existingCar).Return(nil)

	req, _ := http.NewRequest("DELETE", "/api/v1/cars/"+carID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "Car deleted successfully", resp.Message)

	mockService.AssertExpectations(t)
}

func TestDeleteCar_NotFound(t *testing.T) {
	mockService := new(MockCarService)
	controller := NewCarController(mockService)
	router := setupCarRouter(controller)

	carID := uuid.New()

	mockService.On("GetCarByID", mock.Anything, carID).Return(nil, nil)

	req, _ := http.NewRequest("DELETE", "/api/v1/cars/"+carID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)

	mockService.AssertExpectations(t)
}

func TestDeleteCar_InvalidIDFormat(t *testing.T) {
	mockService := new(MockCarService)
	controller := NewCarController(mockService)
	router := setupCarRouter(controller)

	req, _ := http.NewRequest("DELETE", "/api/v1/cars/invalid-uuid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "Invalid car ID format", resp.Error.Message)

	mockService.AssertNotCalled(t, "GetCarByID")
}

func TestDeleteCar_ServiceError(t *testing.T) {
	mockService := new(MockCarService)
	controller := NewCarController(mockService)
	router := setupCarRouter(controller)

	carID := uuid.New()
	existingCar := &models.Car{
		ID:           carID,
		Brand:        "Toyota",
		Model:        "Vios",
		Year:         2023,
		LicensePlate: "B 1234 XYZ",
		Color:        "Black",
		Transmission: models.TransmissionAutomatic,
		Status:       models.CarStatusAvailable,
	}

	mockService.On("GetCarByID", mock.Anything, carID).Return(existingCar, nil)
	mockService.On("DeleteCar", mock.Anything, existingCar).Return(assert.AnError)

	req, _ := http.NewRequest("DELETE", "/api/v1/cars/"+carID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)

	mockService.AssertExpectations(t)
}