package dto

import (
	"time"
)

// DateFormat is the custom date format used for DD/MM/YYYY
const DateFormat = "02/01/2006"

// Date is a custom type for parsing DD/MM/YYYY date format
// swagger:strfmt date
type Date struct {
	time.Time
}

// UnmarshalJSON parses the date from DD/MM/YYYY format
func (d *Date) UnmarshalJSON(data []byte) error {
	// Remove quotes
	str := string(data)
	if str == "null" || str == `""` {
		return nil
	}

	parsed, err := time.Parse(DateFormat, str)
	if err != nil {
		return err
	}
	d.Time = parsed
	return nil
}

// MarshalJSON returns the date in DD/MM/YYYY format
func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.Time.Format(DateFormat) + `"`), nil
}

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

type RegisterRequest struct {
	Name            string    `json:"name" binding:"required"`
	Username        string    `json:"username" binding:"required"`
	Password        string    `json:"password" binding:"required"`
	ConfirmPassword string    `json:"confirmPassword" binding:"required"`
	Email           string    `json:"email" binding:"required,email"`
	PhoneNumber     string    `json:"phoneNumber" binding:"required"`
	DateOfBirth     string      `json:"dateOfBirth" binding:"required" swaggertype:"string" format:"date" example:"1999-08-25"`
	RoleID          uint      `json:"roleId"`
}

type RegisterResponse struct {
	User CreateUserResponse `json:"user"`
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

// VerifyOTPInput is used for POST /auth/verify-otp
type VerifyOTPInput struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp"   binding:"required,len=6"`
}

// ResendOTPInput is used for POST /auth/resend-otp
type ResendOTPInput struct {
	Email string `json:"email" binding:"required,email"`
}
