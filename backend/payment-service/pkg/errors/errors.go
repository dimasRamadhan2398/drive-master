package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Details)
	}
	return e.Message
}

func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func (e *AppError) WithDetails(details string) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: e.Message,
		Details: details,
	}
}

var (
	ErrNotFound          = &AppError{Code: http.StatusNotFound, Message: "Resource not found"}
	ErrUnauthorized      = &AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"}
	ErrForbidden         = &AppError{Code: http.StatusForbidden, Message: "Forbidden"}
	ErrBadRequest        = &AppError{Code: http.StatusBadRequest, Message: "Bad request"}
	ErrInternalServer    = &AppError{Code: http.StatusInternalServerError, Message: "Internal server error"}
	ErrConflict          = &AppError{Code: http.StatusConflict, Message: "Conflict"}
	ErrPaymentFailed     = &AppError{Code: http.StatusPaymentRequired, Message: "Payment failed"}
	ErrInvalidAmount     = &AppError{Code: http.StatusBadRequest, Message: "Invalid amount"}
	ErrPaymentExpired    = &AppError{Code: http.StatusGone, Message: "Payment expired"}
	ErrPaymentCancelled  = &AppError{Code: http.StatusBadRequest, Message: "Payment cancelled"}
)

func IsNotFound(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == http.StatusNotFound
	}
	return false
}

func IsUnauthorized(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == http.StatusUnauthorized
	}
	return false
}

func IsPaymentError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == http.StatusPaymentRequired
	}
	return false
}