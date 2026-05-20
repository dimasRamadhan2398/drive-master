package seeders

import (
	"log"

	"gorm.io/gorm"
)

// SeederRunner runs all seeders in the correct order
type SeederRunner struct {
	db *gorm.DB
}

// NewSeederRunner creates a new seeder runner
func NewSeederRunner(db *gorm.DB) *SeederRunner {
	return &SeederRunner{db: db}
}

// RunAll runs all seeders in order
func (r *SeederRunner) RunAll() error {
	log.Println("Running seeders...")

	// 1. Seed roles first (other seeders depend on role IDs)
	roleSeeder := NewRoleSeeder(r.db)
	if err := roleSeeder.Seed(); err != nil {
		return err
	}
	log.Println("Roles seeded successfully")

	// Get role map for other seeders
	roleMap, err := roleSeeder.GetRoles()
	if err != nil {
		return err
	}

	// 2. Seed users
	userSeeder := NewUserSeeder(r.db)
	if err := userSeeder.Seed(roleMap); err != nil {
		return err
	}
	log.Println("Users seeded successfully")

	// 3. Seed work experiences
	workExpSeeder := NewWorkExperienceSeeder(r.db)
	if err := workExpSeeder.Seed(); err != nil {
		return err
	}
	log.Println("Work experiences seeded successfully")

	// 4. Seed instructor areas
	areaSeeder := NewInstructorAreaSeeder(r.db)
	if err := areaSeeder.Seed(); err != nil {
		return err
	}
	log.Println("Instructor areas seeded successfully")

	log.Println("All seeders completed successfully")
	return nil
}
