package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreatePageRequest struct {
	Title  string `json:"title" binding:"required"`
	Slug   string `json:"slug" binding:"required"`
	Status string `json:"status" binding:"required"`
}

type UpdatePageRequest struct {
	Title    *string `json:"title"`
	Slug     *string `json:"slug"`
	Status   *string `json:"status"`
	Sections *string `json:"sections"`
}

type PageResponse struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Status      string     `json:"status"`
	Sections    string     `json:"sections"`
	LastUpdated string     `json:"lastUpdated"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}
