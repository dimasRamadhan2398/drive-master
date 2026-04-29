package controllers

import (
	"net/http"
	"strconv"

	"user-service/models"
	"user-service/models/dto"
	apperrors "user-service/pkg/errors"
	responseRes "user-service/pkg/response"
	"user-service/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	userService       services.IUserService
	authService       services.IAuthService
	memberService     services.IMemberService
	instructorService services.IInstructorService
	roleService       services.IRoleService
	emailService      services.IMailtrapEmailService
	mediaService      services.IMediaService
}

func NewUserController(
	userService services.IUserService,
	authService services.IAuthService,
	memberService services.IMemberService,
	instructorService services.IInstructorService,
	roleService services.IRoleService,
	emailService services.IMailtrapEmailService,
	mediaService services.IMediaService,
) *UserController {
	return &UserController{
		userService:       userService,
		authService:       authService,
		memberService:     memberService,
		instructorService: instructorService,
		roleService:       roleService,
		emailService:      emailService,
		mediaService:      mediaService,
	}
}

// @Summary Login
// @Description Authenticate user with email/username and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginInput true "Login credentials"
// @Success 200 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Failure 401 {object} responseRes.Response
// @Router /auth/login [post]
func (c *UserController) Login(cxt *gin.Context) {
	var input dto.LoginInput
	if err := cxt.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(cxt, apperrors.ErrBadRequest)
		return
	}

	user, err := c.authService.ValidateCredentials(input.Identifier, input.Password)
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	responseRes.Success(cxt, http.StatusOK, "Login successful", user)
}

// @Summary Register
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "User registration data"
// @Success 201 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Failure 409 {object} responseRes.Response
// @Router /auth/register [post]
func (c *UserController) Register(cxt *gin.Context) {
	var input dto.CreateUserRequest
	if err := cxt.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(cxt, apperrors.ErrBadRequest)
		return
	}

	user, err := c.userService.CreateUser(input)
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	// Send welcome email asynchronously
	go c.emailService.SendWelcomeEmail(cxt.Request.Context(), user.Email, user.Username)

	responseRes.Success(cxt, http.StatusCreated, "User registered successfully", user)
}

// @Summary Forgot Password
// @Description Request a password reset link
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ForgotPasswordInput true "Email address"
// @Success 200 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Router /auth/forgot-password [post]
func (c *UserController) ForgotPassword(cxt *gin.Context) {
	var input dto.ForgotPasswordInput
	if err := cxt.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(cxt, apperrors.ErrBadRequest)
		return
	}

	user, err := c.userService.GetUserByEmail(input.EmailAddress)
	if err != nil {
		// Don't reveal if user exists or not
		responseRes.Success(cxt, http.StatusOK, "If the email exists, a reset link has been sent", nil)
		return
	}

	// TODO: Generate reset token and store it
	resetToken := "placeholder-token"

	// Send password reset email asynchronously
	go c.emailService.SendPasswordResetEmail(cxt.Request.Context(), user.Email, resetToken)

	responseRes.Success(cxt, http.StatusOK, "If the email exists, a reset link has been sent", nil)
}

// @Summary Reset Password
// @Description Reset password using a token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordInput true "Reset password data"
// @Success 200 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Router /auth/reset-password [post]
func (c *UserController) ResetPassword(cxt *gin.Context) {
	var input dto.ResetPasswordInput
	if err := cxt.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(cxt, apperrors.ErrBadRequest)
		return
	}

	// TODO: Validate reset token and reset password
	responseRes.Success(cxt, http.StatusOK, "Password has been reset successfully", nil)
}

// @Summary Create User
// @Description Create a new user (admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "User data"
// @Success 201 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Failure 409 {object} responseRes.Response
// @Router /users [post]
func (c *UserController) CreateUser(cxt *gin.Context) {
	var input dto.CreateUserRequest
	if err := cxt.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(cxt, apperrors.ErrBadRequest)
		return
	}

	user, err := c.userService.CreateUser(input)
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	responseRes.Success(cxt, http.StatusCreated, "User created successfully", user)
}

// @Summary Get All Users
// @Description Get all users with their profiles
// @Tags Users
// @Produce json
// @Success 200 {object} responseRes.Response
// @Failure 500 {object} responseRes.Response
// @Router /users [get]
func (c *UserController) GetAllUsers(cxt *gin.Context) {
	users, err := c.userService.GetAllUsersWithProfiles()
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	responseRes.Success(cxt, http.StatusOK, "Users retrieved successfully", users)
}

// @Summary Get User by ID
// @Description Get a user by their ID with profiles
// @Tags Users
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Failure 404 {object} responseRes.Response
// @Router /users/{id} [get]
func (c *UserController) GetUserByID(cxt *gin.Context) {
	id, err := getUserIDFromPath(cxt, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	user, err := c.userService.GetUserByIDWithProfiles(id)
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	responseRes.Success(cxt, http.StatusOK, "User retrieved successfully", user)
}

// @Summary Update User
// @Description Update user information
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param request body dto.UpdateUserRequest true "Update user data"
// @Success 200 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Failure 404 {object} responseRes.Response
// @Router /users/{id} [put]
func (c *UserController) UpdateUser(cxt *gin.Context) {
	id, err := getUserIDFromPath(cxt, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	var input dto.UpdateUserRequest
	if err := cxt.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(cxt, apperrors.ErrBadRequest)
		return
	}

	user, err := c.userService.GetUserByID(id)
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	// Convert DTO to model for update
	userModel := &models.User{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		EmailAddress: user.Email,
		PhoneNumber:  user.PhoneNumber,
		RoleID:       user.RoleID,
	}

	// Update fields if provided
	if input.Username != nil {
		userModel.Username = *input.Username
	}
	if input.EmailAddress != nil {
		userModel.Email = *input.EmailAddress
		userModel.EmailAddress = *input.EmailAddress
	}
	if input.PhoneNumber != nil {
		userModel.PhoneNumber = *input.PhoneNumber
	}
	if input.DateOfBirth != nil {
		userModel.DateOfBirth = *input.DateOfBirth
	}
	if input.Image != nil {
		userModel.Image = *input.Image
	}
	if input.Address != nil {
		userModel.Address = *input.Address
	}
	if input.RoleID != nil {
		userModel.RoleID = *input.RoleID
	}

	if err := c.userService.UpdateUser(userModel); err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	responseRes.Success(cxt, http.StatusOK, "User updated successfully", user)
}

// @Summary Delete User
// @Description Delete a user by ID
// @Tags Users
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Failure 404 {object} responseRes.Response
// @Router /users/{id} [delete]
func (c *UserController) DeleteUser(cxt *gin.Context) {
	id, err := getUserIDFromPath(cxt, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	user, err := c.userService.GetUserByID(id)
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	// Convert DTO to model for delete
	userModel := &models.User{
		ID: user.ID,
	}

	if err := c.userService.DeleteUser(userModel); err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	responseRes.Success(cxt, http.StatusOK, "User deleted successfully", nil)
}

// @Summary Change Password
// @Description Change user password
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param request body dto.ChangePasswordInput true "Password change data"
// @Success 200 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Failure 401 {object} responseRes.Response
// @Router /users/{id}/password [patch]
func (c *UserController) ChangePassword(cxt *gin.Context) {
	id, err := getUserIDFromPath(cxt, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	var input dto.ChangePasswordInput
	if err := cxt.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(cxt, apperrors.ErrBadRequest)
		return
	}

	if err := c.authService.ChangePassword(id, input.CurrentPassword, input.NewPassword); err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	responseRes.Success(cxt, http.StatusOK, "Password changed successfully", nil)
}

// @Summary Get Member Profile
// @Description Get member profile by user ID
// @Tags Members
// @Produce json
// @Param userId path string true "User ID (UUID)"
// @Success 200 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Failure 404 {object} responseRes.Response
// @Router /members/{userId}/profile [get]
func (c *UserController) GetMemberProfile(cxt *gin.Context) {
	userID, err := getUserIDFromPath(cxt, "userId")
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	profile, err := c.memberService.GetMemberProfile(userID)
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	responseRes.Success(cxt, http.StatusOK, "Member profile retrieved successfully", profile)
}

// @Summary Update Member Profile
// @Description Update member profile
// @Tags Members
// @Accept json
// @Produce json
// @Param userId path string true "User ID (UUID)"
// @Param request body dto.UpdateUserRequest true "Update profile data"
// @Success 200 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Failure 404 {object} responseRes.Response
// @Router /members/{userId}/profile [put]
func (c *UserController) UpdateMemberProfile(cxt *gin.Context) {
	userID, err := getUserIDFromPath(cxt, "userId")
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	profile, err := c.memberService.GetMemberProfile(userID)
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	// Bind request body
	var input struct {
		// Add member profile update fields here
	}
	if err := cxt.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(cxt, apperrors.ErrBadRequest)
		return
	}

	if err := c.memberService.UpdateMemberProfile(profile); err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	responseRes.Success(cxt, http.StatusOK, "Member profile updated successfully", profile)
}

// @Summary Get All Instructors
// @Description Get all instructors with their profiles
// @Tags Instructors
// @Produce json
// @Success 200 {object} responseRes.Response
// @Failure 500 {object} responseRes.Response
// @Router /instructors [get]
func (c *UserController) GetAllInstructors(cxt *gin.Context) {
	users, err := c.userService.GetAllUsersWithProfiles()
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	responseRes.Success(cxt, http.StatusOK, "Instructors retrieved successfully", users)
}

// @Summary Get Instructor Profile
// @Description Get instructor profile by user ID with work experiences
// @Tags Instructors
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Failure 404 {object} responseRes.Response
// @Router /instructors/{id}/profile [get]
func (c *UserController) GetInstructorProfile(cxt *gin.Context) {
	userID, err := getUserIDFromPath(cxt, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	profile, err := c.instructorService.GetInstructorProfile(userID)
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	workExps, err := c.instructorService.GetWorkExperiences(userID)
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	responseRes.Success(cxt, http.StatusOK, "Instructor profile retrieved successfully", gin.H{
		"profile":         profile,
		"workExperiences": workExps,
	})
}

// @Summary Update Instructor Profile
// @Description Update instructor profile
// @Tags Instructors
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param request body dto.UpdateInstructorProfileInput true "Update profile data"
// @Success 200 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Failure 404 {object} responseRes.Response
// @Router /instructors/{id}/profile [put]
func (c *UserController) UpdateInstructorProfile(cxt *gin.Context) {
	userID, err := getUserIDFromPath(cxt, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	var input dto.UpdateInstructorProfileInput
	if err := cxt.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(cxt, apperrors.ErrBadRequest)
		return
	}

	profile, err := c.instructorService.GetInstructorProfile(userID)
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	// Convert DTO to model for update
	profileModel := &models.InstructorProfile{
		UserID:                profile.UserID,
		LicenseNumber:          profile.LicenseNumber,
		YearsOfExperience:      profile.YearsOfExperience,
		Bio:                    profile.Bio,
		LicenseExpiry:          profile.LicenseExpiry,
		PhotoURL:               profile.PhotoURL,
		IsActive:               profile.IsActive,
		NumberOfStudents:       profile.NumberOfStudents,
		SessionsCompleted:     profile.SessionsCompleted,
		AverageRating:         profile.AverageRating,
		BNSPCertificateNumber:  profile.BNSPCertificateNumber,
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

	if err := c.instructorService.UpdateInstructorProfile(profileModel); err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	responseRes.Success(cxt, http.StatusOK, "Instructor profile updated successfully", profile)
}

// @Summary Create Work Experience
// @Description Add work experience for an instructor
// @Tags Instructors
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param request body dto.CreateWorkExperienceRequest true "Work experience data"
// @Success 201 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Failure 404 {object} responseRes.Response
// @Router /instructors/{id}/work-experiences [post]
func (c *UserController) CreateWorkExperience(cxt *gin.Context) {
	userID, err := getUserIDFromPath(cxt, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	var input dto.CreateWorkExperienceRequest
	if err := cxt.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(cxt, apperrors.ErrBadRequest)
		return
	}

	workExp, err := c.instructorService.CreateWorkExperience(userID, input)
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	responseRes.Success(cxt, http.StatusCreated, "Work experience created successfully", workExp)
}

// @Summary Update Work Experience
// @Description Update a work experience
// @Tags Instructors
// @Accept json
// @Produce json
// @Param expId path int true "Work Experience ID"
// @Param request body dto.CreateWorkExperienceRequest true "Update work experience data"
// @Success 200 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Failure 404 {object} responseRes.Response
// @Router /work-experiences/{expId} [put]
func (c *UserController) UpdateWorkExperience(cxt *gin.Context) {
	expID, err := getUintIDFromPath(cxt, "expId")
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	var input dto.CreateWorkExperienceRequest
	if err := cxt.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(cxt, apperrors.ErrBadRequest)
		return
	}

	// TODO: Get work experience by ID and update
	workExp := &struct {
		ID           uint
		InstructorID uuid.UUID
		CompanyName  string
		Role         string
	}{ID: expID}

	workExp.CompanyName = input.CompanyName
	workExp.Role = input.Role

	// This needs the repo to fetch by ID - for now placeholder
	responseRes.Success(cxt, http.StatusOK, "Work experience updated successfully", workExp)
}

// @Summary Delete Work Experience
// @Description Delete a work experience
// @Tags Instructors
// @Produce json
// @Param expId path int true "Work Experience ID"
// @Success 200 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Failure 404 {object} responseRes.Response
// @Router /work-experiences/{expId} [delete]
func (c *UserController) DeleteWorkExperience(cxt *gin.Context) {
	expID, err := getUintIDFromPath(cxt, "expId")
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	if err := c.instructorService.DeleteWorkExperience(expID); err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	responseRes.Success(cxt, http.StatusOK, "Work experience deleted successfully", nil)
}

// @Summary Add Coverage Area
// @Description Add coverage area for an instructor
// @Tags Instructors
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param request body dto.AddCoverageAreaInput true "Coverage area data"
// @Success 201 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Failure 404 {object} responseRes.Response
// @Router /instructors/{id}/coverage-areas [post]
func (c *UserController) AddCoverageArea(cxt *gin.Context) {
	userID, err := getUserIDFromPath(cxt, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	var input dto.AddCoverageAreaInput
	if err := cxt.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(cxt, apperrors.ErrBadRequest)
		return
	}

	// TODO: Look up area ID by name from location service
	// For now, use a placeholder area ID
	areaID := uint(1)

	if err := c.instructorService.AddCoverageArea(userID, areaID); err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	responseRes.Success(cxt, http.StatusCreated, "Coverage area added successfully", nil)
}

// @Summary Remove Coverage Area
// @Description Remove coverage area for an instructor
// @Tags Instructors
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param areaId path int true "Coverage Area ID"
// @Success 200 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Failure 404 {object} responseRes.Response
// @Router /instructors/{id}/coverage-areas/{areaId} [delete]
func (c *UserController) RemoveCoverageArea(cxt *gin.Context) {
	userID, err := getUserIDFromPath(cxt, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	areaID, err := getUintIDFromPath(cxt, "areaId")
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	if err := c.instructorService.RemoveCoverageArea(userID, areaID); err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	responseRes.Success(cxt, http.StatusOK, "Coverage area removed successfully", nil)
}

// @Summary Upload Media
// @Description Upload a file
// @Tags Media
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param fileName formData string false "File name"
// @Param folder formData string false "Folder path"
// @Success 200 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Router /media/upload [post]
func (c *UserController) UploadMedia(cxt *gin.Context) {
	file, _, err := cxt.Request.FormFile("file")
	if err != nil {
		responseRes.ErrorFromAppError(cxt, apperrors.ErrBadRequest)
		return
	}
	defer file.Close()

	// Read file data
	fileData := make([]byte, 1024*1024) // 1MB max
	n, err := file.Read(fileData)
	if err != nil {
		responseRes.ErrorFromAppError(cxt, apperrors.ErrBadRequest)
		return
	}
	fileData = fileData[:n]

	fileName := cxt.PostForm("fileName")
	folder := cxt.PostForm("folder")

	resp, err := c.mediaService.UploadMedia(cxt.Request.Context(), services.UploadMediaInput{
		FileData: fileData,
		FileName: fileName,
		Folder:   folder,
	})
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	responseRes.Success(cxt, http.StatusOK, "Media uploaded successfully", resp)
}

// @Summary Upload Base64 Media
// @Description Upload a file from base64 encoded data
// @Tags Media
// @Accept json
// @Produce json
// @Param request body dto.UploadBase64MediaRequest true "Base64 encoded file data"
// @Success 200 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Router /media/upload-base64 [post]
func (c *UserController) UploadBase64Media(cxt *gin.Context) {
	var input dto.UploadBase64MediaRequest
	if err := cxt.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(cxt, apperrors.ErrBadRequest)
		return
	}

	resp, err := c.mediaService.UploadBase64Media(
		cxt.Request.Context(),
		input.Base64Data,
		input.FileName,
		input.Folder,
	)
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	responseRes.Success(cxt, http.StatusOK, "Media uploaded successfully", resp)
}

// @Summary Delete Media
// @Description Delete a media file by ID
// @Tags Media
// @Produce json
// @Param fileId path string true "File ID"
// @Success 200 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Failure 404 {object} responseRes.Response
// @Router /media/{fileId} [delete]
func (c *UserController) DeleteMedia(cxt *gin.Context) {
	fileID := cxt.Param("fileId")
	if fileID == "" {
		responseRes.ErrorFromAppError(cxt, apperrors.ErrBadRequest)
		return
	}

	if err := c.mediaService.DeleteMedia(cxt.Request.Context(), fileID); err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	responseRes.Success(cxt, http.StatusOK, "Media deleted successfully", nil)
}

// @Summary Get Media Metadata
// @Description Get metadata for a media file
// @Tags Media
// @Produce json
// @Param fileId path string true "File ID"
// @Success 200 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Failure 404 {object} responseRes.Response
// @Router /media/{fileId}/metadata [get]
func (c *UserController) GetMediaMetadata(cxt *gin.Context) {
	fileID := cxt.Param("fileId")
	if fileID == "" {
		responseRes.ErrorFromAppError(cxt, apperrors.ErrBadRequest)
		return
	}

	resp, err := c.mediaService.GetMediaMetadata(cxt.Request.Context(), fileID)
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	responseRes.Success(cxt, http.StatusOK, "Media metadata retrieved successfully", resp)
}

// @Summary Get All Roles
// @Description Get all available roles
// @Tags Roles
// @Produce json
// @Success 200 {object} responseRes.Response
// @Failure 500 {object} responseRes.Response
// @Router /roles [get]
func (c *UserController) GetAllRoles(cxt *gin.Context) {
	roles, err := c.roleService.GetAllRoles()
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	responseRes.Success(cxt, http.StatusOK, "Roles retrieved successfully", roles)
}

// @Summary Get Role by ID
// @Description Get a role by ID
// @Tags Roles
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Failure 404 {object} responseRes.Response
// @Router /roles/{id} [get]
func (c *UserController) GetRole(cxt *gin.Context) {
	id, err := getUintIDFromPath(cxt, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	role, err := c.roleService.GetRole(id)
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	responseRes.Success(cxt, http.StatusOK, "Role retrieved successfully", role)
}

// RegisterRoutes registers all user routes
func (c *UserController) RegisterRoutes(router *gin.Engine) {
	users := router.Group("/users")
	{
		users.POST("", c.CreateUser)
		users.GET("", c.GetAllUsers)
		users.GET("/:id", c.GetUserByID)
		users.PUT("/:id", c.UpdateUser)
		users.DELETE("/:id", c.DeleteUser)
		users.PATCH("/:id/password", c.ChangePassword)
	}

	auth := router.Group("/auth")
	{
		auth.POST("/login", c.Login)
		auth.POST("/register", c.Register)
		auth.POST("/forgot-password", c.ForgotPassword)
		auth.POST("/reset-password", c.ResetPassword)
	}

	members := router.Group("/members")
	{
		members.GET("/:userId/profile", c.GetMemberProfile)
		members.PUT("/:userId/profile", c.UpdateMemberProfile)
	}

	instructors := router.Group("/instructors")
	{
		instructors.GET("", c.GetAllInstructors)
		instructors.GET("/:id/profile", c.GetInstructorProfile)
		instructors.PUT("/:id/profile", c.UpdateInstructorProfile)
		instructors.POST("/:id/work-experiences", c.CreateWorkExperience)
		instructors.PUT("/work-experiences/:expId", c.UpdateWorkExperience)
		instructors.DELETE("/work-experiences/:expId", c.DeleteWorkExperience)
		instructors.POST("/:id/coverage-areas", c.AddCoverageArea)
		instructors.DELETE("/:id/coverage-areas/:areaId", c.RemoveCoverageArea)
	}

	media := router.Group("/media")
	{
		media.POST("/upload", c.UploadMedia)
		media.POST("/upload-base64", c.UploadBase64Media)
		media.DELETE("/:fileId", c.DeleteMedia)
		media.GET("/:fileId/metadata", c.GetMediaMetadata)
	}

	roles := router.Group("/roles")
	{
		roles.GET("", c.GetAllRoles)
		roles.GET("/:id", c.GetRole)
	}
}

// Helper to get user ID from path (returns UUID)
func getUserIDFromPath(c *gin.Context, paramName string) (uuid.UUID, error) {
	idStr := c.Param(paramName)
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, apperrors.ErrBadRequest
	}
	return id, nil
}

// Helper to get uint ID from path
func getUintIDFromPath(c *gin.Context, paramName string) (uint, error) {
	idStr := c.Param(paramName)
	if idStr == "" {
		return 0, apperrors.ErrBadRequest
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, apperrors.ErrBadRequest
	}
	return uint(id), nil
}