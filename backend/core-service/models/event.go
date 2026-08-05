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
	CarID        string `json:"carId"`
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	Year         int    `json:"year"`
	LicensePlate string `json:"licensePlate"`
	Transmission string `json:"transmission"`
	CreatedAt    string `json:"createdAt"`
}

type CarUpdatedEvent struct {
	CarID     string `json:"carId"`
	Brand     string `json:"brand,omitempty"`
	Model     string `json:"model,omitempty"`
	Year      int    `json:"year,omitempty"`
	Status    string `json:"status,omitempty"`
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
	CreatedAt   string  `json:"createdAt"`
}

type PackageUpdatedEvent struct {
	PackageID   string  `json:"packageId"`
	Name        string  `json:"name,omitempty"`
	PackageType string  `json:"packageType,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Status      string  `json:"status,omitempty"`
	UpdatedAt   string  `json:"updatedAt"`
}

type PackageDeletedEvent struct {
	PackageID string `json:"packageId"`
}

type EnrollmentPaidEvent struct {
	EnrollmentID  string                 `json:"enrollment_id,omitempty"`
	UserID        string                 `json:"user_id,omitempty"`
	PackageID     string                 `json:"package_id,omitempty"`
	TotalPrice    float64                `json:"total_price,omitempty"`
	TotalSessions int                    `json:"total_sessions,omitempty"`
	PackageName   string                 `json:"package_name,omitempty"`
	Data          map[string]interface{} `json:"data,omitempty"`
}

// ========== ARTICLE EVENTS ==========

type ArticleCreatedEvent struct {
	ArticleID   string   `json:"articleId"`
	Title       string   `json:"title"`
	Slug        string   `json:"slug"`
	Status      string   `json:"status"`
	AuthorID    string   `json:"authorId"`
	IsFeatured  bool     `json:"isFeatured"`
	IsSpotlight bool     `json:"isSpotlight"`
	CreatedAt   string   `json:"createdAt"`
}

type ArticleUpdatedEvent struct {
	ArticleID   string   `json:"articleId"`
	Title       string   `json:"title,omitempty"`
	Slug        string   `json:"slug,omitempty"`
	Status      string   `json:"status,omitempty"`
	AuthorID    string   `json:"authorId,omitempty"`
	IsFeatured  *bool    `json:"isEnabled,omitempty"`
	IsSpotlight *bool    `json:"isSpotlight,omitempty"`
	UpdatedAt   string   `json:"updatedAt"`
}

type ArticleDeletedEvent struct {
	ArticleID string `json:"articleId"`
	DeletedAt string `json:"deletedAt"`
}

type ArticlePublishedEvent struct {
	ArticleID  string `json:"articleId"`
	Slug       string `json:"slug"`
	PublishedAt string `json:"publishedAt"`
}

type ArticleArchivedEvent struct {
	ArticleID string `json:"articleId"`
	ArchivedAt string `json:"archivedAt"`
}

// ========== USER EVENTS (consumed from user-service) ==========

type UserUpdatedEvent struct {
	UserID   uint   `json:"userId"`
	Email    string `json:"email,omitempty"`
	Name     string `json:"name,omitempty"`
	Phone    string `json:"phone,omitempty"`
	ImageURL string `json:"imageUrl,omitempty"`
	UpdatedAt string `json:"updatedAt"`
}

type UserDeletedEvent struct {
	UserID    uint   `json:"userId"`
	DeletedAt string `json:"deletedAt"`
}

// ========== REGION EVENTS ==========

type RegionProvinceUpdatedEvent struct {
	ProvinceID uint   `json:"provinceId"`
	Name       string `json:"name"`
	UpdatedAt  string `json:"updatedAt"`
}

type RegionRegencyUpdatedEvent struct {
	RegencyID  uint   `json:"regencyId"`
	ProvinceID uint   `json:"provinceId"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	UpdatedAt  string `json:"updatedAt"`
}

type RegionDistrictUpdatedEvent struct {
	DistrictID uint   `json:"districtId"`
	RegencyID  uint   `json:"regencyId"`
	Name       string `json:"name"`
	UpdatedAt  string `json:"updatedAt"`
}

// ========== SYSTEM EVENTS ==========

type CacheInvalidatedEvent struct {
	CacheKey    string `json:"cacheKey"`
	Pattern     string `json:"pattern,omitempty"`
	InvalidatedBy string `json:"invalidatedBy"`
	TTL         int    `json:"ttl"`
}
