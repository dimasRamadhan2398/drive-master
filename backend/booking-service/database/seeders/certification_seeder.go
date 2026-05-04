package seeders

import (
	"log"
	"time"

	"booking-service/models"

	"gorm.io/gorm"
)

// CertificationSeeder seeds sample certifications for testing
type CertificationSeeder struct{}

// Name returns the seeder name
func (s *CertificationSeeder) Name() string {
	return "certifications"
}

// Run executes the certification seeder
func (s *CertificationSeeder) Run(db *gorm.DB) error {
	log.Println("Running certification seeder...")

	// Check if certifications already exist
	var count int64
	db.Model(&models.Certification{}).Count(&count)
	if count > 0 {
		log.Println("Certifications already exist, skipping...")
		return nil
	}

	// Sample certifications
	certifications := []models.Certification{
		{
			Type:      "driving-license-b",
			Recipient: "John Doe",
			IssueDate: time.Now().AddDate(0, -2, 0),
			PackageID: 1,
			Status:    models.CertificationStatusIssued,
		},
		{
			Type:      "defensive-driving",
			Recipient: "Jane Smith",
			IssueDate: time.Now().AddDate(0, -1, 0),
			PackageID: 2,
			Status:    models.CertificationStatusIssued,
		},
		{
			Type:      "driving-license-a",
			Recipient: "Bob Wilson",
			IssueDate: time.Now(),
			PackageID: 3,
			Status:    models.CertificationStatusPending,
		},
	}

	for _, cert := range certifications {
		if err := db.Create(&cert).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d certifications", len(certifications))
	return nil
}
