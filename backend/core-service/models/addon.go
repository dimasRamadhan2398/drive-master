package models

import (
	"time"

	"github.com/google/uuid"
)

// AddOnStatus represents the status of an add-on
type AddOnStatus string

const (
	AddOnStatusActive   AddOnStatus = "active"
	AddOnStatusInactive AddOnStatus = "inactive"
)

// AddOn represents additional purchasable items (e.g., extra sessions, additional features)
type AddOn struct {
	ID          uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title       string       `json:"title" gorm:"size:255;not null"`         // e.g., "Extra Session", "Priority Support"
	Description string       `json:"description" gorm:"type:text"`           // Detailed description
	Price       float64      `json:"price" gorm:"type:decimal(10,2);not null"` // Add-on price
	Sessions    int          `json:"sessions" gorm:"default:1"`               // Number of sessions included (for session-based add-ons)
	Status      AddOnStatus  `json:"status" gorm:"size:20;not null;default:'active'"`
	ImageURL    string       `json:"imageUrl" gorm:"size:500"`               // Image URL
	SortOrder   int          `json:"sortOrder" gorm:"default:0"`            // Display order
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}
