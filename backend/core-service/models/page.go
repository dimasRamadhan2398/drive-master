package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PageStatus string

const (
	PageStatusDraft     PageStatus = "draft"
	PageStatusPublished PageStatus = "published"
)

type Page struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title     string         `json:"title" gorm:"size:255;not null"`
	Slug      string         `json:"slug" gorm:"size:255;uniqueIndex:idx_pages_slug,where:deleted_at IS NULL;not null"`
	Status    PageStatus     `json:"status" gorm:"size:20;default:'draft'"`
	Sections  string         `json:"sections" gorm:"type:text"` // JSON-serialized array of sections
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}
