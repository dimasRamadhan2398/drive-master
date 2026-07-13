package services

import (
	"context"
	"core-service/models"
	"core-service/models/dto"
	"core-service/pkg/base"
	"core-service/pkg/kafka"
	"core-service/repositories"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
)

// IArticleService defines the interface for blog/FAQ service
type IArticleService interface {
	// Blog endpoints
	GetBlogArticles(ctx context.Context, page, limit int, status string) (*dto.BlogPostListResponse, error)
	GetBlogPostByID(ctx context.Context, id uuid.UUID) (*dto.BlogPostResponse, error)
	CreateBlogPost(ctx context.Context, req *dto.CreateBlogPostRequest) (*models.Article, error)
	UpdateBlogPost(ctx context.Context, id uuid.UUID, req *dto.UpdateBlogPostRequest) (*models.Article, error)
	DeleteBlogPost(ctx context.Context, id uuid.UUID) error
	IncrementViewCount(ctx context.Context, id uuid.UUID) error

	// FAQ endpoints
	GetFAQs(ctx context.Context) ([]models.FAQ, error)
	CreateFAQ(ctx context.Context, req *dto.CreateFAQRequest) (*models.FAQ, error)
	DeleteFAQ(ctx context.Context, id uuid.UUID) error
}

// ArticleService implements IArticleService
type ArticleService struct {
	articleRepo    repositories.IArticleRepository
	faqRepo        repositories.IFAQRepository
	eventPublisher *kafka.EventPublisher
	mediaSvc       IMediaService
}

func NewArticleService(articleRepo repositories.IArticleRepository, faqRepo repositories.IFAQRepository, eventPublisher *kafka.EventPublisher, mediaSvc IMediaService) IArticleService {
	return &ArticleService{
		articleRepo:    articleRepo,
		faqRepo:        faqRepo,
		eventPublisher: eventPublisher,
		mediaSvc:       mediaSvc,
	}
}

// GetBlogPostByID retrieves a blog post by ID and returns it as a DTO
func (s *ArticleService) GetBlogPostByID(ctx context.Context, id uuid.UUID) (*dto.BlogPostResponse, error) {
	article, err := s.articleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, fmt.Errorf("article not found")
	}

	// Map to DTO with hardcoded author
	return &dto.BlogPostResponse{
		ID:            article.ID,
		Title:         article.Title,
		Slug:          article.Slug,
		Author:        "Admin", // Hardcoded author for all blog posts
		Content:       article.Content,
		Excerpt:       article.LeadParagraph,
		FeaturedImage:  article.FeaturedImage,
		ViewCount:     article.ViewCount,
		LikeCount:     article.LikeCount,
		ShareCount:    article.ShareCount,
		ReadingTime:   article.ReadingTime,
		CreatedAt:     article.CreatedAt,
		UpdatedAt:     article.UpdatedAt,
		Publishing: &dto.Publishing{
			Status:      string(article.Status),
			PublishedAt: article.PublishedAt,
			ScheduledAt: article.ScheduledAt,
		},
		Attractiveness: &dto.Attractiveness{
			IsFeatured:  article.IsFeatured,
			IsSpotlight: article.IsSpotlight,
			Priority:    article.Priority,
		},
	}, nil
}

// IncrementViewCount increments the view count for a blog post
func (s *ArticleService) IncrementViewCount(ctx context.Context, id uuid.UUID) error {
	article, err := s.articleRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if article == nil {
		return fmt.Errorf("article not found")
	}

	article.ViewCount++
	return s.articleRepo.Update(ctx, article)
}

// GetBlogArticles retrieves blog articles with optional status filter and pagination
func (s *ArticleService) GetBlogArticles(ctx context.Context, page, limit int, status string) (*dto.BlogPostListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit
	opts := base.NewQueryOptions()
	opts.Limit = limit
	opts.Offset = offset
	opts.Order = "published_at DESC, created_at DESC"

	articles, err := s.articleRepo.GetBlogArticles(ctx, opts, status)
	if err != nil {
		return nil, err
	}

	total, err := s.articleRepo.CountBlogArticles(ctx, status)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	// Convert to blog post response
	blogPosts := make([]dto.BlogPostResponse, len(articles))
	for i, a := range articles {
		statusStr := string(a.Status)
		blogPosts[i] = dto.BlogPostResponse{
			ID:            a.ID,
			Title:         a.Title,
			Slug:          a.Slug,
			Author:        "Admin", // Hardcoded author for all blog posts
			Content:       a.Content,
			Excerpt:       a.LeadParagraph,
			FeaturedImage: a.FeaturedImage,
			ViewCount:     a.ViewCount,
			LikeCount:     a.LikeCount,
			ShareCount:    a.ShareCount,
			ReadingTime:   a.ReadingTime,
			CreatedAt:     a.CreatedAt,
			UpdatedAt:     a.UpdatedAt,
			Publishing: &dto.Publishing{
				Status:      statusStr,
				PublishedAt: a.PublishedAt,
				ScheduledAt: a.ScheduledAt,
			},
			Attractiveness: &dto.Attractiveness{
				IsFeatured:  a.IsFeatured,
				IsSpotlight: a.IsSpotlight,
				Priority:    a.Priority,
			},
		}
	}

	return &dto.BlogPostListResponse{
		Data: blogPosts,
		Pagination: dto.PaginationMeta{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		},
	}, nil
}

// DeleteBlogPost soft-deletes a blog post
func (s *ArticleService) DeleteBlogPost(ctx context.Context, id uuid.UUID) error {
	article, err := s.articleRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if article == nil {
		return nil
	}

	return s.articleRepo.Delete(ctx, article)
}

// CreateBlogPost creates a new blog post with featured image upload to ImageKit
func (s *ArticleService) CreateBlogPost(ctx context.Context, req *dto.CreateBlogPostRequest) (*models.Article, error) {
	// Calculate reading time from content
	readingTime := s.calculateReadingTimeFromContent(req.Content)

	status := models.ArticleStatusDraft
	if req.Publishing != nil && req.Publishing.Status != "" {
		status = models.ArticleStatus(req.Publishing.Status)
	} else if req.Status != "" {
		status = models.ArticleStatus(req.Status)
	}

	// Build article model first (without featured image) to get the ID
	article := &models.Article{
		Title:         req.Title,
		Slug:          req.Slug,
		LeadParagraph: req.LeadParagraph,
		Content:       req.Content,
		AuthorID:      req.AuthorID,
		Status:        status,
		ReadingTime:   readingTime,
		IsFeatured:    req.Attractiveness != nil && req.Attractiveness.IsFeatured,
		IsSpotlight:   req.Attractiveness != nil && req.Attractiveness.IsSpotlight,
		Priority:      0,
	}

	if req.Attractiveness != nil {
		article.Priority = req.Attractiveness.Priority
	}

	if req.Publishing != nil {
		article.PublishedAt = req.Publishing.PublishedAt
		article.ScheduledAt = req.Publishing.ScheduledAt
	}

	// Create article first to get the ID
	if err := s.articleRepo.Create(ctx, article); err != nil {
		return nil, err
	}

	// Upload featured image to ImageKit if provided
	if req.FeaturedImage != nil && req.FeaturedImage.Data != "" {
		// Validate file format
		if !dto.IsSupportedImageFormat(req.FeaturedImage.FileName) {
			return nil, fmt.Errorf("unsupported image format: only jpg, jpeg, and png are supported")
		}

		// Upload to ImageKit with folder /blogs/{blogId}
		uploadResp, err := s.mediaSvc.UploadBase64Media(ctx, req.FeaturedImage.Data, req.FeaturedImage.FileName, fmt.Sprintf("blogs/%s", article.ID.String()))
		if err != nil {
			return nil, fmt.Errorf("failed to upload featured image: %w", err)
		}

		// Update article with featured image URL
		article.FeaturedImage = uploadResp.URL
		if err := s.articleRepo.Update(ctx, article); err != nil {
			return nil, fmt.Errorf("failed to update article with featured image: %w", err)
		}
	}

	return article, nil
}

// UpdateBlogPost updates an existing blog post with optional featured image upload
func (s *ArticleService) UpdateBlogPost(ctx context.Context, id uuid.UUID, req *dto.UpdateBlogPostRequest) (*models.Article, error) {
	article, err := s.articleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, fmt.Errorf("article not found")
	}

	// Update basic fields
	if req.Title != "" {
		article.Title = req.Title
	}
	if req.Slug != "" {
		article.Slug = req.Slug
	}
	if req.LeadParagraph != "" {
		article.LeadParagraph = req.LeadParagraph
	}
	if req.Content != "" {
		article.Content = req.Content
		article.ReadingTime = s.calculateReadingTimeFromContent(req.Content)
	}

	// Update publishing fields
	if req.Publishing != nil {
		if req.Publishing.Status != "" {
			article.Status = models.ArticleStatus(req.Publishing.Status)
		}
		if req.Publishing.PublishedAt != nil {
			article.PublishedAt = req.Publishing.PublishedAt
		}
		if req.Publishing.ScheduledAt != nil {
			article.ScheduledAt = req.Publishing.ScheduledAt
		}
	} else if req.Status != "" {
		article.Status = models.ArticleStatus(req.Status)
	}

	// Update attractiveness fields
	if req.Attractiveness != nil {
		article.IsFeatured = req.Attractiveness.IsFeatured
		article.IsSpotlight = req.Attractiveness.IsSpotlight
		article.Priority = req.Attractiveness.Priority
	}

	// Handle featured image upload if provided
	if req.FeaturedImage != nil {
		if req.FeaturedImage.Data != "" {
			// Validate file format
			if !dto.IsSupportedImageFormat(req.FeaturedImage.FileName) {
				return nil, fmt.Errorf("unsupported image format: only jpg, jpeg, and png are supported")
			}

			// Upload new featured image to ImageKit
			uploadResp, err := s.mediaSvc.UploadBase64Media(ctx, req.FeaturedImage.Data, req.FeaturedImage.FileName, fmt.Sprintf("blogs/%s", article.ID.String()))
			if err != nil {
				return nil, fmt.Errorf("failed to upload featured image: %w", err)
			}
			article.FeaturedImage = uploadResp.URL
		}
	}

	if err := s.articleRepo.Update(ctx, article); err != nil {
		return nil, err
	}

	return article, nil
}

// GetFAQs retrieves all FAQs
func (s *ArticleService) GetFAQs(ctx context.Context) ([]models.FAQ, error) {
	return s.faqRepo.FindActive(ctx)
}

// CreateFAQ creates a new FAQ
func (s *ArticleService) CreateFAQ(ctx context.Context, req *dto.CreateFAQRequest) (*models.FAQ, error) {
	category := req.Category
	if category == "" {
		category = "general"
	}

	order := req.Order
	if order == 0 {
		// Get max order and add 1
		faqs, err := s.faqRepo.FindAll(ctx)
		if err == nil && len(faqs) > 0 {
			maxOrder := 0
			for _, f := range faqs {
				if f.Order > maxOrder {
					maxOrder = f.Order
				}
			}
			order = maxOrder + 1
		}
	}

	faq := &models.FAQ{
		Question: req.Question,
		Answer:   req.Answer,
		Category: category,
		Order:    order,
		IsActive: true,
	}

	if err := s.faqRepo.Create(ctx, faq); err != nil {
		return nil, err
	}

	return faq, nil
}

// DeleteFAQ soft-deletes an FAQ
func (s *ArticleService) DeleteFAQ(ctx context.Context, id uuid.UUID) error {
	faq, err := s.faqRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if faq == nil {
		return nil
	}

	return s.faqRepo.DeleteSoft(ctx, faq)
}

// Helper methods

// calculateReadingTimeFromContent calculates reading time from content string (avg 200 words/min)
func (s *ArticleService) calculateReadingTimeFromContent(content string) int {
	if content == "" {
		return 1
	}

	wordCount := len(strings.Fields(content))
	minutes := wordCount / 200
	if minutes < 1 {
		minutes = 1
	}
	return minutes
}

// truncate truncates a string to maxLen characters
func (s *ArticleService) truncate(str string, maxLen int) string {
	if utf8.RuneCountInString(str) <= maxLen {
		return str
	}
	runes := []rune(str)
	return string(runes[:maxLen])
}