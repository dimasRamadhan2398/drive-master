package middlewares

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"user-service/pkg/base"
	"user-service/pkg/config"
	"user-service/pkg/constants"
	apperrors "user-service/pkg/errors"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func HandlePanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("Recovered from panic: %v", r)
				c.JSON(http.StatusInternalServerError, base.Response{
					Status:  apperrors.Error,
					Message: apperrors.ErrInternalServer.Error(),
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

func RateLimiter(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusTooManyRequests, base.Response{
				Status:  apperrors.Error,
				Message: apperrors.ErrTooManyRequests.Error(),
			})
			c.Abort()
		}
		c.Next()
	}
}

func extractBearerToken(token string) string {
	arrayToken := strings.Split(token, " ")
	if len(arrayToken) == 2 {
		return arrayToken[1]
	}
	return ""
}

func responseUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, base.Response{
		Status:  apperrors.Error,
		Message: message,
	})
	c.Abort()
}

func responseError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, base.Response{
		Status:  apperrors.Error,
		Message: message,
	})
	c.Abort()
}

func validateAPIKey(c *gin.Context) error {
	apiKey := c.GetHeader(constants.XApiKey)
	requestAt := c.GetHeader(constants.XRequestAt)
	serviceName := c.GetHeader(constants.XServiceName)

	if apiKey == "" || requestAt == "" || serviceName == "" {
		return apperrors.ErrUnauthorized
	}

	cfg := config.Get()
	validateKey := fmt.Sprintf("%s:%s:%s", serviceName, cfg.App.SignatureKey, requestAt)
	hash := sha256.New()
	hash.Write([]byte(validateKey))
	resultHash := hex.EncodeToString(hash.Sum(nil))

	if apiKey != resultHash {
		return apperrors.ErrUnauthorized
	}
	return nil
}

func contains(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

// Authenticate validates both the API key signature and Bearer token
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(constants.Authorization)
		if token == "" {
			responseUnauthorized(c, apperrors.ErrUnauthorized.Error())
			return
		}

		err := validateAPIKey(c)
		if err != nil {
			responseUnauthorized(c, err.Error())
			return
		}

		tokenString := extractBearerToken(token)
		if tokenString == "" {
			responseUnauthorized(c, apperrors.ErrUnauthorized.Error())
			return
		}

		// Store token in context for downstream handlers
		c.Set(constants.Token, tokenString)
		c.Next()
	}
}

// AuthenticateWithoutToken only validates the API key signature (no Bearer token required)
func AuthenticateWithoutToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := validateAPIKey(c)
		if err != nil {
			responseUnauthorized(c, err.Error())
			return
		}
		c.Next()
	}
}