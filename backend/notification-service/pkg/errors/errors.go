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

// Common errors for notification service
var (
	ErrBadRequest           = New("BAD_REQUEST", "Bad request", http.StatusBadRequest)
	ErrUnauthorized         = New("UNAUTHORIZED", "Unauthorized", http.StatusUnauthorized)
	ErrForbidden            = New("FORBIDDEN", "Forbidden", http.StatusForbidden)
	ErrNotFound             = New("NOT_FOUND", "Resource not found", http.StatusNotFound)
	ErrConflict             = New("CONFLICT", "Resource conflict", http.StatusConflict)
	ErrInternalServer       = New("INTERNAL_ERROR", "Internal server error", http.StatusInternalServerError)
	ErrDatabase             = New("DATABASE_ERROR", "Database error", http.StatusInternalServerError)
	ErrValidation           = New("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest)
	ErrDuplicateEntry       = New("DUPLICATE_ENTRY", "Duplicate entry", http.StatusConflict)
	ErrEmailSendFailed      = New("EMAIL_SEND_FAILED", "Failed to send email", http.StatusInternalServerError)
	ErrWhatsAppSendFailed   = New("WHATSAPP_SEND_FAILED", "Failed to send WhatsApp message", http.StatusInternalServerError)
	ErrTemplateNotFound     = New("TEMPLATE_NOT_FOUND", "Notification template not found", http.StatusNotFound)
	ErrInvalidPhoneNumber   = New("INVALID_PHONE_NUMBER", "Invalid phone number", http.StatusBadRequest)
	ErrSubscriptionExists   = New("SUBSCRIPTION_EXISTS", "Email already subscribed", http.StatusConflict)
)