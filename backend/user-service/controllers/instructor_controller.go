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
	userService            services.IUserService
	authService            services.IAuthService
	memberService          services.IMemberService
	instructorService      services.IInstructorService
	roleService            services.IRoleService
	emailService           services.IMailtrapEmailService
	mediaService           services.IMediaService
	recurringScheduleService services.IRecurringScheduleService
}

type IInstructorController interface {
	GetInstructorProfile(ctx *gin.Context)
	UpdateInstructorProfile(ctx *gin.Context)
	DeleteInstructor(ctx *gin.Context)
	GetInstructorLists(ctx *gin.Context)
	CreateInstructorProfile(ctx *gin.Context)
	RegisterInstructor(ctx *gin.Context)
	UploadProfilePic(ctx *gin.Context)
	DeleteProfilePic(ctx *gin.Context)
	UploadBase64Media(ctx *gin.Context)
	GetMediaMetadata(ctx *gin.Context)
	GetAllInstructorsWithRecurringSchedules(ctx *gin.Context)
}

func NewInstructorController(
	userService services.IUserService,
	authService services.IAuthService,
	memberService services.IMemberService,
	instructorService services.IInstructorService,
	roleService services.IRoleService,
	emailService services.IMailtrapEmailService,
	mediaService services.IMediaService,
	recurringScheduleService services.IRecurringScheduleService,
) IInstructorController {
	return &InstructorController{
		userService:            userService,
		authService:            authService,
		memberService:          memberService,
		instructorService:      instructorService,
		roleService:            roleService,
		emailService:           emailService,
		mediaService:           mediaService,
		recurringScheduleService: recurringScheduleService,
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

	responseRes.Paginated(ctx, http.StatusOK, "Instructors retrieved successfully", result.Data, &result.Pagination)
}

// @Summary Create Instructor Profile
// @Description Create an instructor profile for a user with optional data
// @Tags Instructors
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param request body dto.InstructorProfileRequest false "Instructor profile data (optional)"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /instructors [post]
func (c *InstructorController) CreateInstructorProfile(ctx *gin.Context) {
	userID, err := getUserIDFromPath(ctx, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	var req dto.InstructorProfileRequest
	// BindJSON will succeed even if body is empty (all fields are optional)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// If binding fails completely (e.g., invalid JSON), still allow creation with defaults
		// Only return error if it's a syntax issue
		if ctx.Request.ContentLength > 0 {
			responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
			return
		}
	}

	var result *dto.InstructorProfileResponse

	// Check if any field was provided in the request
	hasInput := req.LicenseNumber != nil || req.LicenseExpiry.Unix() > 0 ||
		req.BNSPCertificateNumber != nil || req.Description != nil ||
		req.YearsOfExperience != nil

	if hasInput {
		result, err = c.instructorService.CreateInstructorProfileWithInput(ctx.Request.Context(), userID, req)
	} else {
		result, err = c.instructorService.CreateInstructorProfile(ctx.Request.Context(), userID)
	}

	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusCreated, "Instructor profile created successfully", result)
}

// @Summary Register Instructor
// @Description Register a new instructor with user account and instructor profile in a single transaction
// @Tags Instructors
// @Accept json
// @Produce json
// @Param request body dto.CreateInstructorWithUserRequest true "Instructor registration data"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /instructors/register [post]
func (c *InstructorController) RegisterInstructor(ctx *gin.Context) {
	var req dto.CreateInstructorWithUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	result, err := c.instructorService.CreateInstructorWithUser(ctx.Request.Context(), req)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusCreated, "Instructor registered successfully", result)
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

// @Summary Delete Instructor
// @Description Soft delete an instructor by changing role to 'member'
// @Tags Instructors
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /instructors/{id} [delete]
func (c *InstructorController) DeleteInstructor(ctx *gin.Context) {
	userID, err := getUserIDFromPath(ctx, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	// Get the instructor profile first to verify it exists
	profile, err := c.instructorService.GetInstructorProfile(ctx.Request.Context(), userID)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	// Delete the instructor profile
	if err := c.instructorService.DeleteInstructorProfile(ctx.Request.Context(), userID); err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	// Get member role and update user role
	memberRole, err := c.roleService.GetRoleByName(ctx, "member")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	// Update user role to member
	if err := c.roleService.UpdateUserRole(ctx.Request.Context(), userID, memberRole.ID); err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Instructor deleted successfully", gin.H{
		"userId":      profile.UserID,
		"newRole":     "member",
		"deletedAt":   time.Now(),
	})
}

// @Summary Upload Media
// @Description Upload a profile picture for an instructor
// @Tags Media
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param file formData file true "File to upload"
// @Param fileName formData string false "File name"
// @Param folder formData string false "Folder path"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /instructors/{id}/media/upload [post]
func (c *InstructorController) UploadProfilePic(ctx *gin.Context) {
	userID, err := getUserIDFromPath(ctx, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

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

	// Update only the photo URL field
	if err := c.instructorService.UpdateInstructorPhotoURL(ctx.Request.Context(), userID, resp.URL); err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Media uploaded successfully", resp)
}

// @Summary Upload Base64 Media
// @Description Upload a file from base64 encoded data for an instructor
// @Tags Media
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param request body dto.UploadBase64MediaRequest true "Base64 encoded file data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /instructors/{id}/media/upload-base64 [post]
func (c *InstructorController) UploadBase64Media(ctx *gin.Context) {
	userID, err := getUserIDFromPath(ctx, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

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

	// Update only the photo URL field
	if err := c.instructorService.UpdateInstructorPhotoURL(ctx.Request.Context(), userID, resp.URL); err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Media uploaded successfully", resp)
}


// @Summary Delete Media
// @Description Delete a media file and clear instructor photo URL
// @Tags Media
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /instructors/{id}/media [delete]
func (c *InstructorController) DeleteProfilePic(ctx *gin.Context) {
	userID, err := getUserIDFromPath(ctx, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	// Get the current profile to get the photo URL
	profile, err := c.instructorService.GetInstructorProfile(ctx.Request.Context(), userID)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	// Delete from media service if there's a photo
	if profile.PhotoURL != "" {
		if err := c.mediaService.DeleteMedia(ctx.Request.Context(), profile.PhotoURL); err != nil {
			responseRes.ErrorFromGeneric(ctx, err)
			return
		}
	}

	// Clear the photo URL in the instructor profile
	profileModel := &models.InstructorProfile{
		UserID:                profile.UserID,
		LicenseNumber:         profile.LicenseNumber,
		YearsOfExperience:     profile.YearsOfExperience,
		Bio:                   profile.Bio,
		LicenseExpiry:         profile.LicenseExpiry,
		PhotoURL:              "",
		IsActive:              profile.IsActive,
		NumberOfStudents:      profile.NumberOfStudents,
		SessionsCompleted:     profile.SessionsCompleted,
		AverageRating:         profile.AverageRating,
		BNSPCertificateNumber: profile.BNSPCertificateNumber,
	}

	if err := c.instructorService.UpdateInstructorProfile(ctx.Request.Context(), profileModel); err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Media deleted successfully", nil)
}

// @Summary Get Media Metadata
// @Description Get metadata for a media file
// @Tags Media
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /instructors/{id}/media/metadata [get]
func (c *InstructorController) GetMediaMetadata(ctx *gin.Context) {
	userID, err := getUserIDFromPath(ctx, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	// Get the current profile to get the photo URL
	profile, err := c.instructorService.GetInstructorProfile(ctx.Request.Context(), userID)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	if profile.PhotoURL == "" {
		responseRes.Error(ctx, http.StatusNotFound, "No media found", "Instructor has no profile photo", "")
		return
	}

	resp, err := c.mediaService.GetMediaMetadata(ctx.Request.Context(), profile.PhotoURL)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Media metadata retrieved successfully", resp)
}

// @Summary Get All Instructors with Recurring Schedules
// @Description Get all instructors with their recurring schedules for schedule generation
// @Tags Instructors
// @Produce json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /instructors/with-schedules [get]
func (c *InstructorController) GetAllInstructorsWithRecurringSchedules(ctx *gin.Context) {
	// Get all instructors
	instructors, err := c.userService.GetAllInstructorsForScheduling(ctx.Request.Context())
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	// Build response with recurring schedules
	type InstructorWithSchedules struct {
		ID                 string                           `json:"id"`
		FirstName          string                           `json:"firstName"`
		LastName           string                           `json:"lastName"`
		Email              string                           `json:"email"`
		RecurringSchedules []dto.RecurringScheduleResponse  `json:"recurringSchedules"`
	}

	var result []InstructorWithSchedules
	for _, instructor := range instructors {
		// Get recurring schedules for this instructor
		schedules, err := c.recurringScheduleService.GetRecurringSchedules(ctx.Request.Context(), instructor.ID)
		if err != nil {
			responseRes.ErrorFromGeneric(ctx, err)
			return
		}

		item := InstructorWithSchedules{
			ID:                 instructor.ID.String(),
			FirstName:          instructor.FirstName,
			LastName:           instructor.LastName,
			Email:              instructor.EmailAddress,
			RecurringSchedules: schedules,
		}
		result = append(result, item)
	}

	responseRes.Success(ctx, http.StatusOK, "Instructors with recurring schedules retrieved successfully", result)
}