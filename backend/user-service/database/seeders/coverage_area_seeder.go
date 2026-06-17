package seeders

import (
	"log"

	"user-service/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CoverageAreaSeeder seeds coverage areas for instructors
type CoverageAreaSeeder struct {
	db *gorm.DB
}

func NewCoverageAreaSeeder(db *gorm.DB) *CoverageAreaSeeder {
	return &CoverageAreaSeeder{db: db}
}

// Seed seeds coverage areas with AreaTypeDistrict and area ID 108
func (s *CoverageAreaSeeder) Seed() error {
	// Get instructor users
	var instructors []models.User
	if err := s.db.Joins("JOIN roles ON users.role_id = roles.id").
		Where("roles.name = ?", "instructor").
		Find(&instructors).Error; err != nil {
		return err
	}

	if len(instructors) == 0 {
		log.Println("No instructors found, skipping coverage area seeding")
		return nil
	}

	// Seed coverage area with AreaTypeDistrict and area ID 108 for each instructor
	for _, instructor := range instructors {
		record := models.InstructorArea{
			ID:           uuid.New(),
			InstructorID: instructor.ID,
			AreaType:     models.AreaTypeDistrict,
			AreaID:       108,
		}

		// Use FirstOrCreate to avoid duplicates
		if err := s.db.FirstOrCreate(&record, models.InstructorArea{
			InstructorID: instructor.ID,
			AreaType:     models.AreaTypeDistrict,
			AreaID:       108,
		}).Error; err != nil {
			return err
		}

		log.Printf("Seeded coverage area (district 108) for instructor %s", instructor.Username)
	}

	return nil
}
