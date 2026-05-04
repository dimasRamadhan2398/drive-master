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
	ResetPassword(ctx *gin.Context)
	ConfirmResetPassword(ctx *gin.Context)
	VerifyOTP(ctx *gin.Context)
	ResendOTP(ctx *gin.Context)
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

// Login authenticates user with credentials and returns JWT token
//
//	@Summary		Login
//	@Description	Authenticate user with email/username and password
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.LoginInput			true	"Login credentials"
//	@Success		200		{object}	response.Response		"Login successful"
//	@Failure		400		{object}	response.Response		"Bad request"
//	@Failure		401		{object}	response.Response		"Invalid credentials"
//	@Router			/auth/login [post]
func (a *AuthController) Login(ctx *gin.Context) {
	var input dto.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	if err := validator.New().Struct(input); err != nil {
		responseRes.Error(ctx, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity), err.Error(), "")
		return
	}

	loginResp, err := a.authService.Login(ctx.Request.Context(), &input)
	if err != nil {
		responseRes.Error(ctx, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), "Invalid username or password", "")
		return
	}

	role, err := a.roleService.GetRole(ctx, loginResp.User.RoleID)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}
	loginResp.User.Role = dto.RoleResponse{ID: role.ID, Name: role.Name}

	responseRes.Success(ctx, http.StatusOK, "Login successful", loginResp)
}

// Register creates new user account and sends OTP for email verification
//
//	@Summary		Register
//	@Description	Create a new user account and send OTP for email verification
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RegisterRequest	true	"User registration data"
//	@Success		201		{object}	response.Response		"Registration successful"
//	@Failure		400		{object}	response.Response		"Bad request"
//	@Failure		409		{object}	response.Response		"Username or email already exists"
//	@Router			/auth/register [post]
func (a *AuthController) Register(ctx *gin.Context) {
	var input dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	if err := validator.New().Struct(input); err != nil {
		responseRes.Error(ctx, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity), err.Error(), "")
		return
	}

	registerResp, err := a.authService.Register(ctx.Request.Context(), &input)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	go a.emailService.SendWelcomeEmail(ctx.Request.Context(), input.Email, input.Username)

	responseRes.Success(ctx, http.StatusCreated, "User registered successfully. Please check your email for OTP verification", registerResp)
}

// ResetPassword sends password reset link to email
//
//	@Summary		Forgot Password
//	@Description	Request password reset link sent to email
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.ForgotPasswordInput	true	"Email address"
//	@Success		200		{object}	response.Response		"Reset link sent if email exists"
//	@Failure		400		{object}	response.Response		"Bad request"
//	@Router			/auth/forgot-password [post]
func (a *AuthController) ResetPassword(ctx *gin.Context) {
	var input dto.ForgotPasswordInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	if err := validator.New().Struct(input); err != nil {
		responseRes.Error(ctx, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity), err.Error(), "")
		return
	}

	user, err := a.userService.GetUserByEmail(ctx.Request.Context(), input.EmailAddress)
	if err != nil {
		responseRes.Success(ctx, http.StatusOK, "If the email exists, a reset link has been sent", nil)
		return
	}

	go a.emailService.SendPasswordResetEmail(ctx.Request.Context(), user.Email, user.Username)

	responseRes.Success(ctx, http.StatusOK, "If the email exists, a reset link has been sent", nil)
}

// ConfirmResetPassword resets password using valid token
//
//	@Summary		Confirm Reset Password
//	@Description	Reset password using token from email
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.ResetPasswordInput	true	"Reset password data"
//	@Success		200		{object}	response.Response		"Password reset successful"
//	@Failure		400		{object}	response.Response		"Bad request"
//	@Router			/auth/confirm-reset-password [post]
func (a *AuthController) ConfirmResetPassword(ctx *gin.Context) {
	var input dto.ResetPasswordInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	if err := validator.New().Struct(input); err != nil {
		responseRes.Error(ctx, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity), err.Error(), "")
		return
	}

	// TODO: Implement token validation and password reset
	responseRes.Success(ctx, http.StatusOK, "Password has been reset successfully", nil)
}

// VerifyOTP verifies the OTP code and marks user email as verified
//
//	@Summary		Verify OTP
//	@Description	Verify OTP sent to email and mark user as verified
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.VerifyOTPInput	true	"OTP verification data"
//	@Success		200		{object}	response.Response		"Email verified successfully"
//	@Failure		400		{object}	response.Response		"Bad request"
//	@Failure		401		{object}	response.Response		"Invalid or expired OTP"
//	@Router			/auth/verify-otp [post]
func (a *AuthController) VerifyOTP(ctx *gin.Context) {
	var input dto.VerifyOTPInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	if err := validator.New().Struct(input); err != nil {
		responseRes.Error(ctx, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity), err.Error(), "")
		return
	}

	if err := a.authService.VerifyOTP(ctx.Request.Context(), input.Email, input.OTP); err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Email verified successfully", nil)
}

// ResendOTP resends verification OTP to email
//
//	@Summary		Resend OTP
//	@Description	Resend OTP to email address for verification
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.ResendOTPInput	true	"Email address"
//	@Success		200		{object}	response.Response		"OTP sent successfully"
//	@Failure		400		{object}	response.Response		"Bad request"
//	@Router			/auth/resend-otp [post]
func (a *AuthController) ResendOTP(ctx *gin.Context) {
	var input dto.ResendOTPInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(ctx, apperrors.ErrBadRequest)
		return
	}

	if err := validator.New().Struct(input); err != nil {
		responseRes.Error(ctx, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity), err.Error(), "")
		return
	}

	if err := a.authService.ResendOTP(ctx.Request.Context(), input.Email); err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "OTP has been sent to your email", nil)
}
