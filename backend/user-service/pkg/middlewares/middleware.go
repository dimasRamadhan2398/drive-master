package middlewares

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"
	"user-service/pkg/base"
	"user-service/pkg/config"
	"user-service/pkg/constants"
	apperrors "user-service/pkg/errors"
	"user-service/pkg/response"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

const (
	DefaultMaxRequests        = 100               // requests per window
	DefaultExpirationTTLSeconds = 60              // 1 minute window
)

type Claims struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// AuthMiddleware validates JWT tokens

type IAuthMiddleware interface {
	Authenticate() gin.HandlerFunc
    AuthenticateWithoutToken() gin.HandlerFunc
}

type AuthMiddleware struct {
	secret string
}

func NewAuthMiddleware(secret string) IAuthMiddleware {
	return &AuthMiddleware{
		secret: secret,
	}
}


// Authenticate validates both the API key signature and Bearer token
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
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
		// Parse and validate token
		newToken, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.secret), nil
		})
		// jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 	return []byte(m.secret), nil
		// })

		if err != nil || !newToken.Valid {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}
		c.Set(constants.Token, tokenString)
		c.Next()
	}
}

// AuthenticateWithoutToken only validates the API key signature (no Bearer token required)
func (m *AuthMiddleware) AuthenticateWithoutToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := validateAPIKey(c)
		if err != nil {
			responseUnauthorized(c, err.Error())
			return
		}
		c.Next()
	}
}

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

// RateLimiter creates a rate limiting middleware using tollbooth
// maxRequests: maximum number of requests allowed per expiration window
// expirationTTL: how long the rate limit window lasts (e.g., 60 seconds)
func RateLimiter(maxRequests float64, expirationTTL time.Duration) gin.HandlerFunc {
	lmt := tollbooth.NewLimiter(maxRequests, &limiter.ExpirableOptions{
		DefaultExpirationTTL: expirationTTL,
	})
	lmt.SetIPLookups([]string{"X-Forwarded-For", "X-Real-IP", "RemoteAddr"})

	// Paths to exclude from rate limiting
	excludedPaths := []string{"/swagger/", "/health"}

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Skip rate limiting for excluded paths
		shouldSkip := false
		for _, excludedPath := range excludedPaths {
			if strings.HasPrefix(path, excludedPath) {
				shouldSkip = true
				break
			}
		}

		if shouldSkip {
			c.Next()
			return
		}

		err := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusTooManyRequests, base.Response{
				Status:  apperrors.Error,
				Message: apperrors.ErrTooManyRequests.Error(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func extractBearerToken(token string) string {
	arrayToken := strings.Split(token, " ")
	if len(arrayToken) == 2 {
		return arrayToken[1]
	}
	// If no "Bearer " prefix, return the token as-is
	return token
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
	cfg := config.Get()
	if cfg.App.AppEnv == "local" || cfg.App.AppEnv == "development" {
		return nil
	}
	apiKey := c.GetHeader(constants.XApiKey)
	requestAt := c.GetHeader(constants.XRequestAt)
	serviceName := c.GetHeader(constants.XServiceName)

	if apiKey == "" || requestAt == "" || serviceName == "" {
		return apperrors.ErrUnauthorized
	}

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