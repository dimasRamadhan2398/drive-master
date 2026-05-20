package models

import (
	"time"

	"github.com/google/uuid"
)

// Role represents the roles table
type Role struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:50;not null;uniqueIndex"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// User represents the users table
type User struct {
	ID            uuid.UUID        `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	FirstName    string    			`json:"firstName" gorm:"required,min=2"`
	LastName     string    			`json:"lastName" gorm:"required,min=2"`
	Username      string           `json:"username" gorm:"size:120;not null;uniqueIndex"`
	PasswordHash  string           `json:"-" gorm:"size:255;not null"`
	EmailAddress         string           `json:"email_address" gorm:"size:190;not null;uniqueIndex"`
	PhoneNumber   string           `json:"phoneNumber" gorm:"size:20"`
	Image         string           `json:"image" gorm:"size:500"`
	DateOfBirth   time.Time        `json:"dateOfBirth" gorm:"type:date"`
	Address       string           `json:"address" gorm:"size:255"`
	IsActive      bool             `json:"isActive" gorm:"default:true"`
	IsVerified    bool             `json:"isVerified" gorm:"default:false"`
	RoleID        uint             `json:"roleId" gorm:"not null"`
	Role          Role             `json:"role" gorm:"foreignKey:RoleID"`
	MemberProfile *MemberProfile    `json:"memberProfile,omitempty" gorm:"foreignKey:UserID"`
	InstructorProfile *InstructorProfile `json:"instructorProfile,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt     time.Time        `json:"createdAt"`
	UpdatedAt     time.Time        `json:"updatedAt"`
}

type UserSession struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID       uuid.UUID  `json:"userId" gorm:"type:uuid;not null;index"`
	RefreshToken string     `json:"-" gorm:"size:500;not null;uniqueIndex"`
	DeviceInfo   string     `json:"deviceInfo" gorm:"size:255"`
	IPAddress    string     `json:"ipAddress" gorm:"size:50"`
	ExpiresAt    time.Time  `json:"expiresAt" gorm:"not null"`
	LastUsedAt   time.Time  `json:"lastUsedAt"`
	CreatedAt    time.Time  `json:"createdAt"`
}

func (s *UserSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// MemberProfile represents the member_profiles table
type MemberProfile struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID            uuid.UUID `json:"userId" gorm:"type:uuid;not null;uniqueIndex"`
	SessionsCompleted int       `json:"sessionsCompleted" gorm:"default:0"`
	TrainingTime      int       `json:"trainingTime" gorm:"default:0"` // in minutes
	AverageRating     float64   `json:"averageRating" gorm:"default:0"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}


// InstructorProfile represents the instructor_profiles table
type InstructorProfile struct {
	ID                       uuid.UUID        `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID                   uuid.UUID        `json:"userId" gorm:"type:uuid;not null;uniqueIndex"`
	LicenseNumber            string           `json:"licenseNumber" gorm:"size:50"`
	BNSPCertificateNumber    string           `json:"bnspCertificateNumber" gorm:"size:50"`
	YearsOfExperience        int              `json:"yearsOfExperience" gorm:"default:0"`
	Bio                      string           `json:"bio" gorm:"type:text"`
	LicenseExpiry            time.Time        `json:"licenseExpiry"`
	PhotoURL                 string           `json:"photoURL" gorm:"size:500"`
	IsActive                 bool             `json:"isActive" gorm:"default:true"`
	NumberOfStudents         int              `json:"numberOfStudents" gorm:"default:0"`
	SessionsCompleted        int              `json:"sessionsCompleted" gorm:"default:0"`
	AverageRating            float64          `json:"averageRating" gorm:"default:0"`
	WorkExperiences          []WorkExperience `json:"workExperiences" gorm:"foreignKey:InstructorID"`
	CreatedAt                time.Time        `json:"createdAt"`
	UpdatedAt                time.Time        `json:"updatedAt"`
}

// WorkExperience represents the work_experiences table
type WorkExperience struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	InstructorID uuid.UUID  `json:"instructorId" gorm:"type:uuid;not null"`
	CompanyName  string     `json:"companyName" gorm:"size:255;not null"`
	Role         string     `json:"role" gorm:"size:100;not null"`
	StartDate    time.Time  `json:"startDate" gorm:"not null"`
	EndDate      *time.Time `json:"endDate"`
	Description  string     `json:"description" gorm:"type:text"`
	IsVerified   bool       `json:"isVerified" gorm:"default:false"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

// AreaType defines the type of coverage area
type AreaType string

const (
	AreaTypeProvince AreaType = "province"
	AreaTypeRegency  AreaType = "regency"
	AreaTypeDistrict AreaType = "district"
)

// InstructorArea represents the instructor_areas table
// Links an instructor to a specific region from core-service
type InstructorArea struct {
	InstructorID uuid.UUID `json:"instructorId" gorm:"type:uuid;not null"`
	AreaType     AreaType  `json:"areaType" gorm:"type:varchar(20);not null"` // province, regency, district
	AreaID       uint      `json:"areaId" gorm:"not null"`                   // ID from core-service (province/regency/district ID)
}

// InstructorAreaWithDetails includes region details for display
type InstructorAreaWithDetails struct {
	InstructorID uuid.UUID `json:"instructorId"`
	AreaType     AreaType  `json:"areaType"`
	AreaID       uint      `json:"areaId"`
	AreaName     string    `json:"areaName"`
}

// CreateUserInput is used internally by services
type CreateUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserCreatedEvent is used for publishing to Kafka
type UserCreatedEvent struct {
	UserID uuid.UUID   `json:"userId"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

// CertificationStatus represents the status of a certification
type CertificationStatus string

const (
	CertificationStatusPending   CertificationStatus = "pending"
	CertificationStatusVerified  CertificationStatus = "verified"
	CertificationStatusExpired   CertificationStatus = "expired"
	CertificationStatusRevoked   CertificationStatus = "revoked"
)

// Certification represents an instructor's certification (e.g., BNSP certificate)
type Certification struct {
	ID              uuid.UUID            `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	InstructorID    uuid.UUID            `json:"instructorId" gorm:"type:uuid;not null;index"`
	CertType        string               `json:"certType" gorm:"size:50;not null"`        // e.g., "BNSP", "SIM", "AWS"
	CertNumber      string               `json:"certNumber" gorm:"size:100;not null"`     // Certificate number
	IssuedBy        string               `json:"issuedBy" gorm:"size:255"`               // Issuing authority
	IssuedDate      time.Time            `json:"issuedDate" gorm:"not null"`              // Date of issue
	ExpiryDate      *time.Time           `json:"expiryDate"`                             // Expiration date (nullable)
	Status          CertificationStatus `json:"status" gorm:"type:varchar(20);default:'pending'"`
	DocumentURL     string               `json:"documentUrl" gorm:"size:500"`            // URL to scanned document
	Notes           string               `json:"notes" gorm:"type:text"`                 // Additional notes
	VerifiedAt      *time.Time           `json:"verifiedAt"`                             // When it was verified
	VerifiedBy      *uuid.UUID           `json:"verifiedBy"`                             // Who verified it
	CreatedAt       time.Time            `json:"createdAt"`
	UpdatedAt       time.Time            `json:"updatedAt"`
}

// EntitlementStatus represents the status of an entitlement
type EntitlementStatus string

const (
	EntitlementStatusActive   EntitlementStatus = "active"
	EntitlementStatusUsed     EntitlementStatus = "used"
	EntitlementStatusExpired  EntitlementStatus = "expired"
	EntitlementStatusRevoked  EntitlementStatus = "revoked"
)

// Entitlement represents a member's entitlement to training sessions
// Created when a member purchases a package from booking-service
type Entitlement struct {
	ID              uuid.UUID           `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	MemberID        uuid.UUID           `json:"memberId" gorm:"type:uuid;not null;index"`
	BookingID       uuid.UUID           `json:"bookingId" gorm:"type:uuid;index"`           // Reference to booking that created this
	PackageID       uuid.UUID           `json:"packageId" gorm:"type:uuid"`                // Package ID from core-service
	PackageName     string              `json:"packageName" gorm:"size:255"`               // Package name (denormalized)
	TotalSessions   int                 `json:"totalSessions" gorm:"default:0"`             // Total sessions in package
	Remaining       int                 `json:"remaining" gorm:"default:0"`                 // Remaining sessions
	UsedSessions    int                 `json:"usedSessions" gorm:"default:0"`              // Sessions already used
	StartDate       time.Time           `json:"startDate"`                                 // When entitlement becomes active
	EndDate         *time.Time          `json:"endDate"`                                  // Expiration date
	Status          EntitlementStatus   `json:"status" gorm:"type:varchar(20);default:'active'"`
	CreatedAt       time.Time           `json:"createdAt"`
	UpdatedAt       time.Time           `json:"updatedAt"`
}