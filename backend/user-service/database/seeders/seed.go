package seeders

import (
	"user-service/models"

	"gorm.io/gorm"
)

// RoleSeeder seeds the initial roles
type RoleSeeder struct {
	db *gorm.DB
}

func NewRoleSeeder(db *gorm.DB) *RoleSeeder {
	return &RoleSeeder{db: db}
}

// Seed inserts initial roles into the database
func (s *RoleSeeder) Seed() error {
	roles := []models.Role{
		{Name: "admin"},
		{Name: "instructor"},
		{Name: "member"},
		{Name: "super_admin"},
	}

	for _, role := range roles {
		// Use FirstOrCreate to avoid duplicates
		if err := s.db.FirstOrCreate(&role, models.Role{Name: role.Name}).Error; err != nil {
			return err
		}
	}

	return nil
}

// Roles returns all seeded roles
func (s *RoleSeeder) GetRoles() (map[string]uint, error) {
	var roles []models.Role
	if err := s.db.Find(&roles).Error; err != nil {
		return nil, err
	}

	roleMap := make(map[string]uint)
	for _, role := range roles {
		roleMap[role.Name] = role.ID
	}
	return roleMap, nil
}
