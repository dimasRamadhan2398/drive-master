package controllers

import (
	"testing"

	"core-service/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// CreateMockCar creates a mock car for testing
func CreateMockCar() *models.Car {
	return &models.Car{
		ID:           uuid.New(),
		Brand:        "Toyota",
		Model:        "Vios",
		Year:         2023,
		LicensePlate: "B 1234 XYZ",
		Color:        "Black",
		Transmission: models.TransmissionAutomatic,
		Status:       models.CarStatusAvailable,
		Mileage:      0,
		ImageURL:     "https://example.com/car.jpg",
		Notes:        "Test car",
	}
}

// CreateMockCarWithStatus creates a mock car with specific status
func CreateMockCarWithStatus(status models.CarStatus) *models.Car {
	car := CreateMockCar()
	car.Status = status
	return car
}

// CreateMockPackage creates a mock package for testing
func CreateMockPackage() *models.Package {
	return &models.Package{
		ID:              uuid.New(),
		Name:            "Basic Package",
		Description:     "A basic driving package",
		PackageType:     models.PackageTypeBronze,
		Price:           500000,
		DiscountPrice:   450000,
		DurationMinutes: 60,
		TotalSessions:   4,
		Status:          models.PackageStatusActive,
		ImageURL:        "https://example.com/package.jpg",
		Benefits: []models.PackageBenefit{
			{
				ID:          uuid.New(),
				Title:       "4 Sessions",
				Description: "4 driving sessions with instructor",
				Icon:        "car",
				SortOrder:   1,
			},
		},
	}
}

// CreateMockPackageWithType creates a mock package with specific type
func CreateMockPackageWithType(packageType models.PackageType) *models.Package {
	pkg := CreateMockPackage()
	pkg.PackageType = packageType
	return pkg
}

// CreateMockProvince creates a mock province for testing
func CreateMockProvince() *models.Province {
	return &models.Province{
		ID:   11,
		Name: "ACEH",
	}
}

// CreateMockRegency creates a mock regency for testing
func CreateMockRegency() *models.Regency {
	return &models.Regency{
		ID:         1101,
		ProvinceID: 11,
		Name:       "KABUPATEN SIMEULUE",
		Type:       "kabupaten",
	}
}

// CreateMockDistrict creates a mock district for testing
func CreateMockDistrict() *models.District {
	return &models.District{
		ID:        110101,
		Name:      "TAMIN",
		RegencyID: 1101,
	}
}

// AssertCarResponse asserts that the car data matches expected values
func AssertCarResponse(t *testing.T, car *models.Car, expectedBrand, expectedModel string) {
	assert.NotNil(t, car)
	assert.Equal(t, expectedBrand, car.Brand)
	assert.Equal(t, expectedModel, car.Model)
}

// AssertPackageResponse asserts that the package data matches expected values
func AssertPackageResponse(t *testing.T, pkg *models.Package, expectedName string) {
	assert.NotNil(t, pkg)
	assert.Equal(t, expectedName, pkg.Name)
}

// CreateValidCreateCarRequest returns a valid create car request JSON string
func CreateValidCreateCarRequest() string {
	return `{
		"brand": "Toyota",
		"model": "Vios",
		"year": 2023,
		"licensePlate": "B 1234 XYZ",
		"color": "Black",
		"transmission": "automatic",
		"imageUrl": "https://example.com/car.jpg",
		"notes": "Test car"
	}`
}

// CreateValidCreatePackageRequest returns a valid create package request JSON string
func CreateValidCreatePackageRequest() string {
	return `{
		"name": "Basic Package",
		"description": "A basic driving package",
		"packageType": "bronze",
		"price": 500000,
		"discountPrice": 450000,
		"durationMinutes": 60,
		"totalSessions": 4,
		"imageUrl": "https://example.com/package.jpg",
		"benefits": [
			{
				"title": "4 Sessions",
				"description": "4 driving sessions with instructor",
				"icon": "car",
				"sortOrder": 1
			}
		]
	}`
}