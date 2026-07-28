package seeders

import (
	"gorm.io/gorm"
)

// RunNamedSeeder runs a specific named seeder
func RunNamedSeeder(db *gorm.DB, name string) error {
	// Add custom seeders here
	switch name {
	case "sample":
		return runSampleSeeder(db)
	default:
		return nil
	}
}

func runSampleSeeder(db *gorm.DB) error {
	// Add sample data seeder here
	return nil
}