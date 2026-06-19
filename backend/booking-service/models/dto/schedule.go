package dto

import (
	"time"

	"github.com/google/uuid"
)

// ScheduleStatus represents the status of a schedule slot
type ScheduleStatus string

const (
	ScheduleStatusAvailable  ScheduleStatus = "available"
	ScheduleStatusBooked     ScheduleStatus = "booked"
	ScheduleStatusInProgress ScheduleStatus = "in-progress"
	ScheduleStatusCompleted  ScheduleStatus = "completed"
	ScheduleStatusBlocked    ScheduleStatus = "blocked"
)

// Schedule represents a time slot in the scheduling calendar.
// It tracks availability for instructors and cars, and links to enrollments.
type Schedule struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Date         time.Time      `json:"date" gorm:"type:date;not null;index"`
	Time         string         `json:"time" gorm:"size:10;not null"` // HH:MM format
	Duration     int            `json:"duration" gorm:"default:60"`   // duration in minutes
	InstructorID uuid.UUID      `json:"instructorId" gorm:"type:uuid;not null;index"` // ref: user-service (UUID)
	CarID        uuid.UUID      `json:"carId" gorm:"type:uuid;not null;index"`   // ref: core-service (car)
	UserID       *uuid.UUID     `json:"userId" gorm:"type:uuid;index"`           // ref: user-service (nullable, assigned when booked)
	EnrollmentID *uuid.UUID     `json:"enrollmentId" gorm:"index"`             // ref: Enrollment (nullable)
	Status       ScheduleStatus `json:"status" gorm:"type:varchar(20);default:'available'"`
	Notes        string         `json:"notes" gorm:"type:text"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
}

type CreateScheduleRequest struct {
	Date         string    `json:"date" binding:"required"`                   // YYYY-MM-DD format
	Time         string    `json:"time" binding:"required"`                  // HH:MM format
	Duration     int       `json:"duration" binding:"omitempty"`             // in minutes, default 90
	InstructorID uuid.UUID `json:"instructorId" binding:"required"`
	CarID        uuid.UUID `json:"carId" binding:"required"`
	Notes        string    `json:"notes"`
}

type UpdateScheduleRequest struct {
	Date         *string    `json:"date" binding:"omitempty"`    // YYYY-MM-DD format
	Time         *string    `json:"time" binding:"omitempty"`     // HH:MM format
	Duration     *int       `json:"duration" binding:"omitempty"`
	InstructorID *uuid.UUID `json:"instructorId" binding:"omitempty"`
	CarID        *uuid.UUID `json:"carId" binding:"omitempty"`
	Notes        *string    `json:"notes" binding:"omitempty"`
	Status       *string    `json:"status" binding:"omitempty"`  // update slot status
}

type ScheduleResponse struct {
	ID             uint      `json:"id"`
	Date           string    `json:"date"`         // YYYY-MM-DD format
	Time           string    `json:"time"`         // HH:MM format
	Duration       int       `json:"duration"`     // in minutes
	InstructorID   uuid.UUID `json:"instructorId"`
	InstructorName string    `json:"instructorName"`
	CarID          uuid.UUID `json:"carId"`
	CarName        string    `json:"carName"`
	UserID         *uuid.UUID `json:"userId"`
	UserName       *string   `json:"userName"`        // nullable
	BookingID      *uint     `json:"bookingId"`     // nullable
	Status         string    `json:"status"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type ScheduleWithDetailsResponse struct {
	ScheduleResponse
	InstructorName string `json:"instructorName,omitempty"`
	CarName        string `json:"carName,omitempty"`
	StudentName    string `json:"studentName,omitempty"`
}

// type ScheduleListResponse struct {
// 	Data       []ScheduleResponse `json:"data"`
// 	Total      int64               `json:"total"`
// 	Page       int                 `json:"page"`
// 	Limit      int                 `json:"limit"`
// 	TotalPages int                 `json:"totalPages"`
// }

type ScheduleListResponse = PagedData[ScheduleResponse]

type ScheduleFilterParams struct {
	ListParams
	Date         string `form:"date"`          // YYYY-MM-DD format
	StartDate    string `form:"startDate"`      // YYYY-MM-DD format
	EndDate      string `form:"endDate"`        // YYYY-MM-DD format
	InstructorID string `form:"instructorId"`
	CarID        string `form:"carId"`
	Status       string `form:"status"`
}

type ScheduleStatsResponse struct {
	AvailableSchedule int64 `json:"availableSchedule"`
	BookedSchedule    int64 `json:"bookedSchedule"`
	InProgressSchedule int64 `json:"inProgressSchedule"`
	CompletedSchedule int64 `json:"completedSchedule"`
	BlockedSchedule    int64 `json:"blockedSchedule"`
}