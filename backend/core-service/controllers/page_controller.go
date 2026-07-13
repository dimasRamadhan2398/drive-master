package controllers

import (
	"core-service/models/dto"
	"core-service/pkg/response"
	"core-service/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IPageController interface {
	GetAllPages(ctx *gin.Context)
	GetPageByID(ctx *gin.Context)
	GetPageBySlug(ctx *gin.Context)
	CreatePage(ctx *gin.Context)
	UpdatePage(ctx *gin.Context)
	DeletePage(ctx *gin.Context)
}

type PageController struct {
	pageService services.IPageService
}

func NewPageController(pageService services.IPageService) IPageController {
	return &PageController{
		pageService: pageService,
	}
}

func (c *PageController) GetAllPages(ctx *gin.Context) {
	pages, err := c.pageService.GetAllPages(ctx.Request.Context())
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch pages: "+err.Error())
		return
	}

	response.OK(ctx, "Pages fetched successfully", pages)
}

func (c *PageController) GetPageByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid page ID format")
		return
	}

	page, err := c.pageService.GetPageByID(ctx.Request.Context(), id)
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch page: "+err.Error())
		return
	}

	if page == nil {
		response.NotFound(ctx, "Page not found")
		return
	}

	response.OK(ctx, "Page fetched successfully", page)
}

func (c *PageController) GetPageBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if slug == "" {
		response.BadRequest(ctx, "Slug parameter is required")
		return
	}

	page, err := c.pageService.GetPageBySlug(ctx.Request.Context(), slug)
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch page by slug: "+err.Error())
		return
	}

	if page == nil {
		response.NotFound(ctx, "Page not found")
		return
	}

	response.OK(ctx, "Page fetched successfully", page)
}

func (c *PageController) CreatePage(ctx *gin.Context) {
	var req dto.CreatePageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request payload: "+err.Error())
		return
	}

	page, err := c.pageService.CreatePage(ctx.Request.Context(), &req)
	if err != nil {
		response.InternalServerError(ctx, "Failed to create page: "+err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Page created successfully",
		"data":    page,
	})
}

func (c *PageController) UpdatePage(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid page ID format")
		return
	}

	var req dto.UpdatePageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request payload: "+err.Error())
		return
	}

	page, err := c.pageService.UpdatePage(ctx.Request.Context(), id, &req)
	if err != nil {
		response.InternalServerError(ctx, "Failed to update page: "+err.Error())
		return
	}

	if page == nil {
		response.NotFound(ctx, "Page not found")
		return
	}

	response.OK(ctx, "Page updated successfully", page)
}

func (c *PageController) DeletePage(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid page ID format")
		return
	}

	err = c.pageService.DeletePage(ctx.Request.Context(), id)
	if err != nil {
		response.InternalServerError(ctx, "Failed to delete page: "+err.Error())
		return
	}

	response.OK(ctx, "Page deleted successfully", nil)
}
