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
	// Custom member IDs for seeding
	customMemberIDs := []uuid.UUID{
		uuid.MustParse("433c4e17-83cb-4133-a4ab-7abaef8a3afe"),
		uuid.MustParse("86f5e950-cedb-495b-9528-48a7dffa6919"),
	}

	// Check if entitlements already exist for these members
	var count int64
	s.db.Model(&models.Entitlement{}).Where("member_id IN ?", customMemberIDs).Count(&count)
	if count > 0 {
		return nil // Already seeded for these members
	}

	// Entitlements for member 1: 433c4e17-83cb-4133-a4ab-7abaef8a3afe
	entitlements := []models.Entitlement{
		{
			ID:            uuid.New(),
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
			ID:            uuid.New(),
			MemberID:      uuid.MustParse("433c4e17-83cb-4133-a4ab-7abaef8a3afe"),
			PackageID:     uuid.MustParse("156e94ec-0b53-4a7f-b552-a5e4e7756726"),
			PackageName:   "Platinum Package",
			IsNightSession: true,
			TotalSessions: 12,
			Remaining:     0,
			UsedSessions:  12,
			StartDate:     time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
			Status:        models.EntitlementStatusUsed,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			ID:            uuid.New(),
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
	}

	// Entitlements for member 2: 86f5e950-cedb-495b-9528-48a7dffa6919
	entitlements = append(entitlements,
		models.Entitlement{
			ID:            uuid.New(),
			MemberID:      uuid.MustParse("86f5e950-cedb-495b-9528-48a7dffa6919"),
			PackageID:     uuid.MustParse("1ae7b1f4-2fbc-41e8-a473-cd1f9e900d5c"),
			PackageName:   "Silver Package",
			TotalSessions: 8,
			Remaining:     4,
			UsedSessions:  4,
			IsNightSession: true,
			IsWeekendSession: true,
			StartDate:     time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC),
			Status:        models.EntitlementStatusActive,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		models.Entitlement{
			ID:            uuid.New(),
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
	)

	for _, entitlement := range entitlements {
		if err := s.db.Create(&entitlement).Error; err != nil {
			return err
		}
	}

	return nil
}