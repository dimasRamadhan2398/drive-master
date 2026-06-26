package seeders

import (
	"log"
	"time"

	"booking-service/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EnrollmentSeeder seeds sample enrollments for testing
type EnrollmentSeeder struct{}

// Name returns the seeder name
func (s *EnrollmentSeeder) Name() string {
	return "enrollments"
}

// Run executes the enrollment seeder
func (s *EnrollmentSeeder) Run(db *gorm.DB) error {
	log.Println("Running enrollment seeder...")

	// Check if enrollments already exist
	var count int64
	db.Model(&models.Enrollment{}).Count(&count)
	if count > 0 {
		log.Println("Enrollments already exist, skipping...")
		return nil
	}

	// Sample enrollments - using UUIDs for PackageID
	enrollments := []models.Enrollment{
		{
			UserID:     uuid.MustParse("0ccae045-9154-41c6-957c-c42f812b809f"),
			PackageID:  uuid.MustParse("11111111-1111-1111-1111-111111111111"), // Package 1 UUID
			Status:     models.EnrollmentStatusPaid,
			TotalPrice: 2500000,
			PaidAt:     func() *time.Time { t := time.Now().AddDate(0, 0, -30); return &t }(),
			ExpiresAt:  time.Now().AddDate(1, 0, 0),
		},
		{
			UserID:     uuid.MustParse("898170a5-08db-4467-b33e-049660a4231c"),
			PackageID:  uuid.MustParse("22222222-2222-2222-2222-222222222222"), // Package 2 UUID
			Status:     models.EnrollmentStatusPaid,
			TotalPrice: 4500000,
			PaidAt:     func() *time.Time { t := time.Now().AddDate(0, 0, -15); return &t }(),
			ExpiresAt:  time.Now().AddDate(1, 0, 0),
		},
		{
			UserID:     uuid.MustParse("9eac3536-bdc4-41b1-aaea-f9acbd231e20"),
			PackageID:  uuid.MustParse("11111111-1111-1111-1111-111111111111"), // Package 1 UUID
			Status:     models.EnrollmentStatusPendingPayment,
			TotalPrice: 2500000,
			PaidAt:     nil,
			ExpiresAt:  time.Now().AddDate(1, 0, 0),
		},
		{
			UserID:     uuid.MustParse("0ccae045-9154-41c6-957c-c42f812b809f"),
			PackageID:  uuid.MustParse("33333333-3333-3333-3333-333333333333"), // Package 3 UUID
			Status:     models.EnrollmentStatusInProgress,
			TotalPrice: 6500000,
			PaidAt:     func() *time.Time { t := time.Now().AddDate(0, 0, -60); return &t }(),
			ExpiresAt:  time.Now().AddDate(0, 6, 0),
		},
	}

	for _, enrollment := range enrollments {
		if err := db.Create(&enrollment).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d enrollments", len(enrollments))
	return nil
}