package controllers

import (
	"core-service/models"
	"core-service/models/dto"
	"core-service/pkg/response"
	"core-service/services"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// PackageController handles package-related HTTP requests
type PackageController struct {
	packageService services.IPackageService
}

// IPackageController defines the interface for package controller
type IPackageController interface {
	GetAllPackages(ctx *gin.Context)
	GetPackageByID(ctx *gin.Context)
	CreatePackage(ctx *gin.Context)
	UpdatePackage(ctx *gin.Context)
	DeletePackage(ctx *gin.Context)
	ToggleStatusPackage(ctx *gin.Context)
	IncrementCount(ctx *gin.Context)
}

// NewPackageController creates a new package controller
func NewPackageController(packageService services.IPackageService) IPackageController {
	return &PackageController{
		packageService: packageService,
	}
}

// ToggleStatusPackage handles PUT /api/v1/packages/toggle-status/:id
// @Summary Toggle package status
// @Description Toggles a package status between active and inactive
// @Tags Packages
// @Produce json
// @Param id path string true "Package ID"
// @Success 200 {object} response.Response
// @Router /packages/toggle-status/{id} [put]
func (c *PackageController) ToggleStatusPackage(ctx *gin.Context) {
	idParam := ctx.Param("id")
	if idParam == "" {
		response.BadRequest(ctx, "Package ID is required")
		return
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid package ID format")
		return
	}

	pkg, err := c.packageService.ToggleStatusPackage(ctx.Request.Context(), id)
	if err != nil {
		response.InternalServerError(ctx, "Failed to toggle package status")
		return
	}

	response.OK(ctx, "Package status toggled successfully", pkg)
}

// GetAllPackages handles GET /api/v1/packages
// @Summary Get all packages
// @Description Retrieves all packages with pagination
// @Tags Packages
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.Response
// @Router /packages [get]
func (c *PackageController) GetAllPackages(ctx *gin.Context) {
	var query struct {
		dto.PaginationQuery
		Status string `form:"status"`
	}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.BadRequest(ctx, "Invalid query parameters")
		return
	}

	// Set defaults
	page := query.Page
	limit := query.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	var packages []models.Package
	var total int64
	var err error

	if query.Status != "" {
		packages, err = c.packageService.GetPackagesByStatus(ctx.Request.Context(), models.PackageStatus(query.Status))
		total = int64(len(packages))
	} else {
		packages, total, err = c.packageService.GetAllPackagesPaginated(ctx.Request.Context(), page, limit)
	}

	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch packages")
		return
	}

	// Build response items
	var items []models.Package
	for _, pkg := range packages {
		items = append(items, pkg)
	}

	pagination := dto.NewPaginationMeta(total, page, limit)
	response.Paginated(ctx, http.StatusOK, "Packages fetched successfully", items, &pagination)
}

// GetPackageByID handles GET /api/v1/packages/:id
// @Summary Get package by ID
// @Description Retrieves a specific package by ID
// @Tags Packages
// @Produce json
// @Param id path string true "Package ID"
// @Success 200 {object} response.Response
// @Router /packages/{id} [get]
func (c *PackageController) GetPackageByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid package ID format")
		return
	}

	pkg, err := c.packageService.GetPackageByID(ctx.Request.Context(), id)
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch package")
		return
	}

	if pkg == nil {
		response.NotFound(ctx, "Package not found")
		return
	}

	response.OK(ctx, "Package fetched successfully", pkg)
}

// CreatePackage handles POST /api/v1/packages
// @Summary Create a new package
// @Description Creates a new package
// @Tags Packages
// @Accept json
// @Produce json
// @Param request body dto.CreatePackageRequest true "Package data"
// @Success 201 {object} response.Response
// @Router /packages [post]
func (c *PackageController) CreatePackage(ctx *gin.Context) {
	var req dto.CreatePackageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body: "+err.Error())
		return
	}

	if err := validator.New().Struct(req); err != nil {
		response.BadRequest(ctx, "Validation failed: "+err.Error())
		return
	}

	// Create benefits
	var benefits []models.PackageBenefit
	for i, b := range req.Benefits {
		switch v := b.(type) {
		case string:
			// Simple string - use as title
			benefits = append(benefits, models.PackageBenefit{
				Title:     v,
				SortOrder: i,
			})
		case map[string]interface{}:
			// Object format
			title, _ := v["title"].(string)
			description, _ := v["description"].(string)
			icon, _ := v["icon"].(string)
			sortOrder := i
			if so, ok := v["sortOrder"].(float64); ok {
				sortOrder = int(so)
			}
			benefits = append(benefits, models.PackageBenefit{
				Title:       title,
				Description: description,
				Icon:        icon,
				SortOrder:   sortOrder,
			})
		}
	}

	pkg := &models.Package{
		Name:            req.Name,
		Description:     req.Description,
		PackageType:     req.PackageType,
		Price:           req.Price,
		DiscountPrice:   req.DiscountPrice,
		DurationMinutes: req.DurationMinutes,
		TotalSessions:   req.TotalSessions,
		ImageURL:        req.ImageURL,
		Benefits:        benefits,
		Status:          models.PackageStatusActive,
	}

	if err := c.packageService.CreatePackage(ctx.Request.Context(), pkg); err != nil {
		response.InternalServerError(ctx, "Failed to create package")
		return
	}

	response.Created(ctx, "Package created successfully", pkg)
}

// UpdatePackage handles PUT /api/v1/packages/:id
// @Summary Update a package
// @Description Updates an existing package
// @Tags Packages
// @Accept json
// @Produce json
// @Param id path string true "Package ID"
// @Param request body dto.UpdatePackageRequest true "Package data"
// @Success 200 {object} response.Response
// @Router /packages/{id} [put]
func (c *PackageController) UpdatePackage(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid package ID format")
		return
	}

	// Get existing package
	pkg, err := c.packageService.GetPackageByID(ctx.Request.Context(), id)
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch package")
		return
	}

	if pkg == nil {
		response.NotFound(ctx, "Package not found")
		return
	}

	var req dto.UpdatePackageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body: "+err.Error())
		return
	}

	// Update fields
	if req.Name != "" {
		pkg.Name = req.Name
	}
	if req.Description != "" {
		pkg.Description = req.Description
	}
	if req.PackageType != "" {
		pkg.PackageType = req.PackageType
	}
	if req.Price != 0 {
		pkg.Price = req.Price
	}
	if req.DiscountPrice != 0 {
		pkg.DiscountPrice = req.DiscountPrice
	}
	if req.DurationMinutes != 0 {
		pkg.DurationMinutes = req.DurationMinutes
	}
	if req.TotalSessions != 0 {
		pkg.TotalSessions = req.TotalSessions
	}
	if req.Status != "" {
		pkg.Status = req.Status
	}
	if req.ImageURL != "" {
		pkg.ImageURL = req.ImageURL
	}

	// Handle Features
	if len(req.Features) > 0 {
		var features models.StringArray
		for _, b := range req.Features {
			switch v := b.(type) {
			case string:
				features = append(features, v)
			case map[string]interface{}:
				if title, ok := v["title"].(string); ok && title != "" {
					features = append(features, title)
				}
			}
		}
		pkg.Features = features
	}

	// Update Highlight
	pkg.Highlight = req.Highlight

	if err := c.packageService.UpdatePackage(ctx.Request.Context(), pkg); err != nil {
		response.InternalServerError(ctx, "Failed to update package")
		return
	}

	response.OK(ctx, "Package updated successfully", pkg)
}

// DeletePackage handles DELETE /api/v1/packages/:id
// @Summary Delete a package
// @Description Deletes a package by ID
// @Tags Packages
// @Produce json
// @Param id path string true "Package ID"
// @Success 200 {object} response.Response
// @Router /packages/{id} [delete]
func (c *PackageController) DeletePackage(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid package ID format")
		return
	}

	pkg, err := c.packageService.GetPackageByID(ctx.Request.Context(), id)
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch package")
		return
	}

	if pkg == nil {
		response.NotFound(ctx, "Package not found")
		return
	}

	if err := c.packageService.DeletePackage(ctx.Request.Context(), pkg); err != nil {
		response.InternalServerError(ctx, "Failed to delete package")
		return
	}

	response.OK(ctx, "Package deleted successfully", nil)
}

// IncrementCount handles POST /api/v1/packages/:id/increment-count
func (c *PackageController) IncrementCount(ctx *gin.Context) {
	idParam := ctx.Param("id")
	if idParam == "" {
		response.BadRequest(ctx, "Package ID is required")
		return
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(ctx, "Invalid package ID format")
		return
	}

	if err := c.packageService.IncrementStudentCount(ctx.Request.Context(), id); err != nil {
		response.InternalServerError(ctx, "Failed to increment package student count: "+err.Error())
		return
	}

	response.OK(ctx, "Package student count incremented successfully", nil)
}
