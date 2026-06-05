package seeders

import (
	"core-service/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RunGeneralSettingsSeeder seeds default general settings
func RunGeneralSettingsSeeder(db *gorm.DB) error {
	settings := models.GeneralSettings{
		ID:              uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		BusinessName:    "Drive",
		Email:           "info@drive.com",
		Phone:           "+62 21 1234 5678",
		Fax:             "+62 21 1234 5679",
		WhatsApp:        "+62 812 3456 7890",
		Address:         "Jl. Sudirman No. 123, Jakarta Selatan",
		HoursMonFri:     "08:00 - 17:00",
		HoursSatSun:     "09:00 - 15:00",
		HoursNightShift: "19:00 - 22:00",
		NotifyEmail:     true,
		NotifySMS:       false,
		NotifyWhatsApp:  false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Only create if not exists (singleton pattern)
	result := db.Where("id = ?", settings.ID).FirstOrCreate(&settings)
	if result.Error != nil {
		return result.Error
	}

	return nil
}