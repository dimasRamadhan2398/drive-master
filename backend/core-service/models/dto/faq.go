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

// Blog Article DTOs

type CreateBlogArticleRequest struct {
	Title         string    `json:"title" binding:"required,max=255"`
	Slug          string    `json:"slug" binding:"required,max=255"`
	LeadParagraph string    `json:"leadParagraph"`
	BodyBlocks    []byte    `json:"bodyBlocks"`
	FeaturedImage string    `json:"featuredImage"`
	CategoryID    uuid.UUID `json:"categoryId"`
	AuthorID      uuid.UUID `json:"authorId" binding:"required"`
	Tags          []string  `json:"tags"`
	Status        string    `json:"status"`
}

type BlogArticleResponse struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Slug          string    `json:"slug"`
	LeadParagraph string    `json:"leadParagraph"`
	FeaturedImage string    `json:"featuredImage"`
	ReadingTime   int       `json:"readingTime"`
	ViewCount     int64     `json:"viewCount"`
	LikeCount     int64     `json:"likeCount"`
	Status        string    `json:"status"`
	PublishedAt   *time.Time `json:"publishedAt"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type BlogArticleListResponse struct {
	Articles   []BlogArticleResponse `json:"articles"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	Limit      int                  `json:"limit"`
	TotalPages int                  `json:"totalPages"`
}