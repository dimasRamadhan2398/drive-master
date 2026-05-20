package models

import "time"

// UserCreatedEvent represents a user creation event
type UserCreatedEvent struct {
	UserID uint   `json:"userId"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

// ProcessedEvent represents a processed event stored in DB
type ProcessedEvent struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	EventType string    `json:"eventType" gorm:"size:80;index;not null"`
	Payload   string    `json:"payload" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"createdAt"`
}

// ========== CAR EVENTS ==========

type CarCreatedEvent struct {
	CarID       string `json:"carId"`
	Brand       string `json:"brand"`
	Model       string `json:"model"`
	Year        int    `json:"year"`
	LicensePlate string `json:"licensePlate"`
	Transmission string `json:"transmission"`
	CreatedAt   string `json:"createdAt"`
}

type CarUpdatedEvent struct {
	CarID  string `json:"carId"`
	Brand  string `json:"brand,omitempty"`
	Model  string `json:"model,omitempty"`
	Year   int    `json:"year,omitempty"`
	Status string `json:"status,omitempty"`
	UpdatedAt string `json:"updatedAt"`
}

type CarDeletedEvent struct {
	CarID string `json:"carId"`
}

// ========== PACKAGE EVENTS ==========

type PackageCreatedEvent struct {
	PackageID   string  `json:"packageId"`
	Name        string  `json:"name"`
	PackageType string  `json:"packageType"`
	Price       float64 `json:"price"`
	CreatedAt   string `json:"createdAt"`
}

type PackageUpdatedEvent struct {
	PackageID   string  `json:"packageId"`
	Name        string  `json:"name,omitempty"`
	PackageType string  `json:"packageType,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Status      string  `json:"status,omitempty"`
	UpdatedAt   string `json:"updatedAt"`
}

type PackageDeletedEvent struct {
	PackageID string `json:"packageId"`
}
