package dto

import "time"

// Booking DTOs

type CreateBookingRequest struct {
	UserID        uint      `json:"userId" binding:"required"`
	InstructorID  uint      `json:"instructorId" binding:"required"`
	EntitlementID uint      `json:"entitlementId" binding:"required"`
	DateOfSession time.Time `json:"dateOfSession" binding:"required"`
	FromTime      time.Time `json:"fromTime" binding:"required"`
	ToTime        time.Time `json:"toTime" binding:"required"`
	CarID         uint      `json:"carId" binding:"required"`
	Area          string    `json:"area"`
	Notes         string    `json:"notes"`
}

type UpdateBookingRequest struct {
	DateOfSession *time.Time `json:"dateOfSession" binding:"omitempty"`
	FromTime      *time.Time `json:"fromTime" binding:"omitempty"`
	ToTime        *time.Time `json:"toTime" binding:"omitempty"`
	CarID         *uint      `json:"carId" binding:"omitempty"`
	Area          *string    `json:"area" binding:"omitempty"`
	Notes         *string    `json:"notes" binding:"omitempty"`
	Status        *string    `json:"status" binding:"omitempty"`
}

type BookingResponse struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"userId"`
	InstructorID  uint      `json:"instructorId"`
	EntitlementID uint      `json:"entitlementId"`
	DateOfSession time.Time `json:"dateOfSession"`
	FromTime      time.Time `json:"fromTime"`
	ToTime        time.Time `json:"toTime"`
	CarID         uint      `json:"carId"`
	Area          string    `json:"area"`
	Notes         string    `json:"notes"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type BookingListResponse struct {
	Data       []BookingResponse `json:"data"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"totalPages"`
}

// Session DTOs

type CreateSessionRequest struct {
	UserID        uint      `json:"userId" binding:"required"`
	InstructorID  uint      `json:"instructorId" binding:"required"`
	DateOfSession time.Time `json:"dateOfSession" binding:"required"`
	Duration      int       `json:"duration" binding:"required"`
	CarID         uint      `json:"carId" binding:"required"`
	Area          string    `json:"area"`
	Notes         string    `json:"notes"`
}

type SessionResponse struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"userId"`
	InstructorID  uint      `json:"instructorId"`
	DateOfSession time.Time `json:"dateOfSession"`
	Duration      int       `json:"duration"`
	CarID         uint      `json:"carId"`
	Area          string    `json:"area"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type SessionListResponse struct {
	Data       []SessionResponse `json:"data"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"totalPages"`
}

// UserEntitlement DTOs

type CreateEntitlementRequest struct {
	UserID            uint      `json:"userId" binding:"required"`
	SourceType        string    `json:"sourceType" binding:"required"`
	SourceID          string    `json:"sourceId" binding:"required"`
	TotalSessions     int       `json:"totalSessions" binding:"required"`
	SessionsRemaining int       `json:"sessionsRemaining" binding:"required"`
	ExpiresAt         time.Time `json:"expiresAt" binding:"required"`
}

type UpdateEntitlementRequest struct {
	SessionsRemaining *int       `json:"sessionsRemaining" binding:"omitempty"`
	ExpiresAt         *time.Time `json:"expiresAt" binding:"omitempty"`
}

type EntitlementResponse struct {
	ID                uint      `json:"id"`
	UserID            uint      `json:"userId"`
	SourceType        string    `json:"sourceType"`
	SourceID          string    `json:"sourceId"`
	TotalSessions     int       `json:"totalSessions"`
	SessionsRemaining int       `json:"sessionsRemaining"`
	ExpiresAt         time.Time `json:"expiresAt"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type EntitlementListResponse struct {
	Data       []EntitlementResponse `json:"data"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
	TotalPages int                    `json:"totalPages"`
}

// Certification DTOs

type CreateCertificationRequest struct {
	Type       string `json:"type" binding:"required"`
	Recipient  string `json:"recipient" binding:"required"`
	IssueDate  time.Time `json:"issueDate" binding:"required"`
	PackageID  uint   `json:"packageId" binding:"required"`
}

type UpdateCertificationRequest struct {
	Status *string `json:"status" binding:"omitempty"`
}

type CertificationResponse struct {
	ID         uint      `json:"id"`
	Type       string    `json:"type"`
	Recipient  string    `json:"recipient"`
	IssueDate  time.Time `json:"issueDate"`
	PackageID  uint      `json:"packageId"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type CertificationListResponse struct {
	Data       []CertificationResponse `json:"data"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
	TotalPages int                    `json:"totalPages"`
}

// UserCertification DTOs

type IssueCertificationRequest struct {
	UserID         uint `json:"userId" binding:"required"`
	CertificationID uint `json:"certificationId" binding:"required"`
}

type UserCertificationResponse struct {
	UserID          uint      `json:"userId"`
	CertificationID uint      `json:"certificationId"`
	IssuedAt        time.Time `json:"issuedAt"`
	Certification   CertificationResponse `json:"certification"`
}

// Common list params
type ListParams struct {
	Page  int `form:"page,default=1"`
	Limit int `form:"limit,default=10"`
}