package models

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// PackageType represents the type of package
type PackageType string

const (
	PackageTypeBronze		 	PackageType = "bronze"
	PackageTypeSilver        	PackageType = "silver"
	PackageTypeGold    			PackageType = "gold"
	PackageTypePlatinum          PackageType = "platinum"
)

// PackageStatus represents the status of a package
type PackageStatus string

const (
	PackageStatusActive   PackageStatus = "active"
	PackageStatusInactive PackageStatus = "inactive"
)

// StringArray is a custom type for PostgreSQL text arrays
type StringArray []string

// Value implements driver.Valuer for database serialization
func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return pq.Array(a).Value()
}

// Scan implements sql.Scanner for database deserialization
func (a *StringArray) Scan(src interface{}) error {
	return (*pq.StringArray)(a).Scan(src)
}

// Package represents the packages table
type Package struct {
	ID              uuid.UUID       `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name            string          `json:"name" gorm:"size:255;not null"`            // e.g., "Basic Package", "Premium Package"
	Description     string          `json:"description" gorm:"type:text"`            // Package description
	PackageType     PackageType     `json:"packageType" gorm:"size:50;not null"`      // Type of package
	Price           float64         `json:"price" gorm:"type:decimal(10,2);not null"` // Package price
	DiscountPrice   float64         `json:"discountPrice" gorm:"type:decimal(10,2);default:0"`
	DurationMinutes int             `json:"durationMinutes" gorm:"default:60"`        // Duration in minutes
	TotalSessions   int             `json:"totalSessions" gorm:"default:1"`           // Number of sessions included
	Features        StringArray     `json:"features" gorm:"type:text[]"`              // Array of feature strings
	Highlight       bool            `json:"highlight" gorm:"default:false"`           // Whether package is highlighted
	Status          PackageStatus   `json:"status" gorm:"size:20;not null;default:'active'"`
	ImageURL        string          `json:"imageUrl" gorm:"size:500"`                // Package image URL
	Benefits        []PackageBenefit `json:"benefits,omitempty" gorm:"-"`           // Excluded from auto-migrate, managed separately
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
}

// PackageBenefit represents the benefits included in a package
type PackageBenefit struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PackageID   uuid.UUID `json:"packageId" gorm:"type:uuid;not null"`
	Title       string    `json:"title" gorm:"size:255;not null"`  // e.g., "10 Sessions with Instructor"
	Description string    `json:"description" gorm:"size:500"`      // Detailed description
	Icon        string    `json:"icon" gorm:"size:100"`            // Icon name or URL
	SortOrder   int       `json:"sortOrder" gorm:"default:0"`       // Display order
	CreatedAt   time.Time `json:"createdAt"`
}