package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application error
type AppError struct {
	Code       string
	Message    string
	StatusCode int
	Err        error
	Meta       map[string]interface{}
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// Is checks if an error is of type AppError
func Is(err error, target *AppError) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == target.Code
	}
	return false
}

// New creates a new AppError
func New(code string, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// Wrap wraps an existing error
func Wrap(code string, message string, statusCode int, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// Common errors
var (
	ErrBadRequest         = New("BAD_REQUEST", "Bad request", http.StatusBadRequest)
	ErrUnauthorized       = New("UNAUTHORIZED", "Unauthorized", http.StatusUnauthorized)
	ErrForbidden          = New("FORBIDDEN", "Forbidden", http.StatusForbidden)
	ErrNotFound           = New("NOT_FOUND", "Resource not found", http.StatusNotFound)
	ErrConflict           = New("CONFLICT", "Resource conflict", http.StatusConflict)
	ErrInternalServer     = New("INTERNAL_ERROR", "Internal server error", http.StatusInternalServerError)
	ErrDatabase           = New("DATABASE_ERROR", "Database error", http.StatusInternalServerError)
	ErrValidation         = New("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest)
	ErrInvalidCredentials = New("INVALID_CREDENTIALS", "Invalid credentials", http.StatusUnauthorized)
	ErrInvalidToken       = New("INVALID_TOKEN", "Invalid token", http.StatusUnauthorized)
	ErrTokenExpired       = New("TOKEN_EXPIRED", "Token expired", http.StatusUnauthorized)
	ErrTooManyRequests    = New("TOO_MANY_REQUESTS", "Too many requests", http.StatusTooManyRequests)
	ErrUnknownDevice      = New("UNKNOWN_DEVICE", "Unknown device", http.StatusUnauthorized)
	ErrUsernameExist 	  = New("USERNAME_EXIST", "Username already exist", http.StatusConflict)
	ErrEmailExist      	  = New("EMAIL_EXIST", "Email already exist", http.StatusConflict)
	ErrPasswordDoesNotMatch = New("PASSWORD_DOES_NOT_MATCH", "Password does not match", http.StatusUnauthorized)
	ErrDuplicateEntry     = New("DUPLICATE_ENTRY", "Duplicate entry", http.StatusConflict)
	ErrGenerateOTP        = New("FAILED_TO_GENERATE_OTP", "Error generating OTP", http.StatusInternalServerError)
	ErrAccountLocked      = New("ACCOUNT_LOCKED", "Account locked", http.StatusUnauthorized)
	ErrIDNotFound 		  = New("ID_NOT_FOUND", "Resource not found. Please write the correct ID", http.StatusNotFound)
	ErrUserNotFound       = New("USER_NOT_FOUND", "User not found", http.StatusNotFound)
	ErrAlreadyVerified    = New("ALREADY_VERIFIED", "User is already verified", http.StatusConflict)
	ErrOTPExpired        = New("OTP_EXPIRED", "OTP has expired", http.StatusUnauthorized)
	ErrInvalidOTP        = New("INVALID_OTP", "Invalid OTP code", http.StatusUnauthorized)
	ErrEmailNotVerified  = New("EMAIL_NOT_VERIFIED", "Email not verified", http.StatusUnauthorized)
)