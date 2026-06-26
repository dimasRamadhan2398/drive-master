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

// IArticleService defines the interface for article service
type IArticleService interface {
	CreateArticle(ctx context.Context, req *dto.CreateArticleRequest) (*models.Article, error)
	GetArticleByID(ctx context.Context, id uuid.UUID) (*models.Article, error)
	GetArticleBySlug(ctx context.Context, slug string) (*models.Article, error)
	GetArticles(ctx context.Context, req *dto.ArticleListRequest) (*dto.ArticleListResponse, error)
	UpdateArticle(ctx context.Context, id uuid.UUID, req *dto.UpdateArticleRequest) (*models.Article, error)
	DeleteArticle(ctx context.Context, id uuid.UUID) error

	// Public endpoints
	IncrementViewCount(ctx context.Context, id uuid.UUID) error
	SearchArticles(ctx context.Context, query string, page, limit int) (*dto.ArticleListResponse, error)
	GetFeaturedArticles(ctx context.Context, limit int) ([]models.Article, error)
	GetSpotlightArticle(ctx context.Context) (*models.Article, error)
	GetRelatedArticles(ctx context.Context, id uuid.UUID, limit int) ([]models.Article, error)
	GetArticlesByTag(ctx context.Context, tag string, page, limit int) (*dto.ArticleListResponse, error)

	// Admin endpoints
	PublishArticle(ctx context.Context, id uuid.UUID) error
	ArchiveArticle(ctx context.Context, id uuid.UUID) error

	// Blog endpoints
	GetBlogArticles(ctx context.Context, page, limit int, status string) (*dto.BlogPostListResponse, error)
	CreateBlogArticle(ctx context.Context, req *dto.CreateBlogArticleRequest) (*models.Article, error)
	CreateBlogPost(ctx context.Context, req *dto.CreateBlogPostRequest) (*models.Article, error)
	UpdateBlogPost(ctx context.Context, id uuid.UUID, req *dto.UpdateBlogPostRequest) (*models.Article, error)
	DeleteBlogArticle(ctx context.Context, id uuid.UUID) error

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

// CreateArticle creates a new article
func (s *ArticleService) CreateArticle(ctx context.Context, req *dto.CreateArticleRequest) (*models.Article, error) {
	// Auto-generate SEO fields if not provided
	metaTitle := req.MetaTitle
	if metaTitle == "" {
		metaTitle = s.truncate(req.Title, 70)
	}
	metaDescription := req.MetaDescription
	if metaDescription == "" {
		metaDescription = s.truncate(req.LeadParagraph, 160)
	}
	ogTitle := req.OGTitle
	if ogTitle == "" {
		ogTitle = s.truncate(req.Title, 95)
	}
	ogDescription := req.OGDescription
	if ogDescription == "" {
		ogDescription = s.truncate(req.LeadParagraph, 200)
	}

	// Calculate reading time from content (avg 200 words/min)
	readingTime := s.calculateReadingTimeFromContent(req.Content)

	status := models.ArticleStatus(req.Status)
	if status == "" {
		status = models.ArticleStatusDraft
	}

	// Handle category ID pointer
	var categoryID *uuid.UUID
	if req.CategoryID != uuid.Nil {
		categoryID = &req.CategoryID
	}

	article := &models.Article{
		Title:           req.Title,
		Slug:            req.Slug,
		LeadParagraph:   req.LeadParagraph,
		Content:         req.Content,
		Footer:          req.Footer,
		FeaturedImage:   req.FeaturedImage,
		CategoryID:      categoryID,
		AuthorID:        req.AuthorID,
		Tags:            req.Tags,
		Status:          status,
		PublishedAt:     req.PublishedAt,
		ScheduledAt:     req.ScheduledAt,
		MetaTitle:       metaTitle,
		MetaDescription: metaDescription,
		MetaKeywords:    req.MetaKeywords,
		CanonicalURL:    req.CanonicalURL,
		OGTitle:         ogTitle,
		OGDescription:   ogDescription,
		OGImage:         req.OGImage,
		ReadingTime:     readingTime,
		IsFeatured:      req.IsFeatured,
		IsSpotlight:     req.IsSpotlight,
		Priority:        req.Priority,
	}

	if err := s.articleRepo.Create(ctx, article); err != nil {
		return nil, err
	}

	// Publish event (async to not block response)
	if s.eventPublisher != nil {
		go s.eventPublisher.PublishArticleCreated(context.Background(), article)
	}

	return article, nil
}

// GetArticleByID retrieves an article by ID
func (s *ArticleService) GetArticleByID(ctx context.Context, id uuid.UUID) (*models.Article, error) {
	return s.articleRepo.FindByID(ctx, id)
}

// GetArticleBySlug retrieves an article by slug
func (s *ArticleService) GetArticleBySlug(ctx context.Context, slug string) (*models.Article, error) {
	return s.articleRepo.FindBySlug(ctx, slug)
}

// GetArticles retrieves articles with filters and pagination
func (s *ArticleService) GetArticles(ctx context.Context, req *dto.ArticleListRequest) (*dto.ArticleListResponse, error) {
	// Build filter
	filter := repositories.ArticleFilter{
		CategoryID:  req.CategoryID,
		Tags:        req.Tags,
		AuthorID:    req.AuthorID,
		Status:      req.Status,
		IsFeatured:  req.IsFeatured,
		IsSpotlight: req.IsSpotlight,
		SearchQuery: req.Search,
	}

	// Build query options
	opts := s.buildQueryOptions(req)

	// Get articles
	articles, err := s.articleRepo.FindAll(ctx, opts, filter)
	if err != nil {
		return nil, err
	}

	// Get total count
	total, err := s.articleRepo.CountAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Calculate pagination
	page := req.Page
	if page < 1 {
		page = 1
	}
	limit := req.Limit
	if limit < 1 {
		limit = 10
	}
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	// Convert to brief response
	articleBriefs := make([]dto.ArticleBrief, len(articles))
	for i, a := range articles {
		articleBriefs[i] = s.toArticleBrief(&a)
	}

	return &dto.ArticleListResponse{
		Articles:   articleBriefs,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// UpdateArticle updates an article
func (s *ArticleService) UpdateArticle(ctx context.Context, id uuid.UUID, req *dto.UpdateArticleRequest) (*models.Article, error) {
	article, err := s.articleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
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
	if req.Footer != "" {
		article.Footer = req.Footer
	}
	if req.FeaturedImage != "" {
		article.FeaturedImage = req.FeaturedImage
	}
	if req.CategoryID != nil {
		article.CategoryID = req.CategoryID
	}
	if len(req.Tags) > 0 {
		article.Tags = req.Tags
	}

	// Update SEO fields
	if req.MetaTitle != "" {
		article.MetaTitle = req.MetaTitle
	}
	if req.MetaDescription != "" {
		article.MetaDescription = req.MetaDescription
	}
	if req.MetaKeywords != "" {
		article.MetaKeywords = req.MetaKeywords
	}
	if req.CanonicalURL != "" {
		article.CanonicalURL = req.CanonicalURL
	}
	if req.OGTitle != "" {
		article.OGTitle = req.OGTitle
	}
	if req.OGDescription != "" {
		article.OGDescription = req.OGDescription
	}
	if req.OGImage != "" {
		article.OGImage = req.OGImage
	}

	// Update status
	if req.Status != "" {
		article.Status = models.ArticleStatus(req.Status)
	}
	if req.PublishedAt != nil {
		article.PublishedAt = req.PublishedAt
	}

	// Update attractiveness fields
	if req.IsFeatured {
		article.IsFeatured = true
	}
	if req.IsSpotlight {
		article.IsSpotlight = true
	}
	if req.Priority != 0 {
		article.Priority = req.Priority
	}

	// Auto-regenerate SEO if title/content changed and SEO not explicitly set
	if req.MetaTitle == "" && req.Title != "" {
		article.MetaTitle = s.truncate(article.Title, 70)
	}

	if err := s.articleRepo.Update(ctx, article); err != nil {
		return nil, err
	}

	// Publish event (async to not block response)
	if s.eventPublisher != nil {
		go s.eventPublisher.PublishArticleUpdated(context.Background(), article)
	}

	return article, nil
}

// DeleteArticle soft-deletes an article
func (s *ArticleService) DeleteArticle(ctx context.Context, id uuid.UUID) error {
	article, err := s.articleRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.articleRepo.Delete(ctx, article); err != nil {
		return err
	}

	// Publish event (async to not block response)
	if s.eventPublisher != nil {
		go s.eventPublisher.PublishArticleDeleted(context.Background(), article.ID.String())
	}

	return nil
}

// IncrementViewCount increments the view count of an article
func (s *ArticleService) IncrementViewCount(ctx context.Context, id uuid.UUID) error {
	return s.articleRepo.IncrementViewCount(ctx, id)
}

// SearchArticles searches articles by query
func (s *ArticleService) SearchArticles(ctx context.Context, query string, page, limit int) (*dto.ArticleListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	articles, err := s.articleRepo.Search(ctx, query, &base.QueryOptions{
		Limit:  limit,
		Offset: offset,
		Order:  "published_at DESC",
	})
	if err != nil {
		return nil, err
	}

	filter := repositories.ArticleFilter{SearchQuery: query}
	total, err := s.articleRepo.CountAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	articleBriefs := make([]dto.ArticleBrief, len(articles))
	for i, a := range articles {
		articleBriefs[i] = s.toArticleBrief(&a)
	}

	return &dto.ArticleListResponse{
		Articles:   articleBriefs,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// GetFeaturedArticles retrieves featured articles
func (s *ArticleService) GetFeaturedArticles(ctx context.Context, limit int) ([]models.Article, error) {
	return s.articleRepo.FindFeatured(ctx, limit)
}

// GetSpotlightArticle retrieves the spotlight article
func (s *ArticleService) GetSpotlightArticle(ctx context.Context) (*models.Article, error) {
	return s.articleRepo.FindSpotlight(ctx)
}

// GetRelatedArticles retrieves related articles based on tags
func (s *ArticleService) GetRelatedArticles(ctx context.Context, id uuid.UUID, limit int) ([]models.Article, error) {
	return s.articleRepo.FindRelated(ctx, id, limit)
}

// GetArticlesByTag retrieves articles by tag
func (s *ArticleService) GetArticlesByTag(ctx context.Context, tag string, page, limit int) (*dto.ArticleListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	articles, err := s.articleRepo.FindByTag(ctx, tag, &base.QueryOptions{
		Limit:  limit,
		Offset: offset,
		Order:  "published_at DESC",
	})
	if err != nil {
		return nil, err
	}

	filter := repositories.ArticleFilter{Tags: []string{tag}}
	total, err := s.articleRepo.CountAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	articleBriefs := make([]dto.ArticleBrief, len(articles))
	for i, a := range articles {
		articleBriefs[i] = s.toArticleBrief(&a)
	}

	return &dto.ArticleListResponse{
		Articles:   articleBriefs,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// PublishArticle publishes an article
func (s *ArticleService) PublishArticle(ctx context.Context, id uuid.UUID) error {
	if err := s.articleRepo.Publish(ctx, id); err != nil {
		return err
	}

	// Publish event (async to not block response)
	if s.eventPublisher != nil {
		article, _ := s.articleRepo.FindByID(ctx, id)
		if article != nil {
			go s.eventPublisher.PublishArticlePublished(context.Background(), article)
		}
	}

	return nil
}

// ArchiveArticle archives an article
func (s *ArticleService) ArchiveArticle(ctx context.Context, id uuid.UUID) error {
	if err := s.articleRepo.Archive(ctx, id); err != nil {
		return err
	}

	// Publish event (async to not block response)
	if s.eventPublisher != nil {
		go s.eventPublisher.PublishArticleArchived(context.Background(), id.String())
	}

	return nil
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
		blogPosts[i] = dto.BlogPostResponse{
			ID:            a.ID,
			Title:         a.Title,
			Slug:          a.Slug,
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
				Status:      string(a.Status),
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

// CreateBlogArticle creates a new blog article
func (s *ArticleService) CreateBlogArticle(ctx context.Context, req *dto.CreateBlogArticleRequest) (*models.Article, error) {
	// Calculate reading time from content
	readingTime := s.calculateReadingTimeFromContent(req.Content)

	status := models.ArticleStatus(req.Status)
	if status == "" {
		status = models.ArticleStatusDraft
	}

	// Handle category ID pointer
	var categoryID *uuid.UUID
	if req.CategoryID != uuid.Nil {
		categoryID = &req.CategoryID
	}

	article := &models.Article{
		Title:         req.Title,
		Slug:          req.Slug,
		LeadParagraph: req.LeadParagraph,
		Content:       req.Content,
		FeaturedImage: req.FeaturedImage,
		CategoryID:    categoryID,
		AuthorID:      req.AuthorID,
		Tags:          req.Tags,
		Status:        status,
		ReadingTime:   readingTime,
	}

	if err := s.articleRepo.Create(ctx, article); err != nil {
		return nil, err
	}

	return article, nil
}

// DeleteBlogArticle soft-deletes a blog article
func (s *ArticleService) DeleteBlogArticle(ctx context.Context, id uuid.UUID) error {
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

// toArticleBrief converts an article to ArticleBrief DTO
func (s *ArticleService) toArticleBrief(article *models.Article) dto.ArticleBrief {
	return dto.ArticleBrief{
		ID:            article.ID,
		Title:         article.Title,
		Slug:          article.Slug,
		LeadParagraph: article.LeadParagraph,
		FeaturedImage: article.FeaturedImage,
		ReadingTime:   article.ReadingTime,
		PublishedAt:   article.PublishedAt,
	}
}

// buildQueryOptions builds QueryOptions from ArticleListRequest
func (s *ArticleService) buildQueryOptions(req *dto.ArticleListRequest) *base.QueryOptions {
	opts := base.NewQueryOptions()

	opts.Limit = req.Limit
	if req.Page > 0 {
		opts.Offset = (req.Page - 1) * req.Limit
	}

	if req.SortBy != "" {
		order := req.SortBy
		if req.SortOrder != "" {
			order += " " + req.SortOrder
		} else {
			order += " DESC"
		}
		opts.Order = order
	}

	return opts
}