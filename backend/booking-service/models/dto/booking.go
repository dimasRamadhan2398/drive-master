package dto

import (
	"time"
)

// Enrollment DTOs (represents package purchase/enrollment)

type CreateEnrollmentRequest struct {
	UserID    uint   `json:"userId" binding:"required"`
	PackageID uint   `json:"packageId" binding:"required"` // ref: core-service (package)
	AddOns    []uint `json:"addOns"`                      // optional add-ons (night driving, weekend, etc.)
}

type UpdateEnrollmentRequest struct {
	ExpiresAt *time.Time `json:"expiresAt" binding:"omitempty"`
	Status    *string    `json:"status" binding:"omitempty"`
}

type EnrollmentResponse struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"userId"`
	PackageID  uint      `json:"packageId"`
	Status     string    `json:"status"`
	TotalPrice float64   `json:"totalPrice"`
	PaidAt     *time.Time `json:"paidAt"`
	ExpiresAt  time.Time `json:"expiresAt"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type EnrollmentListResponse struct {
	Data       []EnrollmentResponse `json:"data"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	Limit      int                  `json:"limit"`
	TotalPages int                  `json:"totalPages"`
}

// Booking DTOs (DEPRECATED: use Enrollment instead)

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
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"totalPages"`
}

// DrivingSession DTOs

type CreateDrivingSessionRequest struct {
	EnrollmentID  uint      `json:"enrollmentId" binding:"required"`
	EntitlementID uint      `json:"entitlementId" binding:"required"`
	UserID        uint      `json:"userId" binding:"required"`
	InstructorID  uint      `json:"instructorId" binding:"required"`
	CarID         uint      `json:"carId" binding:"required"`
	ScheduleID    *uint     `json:"scheduleId"`
	Date          time.Time `json:"date" binding:"required"`
	Time          string    `json:"time" binding:"required"` // HH:MM format
	Duration      int       `json:"duration" binding:"required"`
	Area          string    `json:"area"`
	Notes         string    `json:"notes"`
}

type DrivingSessionResponse struct {
	ID            uint       `json:"id"`
	EnrollmentID  uint       `json:"enrollmentId"`
	EntitlementID uint       `json:"entitlementId"`
	UserID        uint       `json:"userId"`
	InstructorID  uint       `json:"instructorId"`
	CarID         uint       `json:"carId"`
	ScheduleID    *uint      `json:"scheduleId"`
	Date          string     `json:"date"`  // YYYY-MM-DD format
	Time          string     `json:"time"` // HH:MM format
	Duration      int        `json:"duration"`
	Status        string     `json:"status"`
	Area          string     `json:"area"`
	Notes         string     `json:"notes"`
	StartedAt     *time.Time `json:"startedAt"`
	CompletedAt   *time.Time `json:"completedAt"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

type DrivingSessionListResponse struct {
	Data       []DrivingSessionResponse `json:"data"`
	Total      int64                     `json:"total"`
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
	TotalPages int                       `json:"totalPages"`
}

// Session DTOs (alias for DrivingSession, DEPRECATED)

type CreateSessionRequest = CreateDrivingSessionRequest
type SessionResponse = DrivingSessionResponse
type SessionListResponse = DrivingSessionListResponse

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
type BookSlotRequest struct {
	UserID        uint   `json:"userId" binding:"required"`
	EntitlementID uint   `json:"entitlementId" binding:"required"`
	Notes         string `json:"notes"`
}

// Stats Response DTOs

type SessionStatsResponse struct {
	TotalSessions     int64 `json:"totalSessions"`
	ActiveSessions    int64 `json:"activeSessions"`
	CompletedSessions int64 `json:"completedSessions"`
	PendingSessions   int64 `json:"pendingSessions"`
}

type CertificationStatsResponse struct {
	TotalCertifications int64 `json:"totalCertifications"`
	IssuedCertifications int64 `json:"issuedCertifications"`
	ActiveCertifications int64 `json:"activeCertifications"`
	RevokedCertifications int64 `json:"revokedCertifications"`
}