package dto

import (
	"time"

	"github.com/google/uuid"
)

type InstructorProfileRequest struct {
	Specialization       *string  `json:"specialization"       binding:"omitempty,max=255"`
	Description       *string  `json:"description"       binding:"omitempty"`
	LicenseNumber     *string  `json:"licenseNumber"     binding:"omitempty,min=5"`
	LicenseExpiry     time.Time  `json:"licenseExpiry"     binding:"omitempty"`
	BNSPCertificateNumber *string `json:"bnspCertificateNumber" binding:"omitempty,min=10"`
	YearsOfExperience *int     `json:"yearsOfExperience" binding:"omitempty,min=0"`
}

type InstructorProfileResponse struct {
	UserID            uuid.UUID                `json:"userId"`
	BNSPCertificateNumber string               `json:"bnspCertificateNumber"`
	NumberOfStudents  int                      `json:"numberOfStudents"`
	SessionsCompleted int                      `json:"sessionsCompleted"`
	AverageRating     float64                  `json:"averageRating"`
	Description       string                   `json:"description"`
	Specialization    string                   `json:"specialization"`
	LicenseNumber     string                   `json:"licenseNumber"`
	YearsOfExperience int                      `json:"yearsOfExperience"`
	LicenseExpiry     time.Time                `json:"licenseExpiry"`
	WorkExperiences   []WorkExperienceResponse `json:"workExperiences,omitempty"`
	IsActive          bool                     `json:"isActive"`
	PhotoURL          string                   `json:"photoURL"`
	Bio               string                   `json:"bio"`
}

type UpdateInstructorProfileInput struct {
	Description       *string  `json:"description"       binding:"omitempty"`
	LicenseNumber     *string  `json:"licenseNumber"     binding:"omitempty,min=5"`
	LicenseExpiry     *string  `json:"licenseExpiry"     binding:"omitempty"`
	BNSPCertificateNumber *string `json:"bnspCertificateNumber" binding:"omitempty,min=10"`
	YearsOfExperience *int     `json:"yearsOfExperience" binding:"omitempty,min=0"`
}

type InstructorListResponse = PagedData[UserWithProfileResponse]

// CoverageAreaResponse represents a coverage area in responses
type CoverageAreaResponse struct {
	InstructorID uuid.UUID `json:"instructorId"`
	AreaType     string    `json:"areaType"`
	AreaID       uint      `json:"areaId"`
	AreaName     string    `json:"areaName"`
}

// AddCoverageAreaInput is used for POST /instructors/:id/coverage-areas
type AddCoverageAreaInput struct {
	AreaType string `json:"areaType" binding:"required,oneof=province regency district"`
	AreaID   uint   `json:"areaId" binding:"required,min=1"`
}
 
// RemoveCoverageAreaInput is used for DELETE /instructors/:id/coverage-areas/:areaId
// No body needed — IDs come from path params, this is here for documentation clarity
type RemoveCoverageAreaInput struct {
	InstructorID uint `json:"-"`
	AreaID       uint `json:"-"`
}

// ============================================================
// Certification DTOs
// ============================================================

// CreateCertificationInput is used for POST /instructors/:id/certifications
type CreateCertificationInput struct {
	CertType    string `json:"certType" binding:"required,min=2,max=50"`
	CertNumber  string `json:"certNumber" binding:"required,min=3,max=100"`
	IssuedBy    string `json:"issuedBy" binding:"required,min=2,max=255"`
	IssuedDate  string `json:"issuedDate" binding:"required"` // Format: YYYY-MM-DD
	ExpiryDate  string `json:"expiryDate"`                    // Format: YYYY-MM-DD (optional)
	DocumentURL string `json:"documentUrl" binding:"omitempty,url"`
	Notes       string `json:"notes" binding:"omitempty,max=1000"`
}

// UpdateCertificationInput is used for PUT /instructors/:id/certifications/:certId
type UpdateCertificationInput struct {
	CertType    string `json:"certType" binding:"omitempty,min=2,max=50"`
	CertNumber  string `json:"certNumber" binding:"omitempty,min=3,max=100"`
	IssuedBy    string `json:"issuedBy" binding:"omitempty,min=2,max=255"`
	ExpiryDate  string `json:"expiryDate"`
	DocumentURL string `json:"documentUrl" binding:"omitempty,url"`
	Notes       string `json:"notes" binding:"omitempty,max=1000"`
}

// VerifyCertificationInput is used for POST /instructors/:id/certifications/:certId/verify
type VerifyCertificationInput struct {
	Notes string `json:"notes" binding:"omitempty,max=500"`
}

// IssueCertificationInput is used when issuing certification upon entitlement completion
type IssueCertificationInput struct {
	MemberID      uuid.UUID `json:"memberId"`
	EntitlementID uuid.UUID `json:"entitlementId"`
	PackageID     uuid.UUID `json:"packageId"`
	PackageName   string    `json:"packageName"`
	IssuedAt      time.Time `json:"issuedAt"`
}

// CertificationResponse represents a certification in API responses
type CertificationResponse struct {
	ID            uuid.UUID `json:"id"`
	InstructorID  uuid.UUID `json:"instructorId"`
	CertType      string    `json:"certType"`
	CertNumber    string    `json:"certNumber"`
	IssuedBy      string    `json:"issuedBy"`
	IssuedDate    string    `json:"issuedDate"`
	ExpiryDate    *string   `json:"expiryDate,omitempty"`
	Status        string    `json:"status"`
	DocumentURL   string    `json:"documentUrl,omitempty"`
	Notes         string    `json:"notes,omitempty"`
	VerifiedAt    *string   `json:"verifiedAt,omitempty"`
	CreatedAt     string    `json:"createdAt"`
	UpdatedAt     string    `json:"updatedAt"`
}

// CertificationListResponse represents a paginated list of certifications
type CertificationListResponse = PagedData[CertificationResponse]

// ============================================================
// Entitlement DTOs
// ============================================================

// CreateEntitlementInput is used for POST /members/:id/entitlements
type CreateEntitlementInput struct {
	BookingID     uuid.UUID `json:"bookingId" binding:"required"`
	PackageID    uuid.UUID `json:"packageId" binding:"required"`
	PackageName  string    `json:"packageName" binding:"required,min=2,max=255"`
	TotalSessions int      `json:"totalSessions" binding:"required,min=1"`
	StartDate     string   `json:"startDate" binding:"required"` // Format: YYYY-MM-DD
	EndDate       string   `json:"endDate"`                     // Format: YYYY-MM-DD (optional)
}

// UpdateEntitlementInput is used for PUT /members/:id/entitlements/:entId
type UpdateEntitlementInput struct {
	Remaining int    `json:"remaining" binding:"min=0"`
	EndDate   string `json:"endDate"`
	Status    string `json:"status" binding:"omitempty,oneof=active used expired revoked"`
}

// UseSessionInput is used for POST /members/:id/entitlements/:entId/use-session
type UseSessionInput struct {
	SessionNotes string `json:"sessionNotes" binding:"omitempty,max=500"`
}

// EntitlementResponse represents an entitlement in API responses
type EntitlementResponse struct {
	ID            uuid.UUID `json:"id"`
	MemberID      uuid.UUID `json:"memberId"`
	BookingID     uuid.UUID `json:"bookingId"`
	PackageID     uuid.UUID `json:"packageId"`
	PackageName   string    `json:"packageName"`
	TotalSessions int       `json:"totalSessions"`
	Remaining     int       `json:"remaining"`
	UsedSessions  int       `json:"usedSessions"`
	StartDate     string    `json:"startDate"`
	EndDate       *string   `json:"endDate,omitempty"`
	Status        string    `json:"status"`
	CreatedAt     string    `json:"createdAt"`
	UpdatedAt     string    `json:"updatedAt"`
}

// CreateInstructorWithUserRequest is used for POST /instructors/register
// Creates both a user and an instructor profile in a single transaction
type CreateInstructorWithUserRequest struct {
	// User fields
	FirstName    string `json:"firstName" binding:"required,min=2"`
	LastName     string `json:"lastName" binding:"required,min=2"`
	Username     string `json:"username" binding:"required,min=2"`
	Password     string `json:"password" binding:"required,min=8"`
	Email        string `json:"email" binding:"required,email"`
	PhoneNumber  string `json:"phoneNumber" binding:"required,min=10"`
	DateOfBirth  string `json:"dateOfBirth"` // Format: YYYY-MM-DD
	Address      string `json:"address"`

	// Instructor profile fields
	LicenseNumber          *string `json:"licenseNumber" binding:"omitempty,min=5"`
	LicenseExpiry          *string `json:"licenseExpiry"` // Format: DD/MM/YYYY
	BNSPCertificateNumber *string `json:"bnspCertificateNumber" binding:"omitempty,min=10"`
	YearsOfExperience      *int    `json:"yearsOfExperience" binding:"omitempty,min=0"`
	Specialization         *string `json:"specialization" binding:"omitempty,max=255"`
	Description            *string `json:"description"`
}

// CreateInstructorWithUserResponse returns the created user and instructor profile
type CreateInstructorWithUserResponse struct {
	UserID      uuid.UUID `json:"userId"`
    Email       string    `json:"email"`
    Username    string    `json:"username"`
    FirstName   string    `json:"firstName"`
    LastName    string    `json:"lastName"`
    PhoneNumber string    `json:"phoneNumber"`
    DateOfBirth string    `json:"dateOfBirth"`
    RoleID      uint      `json:"roleId"`
	Profile *InstructorProfileResponse `json:"instructorProfile"`
}

// EntitlementListResponse represents a paginated list of entitlements
type EntitlementListResponse = PagedData[EntitlementResponse]