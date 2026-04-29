package seeders

import (
	"time"
	"user-service/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WorkExperienceSeeder seeds sample work experiences for instructors
type WorkExperienceSeeder struct {
	db *gorm.DB
}

func NewWorkExperienceSeeder(db *gorm.DB) *WorkExperienceSeeder {
	return &WorkExperienceSeeder{db: db}
}

// Seed inserts sample work experiences
func (s *WorkExperienceSeeder) Seed() error {
	// Get instructor users
	var instructors []models.User
	if err := s.db.Joins("JOIN roles ON users.role_id = roles.id").
		Where("roles.name = ?", "instructor").
		Find(&instructors).Error; err != nil {
		return err
	}

	if len(instructors) == 0 {
		return nil // No instructors found, skip
	}

	workExperiences := []struct {
		instructorID string
		companyName  string
		role         string
		startDate    time.Time
		endDate      time.Time
		description  string
	}{
		{
			instructorID: instructors[0].ID.String(),
			companyName:  "Premier Driving School",
			role:         "Senior Driving Instructor",
			startDate:    parseTime("2015-01-15"),
			endDate:      parseTime("2020-06-30"),
			description:  "Led driving education programs for new drivers. Specialized in defensive driving techniques.",
		},
		{
			instructorID: instructors[0].ID.String(),
			companyName:  "Safe Wheels Academy",
			role:         "Driving Instructor",
			startDate:    parseTime("2012-03-01"),
			endDate:      parseTime("2014-12-31"),
			description:  "Provided driving lessons for students of all ages and skill levels.",
		},
		{
			instructorID: instructors[0].ID.String(),
			companyName:  "National Traffic Safety Council",
			role:         "Road Safety Trainer",
			startDate:    parseTime("2010-06-01"),
			endDate:      parseTime("2012-02-28"),
			description:  "Conducted road safety training and awareness programs.",
		},
		{
			instructorID: instructors[0].ID.String(),
			companyName:  "Elite Auto School",
			role:         "Head Instructor",
			startDate:    parseTime("2020-07-01"),
			endDate:      time.Time{}, // Current position, no end date
			description:  "Currently leading a team of driving instructors. Focus on quality assurance and curriculum development.",
		},
		{
			instructorID: instructors[0].ID.String(),
			companyName:  "Metro Driving Center",
			role:         "Driving Instructor",
			startDate:    parseTime("2018-09-01"),
			endDate:      parseTime("2023-08-31"),
			description:  "Taught defensive driving and provided practical driving lessons.",
		},
	}

	// Add more experiences for second instructor if exists
	if len(instructors) > 1 {
		moreExperiences := []struct {
			instructorID string
			companyName  string
			role         string
			startDate    time.Time
			endDate      time.Time
			description  string
		}{
			{
				instructorID: instructors[1].ID.String(),
				companyName:  "City Driving Academy",
				role:         "Driving Instructor",
				startDate:    parseTime("2019-02-01"),
				endDate:      time.Time{}, // Current position
				description:  "Providing comprehensive driving lessons with focus on city driving and traffic awareness.",
			},
			{
				instructorID: instructors[1].ID.String(),
				companyName:  "Highway Safety Institute",
				role:         "Defensive Driving Instructor",
				startDate:    parseTime("2017-06-15"),
				endDate:      parseTime("2018-12-31"),
				description:  "Specialized in highway driving techniques and emergency handling.",
			},
		}
		workExperiences = append(workExperiences, moreExperiences...)
	}

	for _, exp := range workExperiences {
		instructorUUID := parseUUID(exp.instructorID)
		record := models.WorkExperience{
			InstructorID: instructorUUID,
			CompanyName:  exp.companyName,
			Role:         exp.role,
			StartDate:    exp.startDate,
			EndDate:      &exp.endDate,
			Description:  exp.description,
		}

		// Use FirstOrCreate to avoid duplicates
		if err := s.db.FirstOrCreate(&record, models.WorkExperience{
			InstructorID: instructorUUID,
			CompanyName:  exp.companyName,
			Role:         exp.role,
		}).Error; err != nil {
			return err
		}
	}

	return nil
}

func parseTime(dateStr string) time.Time {
	t, _ := time.Parse("2006-01-02", dateStr)
	return t
}

func parseUUID(s string) uuid.UUID {
	id, _ := uuid.Parse(s)
	return id
}
