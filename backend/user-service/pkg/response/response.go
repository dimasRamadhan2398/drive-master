package response

import (
	"net/http"
	"strconv"

	apperrors "user-service/pkg/errors"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Success sends a successful response
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, code string, message string, details string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// ErrorFromAppError sends an error response from an AppError
func ErrorFromAppError(c *gin.Context, appErr *apperrors.AppError) {
	c.JSON(appErr.StatusCode, Response{
		Success: false,
		Error: &ErrorDetail{
			Code:    appErr.Code,
			Message: appErr.Message,
		},
	})
}

// ErrorFromGeneric sends an error response from a generic error
func ErrorFromGeneric(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Error: &ErrorDetail{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		},
	})
}

// OK sends a 200 OK response
func OK(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusOK, message, data)
}

// Created sends a 201 Created response
func Created(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusCreated, message, data)
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, "BAD_REQUEST", message, "")
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, "UNAUTHORIZED", message, "")
}

// Forbidden sends a 403 Forbidden response
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, "FORBIDDEN", message, "")
}

// NotFound sends a 404 Not Found response
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, "NOT_FOUND", message, "")
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", message, "")
}

// ParseUint parses a uint from path parameter
func ParseUint(c *gin.Context, paramName string) (uint, error) {
	idStr := c.Param(paramName)
	if idStr == "" {
		return 0, nil
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
