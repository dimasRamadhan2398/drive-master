package models

import (
	"time"

	"github.com/google/uuid"
)

// TransmissionType represents the type of transmission
type TransmissionType string

const (
	TransmissionManual    TransmissionType = "manual"
	TransmissionAutomatic TransmissionType = "automatic"
)

// CarStatus represents the current status of a car
type CarStatus string

const (
	CarStatusAvailable   CarStatus = "available"
	CarStatusInUse       CarStatus = "in_use"
	CarStatusMaintenance CarStatus = "maintenance"
)

// Car represents the cars table
type Car struct {
	ID           uuid.UUID        `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Brand        string           `json:"brand" gorm:"size:100;not null"`         // e.g., "Toyota", "Honda"
	Model        string           `json:"model" gorm:"size:100;not null"`         // e.g., "Vios", "Civic"
	Year         int              `json:"year" gorm:"not null"`                   // Manufacturing year
	LicensePlate string           `json:"licensePlate" gorm:"size:20;not null;uniqueIndex"`
	Color        string           `json:"color" gorm:"size:50"`                    // e.g., "Black", "White"
	Transmission TransmissionType `json:"transmission" gorm:"size:20;not null;default:'manual'"`
	Status       CarStatus        `json:"status" gorm:"size:20;not null;default:'available'"`
	Mileage      int              `json:"mileage" gorm:"default:0"`                // in kilometers
	ImageURL     string           `json:"imageUrl" gorm:"size:500"`               // URL to car image
	Notes        string           `json:"notes" gorm:"type:text"`                 // Additional notes
	CreatedAt    time.Time        `json:"createdAt"`
	UpdatedAt    time.Time        `json:"updatedAt"`
}