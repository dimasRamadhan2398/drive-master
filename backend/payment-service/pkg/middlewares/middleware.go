package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type IAuthMiddleware interface {
	Authenticate() gin.HandlerFunc
	AuthenticateWithoutToken() gin.HandlerFunc
	AuthorizeAdmin() gin.HandlerFunc
}

type AuthMiddleware struct {
	jwtSecret string
}

func NewAuthMiddleware(jwtSecret string) IAuthMiddleware {
	return &AuthMiddleware{jwtSecret: jwtSecret}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Extract Bearer token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			// No "Bearer " prefix found, try using the whole string
			tokenString = authHeader
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Bearer token is required",
			})
			c.Abort()
			return
		}

		// Parse and validate JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Extract claims and set user context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if user, ok := claims["User"].(map[string]interface{}); ok {
				if userID, ok := user["userId"].(string); ok {
					c.Set("user_id", userID)
				}
				if username, ok := user["username"].(string); ok {
					c.Set("username", username)
				}
				if roleID, ok := user["roleId"].(float64); ok {
					c.Set("role_id", uint(roleID))
				}
			}
		}

		c.Next()
	}
}

func (m *AuthMiddleware) AuthenticateWithoutToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only validate API key, skip Bearer token
		c.Next()
	}
}

func (m *AuthMiddleware) AuthorizeAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// First authenticate
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			tokenString = authHeader
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Check for admin role
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if user, ok := claims["User"].(map[string]interface{}); ok {
				if roleID, ok := user["roleId"].(float64); ok && uint(roleID) == 4 {
					// Admin role
					c.Next()
					return
				}
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Admin access required",
		})
		c.Abort()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID, X-Service-Name, X-Api-Key, X-Request-At")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func RateLimiter(maxRequests float64, expirationTTL time.Duration) gin.HandlerFunc {
	lmt := tollbooth.NewLimiter(maxRequests, &limiter.ExpirableOptions{
		DefaultExpirationTTL: expirationTTL,
	})

	return func(c *gin.Context) {
		httpError := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if httpError != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": httpError.Message,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}