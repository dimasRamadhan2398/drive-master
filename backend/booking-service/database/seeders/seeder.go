package seeders

import (
	"gorm.io/gorm"
)

// Seeder interface that all seeders must implement
type Seeder interface {
	Run(db *gorm.DB) error
	Name() string
}

// SeederRunner manages and runs all seeders
type SeederRunner struct {
	db       *gorm.DB
	seeders  []Seeder
}

// NewSeederRunner creates a new SeederRunner with all available seeders
func NewSeederRunner(db *gorm.DB) *SeederRunner {
	return &SeederRunner{
		db: db,
		seeders: []Seeder{
			&EntitlementSeeder{},
			&CertificationSeeder{},
		},
	}
}

// RunAll runs all seeders in order
func (sr *SeederRunner) RunAll() error {
	for _, seeder := range sr.seeders {
		if err := seeder.Run(sr.db); err != nil {
			return err
		}
	}
	return nil
}

// RunSpecific runs only the specified seeders by name
func (sr *SeederRunner) RunSpecific(names []string) error {
	nameMap := make(map[string]bool)
	for _, name := range names {
		nameMap[name] = true
	}

	for _, seeder := range sr.seeders {
		if nameMap[seeder.Name()] {
			if err := seeder.Run(sr.db); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetSeeders returns all registered seeders
func (sr *SeederRunner) GetSeeders() []Seeder {
	return sr.seeders
}
