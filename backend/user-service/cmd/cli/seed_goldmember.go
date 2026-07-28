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

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	EmailAddress string    `gorm:"column:email_address;uniqueIndex;not null"`
}

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
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "host=drive_master_postgres user=admin_drive password=drivemaster123 dbname=drivemaster_user_service port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	var user User
	result := db.Table("users").Where("email_address = ?", "goldmember@example.com").First(&user)
	if result.Error != nil {
		log.Fatalf("Failed to find user: %v", result.Error)
	}

	fmt.Printf("Found user: %s (ID: %s)\n", user.EmailAddress, user.ID)

	entitlement := Entitlement{
		ID:            uuid.New(),
		MemberID:      user.ID,
		BookingID:     uuid.New(),
		PackageID:     uuid.New(), 
		PackageName:   "10x Session Package",
		TotalSessions: 10,
		Remaining:     1,
		UsedSessions:  9,
		StartDate:     time.Now().AddDate(0, 0, -30),
		Status:        "active",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	var existingEntitlement Entitlement
	result = db.Where("member_id = ? AND status = 'active'", user.ID).First(&existingEntitlement)
	if result.Error == nil {
		existingEntitlement.TotalSessions = 10
		existingEntitlement.Remaining = 1
		existingEntitlement.UsedSessions = 9
        existingEntitlement.PackageName = "10x Session Package"
		existingEntitlement.UpdatedAt = time.Now()
		if err := db.Save(&existingEntitlement).Error; err != nil {
			log.Fatalf("Failed to update entitlement: %v", err)
		}
		fmt.Printf("Updated existing entitlement: %s\n", existingEntitlement.ID)
	} else {
		if err := db.Create(&entitlement).Error; err != nil {
			log.Fatalf("Failed to create entitlement: %v", err)
		}
		fmt.Printf("Created new entitlement: %s\n", entitlement.ID)
	}
}
