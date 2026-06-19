package core

import (
	"time"
)

// CarResponse represents a car from core-service
type CarResponse struct {
	ID           uint      `json:"id"`
	Brand        string    `json:"brand"`
	Model        string    `json:"model"`
	Year         int       `json:"year"`
	LicensePlate string    `json:"licensePlate"`
	IsActive     bool      `json:"isActive"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}


type CarInfo struct {
	ID    string `json:"id"`
	Brand string `json:"brand"`
	Model string `json:"model"`
}
