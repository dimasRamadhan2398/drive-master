package controllers

import (
	"encoding/base64"
	"strings"
	"strconv"

	"core-service/models/dto"
	"core-service/pkg/response"
	"core-service/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ArticleController handles blog-related HTTP requests
type ArticleController struct {
	articleService services.IArticleService
}

// IArticleController defines the interface for blog post controller
type IArticleController interface {
	GetBlogArticles(ctx *gin.Context)
	GetBlogPostByID(ctx *gin.Context)
	CreateBlogPost(ctx *gin.Context)
	UpdateBlogPost(ctx *gin.Context)
	DeleteBlogPost(ctx *gin.Context)
	IncrementViewCount(ctx *gin.Context)
}

// NewArticleController creates a new article controller
func NewArticleController(articleService services.IArticleService) IArticleController {
	return &ArticleController{
		articleService: articleService,
	}
}

// GetBlogArticles handles GET /api/v1/articles/blog
func (c *ArticleController) GetBlogArticles(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	status := ctx.Query("status")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	result, err := c.articleService.GetBlogArticles(ctx.Request.Context(), page, limit, status)
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch blog articles: "+err.Error())
		return
	}

	response.OK(ctx, "Blog articles fetched successfully", result)
}

// GetBlogPostByID handles GET /api/v1/articles/blog/:id
func (c *ArticleController) GetBlogPostByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid article ID format")
		return
	}

	article, err := c.articleService.GetBlogPostByID(ctx.Request.Context(), id)
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch blog post: "+err.Error())
		return
	}

	if article == nil {
		response.NotFound(ctx, "Blog post not found")
		return
	}

	// Map to DTO for response consistency if needed, or return models.Article directly
	response.OK(ctx, "Blog post fetched successfully", article)
}

// CreateBlogPost handles POST /api/v1/articles/blog
// Supports both JSON (application/json) and multipart form-data (multipart/form-data)
func (c *ArticleController) CreateBlogPost(ctx *gin.Context) {
	var req dto.CreateBlogPostRequest
	contentType := ctx.GetHeader("Content-Type")

	if strings.Contains(contentType, "multipart/form-data") {
		// Bind text fields from form
		if err := ctx.ShouldBind(&req); err != nil {
			response.BadRequest(ctx, "Invalid form data: "+err.Error())
			return
		}

		// Handle file upload for featuredImage if present
		file, err := ctx.FormFile("featuredImage")
		if err == nil && file != nil {
			openedFile, err := file.Open()
			if err == nil {
				defer openedFile.Close()
				buffer := make([]byte, file.Size)
				_, err = openedFile.Read(buffer)
				if err == nil {
					base64Data := base64.StdEncoding.EncodeToString(buffer)
					req.FeaturedImage = &dto.FileUpload{
						Data:     base64Data,
						FileName: file.Filename,
						MimeType: file.Header.Get("Content-Type"),
					}
				}
			}
		}
	} else {
		// Fallback to JSON binding
		if err := ctx.ShouldBindJSON(&req); err != nil {
			response.BadRequest(ctx, "Invalid request body: "+err.Error())
			return
		}
	}

	article, err := c.articleService.CreateBlogPost(ctx.Request.Context(), &req)
	if err != nil {
		response.InternalServerError(ctx, "Failed to create blog post: "+err.Error())
		return
	}

	response.Created(ctx, "Blog post created successfully", article)
}

// UpdateBlogPost handles PUT /api/v1/articles/blog/:id
// Supports both JSON and multipart form-data
func (c *ArticleController) UpdateBlogPost(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid article ID format")
		return
	}

	var req dto.UpdateBlogPostRequest
	contentType := ctx.GetHeader("Content-Type")

	if strings.Contains(contentType, "multipart/form-data") {
		// Bind text fields from form
		if err := ctx.ShouldBind(&req); err != nil {
			response.BadRequest(ctx, "Invalid form data: "+err.Error())
			return
		}

		// Handle file upload for featuredImage if present
		file, err := ctx.FormFile("featuredImage")
		if err == nil && file != nil {
			openedFile, err := file.Open()
			if err == nil {
				defer openedFile.Close()
				buffer := make([]byte, file.Size)
				_, err = openedFile.Read(buffer)
				if err == nil {
					base64Data := base64.StdEncoding.EncodeToString(buffer)
					req.FeaturedImage = &dto.FileUpload{
						Data:     base64Data,
						FileName: file.Filename,
						MimeType: file.Header.Get("Content-Type"),
					}
				}
			}
		}
	} else {
		// Fallback to JSON binding
		if err := ctx.ShouldBindJSON(&req); err != nil {
			response.BadRequest(ctx, "Invalid request body: "+err.Error())
			return
		}
	}

	article, err := c.articleService.UpdateBlogPost(ctx.Request.Context(), id, &req)
	if err != nil {
		response.InternalServerError(ctx, "Failed to update blog post: "+err.Error())
		return
	}

	response.OK(ctx, "Blog post updated successfully", article)
}

// DeleteBlogPost handles DELETE /api/v1/articles/blog/:id
func (c *ArticleController) DeleteBlogPost(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid article ID format")
		return
	}

	if err := c.articleService.DeleteBlogPost(ctx.Request.Context(), id); err != nil {
		response.InternalServerError(ctx, "Failed to delete blog post: "+err.Error())
		return
	}

	response.OK(ctx, "Blog post deleted successfully", nil)
}

// IncrementViewCount handles POST /api/v1/articles/:id/view
func (c *ArticleController) IncrementViewCount(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid article ID format")
		return
	}

	if err := c.articleService.IncrementViewCount(ctx.Request.Context(), id); err != nil {
		response.InternalServerError(ctx, "Failed to increment view count: "+err.Error())
		return
	}

	response.OK(ctx, "View count incremented successfully", nil)
}
