package models

import (
	"time"
)

// BookingStatus represents the lifecycle of a booking
type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCompleted BookingStatus = "completed"
	BookingStatusCancelled BookingStatus = "cancelled"
)

// CertificationStatus represents certification lifecycle
type CertificationStatus string

const (
	CertificationStatusPending CertificationStatus = "pending"
	CertificationStatusIssued  CertificationStatus = "issued"
	CertificationStatusRevoked  CertificationStatus = "revoked"
)

// Booking owns the scheduling data.
// Cross-service IDs (UserID, InstructorID, CarID) are stored as plain integers.
// They are resolved at runtime via HTTP calls to user-service and catalog-service.
type Booking struct {
	ID             uint          `json:"id" gorm:"primaryKey"`
	UserID         uint          `json:"userId" gorm:"not null;index"`        // ref: user-service
	InstructorID   uint          `json:"instructorId" gorm:"not null;index"`  // ref: user-service (instructor_profiles.instructor_id)
	EntitlementID  uint          `json:"entitlementId" gorm:"not null;index"`
	DateOfSession  time.Time     `json:"dateOfSession" gorm:"not null"`
	FromTime       time.Time     `json:"fromTime" gorm:"not null"`
	ToTime         time.Time     `json:"toTime" gorm:"not null"`
	CarID          uint          `json:"carId" gorm:"not null;index"`         // ref: catalog-service
	Area           string        `json:"area" gorm:"size:150"`
	Notes          string        `json:"notes" gorm:"type:text"`
	Status         BookingStatus `json:"status" gorm:"type:varchar(30);default:'pending'"`
	CreatedAt      time.Time     `json:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt"`

	// Local association — entitlement lives in this service
	Entitlement UserEntitlement `json:"entitlement" gorm:"foreignKey:EntitlementID"`
}

// Session is the completed record of a booking.
// Created after a booking is marked completed.
type Session struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"userId" gorm:"not null;index"`       // ref: user-service
	InstructorID  uint      `json:"instructorId" gorm:"not null;index"` // ref: user-service
	DateOfSession time.Time `json:"dateOfSession" gorm:"not null"`
	Duration      int       `json:"duration"` // in minutes
	CarID         uint      `json:"carId" gorm:"not null;index"`        // ref: catalog-service
	Area          string    `json:"area" gorm:"size:150"`
	Notes         string    `json:"notes" gorm:"type:text"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// UserEntitlement tracks how many sessions a user has remaining.
// source_type can be "package" or "voucher"; source_id is the ID in that service.
type UserEntitlement struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	UserID            uint      `json:"userId" gorm:"not null;index"`        // ref: user-service
	SourceType        string    `json:"sourceType" gorm:"size:50;not null"`  // "package" | "voucher"
	SourceID          string    `json:"sourceId" gorm:"size:100;not null"`   // ID in catalog/voucher service
	TotalSessions     int       `json:"totalSessions"`
	SessionsRemaining int       `json:"sessionsRemaining"`
	ExpiresAt         time.Time `json:"expiresAt"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

// Certification is issued to a user after completing a package.
// PackageID is a reference to catalog-service; UserID to user-service.
type Certification struct {
	ID        uint                `json:"id" gorm:"primaryKey"`
	Type      string              `json:"type" gorm:"size:100;not null"`
	Recipient string              `json:"recipient" gorm:"size:150;not null"`
	IssueDate time.Time           `json:"issueDate"`
	PackageID uint                `json:"packageId" gorm:"not null;index"` // ref: catalog-service
	Status    CertificationStatus `json:"status" gorm:"type:varchar(30);default:'pending'"`
	CreatedAt time.Time           `json:"createdAt"`
	UpdatedAt time.Time           `json:"updatedAt"`
}

// UserCertification is the join table between users and certifications.
// UserID references user-service; CertificationID is local.
type UserCertification struct {
	UserID          uint      `json:"userId" gorm:"primaryKey"`           // ref: user-service
	CertificationID uint      `json:"certificationId" gorm:"primaryKey"`   // local FK
	IssuedAt        time.Time `json:"issuedAt"`

	Certification Certification `json:"certification" gorm:"foreignKey:CertificationID"`
}