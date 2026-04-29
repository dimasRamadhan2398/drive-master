package dto

import "time"

type CreateWorkExperienceRequest struct {
	CompanyName string    `json:"companyName" binding:"required"`
	Role        string    `json:"role" binding:"required"`
	StartDate   time.Time `json:"startDate" binding:"required"`
	EndDate     *time.Time `json:"endDate"`
	Description string    `json:"description"`
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
