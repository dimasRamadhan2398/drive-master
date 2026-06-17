package controllers

import (
	"net/http"
	"strconv"

	"user-service/models/dto"
	apperrors "user-service/pkg/errors"
	responseRes "user-service/pkg/response"
	"user-service/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CertificationController struct {
	certificationService services.ICertificationService
}

type ICertificationController interface {
	CreateCertification(ctx *gin.Context)
	GetCertification(ctx *gin.Context)
	UpdateCertification(ctx *gin.Context)
	DeleteCertification(ctx *gin.Context)
	ListCertifications(ctx *gin.Context)
	VerifyCertification(ctx *gin.Context)
}

func NewCertificationController(certificationService services.ICertificationService) ICertificationController {
	return &CertificationController{
		certificationService: certificationService,
	}
}

// @Summary Create Certification
// @Description Add a new certification for an instructor
// @Tags Instructors
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param request body dto.CreateCertificationInput true "Certification data"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /instructors/{id}/certifications [post]
func (c *CertificationController) CreateCertification(ctx *gin.Context) {
	instructorID, err := parseUUID(ctx, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	var input dto.CreateCertificationInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	resp, err := c.certificationService.CreateCertification(ctx.Request.Context(), instructorID, input)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusCreated, "Certification created successfully", resp)
}

// @Summary Get Certification
// @Description Get a specific certification for an instructor
// @Tags Instructors
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param certId path string true "Certification ID (UUID)"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /instructors/{id}/certifications/{certId} [get]
func (c *CertificationController) GetCertification(ctx *gin.Context) {
	instructorID, err := parseUUID(ctx, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	certID, err := parseUUID(ctx, "certId")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	resp, err := c.certificationService.GetCertification(ctx.Request.Context(), instructorID, certID)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Certification retrieved successfully", resp)
}

// @Summary Update Certification
// @Description Update a certification for an instructor
// @Tags Instructors
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param certId path string true "Certification ID (UUID)"
// @Param request body dto.UpdateCertificationInput true "Certification data"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /instructors/{id}/certifications/{certId} [put]
func (c *CertificationController) UpdateCertification(ctx *gin.Context) {
	instructorID, err := parseUUID(ctx, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	certID, err := parseUUID(ctx, "certId")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	var input dto.UpdateCertificationInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	resp, err := c.certificationService.UpdateCertification(ctx.Request.Context(), instructorID, certID, input)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Certification updated successfully", resp)
}

// @Summary Delete Certification
// @Description Delete a certification for an instructor
// @Tags Instructors
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param certId path string true "Certification ID (UUID)"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /instructors/{id}/certifications/{certId} [delete]
func (c *CertificationController) DeleteCertification(ctx *gin.Context) {
	instructorID, err := parseUUID(ctx, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	certID, err := parseUUID(ctx, "certId")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	if err := c.certificationService.DeleteCertification(ctx.Request.Context(), instructorID, certID); err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Certification deleted successfully", nil)
}

// @Summary List Certifications
// @Description Get all certifications for an instructor
// @Tags Instructors
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.Response
// @Router /instructors/{id}/certifications [get]
func (c *CertificationController) ListCertifications(ctx *gin.Context) {
	instructorID, err := parseUUID(ctx, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	resp, err := c.certificationService.ListCertifications(ctx.Request.Context(), instructorID, page, limit)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Certifications retrieved successfully", resp)
}

// @Summary Verify Certification
// @Description Verify a certification for an instructor (admin action)
// @Tags Instructors
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param certId path string true "Certification ID (UUID)"
// @Param request body dto.VerifyCertificationInput true "Verification data"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /instructors/{id}/certifications/{certId}/verify [post]
func (c *CertificationController) VerifyCertification(ctx *gin.Context) {
	instructorID, err := parseUUID(ctx, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	certID, err := parseUUID(ctx, "certId")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	var input dto.VerifyCertificationInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	// In a real app, verifiedBy would come from the authenticated user
	// For now, we'll use a placeholder or extract from context
	verifiedBy := uuid.New()

	resp, err := c.certificationService.VerifyCertification(ctx.Request.Context(), instructorID, certID, verifiedBy, input)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Certification verified successfully", resp)
}

// parseUUID parses a UUID from a path parameter
func parseUUID(ctx *gin.Context, param string) (uuid.UUID, error) {
	idStr := ctx.Param(param)
	return uuid.Parse(idStr)
}
