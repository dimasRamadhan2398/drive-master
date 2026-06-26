package seeders

import (
	"core-service/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RunAddonSeeder seeds sample addon data
func RunAddonSeeder(db *gorm.DB) error {
	addons := []models.AddOn{
		{
			ID:          uuid.MustParse("22222222-2222-2222-2222-222222222201"),
			Title:       "Extra Session",
			Description: "Additional driving session to enhance your skills. Each purchase adds 1 extra session to your package.",
			Price:       350000,
			Sessions:    1,
			Status:      models.AddOnStatusActive,
			ImageURL:    "",
			SortOrder:   1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.MustParse("22222222-2222-2222-2222-222222222202"),
			Title:       "Night Session Add-on",
			Description: "Add night driving session to your package. Learn to drive safely in low-light conditions.",
			Price:       250000,
			Sessions:    1,
			Status:      models.AddOnStatusActive,
			ImageURL:    "",
			SortOrder:   2,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.MustParse("22222222-2222-2222-2222-222222222203"),
			Title:       "Weekend Session Add-on",
			Description: "Add weekend driving session to your package. Perfect for those with busy weekday schedules.",
			Price:       200000,
			Sessions:    1,
			Status:      models.AddOnStatusActive,
			ImageURL:    "",
			SortOrder:   3,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.MustParse("22222222-2222-2222-2222-222222222204"),
			Title:       "Extended Session (2 Hours)",
			Description: "Upgrade your regular session to 2 hours for more practice time.",
			Price:       500000,
			Sessions:    1,
			Status:      models.AddOnStatusActive,
			ImageURL:    "",
			SortOrder:   4,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Create or update all addons
	for _, a := range addons {
		result := db.Where("id = ?", a.ID).FirstOrCreate(&a)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}
