package dto

import (
	"core-service/models"
	"time"

	"github.com/google/uuid"
)

// ========== CAR DTOs ==========

type CreateCarRequest struct {
	Brand        string                   `json:"brand" binding:"required"`
	Model        string                   `json:"model" binding:"required"`
	Year         int                      `json:"year" binding:"required"`
	Status       models.CarStatus        `json:"status"`        
	LicensePlate string                   `json:"licensePlate"`
	Color        string                   `json:"color"`
	Transmission models.TransmissionType  `json:"transmission"`
	ImageURL     string                   `json:"imageUrl"`
	Notes        string                   `json:"notes"`
}

type UpdateCarRequest struct {
	Brand        string                  `json:"brand"`
	Model        string                  `json:"model"`
	Year         int                     `json:"year"`
	Color        string                  `json:"color"`
	LicensePlate string                   `json:"licensePlate"`
	Transmission models.TransmissionType `json:"transmission"`
	Status       models.CarStatus        `json:"status"`
	Mileage      int                     `json:"mileage"`
	ImageURL     string                  `json:"imageUrl"`
	Notes        string                  `json:"notes"`
}

type CarResponse struct {
	ID           uuid.UUID               `json:"id"`
	Brand        string                  `json:"brand"`
	Model        string                  `json:"model"`
	Year         int                     `json:"year"`
	LicensePlate string                  `json:"licensePlate"`
	Color        string                  `json:"color"`
	Transmission models.TransmissionType `json:"transmission"`
	Status       models.CarStatus        `json:"status"`
	Mileage      int                     `json:"mileage"`
	ImageURL     string                  `json:"imageUrl"`
	Notes        string                  `json:"notes"`
	CreatedAt    time.Time               `json:"createdAt"`
	UpdatedAt    time.Time               `json:"updatedAt"`
}

type CarListResponse struct {
	Cars  []CarResponse `json:"cars"`
	Total int64         `json:"total"`
}