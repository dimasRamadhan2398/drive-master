package repositories

import (
	"context"
	"core-service/models"
	"core-service/pkg/base"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// ArticleFilter holds filter options for article queries
type ArticleFilter struct {
	CategoryID  *uuid.UUID
	Tags        []string
	AuthorID    *uuid.UUID
	Status      string
	IsFeatured  *bool
	IsSpotlight *bool
	SearchQuery string
}

// IArticleRepository defines the interface for article repository
type IArticleRepository interface {
	Create(ctx context.Context, article *models.Article) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Article, error)
	FindBySlug(ctx context.Context, slug string) (*models.Article, error)
	FindAll(ctx context.Context, opts *base.QueryOptions, filter ArticleFilter) ([]models.Article, error)
	CountAll(ctx context.Context, filter ArticleFilter) (int64, error)
	Update(ctx context.Context, article *models.Article) error
	Delete(ctx context.Context, article *models.Article) error
	IncrementViewCount(ctx context.Context, id uuid.UUID) error

	// SEO-specific queries
	FindByTag(ctx context.Context, tag string, opts *base.QueryOptions) ([]models.Article, error)
	FindFeatured(ctx context.Context, limit int) ([]models.Article, error)
	FindSpotlight(ctx context.Context) (*models.Article, error)
	FindRelated(ctx context.Context, id uuid.UUID, limit int) ([]models.Article, error)
	Search(ctx context.Context, query string, opts *base.QueryOptions) ([]models.Article, error)
	CountByCategory(ctx context.Context, categoryID uuid.UUID) (int64, error)
	Publish(ctx context.Context, id uuid.UUID) error
	Archive(ctx context.Context, id uuid.UUID) error

	// Blog-specific queries
	GetBlogArticles(ctx context.Context, opts *base.QueryOptions) ([]models.Article, error)
	CountBlogArticles(ctx context.Context) (int64, error)
}

// ArticleRepository implements IArticleRepository
type ArticleRepository struct {
	*base.BaseRepository
}

func NewArticleRepository(baseRepo *base.BaseRepository) IArticleRepository {
	return &ArticleRepository{
		BaseRepository: baseRepo,
	}
}

// Create creates a new article
func (r *ArticleRepository) Create(ctx context.Context, article *models.Article) error {
	return r.BaseRepository.Create(ctx, article)
}

// FindByID finds an article by ID
func (r *ArticleRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Article, error) {
	var article models.Article
	if err := r.BaseRepository.FindByIDWithPreload(ctx, &article, id, "Category"); err != nil {
		return nil, err
	}
	return &article, nil
}

// FindBySlug finds an article by slug
func (r *ArticleRepository) FindBySlug(ctx context.Context, slug string) (*models.Article, error) {
	var article models.Article
	if err := r.DB.WithContext(ctx).Preload("Category").Where("slug = ?", slug).First(&article).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &article, nil
}

// FindAll finds all articles with filters and pagination
func (r *ArticleRepository) FindAll(ctx context.Context, opts *base.QueryOptions, filter ArticleFilter) ([]models.Article, error) {
	var articles []models.Article
	query := r.DB.WithContext(ctx).Model(&models.Article{})

	// Apply filters
	query = r.ApplyFilters(query, filter)

	// Apply options
	query = r.ApplyQueryOptions(query, opts)

	if err := query.Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

// CountAll counts articles with filters
func (r *ArticleRepository) CountAll(ctx context.Context, filter ArticleFilter) (int64, error) {
	var count int64
	query := r.DB.WithContext(ctx).Model(&models.Article{})

	query = r.ApplyFilters(query, filter)

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Update updates an article
func (r *ArticleRepository) Update(ctx context.Context, article *models.Article) error {
	return r.BaseRepository.Update(ctx, article)
}

// Delete soft-deletes an article
func (r *ArticleRepository) Delete(ctx context.Context, article *models.Article) error {
	return r.BaseRepository.Delete(ctx, article)
}

// IncrementViewCount increments the view count
func (r *ArticleRepository) IncrementViewCount(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Model(&models.Article{}).Where("id = ?", id).
		Update("view_count", gorm.Expr("view_count + 1")).Error
}

// FindByTag finds articles by tag
func (r *ArticleRepository) FindByTag(ctx context.Context, tag string, opts *base.QueryOptions) ([]models.Article, error) {
	var articles []models.Article
	query := r.DB.WithContext(ctx).Model(&models.Article{}).
		Where("tags @> ?", pq.Array([]string{tag}))

	query = r.ApplyQueryOptions(query, opts)

	if err := query.Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

// FindFeatured finds featured articles
func (r *ArticleRepository) FindFeatured(ctx context.Context, limit int) ([]models.Article, error) {
	var articles []models.Article
	err := r.DB.WithContext(ctx).
		Where("is_featured = ? AND status = ?", true, models.ArticleStatusPublished).
		Order("priority DESC, published_at DESC").
		Limit(limit).
		Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

// FindSpotlight finds the spotlight article
func (r *ArticleRepository) FindSpotlight(ctx context.Context) (*models.Article, error) {
	var article models.Article
	err := r.DB.WithContext(ctx).
		Where("is_spotlight = ? AND status = ?", true, models.ArticleStatusPublished).
		Order("priority DESC, published_at DESC").
		First(&article).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &article, nil
}

// FindRelated finds related articles
func (r *ArticleRepository) FindRelated(ctx context.Context, id uuid.UUID, limit int) ([]models.Article, error) {
	var articles []models.Article

	// Get current article tags
	var tags pq.StringArray
	err := r.DB.WithContext(ctx).Model(&models.Article{}).Where("id = ?", id).Pluck("tags", &tags).Error
	if err != nil {
		return nil, err
	}

	if len(tags) == 0 {
		return articles, nil
	}

	// Find related articles by shared tags, excluding current
	err = r.DB.WithContext(ctx).
		Where("id != ? AND status = ? AND tags && ?", id, models.ArticleStatusPublished, pq.Array([]string(tags))).
		Order("view_count DESC").
		Limit(limit).
		Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

// Search performs full-text search on articles
func (r *ArticleRepository) Search(ctx context.Context, query string, opts *base.QueryOptions) ([]models.Article, error) {
	var articles []models.Article
	searchPattern := "%" + query + "%"

	dbQuery := r.DB.WithContext(ctx).Model(&models.Article{}).
		Where("status = ?", models.ArticleStatusPublished).
		Where("title ILIKE ? OR lead_paragraph ILIKE ? OR meta_keywords ILIKE ?",
			searchPattern, searchPattern, searchPattern)

	dbQuery = r.ApplyQueryOptions(dbQuery, opts)

	if err := dbQuery.Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

// CountByCategory counts articles by category
func (r *ArticleRepository) CountByCategory(ctx context.Context, categoryID uuid.UUID) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&models.Article{}).
		Where("category_id = ? AND status = ?", categoryID, models.ArticleStatusPublished).
		Count(&count).Error
	return count, err
}

// Publish publishes an article (sets status to published and published_at)
func (r *ArticleRepository) Publish(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Model(&models.Article{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       models.ArticleStatusPublished,
			"published_at": gorm.Expr("NOW()"),
		}).Error
}

// Archive archives an article
func (r *ArticleRepository) Archive(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Model(&models.Article{}).Where("id = ?", id).
		Update("status", models.ArticleStatusArchived).Error
}

// ApplyFilters applies filter conditions to the query
func (r *ArticleRepository) ApplyFilters(query *gorm.DB, filter ArticleFilter) *gorm.DB {
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.AuthorID != nil {
		query = query.Where("author_id = ?", *filter.AuthorID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if len(filter.Tags) > 0 {
		query = query.Where("tags && ?", pq.Array(filter.Tags))
	}
	if filter.IsFeatured != nil {
		query = query.Where("is_featured = ?", *filter.IsFeatured)
	}
	if filter.IsSpotlight != nil {
		query = query.Where("is_spotlight = ?", *filter.IsSpotlight)
	}
	if filter.SearchQuery != "" {
		searchPattern := "%" + filter.SearchQuery + "%"
		query = query.Where("title ILIKE ? OR lead_paragraph ILIKE ?",
			searchPattern, searchPattern)
	}
	return query
}

// ApplyQueryOptions applies base query options to a GORM query
func (r *ArticleRepository) ApplyQueryOptions(query *gorm.DB, opts *base.QueryOptions) *gorm.DB {
	if opts == nil {
		return query
	}

	// Apply pagination
	if opts.Offset > 0 {
		query = query.Offset(opts.Offset)
	}
	if opts.Limit > 0 {
		query = query.Limit(opts.Limit)
	}

	// Apply Order
	if opts.Order != "" {
		query = query.Order(opts.Order)
	}

	return query
}

// GetBlogArticles retrieves published blog articles
func (r *ArticleRepository) GetBlogArticles(ctx context.Context, opts *base.QueryOptions) ([]models.Article, error) {
	var articles []models.Article
	query := r.DB.WithContext(ctx).Model(&models.Article{}).
		Where("status = ?", models.ArticleStatusPublished)

	query = r.ApplyQueryOptions(query, opts)
	if opts.Order == "" {
		query = query.Order("published_at DESC, created_at DESC")
	}

	if err := query.Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

// CountBlogArticles counts published blog articles
func (r *ArticleRepository) CountBlogArticles(ctx context.Context) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&models.Article{}).
		Where("status = ?", models.ArticleStatusPublished).
		Count(&count).Error
	return count, err
}
