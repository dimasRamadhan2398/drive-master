package seeders

import (
	"core-service/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RunFAQSeeder seeds initial FAQ questions
func RunFAQSeeder(db *gorm.DB) error {
	faqs := []models.FAQ{
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Question:  "What is an EV?",
			Answer:    "An Electric Vehicle (EV) is a vehicle that uses one or more electric motors for propulsion.",
			Order:     1,
			Category:  "general",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			Question:  "How long does it take to charge?",
			Answer:    "Charging time depends on the charger and the vehicle. A typical home charger can fully charge an EV overnight.",
			Order:     2,
			Category:  "general",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, faq := range faqs {
		result := db.Where("id = ?", faq.ID).FirstOrCreate(&faq)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}