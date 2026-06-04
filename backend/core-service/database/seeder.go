package database

import (
	"core-service/database/seeders"

	"gorm.io/gorm"
)

// RunSeeders runs all database seeders
func RunSeeders(db *gorm.DB) error {
	// Run region seeders (from existing seeders)
	if err := seeders.RunProvinceSeeder(db); err != nil {
		return err
	}

	if err := seeders.RunRegencySeeder(db); err != nil {
		return err
	}

	if err := seeders.RunDistrictsSeeder(db); err != nil {
		return err
	}

	// Run car seeder
	if err := seeders.RunCarSeeder(db); err != nil {
		return err
	}

	// Run package seeder
	if err := seeders.RunPackageSeeder(db); err != nil {
		return err
	}

	// Run article seeder
	if err := seeders.RunArticleSeeder(db); err != nil {
		return err
	}

	return nil
}

// RunSeederByName runs a specific seeder by name
func RunSeederByName(db *gorm.DB, name string) error {
	switch name {
	case "provinces":
		return seeders.RunProvinceSeeder(db)
	case "regencies":
		return seeders.RunRegencySeeder(db)
	case "districts":
		return seeders.RunDistrictsSeeder(db)
	case "cars":
		return seeders.RunCarSeeder(db)
	case "packages":
		return seeders.RunPackageSeeder(db)
	case "articles":
		return seeders.RunArticleSeeder(db)
	case "all":
		return RunSeeders(db)
	default:
		return seeders.RunNamedSeeder(db, name)
	}
}