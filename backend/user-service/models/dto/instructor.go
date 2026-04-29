package dto

import (
	"time"

	"github.com/google/uuid"
)

type InstructorProfileRequest struct {
	Description       *string  `json:"description"       binding:"omitempty"`
	LicenseNumber     *string  `json:"licenseNumber"     binding:"omitempty,min=5"`
	LicenseExpiry     time.Time  `json:"licenseExpiry"     binding:"omitempty"`
	BNSPCertificateNumber *string `json:"bnspCertificateNumber" binding:"omitempty,min=10"`
	YearsOfExperience *int     `json:"yearsOfExperience" binding:"omitempty,min=0"`

}

type InstructorProfileResponse struct {
	UserID            uuid.UUID                    `json:"userId"`
	BNSPCertificateNumber string                   `json:"bnspCertificateNumber"`
	NumberOfStudents  int                      `json:"numberOfStudents"`
	SessionsCompleted int                      `json:"sessionsCompleted"`
	AverageRating     float64                  `json:"averageRating"`
	Description       string                   `json:"description"`
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

type InstructorListResponse struct {
	Data       []UserWithProfileResponse `json:"data"`
	Total      int64                     `json:"total"`
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
	TotalPages int                       `json:"totalPages"`
}

type CoverageAreaResponse struct {
	ID           uint   `json:"id"`
	InstructorID uint   `json:"instructorId"`
	AreaName     string `json:"areaName"`
}
 
// AddCoverageAreaInput is used for POST /instructors/:id/coverage-areas
type AddCoverageAreaInput struct {
	AreaName string `json:"areaName" binding:"required,min=2"`
}
 
// RemoveCoverageAreaInput is used for DELETE /instructors/:id/coverage-areas/:areaId
// No body needed — IDs come from path params, this is here for documentation clarity
type RemoveCoverageAreaInput struct {
	InstructorID uint `json:"-"`
	AreaID       uint `json:"-"`
}