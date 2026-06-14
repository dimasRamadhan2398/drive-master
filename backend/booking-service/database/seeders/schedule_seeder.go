package seeders

import (
	"log"
	"time"

	"booking-service/models/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ScheduleSeeder seeds sample schedules for instructors
type ScheduleSeeder struct{}

// Name returns the seeder name
func (s *ScheduleSeeder) Name() string {
	return "schedules"
}

// Run executes the schedule seeder
func (s *ScheduleSeeder) Run(db *gorm.DB) error {
	log.Println("Running schedule seeder...")

	// Check if schedules already exist
	var count int64
	db.Model(&dto.Schedule{}).Count(&count)
	if count > 0 {
		log.Println("Schedules already exist, skipping...")
		return nil
	}

	// Get active instructors from user-service
	var instructorIDs []uuid.UUID
	if err := db.Raw(`
		SELECT ip.user_id
		FROM instructor_profiles ip
		JOIN users u ON u.id = ip.user_id
		WHERE ip.is_active = true
		LIMIT 10
	`).Scan(&instructorIDs).Error; err != nil {
		log.Printf("Warning: Could not fetch instructors, using defaults: %v", err)
		// Fallback: use default instructor UUIDs
		instructorIDs = []uuid.UUID{
			uuid.MustParse("22222222-2222-2222-2222-222222222201"),
			uuid.MustParse("22222222-2222-2222-2222-222222222202"),
			uuid.MustParse("22222222-2222-2222-2222-222222222203"),
		}
	}

	if len(instructorIDs) == 0 {
		log.Println("No instructors found, skipping schedule seeder...")
		return nil
	}

	// Get available cars from core-service
	var carIDs []uint
	if err := db.Raw(`
		SELECT id::integer as id
		FROM cars
		WHERE status = 'available'
		LIMIT 5
	`).Scan(&carIDs).Error; err != nil {
		log.Printf("Warning: Could not fetch cars, using defaults: %v", err)
		carIDs = []uint{1, 2, 3}
	}

	if len(carIDs) == 0 {
		log.Println("No cars found, skipping schedule seeder...")
		return nil
	}

	// Generate schedules for the next 14 days
	today := time.Now().Truncate(24 * time.Hour)
	schedules := []dto.Schedule{}

	// Time slots (Morning and Afternoon sessions)
	timeSlots := []struct {
		Time     string
		Duration int
	}{
		{"08:00", 90},
		{"10:00", 90},
		{"13:00", 90},
		{"15:00", 90},
		{"18:00", 90}, // Night session
	}

	// Generate for each instructor
	for dayOffset := 0; dayOffset < 14; dayOffset++ {
		date := today.AddDate(0, 0, dayOffset)
		dayOfWeek := date.Weekday()

		for _, instructorID := range instructorIDs {
			for _, carID := range carIDs {
				for _, slot := range timeSlots {
					// Skip night sessions on weekends (or keep them based on business rules)
					// Here we allow night sessions but could restrict based on instructor preferences

					schedule := dto.Schedule{
						Date:         date,
						Time:         slot.Time,
						Duration:     slot.Duration,
						InstructorID: instructorID,
						CarID:        carID,
						UserID:       nil,
						EnrollmentID: nil,
						Status:       dto.ScheduleStatusAvailable,
						Notes:        generateScheduleNotes(slot.Time, dayOfWeek),
					}

					// Make some slots booked (simulate bookings)
					// Every 5th slot per instructor per day is booked
					slotIndex := dayOffset%5 + 1
					if slotIndex%3 == 0 {
						userID := uint(1)
						schedule.UserID = &userID
						schedule.Status = dto.ScheduleStatusBooked
					}

					schedules = append(schedules, schedule)
				}
			}
		}
	}

	// Batch insert schedules
	for i := 0; i < len(schedules); i += 100 {
		end := i + 100
		if end > len(schedules) {
			end = len(schedules)
		}
		batch := schedules[i:end]
		if err := db.CreateInBatches(&batch, 100).Error; err != nil {
			log.Printf("Error seeding batch %d-%d: %v", i, end, err)
			return err
		}
	}

	log.Printf("Seeded %d schedules for %d instructors and %d cars over 14 days",
		len(schedules), len(instructorIDs), len(carIDs))
	return nil
}

// generateScheduleNotes generates helpful notes based on time slot
func generateScheduleNotes(timeStr string, dayOfWeek time.Weekday) string {
	switch timeStr {
	case "18:00":
		return "Night session - runs 18:00-19:30"
	default:
		if dayOfWeek == time.Saturday || dayOfWeek == time.Sunday {
			return "Weekend session"
		}
		return "Regular session"
	}
}