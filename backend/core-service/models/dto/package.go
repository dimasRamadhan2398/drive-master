package dto

import (
	"core-service/models"
	"time"

	"github.com/google/uuid"
)

// ========== PACKAGE DTOs ==========

// BenefitInput can be either a string (just title) or an object with full details
type BenefitInput interface{}

// CreatePackageRequest represents the request body for creating a package
type CreatePackageRequest struct {
	Name            string           `json:"name" binding:"required"`
	Description     string           `json:"description"`
	PackageType     models.PackageType `json:"packageType" binding:"required"`
	Price           float64          `json:"price" binding:"required"`
	DiscountPrice   float64          `json:"discountPrice"`
	DurationMinutes int              `json:"durationMinutes"`
	TotalSessions   int              `json:"totalSessions"`
	ImageURL        string           `json:"imageUrl"`
	Benefits        []BenefitInput   `json:"benefits"`
	IsDiscounted    bool             `json:"isDiscounted"`
	Highlight       bool             `json:"highlight"`
}

type CreateBenefitRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	SortOrder   int    `json:"sortOrder"`
}

type UpdatePackageRequest struct {
	Name            string              `json:"name"`
	Description     string              `json:"description"`
	PackageType     models.PackageType `json:"packageType"`
	Price           float64             `json:"price"`
	DiscountPrice   float64             `json:"discountPrice"`
	DurationMinutes int                 `json:"durationMinutes"`
	TotalSessions   int                 `json:"totalSessions"`
	Status          models.PackageStatus `json:"status"`
	ImageURL        string              `json:"imageUrl"`
	Highlight 		bool					`json:"highlight"`
}

type PackageResponse struct {
	ID              uuid.UUID              `json:"id"`
	Name            string                `json:"name"`
	Description     string                `json:"description"`
	PackageType     models.PackageType   `json:"packageType"`
	Price           float64              `json:"price"`
	DiscountPrice   float64              `json:"discountPrice"`
	DurationMinutes int                   `json:"durationMinutes"`
	TotalSessions   int                   `json:"totalSessions"`
	Status          models.PackageStatus  `json:"status"`
	ImageURL        string                `json:"imageUrl"`
	Benefits        []BenefitResponse     `json:"benefits,omitempty"`
	CreatedAt       time.Time             `json:"createdAt"`
	UpdatedAt       time.Time             `json:"updatedAt"`
}

type BenefitResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	SortOrder   int       `json:"sortOrder"`
}

type PackageListResponse struct {
	Packages []PackageResponse `json:"packages"`
	Total    int64             `json:"total"`
}

type StatusPackageRequest struct {
	Status  models.PackageStatus `json:"status"`
}

func (c *CreatePackageRequest) SetDefaults() {
	if c.DurationMinutes == 0 {
		c.DurationMinutes = 60 // Default to 1 hour
	}
	if c.TotalSessions == 0 {
		c.TotalSessions = 1 // Default to 1 session
	}
	if c.PackageType == "" {
		c.PackageType = models.PackageTypeBronze // Default package type
	}
	if c.IsDiscounted == false {
		c.IsDiscounted = false // Default is not discounted
	}
}