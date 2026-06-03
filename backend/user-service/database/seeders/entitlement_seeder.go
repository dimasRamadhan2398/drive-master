package seeders

import (
	"time"

	"user-service/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EntitlementSeeder struct {
	db *gorm.DB
}

func NewEntitlementSeeder(db *gorm.DB) *EntitlementSeeder {
	return &EntitlementSeeder{db: db}
}

func (s *EntitlementSeeder) Seed() error {
	// Check if entitlements already exist
	var count int64
	s.db.Model(&models.Entitlement{}).Count(&count)
	if count > 0 {
		return nil // Already seeded
	}

	// Get member user IDs from the database
	var members []models.User
	if err := s.db.Preload("Role").Where("role_id = ?", 2).Find(&members).Error; err != nil {
		return err
	}

	// Create sample entitlements for members
	for i, member := range members {
		// Create 1-3 entitlements per member
		numEntitlements := (i % 3) + 1

		for j := 0; j < numEntitlements; j++ {
			totalSessions := 10 + (j * 5) // 10, 15, or 20 sessions
			usedSessions := j * 3          // Some used, some not

			entitlement := &models.Entitlement{
				ID:              uuid.New(),
				MemberID:        member.ID,
				PackageID:       uuid.New(),
				PackageName:     "Package " + string(rune('A'+j)),
				TotalSessions:   totalSessions,
				Remaining:       totalSessions - usedSessions,
				UsedSessions:    usedSessions,
				StartDate:       time.Now().AddDate(0, -1, 0), // Started 1 month ago
				Status:          models.EntitlementStatusActive,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			}

			if err := s.db.Create(entitlement).Error; err != nil {
				return err
			}
		}
	}

	return nil
}