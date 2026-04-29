package dto

import "time"

type CreateWorkExperienceRequest struct {
	CompanyName string    `json:"companyName" binding:"required"`
	Role        string    `json:"role" binding:"required"`
	StartDate   time.Time `json:"startDate" binding:"required"`
	EndDate     *time.Time `json:"endDate"`
	Description string    `json:"description"`
}

type UpdateWorkExperienceRequest struct {
    CompanyName    *string    `json:"companyName"    binding:"omitempty,min=2"`
    Role           *string    `json:"role"           binding:"omitempty,min=2"`
    StartDate      *time.Time `json:"startDate"      binding:"omitempty"`
    Description    *string    `json:"description"    binding:"omitempty"`

    // EndDate uses explicit flag pattern
    EndDate        *time.Time `json:"endDate"`        // the value
    ClearEndDate   *bool      `json:"clearEndDate"` 
}

type CreateWorkExperiencesRequest struct {
	WorkExperiences []CreateWorkExperienceRequest `json:"workExperiences" binding:"required,min=1,dive"`
}

type WorkExperienceResponse struct {
	ID           uint       `json:"id"`
	InstructorID uint       `json:"instructorId"`
	CompanyName  string     `json:"companyName"`
	Role         string     `json:"role"`
	StartDate    time.Time  `json:"startDate"`
	EndDate      *time.Time `json:"endDate"` // nil means currently working here
	Description  string     `json:"description"`
	IsVerified   bool       `json:"isVerified"`
}
