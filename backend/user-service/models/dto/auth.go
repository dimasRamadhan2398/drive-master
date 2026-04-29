package dto

type LoginInput struct {
	// accepts either email or username
	Email string `json:"email" binding:"required"`
	Password   string `json:"password"   binding:"required"`
}

// LoginResponse is returned on successful login
type LoginResponse struct {
	AccessToken  string          `json:"accessToken"`
	RefreshToken string          `json:"refreshToken"`
	ExpiresIn    int64           `json:"expiresIn"` // unix timestamp
	User         GetUserResponse `json:"user"`
}

// RefreshTokenInput is used for POST /auth/refresh
type RefreshTokenInput struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// ForgotPasswordInput is used for POST /auth/forgot-password
type ForgotPasswordInput struct {
	EmailAddress string `json:"emailAddress" binding:"required,email"`
}

// ResetPasswordInput is used for POST /auth/reset-password
type ResetPasswordInput struct {
	Token       string `json:"token"       binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=8"`
}

// ChangePasswordInput is used for PATCH /users/:id/password
// Requires the current password to verify ownership
type ChangePasswordInput struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword"     binding:"required,min=8"`
}
