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
	// Entitlements data
	entitlements := []models.Entitlement{
		{
			ID:            uuid.MustParse("d5a4a9c1-40e1-45bd-85b4-f63ee1893c5c"),
			MemberID:      uuid.MustParse("433c4e17-83cb-4133-a4ab-7abaef8a3afe"),
			PackageID:     uuid.MustParse("5b0f238f-72b5-44eb-b975-197dd10549e7"),
			PackageName:   "Gold Package",
			TotalSessions: 10,
			Remaining:     10,
			UsedSessions:  0,
			StartDate:     time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			Status:        models.EntitlementStatusActive,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			ID:             uuid.MustParse("3c3d52c7-0626-4d1a-821b-c743d528b122"),
			MemberID:       uuid.MustParse("433c4e17-83cb-4133-a4ab-7abaef8a3afe"),
			PackageID:      uuid.MustParse("156e94ec-0b53-4a7f-b552-a5e4e7756726"),
			PackageName:    "Platinum Package",
			IsNightSession: true,
			TotalSessions:  12,
			Remaining:      0,
			UsedSessions:   12,
			StartDate:      time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
			Status:         models.EntitlementStatusUsed,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:            uuid.MustParse("6f7f2b18-bbbe-4978-9588-4688cf8622f9"),
			MemberID:      uuid.MustParse("433c4e17-83cb-4133-a4ab-7abaef8a3afe"),
			PackageID:     uuid.MustParse("156e94ec-0b53-4a7f-b552-a5e4e7756726"),
			PackageName:   "Platinum Package",
			TotalSessions: 12,
			Remaining:     0,
			UsedSessions:  12,
			StartDate:     time.Date(2023, 3, 10, 0, 0, 0, 0, time.UTC),
			Status:        models.EntitlementStatusUsed,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			ID:               uuid.MustParse("2317bba4-1188-44fa-9566-ffccb17cc399"),
			MemberID:         uuid.MustParse("86f5e950-cedb-495b-9528-48a7dffa6919"),
			PackageID:        uuid.MustParse("1ae7b1f4-2fbc-41e8-a473-cd1f9e900d5c"),
			PackageName:      "Silver Package",
			TotalSessions:    8,
			Remaining:        4,
			UsedSessions:     4,
			IsNightSession:   true,
			IsWeekendSession: true,
			StartDate:        time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC),
			Status:           models.EntitlementStatusActive,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
		{
			ID:            uuid.MustParse("a81872cc-77aa-4f28-8688-662888cfbc1f"),
			MemberID:      uuid.MustParse("86f5e950-cedb-495b-9528-48a7dffa6919"),
			PackageID:     uuid.MustParse("49e284e0-193e-4707-b226-f7a4ed9f0fcd"),
			PackageName:   "Bronze Package",
			TotalSessions: 6,
			Remaining:     0,
			UsedSessions:  6,
			StartDate:     time.Date(2024, 11, 20, 0, 0, 0, 0, time.UTC),
			Status:        models.EntitlementStatusUsed,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			ID:            uuid.MustParse("88888888-8888-8888-8888-888888888888"),
			MemberID:      uuid.MustParse("99999999-9999-9999-9999-999999999999"),
			PackageID:     uuid.MustParse("11111111-1111-1111-1111-111111111301"),
			PackageName:   "10x Session",
			TotalSessions: 10,
			Remaining:     1,
			UsedSessions:  9,
			StartDate:     time.Now().AddDate(0, 0, -30),
			Status:        models.EntitlementStatusActive,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	for _, entitlement := range entitlements {
		var existing models.Entitlement
		if err := s.db.Where("id = ?", entitlement.ID).First(&existing).Error; err != nil {
			// Entitlement doesn't exist, create it
			if err := s.db.Create(&entitlement).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
