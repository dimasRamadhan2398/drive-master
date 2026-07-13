package migrations

import "gorm.io/gorm"

func RunNamedMigration(db *gorm.DB, name string) error {
	// Add custom migrations here
	switch name {
	case "create_certifications":
		return createCertifications(db)
	default:
		return nil
	}
}

func createCertifications(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS certifications (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			description TEXT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		);
	`).Error
}