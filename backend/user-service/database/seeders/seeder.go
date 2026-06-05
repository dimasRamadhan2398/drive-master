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

	// 5. Seed coverage areas (district 108)
	coverageAreaSeeder := NewCoverageAreaSeeder(r.db)
	if err := coverageAreaSeeder.Seed(); err != nil {
		return err
	}
	log.Println("Coverage areas seeded successfully")

	// 6. Seed entitlements for members
	entitlementSeeder := NewEntitlementSeeder(r.db)
	if err := entitlementSeeder.Seed(); err != nil {
		return err
	}
	log.Println("Entitlements seeded successfully")

	// 8. Seed testimonials
	testimonialSeeder := NewTestimonialSeeder(r.db)
	if err := testimonialSeeder.Seed(); err != nil {
		return err
	}
	log.Println("Testimonials seeded successfully")

	log.Println("All seeders completed successfully")
	return nil
}
