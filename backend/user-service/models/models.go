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
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name          string    `json:"name" gorm:"size:120;not null"`
	Username      string    `json:"username" gorm:"size:120;not null;uniqueIndex"`
	PasswordHash  string    `json:"-" gorm:"size:255;not null"`
	Email         string    `json:"email" gorm:"size:190;not null;uniqueIndex"`
	EmailAddress  string    `json:"emailAddress" gorm:"size:190;not null;uniqueIndex"`
	PhoneNumber   string    `json:"phoneNumber" gorm:"size:20"`
	Image         string    `json:"image" gorm:"size:500"`
	DateOfBirth   time.Time `json:"dateOfBirth"`
	Address       string    `json:"address" gorm:"size:255"`
	IsActive      bool      `json:"isActive" gorm:"default:true"`
	RoleID        uint      `json:"roleId" gorm:"not null"`
	Role          Role      `json:"role" gorm:"foreignKey:RoleID"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type UserSession struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID       uuid.UUID  `json:"userId" gorm:"type:uuid;not null;index"`
	User         User       `json:"-" gorm:"foreignKey:UserID"`
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
	User              User      `json:"user" gorm:"foreignKey:UserID"`
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
	User                     User             `json:"user" gorm:"foreignKey:UserID"`
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
	ID           uint             `json:"id" gorm:"primaryKey"`
	InstructorID uuid.UUID        `json:"instructorId" gorm:"type:uuid;not null"`
	Instructor   InstructorProfile `json:"instructor" gorm:"foreignKey:InstructorID"`
	CompanyName  string           `json:"companyName" gorm:"size:255;not null"`
	Role         string           `json:"role" gorm:"size:100;not null"`
	StartDate    time.Time        `json:"startDate" gorm:"not null"`
	EndDate      *time.Time       `json:"endDate"`
	Description  string           `json:"description" gorm:"type:text"`
	IsVerified   bool             `json:"isVerified" gorm:"default:false"`
	CreatedAt    time.Time        `json:"createdAt"`
	UpdatedAt    time.Time        `json:"updatedAt"`
}

// InstructorArea represents the instructor_areas table
type InstructorArea struct {
	InstructorID uuid.UUID        `json:"instructorId" gorm:"type:uuid;not null"`
	AreaID 	 	uint		      `json:"areaId" gorm:"type:uint;not null"`
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