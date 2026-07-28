package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateArticleRequest struct {
	Title         string    `json:"title" binding:"required,max=255"`
	Slug          string    `json:"slug" binding:"required,max=255,alphanum"`
	LeadParagraph string    `json:"leadParagraph"`
	Content       string    `json:"content"`
	Footer        string    `json:"footer"`
	FeaturedImage string    `json:"featuredImage" binding:"omitempty,url"`

	// SEO
	MetaTitle       string `json:"metaTitle" binding:"max=70"`
	MetaDescription string `json:"metaDescription" binding:"max=160"`
	MetaKeywords    string `json:"metaKeywords"`
	CanonicalURL    string `json:"canonicalUrl" binding:"omitempty,url"`

	// Open Graph
	OGTitle       string `json:"ogTitle" binding:"max=95"`
	OGDescription string `json:"ogDescription" binding:"max=200"`
	OGImage       string `json:"ogImage" binding:"omitempty,url"`

	// Relations
	CategoryID uuid.UUID `json:"categoryId"`
	Tags       []string  `json:"tags"`
	AuthorID   uuid.UUID `json:"authorId" binding:"required"`

	// Publishing
	Status      string     `json:"status"`
	PublishedAt *time.Time `json:"publishedAt"`
	ScheduledAt *time.Time `json:"scheduledAt"`

	// Attractiveness
	IsFeatured  bool `json:"isFeatured"`
	IsSpotlight bool `json:"isSpotlight"`
	Priority    int  `json:"priority"`
}

type ArticleListRequest struct {
	Page        int        `form:"page"`
	Limit       int        `form:"limit"`
	Search      string     `form:"search"`
	CategoryID  *uuid.UUID `form:"categoryId"`
	Tags        []string   `form:"tags"`
	AuthorID    *uuid.UUID `form:"authorId"`
	Status      string     `form:"status"`
	IsFeatured  *bool      `form:"isFeatured"`
	IsSpotlight *bool      `form:"isSpotlight"`
	SortBy      string     `form:"sortBy"`
	SortOrder   string     `form:"sortOrder"`
}

// Response DTOs
type ArticleResponse struct {
	ID            uuid.UUID          `json:"id"`
	Title         string             `json:"title"`
	Slug          string             `json:"slug"`
	LeadParagraph string             `json:"leadParagraph"`
	Content       string             `json:"content"`
	Footer        string             `json:"footer"`
	FeaturedImage string             `json:"featuredImage"`

	// SEO
	MetaTitle       string `json:"metaTitle"`
	MetaDescription string `json:"metaDescription"`
	MetaKeywords    string `json:"metaKeywords"`
	CanonicalURL    string `json:"canonicalUrl"`

	// OG
	OGTitle       string `json:"ogTitle"`
	OGDescription string `json:"ogDescription"`
	OGImage       string `json:"ogImage"`

	// Relations
	Category        *CategoryResponse `json:"category,omitempty"`
	Tags            []string          `json:"tags"`
	AuthorID        uuid.UUID         `json:"authorId"`
	Author          *AuthorResponse   `json:"author,omitempty"`
	RelatedArticles []ArticleBrief    `json:"relatedArticles,omitempty"`

	// Status
	Status      string     `json:"status"`
	PublishedAt *time.Time `json:"publishedAt"`

	// Metrics
	ViewCount   int64 `json:"viewCount"`
	LikeCount   int64 `json:"likeCount"`
	ShareCount  int64 `json:"shareCount"`
	ReadingTime int   `json:"readingTime"`

	// Attractiveness
	IsFeatured  bool `json:"isFeatured"`
	IsSpotlight bool `json:"isSpotlight"`
	Priority    int  `json:"priority"`

	// Timestamps
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ArticleBrief struct {
	ID            uuid.UUID  `json:"id"`
	Title         string    `json:"title"`
	Slug          string    `json:"slug"`
	LeadParagraph string    `json:"leadParagraph"`
	FeaturedImage string    `json:"featuredImage"`
	ReadingTime   int       `json:"readingTime"`
	PublishedAt   *time.Time `json:"publishedAt"`
}

type ArticleListResponse struct {
	Articles   []ArticleBrief `json:"articles"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"totalPages"`
}

// SEO Preview DTO
type ArticleSEOPreview struct {
	Title        string `json:"title"`
	Slug         string `json:"slug"`
	MetaTitle    string `json:"metaTitle"`
	MetaDescription string `json:"metaDescription"`
	OGTitle      string `json:"ogTitle"`
	OGImage      string `json:"ogImage"`
	CanonicalURL string `json:"canonicalUrl"`
}

// Update DTO
type UpdateArticleRequest struct {
	Title         string     `json:"title" binding:"max=255"`
	Slug          string     `json:"slug" binding:"max=255"`
	LeadParagraph string     `json:"leadParagraph"`
	Content       string     `json:"content"`
	Footer        string     `json:"footer"`
	FeaturedImage string     `json:"featuredImage"`

	// SEO
	MetaTitle       string `json:"metaTitle" binding:"max=70"`
	MetaDescription string `json:"metaDescription" binding:"max=160"`
	MetaKeywords    string `json:"metaKeywords"`
	CanonicalURL    string `json:"canonicalUrl"`

	// Open Graph
	OGTitle       string `json:"ogTitle" binding:"max=95"`
	OGDescription string `json:"ogDescription" binding:"max=200"`
	OGImage       string `json:"ogImage"`

	// Relations
	CategoryID *uuid.UUID `json:"categoryId"`
	Tags       []string   `json:"tags"`

	// Publishing
	Status      string     `json:"status"`
	PublishedAt *time.Time `json:"publishedAt"`

	// Attractiveness
	IsFeatured  bool `json:"isFeatured"`
	IsSpotlight bool `json:"isSpotlight"`
	Priority    int  `json:"priority"`
}

// CategoryResponse DTO
type CategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	ParentID    *uuid.UUID `json:"parentId"`
	Order       int       `json:"order"`
}

// AuthorResponse DTO
type AuthorResponse struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	Avatar string    `json:"avatar"`
	Bio    string    `json:"bio"`
}