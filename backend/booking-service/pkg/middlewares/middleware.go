package middlewares

import (
	"booking-service/pkg/base"
	"booking-service/pkg/config"
	"booking-service/pkg/constants"
	apperrors "booking-service/pkg/errors"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	DefaultMaxRequests          = 100 // requests per window
	DefaultExpirationTTLSeconds = 60  // 1 minute window
)

type IAuthMiddleware interface {
	Authenticate() gin.HandlerFunc
	AuthenticateWithoutToken() gin.HandlerFunc
}

type AuthMiddleware struct {
	secret string
}

// Authenticate implements [IAuthMiddleware].
func (a *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader(constants.Authorization)
		if token == "" {
			ctx.AbortWithStatusJSON(401, gin.H{
				"message": "Unauthorized",
			})
			return
		}
		ctx.Next()
	}
}

// AuthenticateWithoutToken implements [IAuthMiddleware].
func (a *AuthMiddleware) AuthenticateWithoutToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := validateAPIKey(c)
		if err != nil {
			responseUnauthorized(c, err.Error())
			return
		}
		c.Next()
	}
}

func NewAuthMiddleware(secret string) IAuthMiddleware {
	return &AuthMiddleware{
		secret: secret,
	}
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

func responseUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, base.Response{
		Status:  apperrors.Error,
		Message: message,
	})
	c.Abort()
}