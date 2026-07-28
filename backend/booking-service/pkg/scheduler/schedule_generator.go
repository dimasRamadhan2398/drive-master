package scheduler

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"booking-service/clients/user"
	"booking-service/models/dto"
	"booking-service/repositories"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

// ScheduleGenerator handles automatic generation of schedule slots from recurring schedules
type ScheduleGenerator struct {
	scheduleRepo   repositories.IScheduleRepository
	userClient     user.IUserClient
	generationDays int // How many days ahead to generate schedules
	cronScheduler  *cron.Cron
}

// NewScheduleGenerator creates a new schedule generator
func NewScheduleGenerator(
	scheduleRepo repositories.IScheduleRepository,
	userClient user.IUserClient,
) *ScheduleGenerator {
	return &ScheduleGenerator{
		scheduleRepo:   scheduleRepo,
		userClient:      userClient,
		cronScheduler:   cron.New(),
		generationDays: 7, // Generate schedules 7 days ahead by default
	}
}

// Start begins the cron scheduler for automatic schedule generation
// cronExpr: standard cron expression, e.g., "0 0 * * *" (midnight daily)
// Default is "5 0 * * *" (00:05 AM every day)
func (sg *ScheduleGenerator) Start(cronExpr string) error {
	if cronExpr == "" {
		cronExpr = "5 0 * * *" // 00:05 AM every day
	}

	_, err := sg.cronScheduler.AddFunc(cronExpr, func() {
		ctx := context.Background()
		if err := sg.GenerateSchedules(ctx); err != nil {
			log.Printf("Schedule generation failed: %v", err)
		}
	})
	if err != nil {
		return fmt.Errorf("failed to add cron job: %w", err)
	}

	sg.cronScheduler.Start()
	log.Printf("Schedule generator started with cron: %s", cronExpr)
	return nil
}

// Stop gracefully stops the scheduler
func (sg *ScheduleGenerator) Stop() {
	ctx := sg.cronScheduler.Stop()
	<-ctx.Done()
	log.Println("Schedule generator stopped")
}

// GenerateSchedules generates schedule slots for all instructors based on their recurring schedules
func (sg *ScheduleGenerator) GenerateSchedules(ctx context.Context) error {
	log.Println("Starting daily schedule generation...")

	// Get all instructors with their recurring schedules
	instructors, err := sg.userClient.GetAllInstructors(ctx)
	if err != nil {
		log.Printf("Failed to get instructors: %v", err)
		return err
	}

	generatedCount := 0
	skippedCount := 0

	for _, instructor := range instructors {
		// Filter only active recurring schedules
		var activeSchedules []user.RecurringScheduleDTO
		for _, schedule := range instructor.RecurringSchedules {
			if schedule.IsActive {
				activeSchedules = append(activeSchedules, schedule)
			}
		}

		if len(activeSchedules) == 0 {
			continue
		}

		count, err := sg.generateSchedulesForInstructor(ctx, instructor.ID, activeSchedules)
		if err != nil {
			log.Printf("Failed to generate schedules for instructor %s: %v", instructor.ID, err)
			continue
		}
		generatedCount += count
		if count == 0 {
			skippedCount++
		}
	}

	log.Printf("Daily schedule generation completed. Generated %d new slots for %d instructors (skipped: %d).",
		generatedCount, len(instructors)-skippedCount, skippedCount)
	return nil
}

// generateSchedulesForInstructor generates schedule slots for a single instructor
func (sg *ScheduleGenerator) generateSchedulesForInstructor(
	ctx context.Context,
	instructorID uuid.UUID,
	recurringSchedules []user.RecurringScheduleDTO,
) (int, error) {
	generatedCount := 0

	// Generate schedules for the next N days
	for dayOffset := 0; dayOffset < sg.generationDays; dayOffset++ {
		targetDate := time.Now().AddDate(0, 0, dayOffset)
		targetDayOfWeek := int(targetDate.Weekday())

		// Find recurring schedules that match this day of week
		for _, recurring := range recurringSchedules {
			if recurring.DayOfWeek != targetDayOfWeek {
				continue
			}

			// Generate slots for this recurring schedule
			count, err := sg.generateSlotsForRecurringSchedule(ctx, instructorID, targetDate, recurring)
			if err != nil {
				log.Printf("Failed to generate slots for instructor %s on %s: %v",
					instructorID, targetDate.Format("2006-01-02"), err)
				continue
			}
			generatedCount += count
		}
	}

	return generatedCount, nil
}

// generateSlotsForRecurringSchedule generates time slots for a single recurring schedule on a specific date
func (sg *ScheduleGenerator) generateSlotsForRecurringSchedule(
	ctx context.Context,
	instructorID uuid.UUID,
	date time.Time,
	recurring user.RecurringScheduleDTO,
) (int, error) {
	// Parse start and end times (HH:MM format)
	startMinutes, err := parseTimeToMinutes(recurring.StartTime)
	if err != nil {
		return 0, fmt.Errorf("invalid start time %s: %w", recurring.StartTime, err)
	}

	endMinutes, err := parseTimeToMinutes(recurring.EndTime)
	if err != nil {
		return 0, fmt.Errorf("invalid end time %s: %w", recurring.EndTime, err)
	}

	// Generate slots with default duration (60 minutes each)
	slotDuration := 60
	generatedCount := 0

	currentMinutes := startMinutes
	for currentMinutes+slotDuration <= endMinutes {
		// Format time string
		timeStr := formatMinutesToTime(currentMinutes)

		// Check if slot already exists
		exists, err := sg.scheduleRepo.ExistsForInstructorAndDateTime(ctx, instructorID, date, timeStr)
		if err != nil {
			log.Printf("Failed to check existing slot: %v", err)
			continue
		}

		if !exists {
			// Create the schedule slot (status: available, car_id: nil - to be assigned later)
			schedule := &dto.Schedule{
				Date:         date,
				Time:         timeStr,
				Duration:     slotDuration,
				InstructorID: instructorID,
				CarID:        uuid.Nil, // Will be assigned when booked
				Status:       dto.ScheduleStatusAvailable,
			}

			if err := sg.scheduleRepo.Create(ctx, schedule); err != nil {
				log.Printf("Failed to create schedule slot: %v", err)
			} else {
				generatedCount++
				log.Printf("Generated schedule slot: Instructor=%s, Date=%s, Time=%s",
					instructorID, date.Format("2006-01-02"), timeStr)
			}
		}

		// Move to next slot
		currentMinutes += slotDuration
	}

	return generatedCount, nil
}

// RunOnce generates schedules without starting the cron scheduler (for manual trigger)
func (sg *ScheduleGenerator) RunOnce(ctx context.Context) error {
	return sg.GenerateSchedules(ctx)
}

// SetGenerationDays sets how many days ahead to generate schedules
func (sg *ScheduleGenerator) SetGenerationDays(days int) {
	sg.generationDays = days
}

// parseTimeToMinutes converts "HH:MM" string to minutes since midnight
func parseTimeToMinutes(timeStr string) (int, error) {
	if len(timeStr) != 5 || timeStr[2] != ':' {
		return 0, fmt.Errorf("invalid time format: %s (expected HH:MM)", timeStr)
	}

	hourStr := timeStr[0:2]
	minuteStr := timeStr[3:5]

	hour, err := strconv.Atoi(hourStr)
	if err != nil {
		return 0, fmt.Errorf("invalid hour: %s", hourStr)
	}

	minute, err := strconv.Atoi(minuteStr)
	if err != nil {
		return 0, fmt.Errorf("invalid minute: %s", minuteStr)
	}

	if hour < 0 || hour > 23 {
		return 0, fmt.Errorf("hour out of range: %d", hour)
	}

	if minute < 0 || minute > 59 {
		return 0, fmt.Errorf("minute out of range: %d", minute)
	}

	return hour*60 + minute, nil
}

// formatMinutesToTime converts minutes since midnight to "HH:MM" format
func formatMinutesToTime(minutes int) string {
	hour := minutes / 60
	minute := minutes % 60
	return fmt.Sprintf("%02d:%02d", hour, minute)
}
