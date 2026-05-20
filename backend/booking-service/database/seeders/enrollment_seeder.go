package seeders

import (
	"log"
	"time"

	"booking-service/models"

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

	// Sample enrollments
	enrollments := []models.Enrollment{
		{
			UserID:     1,
			PackageID:  1,
			Status:     models.EnrollmentStatusPaid,
			TotalPrice: 2500000,
			PaidAt:     func() *time.Time { t := time.Now().AddDate(0, 0, -30); return &t }(),
			ExpiresAt:  time.Now().AddDate(1, 0, 0),
		},
		{
			UserID:     2,
			PackageID:  2,
			Status:     models.EnrollmentStatusPaid,
			TotalPrice: 4500000,
			PaidAt:     func() *time.Time { t := time.Now().AddDate(0, 0, -15); return &t }(),
			ExpiresAt:  time.Now().AddDate(1, 0, 0),
		},
		{
			UserID:     3,
			PackageID:  1,
			Status:     models.EnrollmentStatusPendingPayment,
			TotalPrice: 2500000,
			PaidAt:     nil,
			ExpiresAt:  time.Now().AddDate(1, 0, 0),
		},
		{
			UserID:     1,
			PackageID:  3,
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