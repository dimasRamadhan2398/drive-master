package seeders

import (
	"user-service/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InstructorAreaSeeder seeds sample coverage areas for instructors
type InstructorAreaSeeder struct {
	db *gorm.DB
}

func NewInstructorAreaSeeder(db *gorm.DB) *InstructorAreaSeeder {
	return &InstructorAreaSeeder{db: db}
}

// Seed inserts sample instructor coverage areas
// Note: Area IDs reference external Area/Location service
func (s *InstructorAreaSeeder) Seed() error {
	// Get instructor users
	var instructors []models.User
	if err := s.db.Joins("JOIN roles ON users.role_id = roles.id").
		Where("roles.name = ?", "instructor").
		Find(&instructors).Error; err != nil {
		return err
	}

	if len(instructors) == 0 {
		return nil // No instructors found, skip
	}

	// Sample coverage areas (these would reference real area IDs from the Area service)
	// For demonstration, we'll use placeholder area IDs
	type instructorArea struct {
		instructorID uuid.UUID
		areaID      uint
	}

	areas := []instructorArea{
		{instructors[0].ID, 1}, // Jakarta Selatan
		{instructors[0].ID, 2}, // Jakarta Pusat
		{instructors[0].ID, 3}, // Jakarta Barat
		{instructors[0].ID, 4}, // Bandung
		{instructors[0].ID, 5}, // Bekasi
	}

	if len(instructors) > 1 {
		areas = append(areas,
			instructorArea{instructors[1].ID, 6},  // Surabaya
			instructorArea{instructors[1].ID, 7},   // Sidoarjo
			instructorArea{instructors[1].ID, 8},   // Gresik
		)
	}

	for _, area := range areas {
		record := models.InstructorArea{
			InstructorID: area.instructorID,
			AreaID:       area.areaID,
		}

		// Use FirstOrCreate to avoid duplicates
		if err := s.db.FirstOrCreate(&record, models.InstructorArea{
			InstructorID: area.instructorID,
			AreaID:       area.areaID,
		}).Error; err != nil {
			return err
		}
	}

	return nil
}
