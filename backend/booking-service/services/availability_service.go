package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"booking-service/clients/user"

	"github.com/google/uuid"
)

// IAvailabilityService defines the interface for availability checking
type IAvailabilityService interface {
	CheckAvailability(ctx context.Context, instructorID uuid.UUID, date time.Time, startTime string, duration int) error
}

// AvailabilityService checks instructor availability against recurring schedules
type AvailabilityService struct {
	userClient user.IUserClient
}

// NewAvailabilityService creates a new availability service
func NewAvailabilityService(userClient user.IUserClient) IAvailabilityService {
	return &AvailabilityService{
		userClient: userClient,
	}
}

// CheckAvailability validates that an instructor is available at the given time
// It checks the instructor's recurring schedules in user-service
func (s *AvailabilityService) CheckAvailability(ctx context.Context, instructorID uuid.UUID, date time.Time, startTime string, duration int) error {
	// Get the day of week (0=Sunday, 1=Monday, ..., 6=Saturday)
	dayOfWeek := int(date.Weekday())

	// Call user-service to get instructor's recurring schedules
	schedules, err := s.userClient.GetInstructorRecurringSchedules(ctx, instructorID)
	if err != nil {
		// If we can't reach user-service, allow the booking (fail open)
		// In production, you might want to fail closed instead
		return fmt.Errorf("failed to check instructor availability: %w", err)
	}

	// Calculate end time
	endTime, err := calculateEndTime(startTime, duration)
	if err != nil {
		return err
	}

	// Check if any active schedule matches the requested day and time
	for _, schedule := range schedules {
		// Only check active schedules
		if !schedule.IsActive {
			continue
		}

		// Check if day of week matches
		if schedule.DayOfWeek != dayOfWeek {
			continue
		}

		// Check if time overlaps
		if timesOverlap(startTime, endTime, schedule.StartTime, schedule.EndTime) {
			return nil // Availability confirmed
		}
	}

	// No matching schedule found - instructor not available at this time
	return errors.New("instructor is not available at this time")
}

// calculateEndTime computes the end time given start time and duration in minutes
func calculateEndTime(startTime string, durationMinutes int) (string, error) {
	// Parse start time (HH:MM format)
	startParts, err := parseTime(startTime)
	if err != nil {
		return "", fmt.Errorf("invalid start time format: %w", err)
	}

	// Add duration
	totalMinutes := startParts.hour*60 + startParts.minute + durationMinutes
	endHour := totalMinutes / 60
	endMinute := totalMinutes % 60

	// Handle overflow past midnight (shouldn't happen for normal durations)
	if endHour >= 24 {
		endHour = endHour % 24
	}

	return fmt.Sprintf("%02d:%02d", endHour, endMinute), nil
}

// timesOverlap checks if two time ranges overlap
// times are in HH:MM format
func timesOverlap(start1, end1, start2, end2 string) bool {
	// Parse all times
	s1, _ := parseTime(start1)
	e1, _ := parseTime(end1)
	s2, _ := parseTime(start2)
	e2, _ := parseTime(end2)

	// Convert to minutes for comparison
	m1Start := s1.hour*60 + s1.minute
	m1End := e1.hour*60 + e1.minute
	m2Start := s2.hour*60 + s2.minute
	m2End := e2.hour*60 + e2.minute

	// Two ranges overlap if:
	// range1 starts before range2 ends AND range1 ends after range2 starts
	return m1Start < m2End && m1End > m2Start
}

type timeParts struct {
	hour   int
	minute int
}

// parseTime parses HH:MM time string into hour and minute
func parseTime(timeStr string) (timeParts, error) {
	if len(timeStr) != 5 || timeStr[2] != ':' {
		return timeParts{}, errors.New("invalid time format, expected HH:MM")
	}

	var hour, minute int
	_, err := fmt.Sscanf(timeStr, "%02d:%02d", &hour, &minute)
	if err != nil {
		return timeParts{}, errors.New("invalid time format")
	}

	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return timeParts{}, errors.New("time out of range")
	}

	return timeParts{hour: hour, minute: minute}, nil
}

// FormatScheduleTime formats time parts back to HH:MM string
func FormatScheduleTime(hour, minute int) string {
	return fmt.Sprintf("%02d:%02d", hour, minute)
}