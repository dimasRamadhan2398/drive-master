package response

import (
	"errors"
	"net/http"
	"strconv"

	"user-service/models/dto"
	apperrors "user-service/pkg/errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Response struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Pagination *dto.PaginationMeta `json:"pagination,omitempty"`
	Error      *ErrorDetail `json:"error,omitempty"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Paginated sends a successful response with data and pagination at same level
func Paginated(c *gin.Context, statusCode int, message string, data interface{}, pagination *dto.PaginationMeta) {
	c.JSON(statusCode, Response{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

// Success sends a successful response (without pagination)
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

// ErrorFromValidationErrors sends an error response from validation errors (map[string][]string)
func ErrorFromValidationErrors(c *gin.Context, errs map[string][]string) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Error: &ErrorDetail{
			Code:    "VALIDATION_ERROR",
			Message: "Validation failed",
			Details: formatValidationErrors(errs),
		},
	})
}

// formatValidationErrors converts validation errors map to a human-readable string
func formatValidationErrors(errs map[string][]string) string {
	if len(errs) == 0 {
		return ""
	}
	var messages []string
	for field, fieldErrs := range errs {
		for _, msg := range fieldErrs {
			messages = append(messages, field+": "+msg)
		}
	}
	return joinStrings(messages, "; ")
}

// joinStrings joins strings with a separator
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
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
	if err == nil {
		return
	}

	// Check if it's an AppError directly
	if appErr, ok := err.(*apperrors.AppError); ok {
		c.JSON(appErr.StatusCode, Response{
			Success: false,
			Error: &ErrorDetail{
				Code:    appErr.Code,
				Message: appErr.Message,
			},
		})
		return
	}

	// Try to unwrap the error to check for AppError
	unwrappedErr := err
	for {
		unwrappedErr = errors.Unwrap(unwrappedErr)
		if unwrappedErr == nil {
			break
		}
		if appErr, ok := unwrappedErr.(*apperrors.AppError); ok {
			c.JSON(appErr.StatusCode, Response{
				Success: false,
				Error: &ErrorDetail{
					Code:    appErr.Code,
					Message: appErr.Message,
				},
			})
			return
		}
	}

	// Check for GORM errors
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error: &ErrorDetail{
				Code:    "NOT_FOUND",
				Message: "Resource not found",
			},
		})
		return
	}

	// Fallback for unknown errors
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