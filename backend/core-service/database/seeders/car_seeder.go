package seeders

import (
	"core-service/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CarSeeder seeds sample car data
func RunCarSeeder(db *gorm.DB) error {
	cars := []models.Car{
		{
			ID:           uuid.New(),
			Brand:        "Toyota",
			Model:        "Vios",
			Year:         2023,
			LicensePlate: "B 1234 ABC",
			Color:        "White",
			Transmission: models.TransmissionAutomatic,
			Status:       models.CarStatusAvailable,
			Mileage:      5000,
			ImageURL:     "https://example.com/toyota-vios.jpg",
			Notes:        "Clean condition, AC working",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			Brand:        "Honda",
			Model:        "Civic",
			Year:         2022,
			LicensePlate: "B 5678 DEF",
			Color:        "Black",
			Transmission: models.TransmissionAutomatic,
			Status:       models.CarStatusAvailable,
			Mileage:      12000,
			ImageURL:     "https://example.com/honda-civic.jpg",
			Notes:        "Full option, leather seats",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			Brand:        "Suzuki",
			Model:        "Baleno",
			Year:         2023,
			LicensePlate: "D 9012 GHI",
			Color:        "Silver",
			Transmission: models.TransmissionManual,
			Status:       models.CarStatusAvailable,
			Mileage:      3000,
			ImageURL:     "https://example.com/suzuki-baleno.jpg",
			Notes:        "Fuel efficient, good for beginner",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			Brand:        "Toyota",
			Model:        "Yaris",
			Year:         2021,
			LicensePlate: "F 3456 JKL",
			Color:        "Red",
			Transmission: models.TransmissionAutomatic,
			Status:       models.CarStatusMaintenance,
			Mileage:      25000,
			ImageURL:     "https://example.com/toyota-yaris.jpg",
			Notes:        "Under maintenance - scheduled service",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			Brand:        "Honda",
			Model:        "City",
			Year:         2024,
			LicensePlate: "B 7890 MNO",
			Color:        "Gray",
			Transmission: models.TransmissionAutomatic,
			Status:       models.CarStatusInUse,
			Mileage:      1000,
			ImageURL:     "https://example.com/honda-city.jpg",
			Notes:        "Brand new, currently in use for class",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			Brand:        "Mitsubishi",
			Model:        "Mirage",
			Year:         2022,
			LicensePlate: "L 1234 PQR",
			Color:        "Blue",
			Transmission: models.TransmissionManual,
			Status:       models.CarStatusAvailable,
			Mileage:      15000,
			ImageURL:     "https://example.com/mitsubishi-mirage.jpg",
			Notes:        "Compact, easy to park",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			Brand:        "Toyota",
			Model:        "Avanza",
			Year:         2023,
			LicensePlate: "B 5678 STU",
			Color:        "White",
			Transmission: models.TransmissionManual,
			Status:       models.CarStatusAvailable,
			Mileage:      8000,
			ImageURL:     "https://example.com/toyota-avanza.jpg",
			Notes:        "Family car, 7 seats",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			Brand:        "Honda",
			Model:        "HR-V",
			Year:         2024,
			LicensePlate: "D 9012 VWX",
			Color:        "Black",
			Transmission: models.TransmissionAutomatic,
			Status:       models.CarStatusAvailable,
			Mileage:      500,
			ImageURL:     "https://example.com/honda-hr-v.jpg",
			Notes:        "SUV, great for long distance",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	for _, car := range cars {
		result := db.Where("license_plate = ?", car.LicensePlate).FirstOrCreate(&car)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}