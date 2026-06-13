package dto

import (
	"time"

	"github.com/google/uuid"
)

// BlogPost DTOs

// Media represents media attached to a blog post
type BlogPostMedia struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Size      string    `json:"size"`
	URL       string    `json:"url"`
	FileType  string    `json:"fileType"` // "image" or "video"
	Order     int       `json:"order"`
	CreatedAt time.Time `json:"createdAt"`
}

// Publishing contains publishing-related fields
type Publishing struct {
	Status      string     `json:"status"` // "draft", "published", "archived"
	PublishedAt *time.Time `json:"publishedAt"`
	ScheduledAt *time.Time `json:"scheduledAt"`
}

// Attractiveness contains visibility/promotion-related fields
type Attractiveness struct {
	IsFeatured  bool `json:"isFeatured"`  // Show in featured section
	IsSpotlight bool `json:"isSpotlight"` // Show as spotlight/highlight
	Priority    int  `json:"priority"`    // Display priority (higher = more prominent)
	Highlight   bool `json:"highlight"`   // Mark as highlighted/popular
}

// CreateBlogPostRequest is the DTO for creating a new blog post
type CreateBlogPostRequest struct {
	Title   string `json:"title" binding:"required,max=255"`
	Slug    string `json:"slug" binding:"max=255"`
	Author  string `json:"author" binding:"max=100"`
	Content string `json:"content"`

	// Media
	Media []BlogPostMedia `json:"media"`

	// Publishing
	Publishing *Publishing `json:"publishing"`

	// Attractiveness
	Attractiveness *Attractiveness `json:"attractiveness"`
}

// UpdateBlogPostRequest is the DTO for updating an existing blog post
type UpdateBlogPostRequest struct {
	Title   string `json:"title" binding:"max=255"`
	Slug    string `json:"slug" binding:"max=255"`
	Author  string `json:"author" binding:"max=100"`
	Content string `json:"content"`

	// Media
	Media []BlogPostMedia `json:"media"`

	// Publishing
	Publishing *Publishing `json:"publishing"`

	// Attractiveness
	Attractiveness *Attractiveness `json:"attractiveness"`
}

// BlogPostResponse is the DTO for blog post responses
type BlogPostResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Author    string    `json:"author"`
	Content   string    `json:"content"`
	Excerpt   string    `json:"excerpt"`

	// Media
	Media []BlogPostMedia `json:"media"`

	// Publishing
	Publishing *Publishing `json:"publishing"`

	// Attractiveness
	Attractiveness *Attractiveness `json:"attractiveness"`

	// Metrics
	ViewCount   int64 `json:"viewCount"`
	LikeCount   int64 `json:"likeCount"`
	ShareCount  int64 `json:"shareCount"`
	ReadingTime int   `json:"readingTime"`

	// Timestamps
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// BlogPostBrief is a condensed response for lists
type BlogPostBrief struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Slug          string    `json:"slug"`
	Author        string    `json:"author"`
	Excerpt       string    `json:"excerpt"`
	FeaturedImage string    `json:"featuredImage"`
	ReadingTime   int       `json:"readingTime"`
	Status        string    `json:"status"`
	PublishedAt   *time.Time `json:"publishedAt"`
	ViewCount     int64     `json:"viewCount"`
	CreatedAt     time.Time `json:"createdAt"`
}

// BlogPostListResponse is the paginated list response
type BlogPostListResponse struct {
	Posts      []BlogPostBrief `json:"posts"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"totalPages"`
}

// BlogPostFilterRequest is for filtering blog posts
type BlogPostFilterRequest struct {
	Page       int        `form:"page"`
	Limit      int        `form:"limit"`
	Search     string     `form:"search"`
	Status     string     `form:"status"`
	Author     string     `form:"author"`
	IsFeatured *bool      `form:"isFeatured"`
	IsSpotlight *bool     `form:"isSpotlight"`
	SortBy     string     `form:"sortBy"`    // "createdAt", "updatedAt", "viewCount", "title"
	SortOrder  string     `form:"sortOrder"` // "asc", "desc"
}
