package seeders

import (
	"log"
	"time"

	"booking-service/models"

	"gorm.io/gorm"
)

// ScheduleSeeder seeds sample schedules for testing
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
	db.Model(&models.Schedule{}).Count(&count)
	if count > 0 {
		log.Println("Schedules already exist, skipping...")
		return nil
	}

	// Generate schedules for the next 7 days
	today := time.Now()
	schedules := []models.Schedule{}

	// Time slots
	timeSlots := []string{"08:00", "09:00", "10:00", "11:00", "13:00", "14:00", "15:00", "16:00"}

	// Generate for multiple instructors and cars
	instructors := []uint{1, 2, 3}
	cars := []uint{1, 2, 3}

	for dayOffset := 0; dayOffset < 7; dayOffset++ {
		date := today.AddDate(0, 0, dayOffset)
		for _, instructorID := range instructors {
			for _, carID := range cars {
				for _, timeSlot := range timeSlots {
					// Randomly make some slots available and some booked
					schedule := models.Schedule{
						Date:         date,
						Time:         timeSlot,
						Duration:     60,
						InstructorID: instructorID,
						CarID:        carID,
						UserID:       nil,
						Status:       models.ScheduleStatusAvailable,
						Notes:        "",
					}

					// Make some slots booked (every 4th slot)
					if (dayOffset+int(instructorID)+int(carID))%4 == 0 {
						userID := uint(1)
						schedule.UserID = &userID
						schedule.Status = models.ScheduleStatusBooked
					}

					schedules = append(schedules, schedule)
				}
			}
		}
	}

	for _, schedule := range schedules {
		if err := db.Create(&schedule).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d schedules", len(schedules))
	return nil
}