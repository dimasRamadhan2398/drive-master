package controllers

import (
	"net/http"
	"time"

	"user-service/models"
	"user-service/models/dto"
	apperrors "user-service/pkg/errors"
	responseRes "user-service/pkg/response"
	"user-service/services"

	"github.com/gin-gonic/gin"
)

type InstructorController struct {
	userService       services.IUserService
	authService       services.IAuthService
	memberService     services.IMemberService
	instructorService services.IInstructorService
	roleService       services.IRoleService
	emailService      services.IMailtrapEmailService
	mediaService      services.IMediaService
}

type IInstructorController interface {
	GetInstructorProfile(ctx *gin.Context)
	UpdateInstructorProfile(ctx *gin.Context)
	GetInstructorLists(ctx *gin.Context)
	UploadProfilePic(ctx *gin.Context)
	DeleteProfilePic(ctx *gin.Context)
	UploadBase64Media(ctx *gin.Context)
	GetMediaMetadata(ctx *gin.Context)
}

func NewInstructorController(
	userService services.IUserService,
	authService services.IAuthService,
	memberService services.IMemberService,
	instructorService services.IInstructorService,
	roleService services.IRoleService,
	emailService services.IMailtrapEmailService,
	mediaService services.IMediaService,
) IInstructorController {
	return &InstructorController{
		userService:       userService,
		authService:       authService,
		memberService:     memberService,
		instructorService: instructorService,
		roleService:       roleService,
		emailService:      emailService,
		mediaService:      mediaService,
	}
}

// @Summary Get All Instructors
// @Description Get all instructors with their profiles (paginated)
// @Tags Instructors
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10, max: 100)"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /instructors [get]
func (c *InstructorController) GetInstructorLists(ctx *gin.Context) {
	var query dto.PaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	// Set defaults
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 10
	}
	if query.Limit > 100 {
		query.Limit = 100
	}

	result, err := c.userService.GetInstructorsWithPagination(ctx.Request.Context(), query.Page, query.Limit)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Instructors retrieved successfully", result)
}

// @Summary Get Instructor Profile
// @Description Get instructor profile by user ID
// @Tags Instructors
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /instructors/{id}/profile [get]
func (c *InstructorController) GetInstructorProfile(ctx *gin.Context) {
	userID, err := getUserIDFromPath(ctx, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	profile, err := c.instructorService.GetInstructorProfile(ctx.Request.Context(), userID)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Instructor profile retrieved successfully", profile)
}

// @Summary Update Instructor Profile
// @Description Update instructor profile
// @Tags Instructors
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param request body dto.UpdateInstructorProfileInput true "Update profile data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /instructors/{id}/profile [put]
func (c *InstructorController) UpdateInstructorProfile(ctx *gin.Context) {
	userID, err := getUserIDFromPath(ctx, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	var input dto.UpdateInstructorProfileInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	profile, err := c.instructorService.GetInstructorProfile(ctx.Request.Context(), userID)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	// Convert DTO to model for update
	profileModel := &models.InstructorProfile{
		UserID:                profile.UserID,
		LicenseNumber:         profile.LicenseNumber,
		YearsOfExperience:     profile.YearsOfExperience,
		Bio:                   profile.Bio,
		LicenseExpiry:         profile.LicenseExpiry,
		PhotoURL:              profile.PhotoURL,
		IsActive:              profile.IsActive,
		NumberOfStudents:      profile.NumberOfStudents,
		SessionsCompleted:     profile.SessionsCompleted,
		AverageRating:         profile.AverageRating,
		BNSPCertificateNumber: profile.BNSPCertificateNumber,
	}

	if input.Description != nil {
		profileModel.Bio = *input.Description
	}
	if input.LicenseNumber != nil {
		profileModel.LicenseNumber = *input.LicenseNumber
	}
	if input.YearsOfExperience != nil {
		profileModel.YearsOfExperience = *input.YearsOfExperience
	}
	if input.LicenseExpiry != nil {
		// Parse the string date in DD/MM/YYYY format
		parsedTime, err := time.Parse("02/01/2006", *input.LicenseExpiry)
		if err != nil {
			responseRes.Error(ctx, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid date format for license expiry. Use DD/MM/YYYY", "")
			return
		}
		profileModel.LicenseExpiry = parsedTime
	}
	if input.BNSPCertificateNumber != nil {
		profileModel.BNSPCertificateNumber = *input.BNSPCertificateNumber
	}

	if err := c.instructorService.UpdateInstructorProfile(ctx.Request.Context(), profileModel); err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Instructor profile updated successfully", profile)
}

// @Summary Upload Media
// @Description Upload a file
// @Tags Media
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param fileName formData string false "File name"
// @Param folder formData string false "Folder path"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /media/upload [post]
func (c *InstructorController) UploadProfilePic(ctx *gin.Context) {
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}
	defer file.Close()

	// Read file data
	fileData := make([]byte, 1024*1024) // 1MB max
	n, err := file.Read(fileData)
	if err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}
	fileData = fileData[:n]

	fileName := ctx.PostForm("fileName")
	folder := ctx.PostForm("folder")

	resp, err := c.mediaService.UploadMedia(ctx.Request.Context(), services.UploadMediaInput{
		FileData: fileData,
		FileName: fileName,
		Folder:   folder,
	})
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Media uploaded successfully", resp)
}

// @Summary Upload Base64 Media
// @Description Upload a file from base64 encoded data
// @Tags Media
// @Accept json
// @Produce json
// @Param request body dto.UploadBase64MediaRequest true "Base64 encoded file data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /media/upload-base64 [post]
func (c *InstructorController) UploadBase64Media(ctx *gin.Context) {
	var input dto.UploadBase64MediaRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	resp, err := c.mediaService.UploadBase64Media(
		ctx.Request.Context(),
		input.Base64Data,
		input.FileName,
		input.Folder,
	)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Media uploaded successfully", resp)
}

// @Summary Delete Media
// @Description Delete a media file by ID
// @Tags Media
// @Produce json
// @Param fileId path string true "File ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /media/{fileId} [delete]
func (c *InstructorController) DeleteProfilePic(ctx *gin.Context) {
	fileID := ctx.Param("fileId")
	if fileID == "" {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	if err := c.mediaService.DeleteMedia(ctx.Request.Context(), fileID); err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Media deleted successfully", nil)
}

// @Summary Get Media Metadata
// @Description Get metadata for a media file
// @Tags Media
// @Produce json
// @Param fileId path string true "File ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /media/{fileId}/metadata [get]
func (c *InstructorController) GetMediaMetadata(ctx *gin.Context) {
	fileID := ctx.Param("fileId")
	if fileID == "" {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	resp, err := c.mediaService.GetMediaMetadata(ctx.Request.Context(), fileID)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Media metadata retrieved successfully", resp)
}