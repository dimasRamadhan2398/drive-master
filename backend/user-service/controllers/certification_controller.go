package controllers

import (
	"net/http"

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
	IssueCertificate(ctx *gin.Context)
	RevokeCertificate(ctx *gin.Context)
	GetCertificate(ctx *gin.Context)
	GetMemberCertificates(ctx *gin.Context)
	DownloadCertificatePDF(ctx *gin.Context)
	GetCertificateStats(ctx *gin.Context)
	GetAllCertificates(ctx *gin.Context)
}

func NewCertificationController(certificationService services.ICertificationService) ICertificationController {
	return &CertificationController{
		certificationService: certificationService,
	}
}

// IssueCertificate godoc
// @Summary Issue a certificate to a member
// @Description Admin issues a certificate to a member who completed a package
// @Tags Certificates
// @Accept json
// @Produce json
// @Param request body dto.IssueMemberCertificateInput true "Certificate issuance data"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/certificates [post]
func (c *CertificationController) IssueCertificate(ctx *gin.Context) {
	var input dto.IssueMemberCertificateInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	resp, err := c.certificationService.IssueCertificate(ctx.Request.Context(), input)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusCreated, "Certificate issued successfully", resp)
}

// RevokeCertificate godoc
// @Summary Revoke a certificate
// @Description Admin revokes a certificate
// @Tags Certificates
// @Produce json
// @Param id path string true "Certificate ID (UUID)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/certificates/{id} [delete]
func (c *CertificationController) RevokeCertificate(ctx *gin.Context) {
	certID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	if err := c.certificationService.RevokeCertificate(ctx.Request.Context(), certID); err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Certificate revoked successfully", nil)
}

// GetCertificate godoc
// @Summary Get certificate details
// @Description Get a specific certificate by ID
// @Tags Certificates
// @Produce json
// @Param id path string true "Certificate ID (UUID)"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/certificates/{id} [get]
func (c *CertificationController) GetCertificate(ctx *gin.Context) {
	certID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	resp, err := c.certificationService.GetCertificate(ctx.Request.Context(), certID)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Certificate retrieved successfully", resp)
}

// GetMemberCertificates godoc
// @Summary Get certificates for a member
// @Description Admin views all certificates for a specific member
// @Tags Certificates
// @Produce json
// @Param memberId path string true "Member ID (UUID)"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/certificates/member/{memberId} [get]
func (c *CertificationController) GetMemberCertificates(ctx *gin.Context) {
	memberID, err := uuid.Parse(ctx.Param("memberId"))
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	resp, err := c.certificationService.GetMemberCertificates(ctx.Request.Context(), memberID)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Member certificates retrieved successfully", resp)
}

// DownloadCertificatePDF godoc
// @Summary Download certificate as PDF
// @Description Generate and download a certificate PDF
// @Tags Certificates
// @Produce application/pdf
// @Param id path string true "Certificate ID (UUID)"
// @Success 200 {file} binary
// @Failure 404 {object} response.Response
// @Router /api/v1/certificates/{id}/pdf [get]
func (c *CertificationController) DownloadCertificatePDF(ctx *gin.Context) {
	certID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	pdfBytes, filename, err := c.certificationService.DownloadCertificatePDF(ctx.Request.Context(), certID)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Header("Content-Type", "application/pdf")
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// GetCertificateStats godoc
// @Summary Get certificate statistics
// @Description Get certificate statistics with growth compared to previous month
// @Tags Certificates
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/certificates/stats [get]
func (c *CertificationController) GetCertificateStats(ctx *gin.Context) {
	resp, err := c.certificationService.GetCertificateStats(ctx.Request.Context())
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Certificate stats retrieved successfully", resp)
}

// GetAllCertificates godoc
// @Summary Get all certificates
// @Description Admin views all certificates
// @Tags Certificates
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/certificates [get]
func (c *CertificationController) GetAllCertificates(ctx *gin.Context) {
	resp, err := c.certificationService.GetAllCertificates(ctx.Request.Context())
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Certificates retrieved successfully", resp)
}
