package models

import (
	"time"

	"github.com/google/uuid"
)

type TestimonialStatus string

const (
	TestimonialStatusDraft     TestimonialStatus = "draft"
	TestimonialStatusPublished TestimonialStatus = "published"
	TestimonialStatusArchived  TestimonialStatus = "archived"
	TestimonialStatusPending   TestimonialStatus = "pending"
)

type Testimonial struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`

	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	UserName  string    `json:"user_name" gorm:"size:150;not null"`
	UserImage string    `json:"user_image" gorm:"size:255"`
	UserRole  string    `json:"user_role" gorm:"size:50"`

	Content string  `json:"content" gorm:"type:text;not null"`
	Rating  float64 `json:"rating" gorm:"type:decimal(2,1);check:rating >= 1 AND rating <= 5"`
	Tags    string  `json:"tags" gorm:"size:255"`

	Status     TestimonialStatus `json:"status" gorm:"type:varchar(20);default:'draft'"`
	IsFeatured bool              `json:"is_featured" gorm:"default:false"`
	AddedBy    uuid.UUID         `json:"added_by" gorm:"type:uuid;not null"`
	AddedAt    time.Time         `json:"added_at"`

	SortOrder int `json:"sort_order" gorm:"default:0"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TestimonialMedia allows multiple photos per testimonial (optional)
type TestimonialMedia struct {
	ID            uuid.UUID `json:"id"             gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TestimonialID uuid.UUID `json:"testimonial_id" gorm:"type:uuid;not null;index"`
	URL           string    `json:"url"            gorm:"size:255;not null"` // ImageKit URL
	MediaType     string    `json:"media_type"    gorm:"size:20"`            // "image" | "video"
	SortOrder     int       `json:"sort_order"    gorm:"default:0"`
	CreatedAt     time.Time `json:"created_at"`
}
