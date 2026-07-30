package controllers

import (
	"net/http"
	"time"
	"user-service/models/dto"
	apperrors "user-service/pkg/errors"
	responseRes "user-service/pkg/response"
	"user-service/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MemberController struct {
	userService        services.IUserService
	authService        services.IAuthService
	memberService      services.IMemberService
	certService       services.ICertificationService
	roleService        services.IRoleService
	emailService       services.IMailtrapEmailService
	mediaService       services.IMediaService
}

func NewMemberController(
	userService services.IUserService,
	authService services.IAuthService,
	memberService services.IMemberService,
	certService services.ICertificationService,
	roleService services.IRoleService,
	emailService services.IMailtrapEmailService,
	mediaService services.IMediaService,
) IMemberController {
	return &MemberController{
		userService:        userService,
		authService:        authService,
		memberService:      memberService,
		certService:       certService,
		roleService:        roleService,
		emailService:       emailService,
		mediaService:       mediaService,
	}
}

type IMemberController interface {
	GetMemberProfile(ctx *gin.Context)
	UpdateMemberProfile(ctx *gin.Context)
	GetMemberLists(ctx *gin.Context)
	FindRecentRegistrations(ctx *gin.Context)
	SearchMembersWithPagination(ctx *gin.Context)
	GetMemberCertificates(ctx *gin.Context)
	DownloadMemberCertificate(ctx *gin.Context)
	CreateMember(ctx *gin.Context)
}

// @Summary Get Member Profile
// @Description Get member profile by user ID
// @Tags Members
// @Produce json
// @Param userId path string true "User ID (UUID)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /members/{userId}/profile [get]
func (m *MemberController) GetMemberProfile(ctx *gin.Context) {
	userID, err := getUserIDFromPath(ctx, "userId")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	profile, err := m.memberService.GetMemberProfile(ctx, userID)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Member profile retrieved successfully", profile)
}

// @Summary Update Member Profile
// @Description Update member profile
// @Tags Members
// @Accept json
// @Produce json
// @Param userId path string true "User ID (UUID)"
// @Param request body dto.UpdateUserRequest true "Update profile data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /members/{userId}/profile [put]
func (m *MemberController) UpdateMemberProfile(ctx *gin.Context) {
	userID, err := getUserIDFromPath(ctx, "userId")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	profile, err := m.memberService.GetRawMemberProfile(ctx, userID)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	// Bind request body
	var input dto.UpdateMemberProfileRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	// Only update fields that are provided (non-empty)
	if input.IdentityFullname != "" {
		profile.IdentityFullname = input.IdentityFullname
	}

	if err := m.memberService.UpdateMemberProfile(ctx, profile); err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Member profile updated successfully", profile)
}

// @Summary Get All Members
// @Description Get all members with their profiles (paginated)
// @Tags Members
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10, max: 100)"
// @Param search query string false "Search query"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /members [get]
func (m *MemberController) GetMemberLists(ctx *gin.Context) {
	var query struct {
		Page   int    `form:"page,default=1"`
		Limit  int    `form:"limit,default=10"`
		Search string `form:"search"`
	}
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

	// Pass search query as pointer (nil if empty)
	var searchQuery *string
	if query.Search != "" {
		searchQuery = &query.Search
	}

	result, err := m.memberService.GetMembersWithPagination(ctx.Request.Context(), query.Page, query.Limit, searchQuery)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	// Use Paginated helper for consistent response format
	responseRes.Paginated(ctx, http.StatusOK, "Members retrieved successfully", result.Data, &result.Pagination)
}

// FindRecentRegistrations implements [IMemberController].
func (m *MemberController) FindRecentRegistrations(ctx *gin.Context) {
	var query struct {
		Limit    int        `form:"limit,default=10"`
		FromDate *time.Time `form:"fromDate" binding:"omitempty"`
		ToDate   *time.Time `form:"toDate" binding:"omitempty"`
	}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	if query.Limit < 1 {
		query.Limit = 10
	}
	if query.Limit > 100 {
		query.Limit = 100
	}

	filters := &dto.RegistrationFilters{
		FromDate: query.FromDate,
		ToDate:   query.ToDate,
	}

	users, err := m.userService.FindRecentRegistrations(ctx.Request.Context(), query.Limit, filters)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Recent registrations retrieved successfully", users)
}

// SearchMembersWithPagination implements [IMemberController].
// @Summary Search Members with Pagination
// @Description Search members by name, email, or username with pagination
// @Tags Members
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10, max: 100)"
// @Param search query string true "Search query"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /members/search [get]
func (m *MemberController) SearchMembersWithPagination(ctx *gin.Context) {
	var query struct {
		Page   int    `form:"page,default=1" binding:"required"`
		Limit  int    `form:"limit,default=10"`
		Search string `form:"search" binding:"required"`
	}
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

	// Pass search query as pointer
	searchQuery := &query.Search

	result, err := m.memberService.GetMembersWithPagination(ctx.Request.Context(), query.Page, query.Limit, searchQuery)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	// Use Paginated helper for consistent response format
	responseRes.Paginated(ctx, http.StatusOK, "Members retrieved successfully", result.Data, &result.Pagination)
}

// @Summary Get Member Certificates
// @Description Get all certificates for a member
// @Tags Members
// @Produce json
// @Param userId path string true "User ID (UUID)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /members/{userId}/certificates [get]
func (m *MemberController) GetMemberCertificates(ctx *gin.Context) {
	userID, err := getUserIDFromPath(ctx, "userId")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	certificates, err := m.certService.GetMemberCertificates(ctx.Request.Context(), userID)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Member certificates retrieved successfully", certificates)
}

// @Summary Download Member Certificate
// @Description Download certificate PDF for a member
// @Tags Members
// @Produce octet-stream
// @Param userId path string true "User ID (UUID)"
// @Param certId path string true "Certificate ID (UUID)"
// @Success 200 {file} binary
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /members/{userId}/certificates/{certId}/download [get]
func (m *MemberController) DownloadMemberCertificate(ctx *gin.Context) {
	certID, err := uuid.Parse(ctx.Param("certId"))
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	pdfData, filename, err := m.certService.DownloadCertificatePDF(ctx.Request.Context(), certID)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	// Set headers for file download
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Header("Content-Type", "application/pdf")
	ctx.Data(http.StatusOK, "application/pdf", pdfData)
}

// CreateMember creates a new member (student) manually (admin only)
func (m *MemberController) CreateMember(ctx *gin.Context) {
	var input struct {
		Email       string `json:"email" binding:"required,email"`
		FirstName   string `json:"firstName" binding:"required"`
		LastName    string `json:"lastName"`
		PhoneNumber string `json:"phoneNumber"`
		Password    string `json:"password" binding:"required"`
		DateOfBirth string `json:"dateOfBirth"`
		Address     string `json:"address"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	// Find the member role
	role, err := m.roleService.GetRoleByName(ctx.Request.Context(), "member")
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	dob := input.DateOfBirth
	if dob == "" {
		dob = "2000-01-01" // Default date of birth
	}

	registerReq := dto.RegisterRequest{
		FirstName:       input.FirstName,
		LastName:        input.LastName,
		Username:        input.Email, // Use email as username
		Password:        input.Password,
		ConfirmPassword: input.Password,
		Email:           input.Email,
		PhoneNumber:     input.PhoneNumber,
		DateOfBirth:     dob,
		RoleID:          role.ID,
	}

	registerResp, err := m.authService.Register(ctx.Request.Context(), &registerReq)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	// Update address if provided
	if input.Address != "" {
		userModel, err := m.userService.GetUserByID(ctx.Request.Context(), registerResp.User.UserID)
		if err == nil {
			userModel.Address = input.Address
			_ = m.userService.UpdateUser(ctx.Request.Context(), userModel)
		}
	}

	// Fetch complete user with profiles to return to the frontend
	createdUser, err := m.userService.GetUserByIDWithProfiles(ctx.Request.Context(), registerResp.User.UserID)
	if err != nil {
		// Fallback to returning basic user info if load fails
		responseRes.Success(ctx, http.StatusCreated, "Member created successfully", registerResp.User)
		return
	}

	responseRes.Success(ctx, http.StatusCreated, "Member created successfully", createdUser)
}
