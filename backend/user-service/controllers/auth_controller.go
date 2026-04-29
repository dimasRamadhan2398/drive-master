package controllers

import (
	"net/http"
	"user-service/models/dto"
	apperrors "user-service/pkg/errors"
	responseRes "user-service/pkg/response"
	"user-service/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthController struct {
	authService  services.IAuthService
	userService  services.IUserService
	emailService services.IMailtrapEmailService
	roleService  services.IRoleService
}


type IAuthController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
	ConfirmForgotPassword(ctx *gin.Context)
	ResetPassword(ctx *gin.Context)
}

func NewAuthController(
	userService services.IUserService,
	authService services.IAuthService,
	emailService services.IMailtrapEmailService,
	roleService services.IRoleService,
) IAuthController {
	return &AuthController{
		userService:  userService,
		authService:  authService,
		emailService: emailService,
		roleService:  roleService,
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
func (a *AuthController) Login(ctx *gin.Context) {
	input := &dto.LoginInput{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		responseRes.Error(ctx, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity), err.Error(), "")
	}

	user, err := a.authService.ValidateCredentials(input.Email, input.Password)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	// Get user role
	role, err := a.roleService.GetRole(user.RoleID)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	// Build login response
	loginResp := dto.LoginResponse{
		User: dto.GetUserResponse{
			UserID:      user.ID,
			Email:       user.Email,
			Username:    user.Username,
			PhoneNumber: user.PhoneNumber,
			Image:       user.Image,
			DateOfBirth: user.DateOfBirth,
			Address:     user.Address,
			RoleID:      user.RoleID,
			Role: dto.RoleResponse{
				ID:   role.ID,
				Name: role.Name,
			},
		},
	}

	responseRes.Success(ctx, http.StatusOK, "Login successful", loginResp)
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
func (a *AuthController) Register(ctx *gin.Context) {
	var input dto.CreateUserRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		responseRes.Error(ctx, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity), err.Error(), "")
	}

	user, err := a.userService.CreateUser(input)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	// Send welcome email asynchronously
	go a.emailService.SendWelcomeEmail(ctx.Request.Context(), user.Email, user.Username)

	responseRes.Success(ctx, http.StatusCreated, "User registered successfully", user)
}

// @Summary Confirm Forgot Password
// @Description Confirm Forgot password using a token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ForgotPasswordInput true "Confirm Forgot password data"
// @Success 200 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Router /auth/confirm-forgot-password [post]
func (a *AuthController) ConfirmForgotPassword(ctx *gin.Context) {
	var input dto.ResetPasswordInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		responseRes.Error(ctx, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity), err.Error(), "")
	}

	// TODO: Validate reset token and reset password
	responseRes.Success(ctx, http.StatusOK, "Password has been reset successfully", nil)
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
func (a *AuthController) ResetPassword(ctx *gin.Context) {
	var input dto.ForgotPasswordInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		responseRes.Error(ctx, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity), err.Error(), "")
	}

	user, err := a.userService.GetUserByEmail(input.EmailAddress)
	if err != nil {
		// Don't reveal if user exists or not
		responseRes.Success(ctx, http.StatusOK, "If the email exists, a reset link has been sent", nil)
		return
	}

	// TODO: Generate reset token and store it
	resetToken := "placeholder-token"

	// Send password reset email asynchronously
	go a.emailService.SendPasswordResetEmail(ctx.Request.Context(), user.Email, resetToken)

	responseRes.Success(ctx, http.StatusOK, "If the email exists, a reset link has been sent", nil)
}