package dto

import (
	"time"

	"github.com/google/uuid"
)

// FAQ DTOs

type CreateFAQRequest struct {
	Question string `json:"question" binding:"required,max=500"`
	Answer   string `json:"answer" binding:"required"`
	Category string `json:"category" binding:"max=100"`
	Order    int    `json:"order"`
}

type UpdateFAQRequest struct {
	Question string `json:"question" binding:"max=500"`
	Answer   string `json:"answer"`
	Category string `json:"category" binding:"max=100"`
	Order    int    `json:"order"`
	IsActive *bool  `json:"isActive"`
}

type FAQResponse struct {
	ID        uuid.UUID `json:"id"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	Order     int       `json:"order"`
	Category  string    `json:"category"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type FAQListResponse struct {
	FAQs  []FAQResponse `json:"faqs"`
	Total int64         `json:"total"`
}

type ReorderFAQRequest struct {
	NewOrder int `json:"newOrder" binding:"required,min=0"`
}