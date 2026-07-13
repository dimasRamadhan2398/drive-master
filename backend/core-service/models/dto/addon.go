package dto

import (
	"core-service/models"
	"time"

	"github.com/google/uuid"
)

// ========== ADDON DTOs ==========

// CreateAddOnRequest represents the request body for creating an add-on
type CreateAddOnRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Price       float64 `json:"price" binding:"required"`
	Sessions    int    `json:"sessions"`
	ImageURL    string `json:"imageUrl"`
	SortOrder   int    `json:"sortOrder"`
}

type UpdateAddOnRequest struct {
	Title       string             `json:"title" binding:"omitempty,min=1,max=255"`
	Description string             `json:"description" binding:"max=2000"`
	Price       float64            `json:"price" binding:"omitempty,gte=0"`
	Sessions    int                `json:"sessions" binding:"omitempty,gte=0"`
	Status      models.AddOnStatus `json:"status" binding:"omitempty,oneof=active inactive"`
	ImageURL    string             `json:"imageUrl" binding:"omitempty,url,max=500"`
	SortOrder   int                `json:"sortOrder"`
}

type AddOnResponse struct {
	ID          uuid.UUID          `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Price       float64            `json:"price"`
	Sessions    int                `json:"sessions"`
	Status      models.AddOnStatus `json:"status"`
	ImageURL    string             `json:"imageUrl"`
	SortOrder   int                `json:"sortOrder"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

type AddOnListResponse struct {
	AddOns []AddOnResponse `json:"addOns"`
	Total  int64           `json:"total"`
}

type AddOnSelection struct {
	AddOnID uuid.UUID `json:"addOnId" binding:"required"`
	Quantity int      `json:"quantity" binding:"required,min=1"`
}

func (c *CreateAddOnRequest) SetDefaults() {
	if c.Sessions == 0 {
		c.Sessions = 1 // Default to 1 session
	}
}
