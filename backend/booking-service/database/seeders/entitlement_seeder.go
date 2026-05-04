package seeders

import (
	"log"
	"time"

	"booking-service/models"

	"gorm.io/gorm"
)

// EntitlementSeeder seeds sample entitlements for testing
type EntitlementSeeder struct{}

// Name returns the seeder name
func (s *EntitlementSeeder) Name() string {
	return "entitlements"
}

// Run executes the entitlement seeder
func (s *EntitlementSeeder) Run(db *gorm.DB) error {
	log.Println("Running entitlement seeder...")

	// Check if entitlements already exist
	var count int64
	db.Model(&models.UserEntitlement{}).Count(&count)
	if count > 0 {
		log.Println("Entitlements already exist, skipping...")
		return nil
	}

	// Sample entitlements
	entitlements := []models.UserEntitlement{
		{
			UserID:            1,
			SourceType:        "package",
			SourceID:          "pkg-001",
			TotalSessions:     10,
			SessionsRemaining: 10,
			ExpiresAt:         time.Now().AddDate(1, 0, 0),
		},
		{
			UserID:            1,
			SourceType:        "voucher",
			SourceID:          "vch-001",
			TotalSessions:     5,
			SessionsRemaining: 5,
			ExpiresAt:         time.Now().AddDate(0, 6, 0),
		},
		{
			UserID:            2,
			SourceType:        "package",
			SourceID:          "pkg-002",
			TotalSessions:     20,
			SessionsRemaining: 20,
			ExpiresAt:         time.Now().AddDate(1, 0, 0),
		},
		{
			UserID:            3,
			SourceType:        "package",
			SourceID:          "pkg-003",
			TotalSessions:     15,
			SessionsRemaining: 15,
			ExpiresAt:         time.Now().AddDate(1, 0, 0),
		},
	}

	for _, ent := range entitlements {
		if err := db.Create(&ent).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d entitlements", len(entitlements))
	return nil
}
