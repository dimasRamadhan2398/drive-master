package user

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// APIResponse is the standard response format from user-service
type APIResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// RecurringScheduleResponse represents a recurring schedule slot from user-service
type RecurringScheduleResponse struct {
	ID           uuid.UUID `json:"id"`
	InstructorID uuid.UUID `json:"instructorId"`
	DayOfWeek    int       `json:"dayOfWeek"` // 0=Sunday, 1=Monday...6=Saturday
	DayName      string    `json:"dayName"`  // e.g., "Monday"
	StartTime    string    `json:"startTime"` // HH:MM format
	EndTime      string    `json:"endTime"`   // HH:MM format
	IsActive     bool      `json:"isActive"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}