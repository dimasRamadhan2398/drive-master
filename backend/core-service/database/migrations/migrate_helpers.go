package migrations

import (
	"gorm.io/gorm"
)

// RunNamedMigration runs a specific named migration
func RunNamedMigration(db *gorm.DB, name string) error {
	// Add custom migrations here
	switch name {
	case "create_regions":
		return createRegions(db)
	case "create_cars":
		return createCars(db)
	case "create_packages":
		return createPackages(db)
	default:
		return nil
	}
}

func createRegions(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS provinces (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL
		);

		CREATE TABLE IF NOT EXISTS regencies (
			id SERIAL PRIMARY KEY,
			province_id INTEGER NOT NULL REFERENCES provinces(id),
			name VARCHAR(100) NOT NULL,
			type VARCHAR(100) NOT NULL
		);

		CREATE TABLE IF NOT EXISTS districts (
			id SERIAL PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(100) NOT NULL,
			regency_id INTEGER NOT NULL REFERENCES regencies(id)
		);
	`).Error
}

func createCars(db *gorm.DB) error {
	return db.Exec(`
		CREATE EXTENSION IF NOT EXISTS "pgcrypto";

		CREATE TABLE IF NOT EXISTS cars (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			brand VARCHAR(100) NOT NULL,
			model VARCHAR(100) NOT NULL,
			year INTEGER NOT NULL,
			license_plate VARCHAR(20) NOT NULL UNIQUE,
			color VARCHAR(50),
			transmission VARCHAR(20) NOT NULL DEFAULT 'manual',
			status VARCHAR(20) NOT NULL DEFAULT 'available',
			mileage INTEGER DEFAULT 0,
			image_url VARCHAR(500),
			notes TEXT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		);

		CREATE INDEX IF NOT EXISTS idx_cars_status ON cars(status);
		CREATE INDEX IF NOT EXISTS idx_cars_transmission ON cars(transmission);
	`).Error
}

func createPackages(db *gorm.DB) error {
	return db.Exec(`
		CREATE EXTENSION IF NOT EXISTS "pgcrypto";

		CREATE TABLE IF NOT EXISTS packages (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			description TEXT,
			package_type VARCHAR(50) NOT NULL,
			price DECIMAL(10,2) NOT NULL,
			discount_price DECIMAL(10,2) DEFAULT 0,
			duration_minutes INTEGER DEFAULT 60,
			total_sessions INTEGER DEFAULT 1,
			status VARCHAR(20) NOT NULL DEFAULT 'active',
			image_url VARCHAR(500),
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS package_benefits (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			package_id UUID NOT NULL REFERENCES packages(id) ON DELETE CASCADE,
			title VARCHAR(255) NOT NULL,
			description VARCHAR(500),
			icon VARCHAR(100),
			sort_order INTEGER DEFAULT 0,
			created_at TIMESTAMP
		);

		CREATE INDEX IF NOT EXISTS idx_packages_type ON packages(package_type);
		CREATE INDEX IF NOT EXISTS idx_packages_status ON packages(status);
		CREATE INDEX IF NOT EXISTS idx_package_benefits_package_id ON package_benefits(package_id);
	`).Error
}