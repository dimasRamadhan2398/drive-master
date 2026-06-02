package dto

import "github.com/google/uuid"

type MemberProfileResponse struct {
	UserID                 uuid.UUID `json:"userId"`
	SessionsCompleted      int       `json:"sessionsCompleted"`
	TrainingTime           int       `json:"trainingTime"` // in minutes
	AverageRating          float64   `json:"averageRating"`
	TotalAvailableSessions int       `json:"totalAvailableSessions"`
}

type MemberListResponse struct {
	Data       []UserWithProfileResponse `json:"data"`
	Total      int64                     `json:"total"`
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
	TotalPages int                       `json:"totalPages"`
}