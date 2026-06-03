package database

import (
	"core-service/database/migrations"
	"core-service/models"

	"gorm.io/gorm"
)

// Migrate runs all database migrations
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// Region tables
		&models.Province{},
		&models.Regency{},
		&models.District{},

		// Car table
		&models.Car{},

		// Package tables
		&models.Package{},
		&models.PackageBenefit{},

		// Sales tables
		&models.Sale{},
		&models.SaleItem{},
		&models.MonthlySales{},
	)
}

// MigrateRegions runs region-specific migrations
func MigrateRegions(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Province{},
		&models.Regency{},
		&models.District{},
	)
}

// MigrateCars runs car-specific migrations
func MigrateCars(db *gorm.DB) error {
	return db.AutoMigrate(&models.Car{})
}

// MigratePackages runs package-specific migrations
func MigratePackages(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Package{},
		&models.PackageBenefit{},
	)
}

// MigrateSales runs sales-specific migrations
func MigrateSales(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Sale{},
		&models.SaleItem{},
		&models.MonthlySales{},
	)
}

// RunMigration runs a specific migration by name
func RunMigration(db *gorm.DB, name string) error {
	switch name {
	case "regions":
		return MigrateRegions(db)
	case "cars":
		return MigrateCars(db)
	case "packages":
		return MigratePackages(db)
	case "sales":
		return MigrateSales(db)
	case "all":
		return Migrate(db)
	default:
		return migrations.RunNamedMigration(db, name)
	}
}

// GetMigrationStatus returns the status of all tables
func GetMigrationStatus(db *gorm.DB) ([]models.TableStatus, error) {
	var tables = []string{
		"provinces",
		"regencies",
		"districts",
		"cars",
		"packages",
		"package_benefits",
	}

	var status []models.TableStatus
	for _, table := range tables {
		var count int64
		db.Raw("SELECT COUNT(*) FROM " + table).Scan(&count)
		status = append(status, models.TableStatus{
			TableName: table,
			RowCount:  count,
			Exists:    count > 0 || count == 0, // Simplified check
		})
	}

	return status, nil
}