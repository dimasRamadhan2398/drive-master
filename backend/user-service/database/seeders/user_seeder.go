package seeders

import (
	"time"

	"user-service/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserSeeder seeds sample users
type UserSeeder struct {
	db *gorm.DB
}

func NewUserSeeder(db *gorm.DB) *UserSeeder {
	return &UserSeeder{db: db}
}

// Seed inserts sample users into the database
func (s *UserSeeder) Seed(roleMap map[string]uint) error {
	// Hash password for all users
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	users := []struct {
		user              models.User
		memberProfile     *models.MemberProfile
		instructorProfile  *models.InstructorProfile
	}{
		{
			user: models.User{
				Username:     "admin",
				PasswordHash: string(hashedPassword),
				Email:        "admin@example.com",
				EmailAddress: "admin@example.com",
				PhoneNumber:  "+6281234567890",
				DateOfBirth:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				RoleID:       roleMap["admin"],
				Image:        "https://example.com/images/admin.jpg",
				Address:      "Jakarta, Indonesia",
			},
		},
		{
			user: models.User{
				Username:     "super_admin",
				PasswordHash: string(hashedPassword),
				Email:        "superadmin@example.com",
				EmailAddress: "superadmin@example.com",
				PhoneNumber:  "+6281234567891",
				DateOfBirth:  time.Date(1985, 6, 15, 0, 0, 0, 0, time.UTC),
				RoleID:       roleMap["super_admin"],
				Image:        "https://example.com/images/super_admin.jpg",
				Address:      "Jakarta, Indonesia",
			},
		},
		{
			user: models.User{
				Username:     "instructor1",
				PasswordHash: string(hashedPassword),
				Email:        "instructor1@example.com",
				EmailAddress: "instructor1@example.com",
				PhoneNumber:  "+6281234567892",
				DateOfBirth:  time.Date(1988, 3, 20, 0, 0, 0, 0, time.UTC),
				RoleID:       roleMap["instructor"],
				Image:        "https://example.com/images/instructor1.jpg",
				Address:      "Bandung, Indonesia",
			},
			instructorProfile: &models.InstructorProfile{
				UserID:                 uuid.Nil, // Will be set after user creation
				LicenseNumber:          "DRIVING-12345",
				LicenseExpiry:          time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
				BNSPCertificateNumber:   "BNSP-2024-123456",
				Bio:                    "Experienced driving instructor with 10+ years of experience",
				PhotoURL:               "https://example.com/images/instructor1-profile.jpg",
				IsActive:               true,
				NumberOfStudents:       150,
				YearsOfExperience:      10,
				SessionsCompleted:      500,
				AverageRating:          4.8,
			},
		},
		{
			user: models.User{
				Username:     "instructor2",
				PasswordHash: string(hashedPassword),
				Email:        "instructor2@example.com",
				EmailAddress: "instructor2@example.com",
				PhoneNumber:  "+6281234567893",
				DateOfBirth:  time.Date(1992, 7, 10, 0, 0, 0, 0, time.UTC),
				RoleID:       roleMap["instructor"],
				Image:        "https://example.com/images/instructor2.jpg",
				Address:      "Surabaya, Indonesia",
			},
			instructorProfile: &models.InstructorProfile{
				UserID:                 uuid.Nil,
				LicenseNumber:          "DRIVING-67890",
				LicenseExpiry:          time.Date(2027, 6, 30, 0, 0, 0, 0, time.UTC),
				BNSPCertificateNumber:  "BNSP-2024-789012",
				Bio:                    "Professional driving instructor specializing in defensive driving",
				PhotoURL:               "https://example.com/images/instructor2-profile.jpg",
				IsActive:               true,
				NumberOfStudents:      85,
				YearsOfExperience:     5,
				SessionsCompleted:      300,
				AverageRating:         4.6,
			},
		},
		{
			user: models.User{
				Username:     "member1",
				PasswordHash: string(hashedPassword),
				Email:        "member1@example.com",
				EmailAddress: "member1@example.com",
				PhoneNumber:  "+6281234567894",
				DateOfBirth:  time.Date(2000, 11, 25, 0, 0, 0, 0, time.UTC),
				RoleID:       roleMap["member"],
				Image:        "https://example.com/images/member1.jpg",
				Address:      "Jakarta, Indonesia",
			},
			memberProfile: &models.MemberProfile{
				UserID:             uuid.Nil, // Will be set after user creation
				SessionsCompleted: 12,
				TrainingTime:      600,
				AverageRating:     4.5,
			},
		},
		{
			user: models.User{
				Username:     "member2",
				PasswordHash: string(hashedPassword),
				Email:        "member2@example.com",
				EmailAddress: "member2@example.com",
				PhoneNumber:  "+6281234567895",
				DateOfBirth:  time.Date(1998, 4, 8, 0, 0, 0, 0, time.UTC),
				RoleID:       roleMap["member"],
				Image:        "https://example.com/images/member2.jpg",
				Address:      "Bekasi, Indonesia",
			},
			memberProfile: &models.MemberProfile{
				UserID:             uuid.Nil,
				SessionsCompleted: 8,
				TrainingTime:       400,
				AverageRating:      4.2,
			},
		},
	}

	for _, u := range users {
		// Check if user already exists
		var existingUser models.User
		if err := s.db.Where("email_address = ?", u.user.EmailAddress).First(&existingUser).Error; err == nil {
			continue // User already exists, skip
		}

		// Create user
		if err := s.db.Create(&u.user).Error; err != nil {
			return err
		}

		// Create member profile if exists
		if u.memberProfile != nil {
			u.memberProfile.UserID = u.user.ID
			if err := s.db.Create(u.memberProfile).Error; err != nil {
				return err
			}
		}

		// Create instructor profile if exists
		if u.instructorProfile != nil {
			u.instructorProfile.UserID = u.user.ID
			if err := s.db.Create(u.instructorProfile).Error; err != nil {
				return err
			}
		}
	}

	return nil
}