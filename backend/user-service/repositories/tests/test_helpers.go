package tests

import (
	"testing"
	"time"

	"user-service/models"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB creates an in-memory SQLite database for testing
func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	return db
}

// CreateMockUser creates a user with a random UUID for testing
func CreateMockUser() *models.User {
	id := uuid.New()
	return &models.User{
		ID:           id,
		Username:     "testuser",
		PasswordHash: "$2a$10$hashedpassword",
		EmailAddress: "test@example.com",
		PhoneNumber:  "+1234567890",
		RoleID:       1,
	}
}

// CreateMockUserWithRole creates a user with a specific role ID
func CreateMockUserWithRole(roleID uint) *models.User {
	user := CreateMockUser()
	user.RoleID = roleID
	return user
}

// CreateMockRole creates a role for testing
func CreateMockRole(id uint, name string) *models.Role {
	return &models.Role{
		ID:   id,
		Name: name,
	}
}

// CreateMockMemberProfile creates a member profile for testing
func CreateMockMemberProfile(userID uuid.UUID) *models.MemberProfile {
	return &models.MemberProfile{
		UserID:            userID,
		SessionsCompleted: 10,
		TrainingTime:      600,
		AverageRating:     4.5,
	}
}

// CreateMockInstructorProfile creates an instructor profile for testing
func CreateMockInstructorProfile(userID uuid.UUID) *models.InstructorProfile {
	return &models.InstructorProfile{
		UserID:               userID,
		LicenseNumber:        "LIC123456",
		BNSPCertificateNumber: "BNSP123456",
		Bio:                  "Experienced instructor",
		IsActive:             true,
		NumberOfStudents:     25,
		YearsOfExperience:    5,
		SessionsCompleted:    150,
		AverageRating:        4.8,
	}
}

// CreateMockWorkExperience creates a work experience for testing
func CreateMockWorkExperience(instructorID uuid.UUID) *models.WorkExperience {
	return &models.WorkExperience{
		ID:           1,
		InstructorID: instructorID,
		CompanyName:  "Test Company",
		Role:         "Senior Instructor",
		Description:  "Teaching driving lessons",
	}
}

// Now returns current timestamp
func Now() time.Time {
	return time.Now()
}