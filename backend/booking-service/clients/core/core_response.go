package core

import (
	"time"

	"github.com/google/uuid"
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

// PackageResponse represents a package from core-service
type PackageResponse struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Price           float64   `json:"price"`
	DiscountPrice   float64   `json:"discountPrice"`
	Sessions        int       `json:"sessions"`        // Number of sessions included
	TotalSessions   int       `json:"totalSessions"`   // JSON field returned by core-service
	Duration        int       `json:"duration"`        // Duration per session in minutes
	ValidityDays    int       `json:"validityDays"`   // Validity period in days
	IsActive        bool      `json:"isActive"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func (p *PackageResponse) GetSessions() int {
	if p.TotalSessions > 0 {
		return p.TotalSessions
	}
	if p.Sessions > 0 {
		return p.Sessions
	}
	return 0
}

// AddOnResponse represents an add-on from core-service
type AddOnResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Sessions    int       `json:"sessions"`     // Number of sessions included
	Status      string    `json:"status"`      // active or inactive
	ImageURL    string    `json:"imageUrl"`
	SortOrder   int       `json:"sortOrder"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
