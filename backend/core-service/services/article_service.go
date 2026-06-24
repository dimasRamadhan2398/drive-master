package services

import (
	"context"
	"core-service/models"
	"core-service/models/dto"
	"core-service/pkg/base"
	"core-service/pkg/kafka"
	"core-service/repositories"
	"encoding/json"
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
	GetBlogArticles(ctx context.Context, page, limit int, status string) (*dto.BlogArticleListResponse, error)
	CreateBlogArticle(ctx context.Context, req *dto.CreateBlogArticleRequest) (*models.Article, error)
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
	mediaSvc IMediaService
}

func NewArticleService(articleRepo repositories.IArticleRepository, faqRepo repositories.IFAQRepository, eventPublisher *kafka.EventPublisher) IArticleService {
	return &ArticleService{
		articleRepo:    articleRepo,
		faqRepo:        faqRepo,
		eventPublisher: eventPublisher,
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

	// Calculate reading time from body blocks (avg 200 words/min)
	readingTime := s.calculateReadingTimeFromBlocks(req.BodyBlocks)

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
		BodyBlocks:      req.BodyBlocks,
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
	if len(req.BodyBlocks) > 0 {
		article.BodyBlocks = req.BodyBlocks
		article.ReadingTime = s.calculateReadingTimeFromBlocks(req.BodyBlocks)
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
func (s *ArticleService) GetBlogArticles(ctx context.Context, page, limit int, status string) (*dto.BlogArticleListResponse, error) {
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

	// Convert to blog article response
	blogArticles := make([]dto.BlogArticleResponse, len(articles))
	for i, a := range articles {
		blogArticles[i] = dto.BlogArticleResponse{
			ID:            a.ID,
			Title:         a.Title,
			Slug:          a.Slug,
			LeadParagraph: a.LeadParagraph,
			FeaturedImage: a.FeaturedImage,
			ReadingTime:   a.ReadingTime,
			ViewCount:     a.ViewCount,
			LikeCount:     a.LikeCount,
			Status:        string(a.Status),
			PublishedAt:   a.PublishedAt,
			CreatedAt:     a.CreatedAt,
			UpdatedAt:     a.UpdatedAt,
		}
	}

	return &dto.BlogArticleListResponse{
		Data:   blogArticles,
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
	// Calculate reading time from body blocks
	readingTime := s.calculateReadingTimeFromBlocks(req.BodyBlocks)

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
		BodyBlocks:    req.BodyBlocks,
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

// calculateReadingTimeFromBlocks calculates reading time from body blocks (avg 200 words/min)
func (s *ArticleService) calculateReadingTimeFromBlocks(bodyBlocks []byte) int {
	if len(bodyBlocks) == 0 {
		return 1
	}

	// Try to extract text content from JSON blocks
	var blocks []map[string]interface{}
	if err := json.Unmarshal(bodyBlocks, &blocks); err != nil {
		return 1
	}

	var textContent strings.Builder
	for _, block := range blocks {
		if content, ok := block["content"].(string); ok {
			textContent.WriteString(content)
			textContent.WriteString(" ")
		}
	}

	wordCount := len(strings.Fields(textContent.String()))
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
