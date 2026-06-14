package controllers

import (
	"net/http"
	"strconv"

	"booking-service/models/dto"
	"booking-service/pkg/base"
	"booking-service/services"

	"github.com/gin-gonic/gin"
)

type CertificationController struct {
	certificationService services.ICertificationService
}

func NewCertificationController(certificationService services.ICertificationService) ICertificationController {
	return &CertificationController{certificationService: certificationService}
}

type ICertificationController interface {
	CreateCertification(c *gin.Context)
	GetCertification(c *gin.Context)
	UpdateCertificationStatus(c *gin.Context)
	ListCertifications(c *gin.Context)
	IssueCertification(c *gin.Context)
	RevokeCertification(c *gin.Context)
	GetUserCertifications(c *gin.Context)
	GetCertificationsByPackage(c *gin.Context)
}

// CreateCertification godoc
// @Summary Create a new certification
// @Description Creates a new certification with the provided details
// @Tags certifications
// @Accept json
// @Produce json
// @Param certification body dto.CreateCertificationRequest true "Certification data"
// @Success 201 {object} dto.CertificationResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /certifications [post]
func (c *CertificationController) CreateCertification(ctx *gin.Context) {
	var req dto.CreateCertificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.certificationService.CreateCertification(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// GetCertification godoc
// @Summary Get a certification by ID
// @Description Retrieves a certification by its ID
// @Tags certifications
// @Accept json
// @Produce json
// @Param id path int true "Certification ID"
// @Success 200 {object} dto.CertificationResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /certifications/{id} [get]
func (c *CertificationController) GetCertification(ctx *gin.Context) {
	id, err := base.GetUintIDFromPath(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid certification id"})
		return
	}

	resp, err := c.certificationService.GetCertification(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateCertificationStatus godoc
// @Summary Update certification status
// @Description Updates the status of a certification
// @Tags certifications
// @Accept json
// @Produce json
// @Param id path int true "Certification ID"
// @Param certification body dto.UpdateCertificationRequest true "Status update data"
// @Success 200 {object} dto.CertificationResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /certifications/{id}/status [put]
func (c *CertificationController) UpdateCertificationStatus(ctx *gin.Context) {
	id, err := base.GetUintIDFromPath(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid certification id"})
		return
	}

	var req dto.UpdateCertificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Status == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "status is required"})
		return
	}

	resp, err := c.certificationService.UpdateCertificationStatus(ctx.Request.Context(), id, *req.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// ListCertifications godoc
// @Summary List all certifications
// @Description Retrieves a paginated list of certifications
// @Tags certifications
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} dto.CertificationListResponse
// @Failure 500 {object} map[string]string
// @Router /certifications [get]
func (c *CertificationController) ListCertifications(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	resp, err := c.certificationService.ListCertifications(ctx.Request.Context(), page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// IssueCertification godoc
// @Summary Issue a certification to a user
// @Description Issues a certification to a specific user
// @Tags certifications
// @Accept json
// @Produce json
// @Param issue body dto.IssueCertificationRequest true "Issue certification data"
// @Success 201 {object} dto.UserCertificationResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /certifications/issue [post]
func (c *CertificationController) IssueCertification(ctx *gin.Context) {
	var req dto.IssueCertificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.certificationService.IssueCertification(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// RevokeCertification godoc
// @Summary Revoke a user's certification
// @Description Revokes a certification from a specific user
// @Tags certifications
// @Accept json
// @Produce json
// @Param revoke body dto.IssueCertificationRequest true "Revoke certification data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /certifications/revoke [post]
func (c *CertificationController) RevokeCertification(ctx *gin.Context) {
	var req dto.IssueCertificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.certificationService.RevokeCertification(ctx.Request.Context(), req.UserID, req.CertificationID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "certification revoked"})
}

// GetUserCertifications godoc
// @Summary Get certifications by user ID
// @Description Retrieves all certifications for a specific user
// @Tags certifications
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Success 200 {object} dto.UserCertificationResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /certifications/user/{userId} [get]
func (c *CertificationController) GetUserCertifications(ctx *gin.Context) {
	userID, err := base.GetUintIDFromPath(ctx, "userId")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	resp, err := c.certificationService.GetUserCertifications(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetCertificationsByPackage godoc
// @Summary Get certifications by package ID
// @Description Retrieves all certifications for a specific package
// @Tags certifications
// @Accept json
// @Produce json
// @Param packageId path int true "Package ID"
// @Success 200 {object} dto.CertificationListResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /certifications/package/{packageId} [get]
func (c *CertificationController) GetCertificationsByPackage(ctx *gin.Context) {
	packageID, err := base.GetUintIDFromPath(ctx, "packageId")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid package id"})
		return
	}

	resp, err := c.certificationService.GetCertificationsByPackage(ctx.Request.Context(), packageID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}