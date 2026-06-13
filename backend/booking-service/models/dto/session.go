package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateScheduleRequest struct {
	Date         string `json:"date" binding:"required"`                   // YYYY-MM-DD format
	Time         string `json:"time" binding:"required"`                  // HH:MM format
	Duration     int    `json:"duration" binding:"omitempty"`             // in minutes, default 90
	InstructorID uuid.UUID `json:"instructorId" binding:"required"`
	CarID        uint   `json:"carId" binding:"required"`
	Notes        string `json:"notes"`
}

type UpdateScheduleRequest struct {
	Date         *string `json:"date" binding:"omitempty"`    // YYYY-MM-DD format
	Time         *string `json:"time" binding:"omitempty"`     // HH:MM format
	Duration     *int    `json:"duration" binding:"omitempty"`
	InstructorID *uuid.UUID `json:"instructorId" binding:"omitempty"`
	CarID        *uint   `json:"carId" binding:"omitempty"`
	Notes        *string `json:"notes" binding:"omitempty"`
	Status       *string `json:"status" binding:"omitempty"`  // update slot status
}

type ScheduleResponse struct {
	ID           uint      `json:"id"`
	Date         string    `json:"date"`         // YYYY-MM-DD format
	Time         string    `json:"time"`         // HH:MM format
	Duration     int       `json:"duration"`     // in minutes
	InstructorID uuid.UUID `json:"instructorId"`
	CarID        uint      `json:"carId"`
	UserID       *uint     `json:"userId"`        // nullable
	BookingID    *uint     `json:"bookingId"`     // nullable
	Status       string    `json:"status"`
	Notes        string    `json:"notes"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type ScheduleWithDetailsResponse struct {
	ScheduleResponse
	InstructorName string `json:"instructorName,omitempty"`
	CarName        string `json:"carName,omitempty"`
	StudentName    string `json:"studentName,omitempty"`
}

type ScheduleListResponse struct {
	Data       []ScheduleResponse         `json:"data"`
	Total      int64                       `json:"total"`
	Page       int                         `json:"page"`
	Limit      int                         `json:"limit"`
	TotalPages int                         `json:"totalPages"`
}

type ScheduleFilterParams struct {
	ListParams
	Date         string `form:"date"`          // YYYY-MM-DD format
	StartDate    string `form:"startDate"`      // YYYY-MM-DD format
	EndDate      string `form:"endDate"`        // YYYY-MM-DD format
	InstructorID string `form:"instructorId"`
	CarID        uint   `form:"carId"`
	Status       string `form:"status"`
}