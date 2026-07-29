package models

import (
	"time"

	"booking-service/models/dto"

	"github.com/google/uuid"
)

// EnrollmentStatus represents the lifecycle of an enrollment (package purchase)
type EnrollmentStatus string

const (
	EnrollmentStatusPendingPayment EnrollmentStatus = "pending"
	EnrollmentStatusPaid           EnrollmentStatus = "paid"
	EnrollmentStatusInProgress     EnrollmentStatus = "in_progress"
	EnrollmentStatusCompleted      EnrollmentStatus = "completed"
	EnrollmentStatusCancelled      EnrollmentStatus = "cancelled"
)

// CertificationStatus represents certification lifecycle
type CertificationStatus string

const (
	CertificationStatusPending CertificationStatus = "pending"
	CertificationStatusIssued  CertificationStatus = "issued"
	CertificationStatusRevoked CertificationStatus = "revoked"
)

// ScheduleStatus is an alias for dto.ScheduleStatus for backward compatibility
// Deprecated: Use dto.ScheduleStatus instead
type ScheduleStatus = dto.ScheduleStatus

// Deprecated: Use dto.ScheduleStatus constants instead
const (
	ScheduleStatusAvailable  ScheduleStatus = dto.ScheduleStatusAvailable
	ScheduleStatusBooked     ScheduleStatus = dto.ScheduleStatusBooked
	ScheduleStatusInProgress ScheduleStatus = dto.ScheduleStatusInProgress
	ScheduleStatusCompleted  ScheduleStatus = dto.ScheduleStatusCompleted
	ScheduleStatusBlocked    ScheduleStatus = dto.ScheduleStatusBlocked
)

// Enrollment represents a user's enrollment/purchase of a driving package.
// This is created when a user pays for a package (Bronze/Silver/Gold/Platinum).
type Enrollment struct {
	ID         uuid.UUID       `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID     uuid.UUID       `json:"userId" gorm:"type:uuid;not null;index"` // ref: user-service
	PackageID  uuid.UUID       `json:"packageId" gorm:"type:uuid;not null;index"`  // ref: core-service (package)
	Status     EnrollmentStatus `json:"status" gorm:"type:varchar(30);default:'pending_payment'"`
	TotalPrice float64         `json:"totalPrice"`                         // base price + add-ons
	PaidAt     *time.Time      `json:"paidAt"`                             // when payment was confirmed
	ExpiresAt  time.Time       `json:"expiresAt"`                         // when the enrollment expires (usually package validity)
	AnonymizedAt *time.Time    `json:"anonymizedAt" gorm:"index"`          // when user was deleted
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt"`
}



// DrivingSession is the record of an actual driving lesson/session.
// Created when a session is scheduled and attended.
type DrivingSession struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	EnrollmentID  uuid.UUID `json:"enrollmentId" gorm:"type:uuid;not null;index"` // ref: Enrollment
	EntitlementID uuid.UUID `json:"entitlementId" gorm:"type:uuid;not null;index"` // ref: UserEntitlement
	UserID        uuid.UUID `json:"userId" gorm:"type:uuid;not null;index"`        // ref: user-service
	InstructorID  uuid.UUID `json:"instructorId" gorm:"type:uuid;not null;index"`  // ref: user-service
	CarID         uuid.UUID `json:"carId" gorm:"type:uuid;not null;index"`         // ref: core-service (car)
	ScheduleID    *uint     `json:"scheduleId" gorm:"index"`             // ref: Schedule (optional)
	Date          time.Time `json:"date" gorm:"type:date;not null"`
	Time          string    `json:"time" gorm:"size:10;not null"`   // HH:MM format
	Duration      int       `json:"duration" gorm:"default:60"`    // in minutes
	Status        string    `json:"status" gorm:"type:varchar(20);default:'scheduled'"` // scheduled | in_progress | completed | cancelled
	Area          string    `json:"area" gorm:"size:150"`
	Notes         string    `json:"notes" gorm:"type:text"`
	AnonymizedAt  *time.Time `json:"anonymizedAt" gorm:"index"`       // when user was deleted
	StartedAt     *time.Time `json:"startedAt"`  // when the session actually started
	CompletedAt   *time.Time `json:"completedAt"` // when the session was completed
	EndTime       *time.Time `json:"endTime"`     // explicit end time set by admin; empty until session ends
	IsEndedByAdmin bool      `json:"isEndedByAdmin" gorm:"default:false"` // true when admin force-completes
	Rating        *float64   `json:"rating" gorm:"type:decimal(2,1)"`    // star rating (1.0 to 5.0)
	Feedback      string     `json:"feedback" gorm:"type:text"`          // student review feedback
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// Session is an alias for DrivingSession for backward compatibility (DEPRECATED)
type Session = DrivingSession

// Certification is issued to a user after completing a package.
// PackageID is a reference to catalog-service; UserID to user-service.
type Certification struct {
	ID        uuid.UUID           `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Type      string              `json:"type" gorm:"size:100;not null"`
	Recipient string              `json:"recipient" gorm:"size:150;not null"`
	IssueDate time.Time           `json:"issueDate"`
	PackageID uuid.UUID           `json:"packageId" gorm:"type:uuid;not null;index"` // ref: catalog-service
	Status    CertificationStatus `json:"status" gorm:"type:varchar(30);default:'pending'"`
	CreatedAt time.Time           `json:"createdAt"`
	UpdatedAt time.Time           `json:"updatedAt"`
}

// UserCertification is the join table between users and certifications.
// UserID references user-service; CertificationID is local.
type UserCertification struct {
	UserID          uuid.UUID    `json:"userId" gorm:"type:uuid;primaryKey"`
	CertificationID uuid.UUID    `json:"certificationId" gorm:"type:uuid;primaryKey"`
	IssuedAt        time.Time    `json:"issuedAt"`

	Certification Certification `json:"certification" gorm:"foreignKey:CertificationID"`
}

// Schedule is an alias for dto.Schedule for backward compatibility
// Deprecated: Use dto.Schedule instead
type Schedule = dto.Schedule

// Payment represents a payment transaction for an enrollment
type Payment struct {
	ID            uint            `json:"id" gorm:"primaryKey"`
	EnrollmentID  uuid.UUID       `json:"enrollmentId" gorm:"type:uuid;not null;index"` // ref: Enrollment
	UserID        uuid.UUID      `json:"userId" gorm:"type:uuid;not null;index"`      // ref: user-service
	OrderID       string          `json:"orderId" gorm:"size:100;uniqueIndex"` // Midtrans order ID
	Amount        float64         `json:"amount" gorm:"not null"`
	PaymentMethod dto.PaymentMethod `json:"paymentMethod" gorm:"type:varchar(30)"`
	Status        dto.PaymentStatus `json:"status" gorm:"type:varchar(30);default:'pending'"`
	PaymentURL    string          `json:"paymentUrl" gorm:"type:text"`       // Snap redirect URL
	TransactionID string          `json:"transactionId" gorm:"size:100"`     // Midtrans transaction ID
	PaidAt        *time.Time      `json:"paidAt"`
	ExpiresAt     *time.Time      `json:"expiresAt"`
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
}
