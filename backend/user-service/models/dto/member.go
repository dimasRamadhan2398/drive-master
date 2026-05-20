package dto

import "github.com/google/uuid"

type MemberProfileResponse struct {
		UserID            uuid.UUID    `json:"userId"`
		SessionsCompleted int     `json:"sessionsCompleted"`
		TrainingTime      int     `json:"trainingTime"` // in minutes
		AverageRating     float64 `json:"averageRating"`
}