// content-service/models/testimonial.go
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
)

type Testimonial struct {
    ID          uint              `json:"id"          gorm:"primaryKey"`

    UserID      uuid.UUID         `json:"userId"      gorm:"type:uuid;not null;index"`
    UserName    string            `json:"userName"    gorm:"size:150;not null"` 
    UserImage   string            `json:"userImage"   gorm:"size:255"`          
    UserRole    string            `json:"userRole"    gorm:"size:50"`           

    Content     string            `json:"content"     gorm:"type:text;not null"`
    Rating      float64           `json:"rating"      gorm:"type:decimal(2,1);check:rating >= 1 AND rating <= 5"`
    Tags        string            `json:"tags"        gorm:"size:255"` 

    Status      TestimonialStatus `json:"status"      gorm:"type:varchar(20);default:'draft'"`
    IsFeatured  bool              `json:"isFeatured"  gorm:"default:false"` 
    AddedBy     uuid.UUID         `json:"addedBy"     gorm:"type:uuid;not null"`
    AddedAt     time.Time         `json:"addedAt"`

    SortOrder   int               `json:"sortOrder"   gorm:"default:0"`

    CreatedAt   time.Time         `json:"createdAt"`
    UpdatedAt   time.Time         `json:"updatedAt"`
}

// TestimonialMedia allows multiple photos per testimonial (optional)
type TestimonialMedia struct {
    ID            uint      `json:"id"           gorm:"primaryKey"`
    TestimonialID uint      `json:"testimonialId" gorm:"not null;index"`
    URL           string    `json:"url"           gorm:"size:255;not null"` // ImageKit URL
    MediaType     string    `json:"mediaType"     gorm:"size:20"`           // "image" | "video"
    SortOrder     int       `json:"sortOrder"     gorm:"default:0"`
    CreatedAt     time.Time `json:"createdAt"`
}