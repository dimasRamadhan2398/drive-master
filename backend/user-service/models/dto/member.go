package dto

import (
	"user-service/models"

	"github.com/google/uuid"
)

type MemberProfileResponse struct {
	UserID                 uuid.UUID            `json:"userId"`
	SessionsCompleted      int                  `json:"sessionsCompleted"`
	TrainingTime           int                  `json:"trainingTime"` // in minutes
	AverageRating          float64              `json:"averageRating"`
	TotalAvailableSessions int                  `json:"totalAvailableSessions"`
	Entitlements           []models.Entitlement `json:"entitlements"`
}

type MemberListResponse = PagedData[UserWithProfileResponse]
