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
