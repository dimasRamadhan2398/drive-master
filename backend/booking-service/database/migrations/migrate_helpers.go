package migrations

import "gorm.io/gorm"

func RunNamedMigration(db *gorm.DB, name string) error {
	// Add custom migrations here
	switch name {
	case "create_certifications":
		return createCertifications(db)
	case "create_entitlements":
		return createEntitlements(db)
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

func createEntitlements(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS entitlements (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL,
			source_type VARCHAR(255) NOT NULL,
			source_id VARCHAR(255) NOT NULL,
			total_sessions INT NOT NULL,
			sessions_remaining INT NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		);
	`).Error
}