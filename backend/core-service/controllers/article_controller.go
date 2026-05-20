package controllers

import (
	"core-service/models/dto"
	"core-service/pkg/response"
	"core-service/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ArticleController handles article-related HTTP requests
type ArticleController struct {
	articleService services.IArticleService
}

// IArticleController defines the interface for article controller
type IArticleController interface {
	GetAllArticles(ctx *gin.Context)
	GetArticleByID(ctx *gin.Context)
	GetArticleBySlug(ctx *gin.Context)
	CreateArticle(ctx *gin.Context)
	UpdateArticle(ctx *gin.Context)
	DeleteArticle(ctx *gin.Context)

	// Public endpoints
	SearchArticles(ctx *gin.Context)
	GetFeaturedArticles(ctx *gin.Context)
	GetSpotlightArticle(ctx *gin.Context)
	GetRelatedArticles(ctx *gin.Context)
	GetArticlesByTag(ctx *gin.Context)
	IncrementViewCount(ctx *gin.Context)

	// Admin endpoints
	PublishArticle(ctx *gin.Context)
	ArchiveArticle(ctx *gin.Context)
}

// NewArticleController creates a new article controller
func NewArticleController(articleService services.IArticleService) IArticleController {
	return &ArticleController{
		articleService: articleService,
	}
}

// GetAllArticles handles GET /api/v1/articles
// @Summary Get all articles
// @Description Retrieves all articles with filters and pagination
// @Tags Articles
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param search query string false "Search query"
// @Param categoryId query string false "Category ID"
// @Param status query string false "Status (draft, published, archived)"
// @Success 200 {object} response.Response
// @Router /articles [get]
func (c *ArticleController) GetAllArticles(ctx *gin.Context) {
	var req dto.ArticleListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "Invalid query parameters: "+err.Error())
		return
	}

	// Set defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	result, err := c.articleService.GetArticles(ctx.Request.Context(), &req)
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch articles: "+err.Error())
		return
	}

	response.OK(ctx, "Articles fetched successfully", result)
}

// GetArticleByID handles GET /api/v1/articles/:id
// @Summary Get article by ID
// @Description Retrieves a specific article by ID
// @Tags Articles
// @Produce json
// @Param id path string true "Article ID"
// @Success 200 {object} response.Response
// @Router /articles/{id} [get]
func (c *ArticleController) GetArticleByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid article ID format")
		return
	}

	article, err := c.articleService.GetArticleByID(ctx.Request.Context(), id)
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch article")
		return
	}

	if article == nil {
		response.NotFound(ctx, "Article not found")
		return
	}

	response.OK(ctx, "Article fetched successfully", article)
}

// GetArticleBySlug handles GET /api/v1/articles/slug/:slug
// @Summary Get article by slug
// @Description Retrieves a specific article by slug
// @Tags Articles
// @Produce json
// @Param slug path string true "Article slug"
// @Success 200 {object} response.Response
// @Router /articles/slug/{slug} [get]
func (c *ArticleController) GetArticleBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if slug == "" {
		response.BadRequest(ctx, "Slug is required")
		return
	}

	article, err := c.articleService.GetArticleBySlug(ctx.Request.Context(), slug)
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch article")
		return
	}

	if article == nil {
		response.NotFound(ctx, "Article not found")
		return
	}

	response.OK(ctx, "Article fetched successfully", article)
}

// CreateArticle handles POST /api/v1/articles
// @Summary Create a new article
// @Description Creates a new article
// @Tags Articles
// @Accept json
// @Produce json
// @Param request body dto.CreateArticleRequest true "Article data"
// @Success 201 {object} response.Response
// @Router /articles [post]
func (c *ArticleController) CreateArticle(ctx *gin.Context) {
	var req dto.CreateArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body: "+err.Error())
		return
	}

	article, err := c.articleService.CreateArticle(ctx.Request.Context(), &req)
	if err != nil {
		response.InternalServerError(ctx, "Failed to create article: "+err.Error())
		return
	}

	response.Created(ctx, "Article created successfully", article)
}

// UpdateArticle handles PUT /api/v1/articles/:id
// @Summary Update an article
// @Description Updates an existing article
// @Tags Articles
// @Accept json
// @Produce json
// @Param id path string true "Article ID"
// @Param request body dto.UpdateArticleRequest true "Article data"
// @Success 200 {object} response.Response
// @Router /articles/{id} [put]
func (c *ArticleController) UpdateArticle(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid article ID format")
		return
	}

	var req dto.UpdateArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body: "+err.Error())
		return
	}

	article, err := c.articleService.UpdateArticle(ctx.Request.Context(), id, &req)
	if err != nil {
		response.InternalServerError(ctx, "Failed to update article: "+err.Error())
		return
	}

	response.OK(ctx, "Article updated successfully", article)
}

// DeleteArticle handles DELETE /api/v1/articles/:id
// @Summary Delete an article
// @Description Soft-deletes an article by ID
// @Tags Articles
// @Produce json
// @Param id path string true "Article ID"
// @Success 200 {object} response.Response
// @Router /articles/{id} [delete]
func (c *ArticleController) DeleteArticle(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid article ID format")
		return
	}

	if err := c.articleService.DeleteArticle(ctx.Request.Context(), id); err != nil {
		response.InternalServerError(ctx, "Failed to delete article")
		return
	}

	response.OK(ctx, "Article deleted successfully", nil)
}

// SearchArticles handles GET /api/v1/articles/search
// @Summary Search articles
// @Description Searches articles by query string
// @Tags Articles
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} response.Response
// @Router /articles/search [get]
func (c *ArticleController) SearchArticles(ctx *gin.Context) {
	query := ctx.Query("q")
	if query == "" {
		response.BadRequest(ctx, "Search query is required")
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	result, err := c.articleService.SearchArticles(ctx.Request.Context(), query, page, limit)
	if err != nil {
		response.InternalServerError(ctx, "Failed to search articles")
		return
	}

	response.OK(ctx, "Articles found", result)
}

// GetFeaturedArticles handles GET /api/v1/articles/featured
// @Summary Get featured articles
// @Description Retrieves featured articles
// @Tags Articles
// @Produce json
// @Param limit query int false "Number of articles"
// @Success 200 {object} response.Response
// @Router /articles/featured [get]
func (c *ArticleController) GetFeaturedArticles(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "5"))
	if limit < 1 {
		limit = 5
	}
	if limit > 20 {
		limit = 20
	}

	articles, err := c.articleService.GetFeaturedArticles(ctx.Request.Context(), limit)
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch featured articles")
		return
	}

	response.OK(ctx, "Featured articles fetched successfully", articles)
}

// GetSpotlightArticle handles GET /api/v1/articles/spotlight
// @Summary Get spotlight article
// @Description Retrieves the spotlight article
// @Tags Articles
// @Produce json
// @Success 200 {object} response.Response
// @Router /articles/spotlight [get]
func (c *ArticleController) GetSpotlightArticle(ctx *gin.Context) {
	article, err := c.articleService.GetSpotlightArticle(ctx.Request.Context())
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch spotlight article")
		return
	}

	if article == nil {
		response.NotFound(ctx, "No spotlight article found")
		return
	}

	response.OK(ctx, "Spotlight article fetched successfully", article)
}

// GetRelatedArticles handles GET /api/v1/articles/:id/related
// @Summary Get related articles
// @Description Retrieves articles related to the specified article
// @Tags Articles
// @Produce json
// @Param id path string true "Article ID"
// @Param limit query int false "Number of articles"
// @Success 200 {object} response.Response
// @Router /articles/{id}/related [get]
func (c *ArticleController) GetRelatedArticles(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid article ID format")
		return
	}

	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "5"))
	if limit < 1 {
		limit = 5
	}
	if limit > 10 {
		limit = 10
	}

	articles, err := c.articleService.GetRelatedArticles(ctx.Request.Context(), id, limit)
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch related articles")
		return
	}

	response.OK(ctx, "Related articles fetched successfully", articles)
}

// GetArticlesByTag handles GET /api/v1/articles/tag/:tag
// @Summary Get articles by tag
// @Description Retrieves articles with a specific tag
// @Tags Articles
// @Produce json
// @Param tag path string true "Tag name"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} response.Response
// @Router /articles/tag/{tag} [get]
func (c *ArticleController) GetArticlesByTag(ctx *gin.Context) {
	tag := ctx.Param("tag")
	if tag == "" {
		response.BadRequest(ctx, "Tag is required")
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	result, err := c.articleService.GetArticlesByTag(ctx.Request.Context(), tag, page, limit)
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch articles by tag")
		return
	}

	response.OK(ctx, "Articles found", result)
}

// IncrementViewCount handles POST /api/v1/articles/:id/view
// @Summary Increment view count
// @Description Increments the view count of an article
// @Tags Articles
// @Produce json
// @Param id path string true "Article ID"
// @Success 200 {object} response.Response
// @Router /articles/{id}/view [post]
func (c *ArticleController) IncrementViewCount(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid article ID format")
		return
	}

	if err := c.articleService.IncrementViewCount(ctx.Request.Context(), id); err != nil {
		response.InternalServerError(ctx, "Failed to increment view count")
		return
	}

	response.OK(ctx, "View count incremented", nil)
}

// PublishArticle handles POST /api/v1/articles/:id/publish
// @Summary Publish an article
// @Description Sets article status to published
// @Tags Articles
// @Produce json
// @Param id path string true "Article ID"
// @Success 200 {object} response.Response
// @Router /articles/{id}/publish [post]
func (c *ArticleController) PublishArticle(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid article ID format")
		return
	}

	if err := c.articleService.PublishArticle(ctx.Request.Context(), id); err != nil {
		response.InternalServerError(ctx, "Failed to publish article")
		return
	}

	response.OK(ctx, "Article published successfully", nil)
}

// ArchiveArticle handles POST /api/v1/articles/:id/archive
// @Summary Archive an article
// @Description Sets article status to archived
// @Tags Articles
// @Produce json
// @Param id path string true "Article ID"
// @Success 200 {object} response.Response
// @Router /articles/{id}/archive [post]
func (c *ArticleController) ArchiveArticle(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid article ID format")
		return
	}

	if err := c.articleService.ArchiveArticle(ctx.Request.Context(), id); err != nil {
		response.InternalServerError(ctx, "Failed to archive article")
		return
	}

	response.OK(ctx, "Article archived successfully", nil)
}