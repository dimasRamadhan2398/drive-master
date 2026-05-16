package dto

import (
	"time"

	"github.com/google/uuid"
)

// ========== CAR DTOs ==========

type TransmissionType string

const (
	TransmissionManual    TransmissionType = "manual"
	TransmissionAutomatic TransmissionType = "automatic"
)

type CarStatus string

const (
	CarStatusAvailable   CarStatus = "available"
	CarStatusInUse       CarStatus = "in_use"
	CarStatusMaintenance CarStatus = "maintenance"
)

type CreateCarRequest struct {
	Brand        string           `json:"brand" binding:"required"`
	Model        string           `json:"model" binding:"required"`
	Year         int              `json:"year" binding:"required"`
	LicensePlate string           `json:"licensePlate" binding:"required"`
	Color        string           `json:"color"`
	Transmission TransmissionType `json:"transmission"`
	ImageURL     string           `json:"imageUrl"`
	Notes        string           `json:"notes"`
}

type UpdateCarRequest struct {
	Brand        string           `json:"brand"`
	Model        string           `json:"model"`
	Year         int              `json:"year"`
	Color        string           `json:"color"`
	Transmission TransmissionType `json:"transmission"`
	Status       CarStatus        `json:"status"`
	Mileage      int              `json:"mileage"`
	ImageURL     string           `json:"imageUrl"`
	Notes        string           `json:"notes"`
}

type CarResponse struct {
	ID           uuid.UUID        `json:"id"`
	Brand        string           `json:"brand"`
	Model        string           `json:"model"`
	Year         int              `json:"year"`
	LicensePlate string           `json:"licensePlate"`
	Color        string           `json:"color"`
	Transmission TransmissionType `json:"transmission"`
	Status       CarStatus        `json:"status"`
	Mileage      int              `json:"mileage"`
	ImageURL     string           `json:"imageUrl"`
	Notes        string           `json:"notes"`
	CreatedAt    time.Time        `json:"createdAt"`
	UpdatedAt    time.Time        `json:"updatedAt"`
}

type CarListResponse struct {
	Cars  []CarResponse `json:"cars"`
	Total int64         `json:"total"`
}

// ========== PACKAGE DTOs ==========

type PackageType string

const (
	PackageTypeBronze    PackageType = "bronze"
	PackageTypeSilver   PackageType = "silver"
	PackageTypeGold     PackageType = "gold"
	PackageTypePlatinum PackageType = "platinum"
)

type PackageStatus string

const (
	PackageStatusActive   PackageStatus = "active"
	PackageStatusInactive PackageStatus = "inactive"
)

type CreateBenefitRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	SortOrder   int    `json:"sortOrder"`
}

type CreatePackageRequest struct {
	Name            string                `json:"name" binding:"required"`
	Description     string                `json:"description"`
	PackageType     PackageType          `json:"packageType" binding:"required"`
	Price           float64              `json:"price" binding:"required"`
	DiscountPrice   float64              `json:"discountPrice"`
	DurationMinutes int                   `json:"durationMinutes"`
	TotalSessions   int                   `json:"totalSessions"`
	ImageURL        string               `json:"imageUrl"`
	Benefits        []CreateBenefitRequest `json:"benefits"`
}

type UpdatePackageRequest struct {
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	PackageType     PackageType   `json:"packageType"`
	Price           float64        `json:"price"`
	DiscountPrice   float64        `json:"discountPrice"`
	DurationMinutes int            `json:"durationMinutes"`
	TotalSessions   int            `json:"totalSessions"`
	Status          PackageStatus  `json:"status"`
	ImageURL        string         `json:"imageUrl"`
}

type BenefitResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	SortOrder   int       `json:"sortOrder"`
}

type PackageResponse struct {
	ID              uuid.UUID        `json:"id"`
	Name            string           `json:"name"`
	Description     string           `json:"description"`
	PackageType     PackageType     `json:"packageType"`
	Price           float64         `json:"price"`
	DiscountPrice   float64         `json:"discountPrice"`
	DurationMinutes int              `json:"durationMinutes"`
	TotalSessions   int              `json:"totalSessions"`
	Status          PackageStatus   `json:"status"`
	ImageURL        string           `json:"imageUrl"`
	Benefits        []BenefitResponse `json:"benefits,omitempty"`
	CreatedAt       time.Time        `json:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt"`
}

type PackageListResponse struct {
	Packages []PackageResponse `json:"packages"`
	Total    int64             `json:"total"`
}

// ========== REGION DTOs ==========

type ProvinceResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type RegencyResponse struct {
	ID         uint   `json:"id"`
	ProvinceID uint   `json:"provinceId"`
	Name       string `json:"name"`
	Type       string `json:"type"`
}

type DistrictResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	RegencyID uint   `json:"regencyId"`
}