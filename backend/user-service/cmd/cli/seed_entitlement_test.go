//go:build ignore

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Entitlement struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	MemberID      uuid.UUID `gorm:"type:uuid;not null;index"`
	BookingID     uuid.UUID `gorm:"type:uuid;index"`
	PackageID     uuid.UUID `gorm:"type:uuid"`
	PackageName   string    `gorm:"size:255"`
	TotalSessions int       `gorm:"default:0"`
	Remaining     int       `gorm:"default:0"`
	UsedSessions  int       `gorm:"default:0"`
	StartDate     time.Time
	Status        string `gorm:"type:varchar(20);default:'active'"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func main() {
	// Get database URL from environment or use default
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "host=localhost user=admin_drive password=drivemaster123 dbname=drivemaster_user_service port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// User ID provided by user
	memberID := uuid.MustParse("0ccae045-9154-41c6-957c-c42f812b809f")
	packageID := uuid.MustParse("5b0f238f-72b5-44eb-b975-197dd10549e7") // Gold Package

	// Create entitlement with 9 sessions completed out of 10
	entitlement := Entitlement{
		ID:            uuid.New(),
		MemberID:      memberID,
		BookingID:     uuid.New(),
		PackageID:     packageID,
		PackageName:   "Gold Package",
		TotalSessions: 10,
		Remaining:     1,
		UsedSessions:  9,
		StartDate:     time.Now().AddDate(0, 0, -30),
		Status:        "active",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Check if user already has an active entitlement for Gold Package
	var existingEntitlement Entitlement
	result := db.Where("member_id = ? AND package_name = ? AND status = 'active'", memberID, "Gold Package").First(&existingEntitlement)
	if result.Error == nil {
		// Update existing entitlement
		existingEntitlement.TotalSessions = 10
		existingEntitlement.Remaining = 1
		existingEntitlement.UsedSessions = 9
		existingEntitlement.UpdatedAt = time.Now()
		if err := db.Save(&existingEntitlement).Error; err != nil {
			log.Fatalf("Failed to update entitlement: %v", err)
		}
		fmt.Printf("Updated existing entitlement: %s\n", existingEntitlement.ID)
	} else {
		// Create new entitlement
		if err := db.Create(&entitlement).Error; err != nil {
			log.Fatalf("Failed to create entitlement: %v", err)
		}
		fmt.Printf("Created new entitlement: %s\n", entitlement.ID)
	}

	// Verify the entitlement
	var entitlements []Entitlement
	db.Where("member_id = ?", memberID).Find(&entitlements)

	fmt.Println("\nEntitlements for user 0ccae045-9154-41c6-957c-c42f812b809f:")
	fmt.Println("ID | Package | Total | Remaining | Used | Status")
	fmt.Println("---|---------|-------|-----------|------|--------")
	for _, e := range entitlements {
		fmt.Printf("%s | %s | %d | %d | %d | %s\n",
			e.ID.String()[:8], e.PackageName, e.TotalSessions, e.Remaining, e.UsedSessions, e.Status)
	}
}
