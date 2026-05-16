package controllers

import (
	"core-service/models"
	"core-service/models/dto"
	"core-service/pkg/response"
	"core-service/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/go-playground/validator/v10"
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
}

// NewPackageController creates a new package controller
func NewPackageController(packageService services.IPackageService) IPackageController {
	return &PackageController{
		packageService: packageService,
	}
}

// GetAllPackages handles GET /api/v1/packages
// @Summary Get all packages
// @Description Retrieves all packages
// @Tags Packages
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/packages [get]
func (c *PackageController) GetAllPackages(ctx *gin.Context) {
	packages, err := c.packageService.GetAllPackages(ctx.Request.Context())
	if err != nil {
		response.InternalServerError(ctx, "Failed to fetch packages")
		return
	}

	response.OK(ctx, "Packages fetched successfully", packages)
}

// GetPackageByID handles GET /api/v1/packages/:id
// @Summary Get package by ID
// @Description Retrieves a specific package by ID
// @Tags Packages
// @Produce json
// @Param id path string true "Package ID"
// @Success 200 {object} response.Response
// @Router /api/v1/packages/{id} [get]
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
// @Router /api/v1/packages [post]
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
	for _, b := range req.Benefits {
		benefits = append(benefits, models.PackageBenefit{
			Title:       b.Title,
			Description: b.Description,
			Icon:        b.Icon,
			SortOrder:   b.SortOrder,
		})
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
// @Router /api/v1/packages/{id} [put]
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
// @Router /api/v1/packages/{id} [delete]
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