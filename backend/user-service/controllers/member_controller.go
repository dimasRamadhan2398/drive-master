package controllers

import (
	"net/http"
	apperrors "user-service/pkg/errors"
	responseRes "user-service/pkg/response"
	"user-service/services"

	"github.com/gin-gonic/gin"
)

type MemberController struct {
	userService   services.IUserService
	authService   services.IAuthService
	memberService services.IMemberService
	roleService   services.IRoleService
	emailService  services.IMailtrapEmailService
	mediaService  services.IMediaService
}

func NewMemberController(
	userService services.IUserService,
	authService services.IAuthService,
	memberService services.IMemberService,
	roleService services.IRoleService,
	emailService services.IMailtrapEmailService,
	mediaService services.IMediaService,
) IMemberController {
	return &MemberController{
		userService:   userService,
		authService:   authService,
		memberService: memberService,
		roleService:   roleService,
		emailService:  emailService,
		mediaService:  mediaService,
	}
}

type IMemberController interface {
	GetMemberProfile(ctx *gin.Context)
	UpdateMemberProfile(ctx *gin.Context)
	GetMemberLists(ctx *gin.Context)
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

	profile, err := m.memberService.GetMemberProfile(ctx, userID)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	// Bind request body
	var input struct {
		// Add member profile update fields here
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
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
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /members [get]
func (m *MemberController) GetMemberLists(ctx *gin.Context) {
	var query struct {
		Page  int `form:"page,default=1"`
		Limit int `form:"limit,default=10"`
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

	result, err := m.memberService.GetMembersWithPagination(ctx.Request.Context(), query.Page, query.Limit)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	// Use Paginated helper for consistent response format
	responseRes.Paginated(ctx, http.StatusOK, "Members retrieved successfully", result.Data, &result.Pagination)
}
