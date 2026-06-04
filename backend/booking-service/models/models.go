package models

import (
	"time"

	"github.com/google/uuid"
)

// EnrollmentStatus represents the lifecycle of an enrollment (package purchase)
type EnrollmentStatus string

const (
	EnrollmentStatusPendingPayment EnrollmentStatus = "pending_payment"
	EnrollmentStatusPaid           EnrollmentStatus = "paid"
	EnrollmentStatusInProgress     EnrollmentStatus = "in_progress"
	EnrollmentStatusCompleted      EnrollmentStatus = "completed"
	EnrollmentStatusCancelled      EnrollmentStatus = "cancelled"
)

// BookingStatus represents the lifecycle of a booking (DEPRECATED: use EnrollmentStatus)
type BookingStatus string

// Booking is an alias for Enrollment for backward compatibility (DEPRECATED: use Enrollment)
type Booking = Enrollment

// Deprecated: Use EnrollmentStatus instead
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
	CertificationStatusRevoked CertificationStatus = "revoked"
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

// Enrollment represents a user's enrollment/purchase of a driving package.
// This is created when a user pays for a package (Bronze/Silver/Gold/Platinum).
type Enrollment struct {
	ID         uint             `json:"id" gorm:"primaryKey"`
	UserID     uint             `json:"userId" gorm:"not null;index"`     // ref: user-service
	PackageID  uint             `json:"packageId" gorm:"not null;index"`  // ref: core-service (package)
	Status     EnrollmentStatus `json:"status" gorm:"type:varchar(30);default:'pending_payment'"`
	TotalPrice float64          `json:"totalPrice"`                       // base price + add-ons
	PaidAt     *time.Time       `json:"paidAt"`                           // when payment was confirmed
	ExpiresAt  time.Time        `json:"expiresAt"`                       // when the enrollment expires (usually package validity)
	CreatedAt  time.Time        `json:"createdAt"`
	UpdatedAt  time.Time        `json:"updatedAt"`

	// Local associations
	Entitlements []UserEntitlement `json:"entitlements" gorm:"foreignKey:EnrollmentID"`
}

// UserEntitlement tracks how many sessions a user has remaining within an enrollment.
// A user can have multiple entitlements if they purchased add-ons.
type UserEntitlement struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	EnrollmentID  uint      `json:"enrollmentId" gorm:"not null;index"` // ref: Enrollment
	UserID        uint      `json:"userId" gorm:"not null;index"`       // ref: user-service
	SourceType    string    `json:"sourceType" gorm:"size:50;not null"` // "package" | "addon" | "voucher"
	SourceID      string    `json:"sourceId" gorm:"size:100;not null"` // ID in core-service
	TotalSessions int       `json:"totalSessions"`
	UsedSessions  int       `json:"usedSessions"`
	ExpiresAt     time.Time `json:"expiresAt"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// DrivingSession is the record of an actual driving lesson/session.
// Created when a session is scheduled and attended.
type DrivingSession struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	EnrollmentID  uint      `json:"enrollmentId" gorm:"not null;index"`  // ref: Enrollment
	EntitlementID uint      `json:"entitlementId" gorm:"not null;index"`  // ref: UserEntitlement
	UserID        uint      `json:"userId" gorm:"not null;index"`        // ref: user-service
	InstructorID  uint      `json:"instructorId" gorm:"not null;index"`  // ref: user-service
	CarID         uint      `json:"carId" gorm:"not null;index"`         // ref: core-service (car)
	ScheduleID    *uint     `json:"scheduleId" gorm:"index"`             // ref: Schedule (optional)
	Date          time.Time `json:"date" gorm:"type:date;not null"`
	Time          string    `json:"time" gorm:"size:10;not null"`   // HH:MM format
	Duration      int       `json:"duration" gorm:"default:60"`    // in minutes
	Status        string    `json:"status" gorm:"type:varchar(20);default:'scheduled'"` // scheduled | in_progress | completed | cancelled
	Area          string    `json:"area" gorm:"size:150"`
	Notes         string    `json:"notes" gorm:"type:text"`
	StartedAt     *time.Time `json:"startedAt"`  // when the session actually started
	CompletedAt   *time.Time `json:"completedAt"` // when the session was completed
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// Session is an alias for DrivingSession for backward compatibility (DEPRECATED)
type Session = DrivingSession

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
	UserID          uint      `json:"userId" gorm:"primaryKey"`           
	CertificationID uint      `json:"certificationId" gorm:"primaryKey"`  
	IssuedAt        time.Time `json:"issuedAt"`

	Certification Certification `json:"certification" gorm:"foreignKey:CertificationID"`
}

// Schedule represents a time slot in the scheduling calendar.
// It tracks availability for instructors and cars, and links to enrollments.
type Schedule struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Date         time.Time      `json:"date" gorm:"type:date;not null;index"`
	Time         string         `json:"time" gorm:"size:10;not null"` // HH:MM format
	Duration     int            `json:"duration" gorm:"default:90"`   // duration in minutes
	InstructorID uuid.UUID      `json:"instructorId" gorm:"type:uuid;not null;index"` // ref: user-service (UUID)
	CarID        uint           `json:"carId" gorm:"not null;index"`             // ref: core-service (car)
	UserID       *uint          `json:"userId" gorm:"index"`                    // ref: user-service (nullable, assigned when booked)
	EnrollmentID *uint          `json:"enrollmentId" gorm:"index"`             // ref: Enrollment (nullable)
	Status       ScheduleStatus `json:"status" gorm:"type:varchar(20);default:'available'"`
	Notes        string         `json:"notes" gorm:"type:text"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
}