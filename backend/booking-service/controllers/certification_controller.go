package controllers

import (
	"net/http"
	"strconv"

	"booking-service/models/dto"
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

func (c *CertificationController) GetCertification(ctx *gin.Context) {
	id, err := getUintIDFromPath(ctx, "id")
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

func (c *CertificationController) UpdateCertificationStatus(ctx *gin.Context) {
	id, err := getUintIDFromPath(ctx, "id")
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

func (c *CertificationController) GetUserCertifications(ctx *gin.Context) {
	userID, err := getUintIDFromPath(ctx, "userId")
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

func (c *CertificationController) GetCertificationsByPackage(ctx *gin.Context) {
	packageID, err := getUintIDFromPath(ctx, "packageId")
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