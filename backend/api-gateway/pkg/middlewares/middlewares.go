package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// must match exactly what user-service puts in the JWT
type JWTClaims struct {
    User struct {
        UserID      string `json:"userId"`
        Email       string `json:"email"`
        Username    string `json:"username"`
        PhoneNumber string `json:"phoneNumber"`
        RoleID      int    `json:"roleId"`
        Role        struct {
            ID   int    `json:"id"`
            Name string `json:"name"`
        } `json:"role"`
    } `json:"User"`
    jwt.RegisteredClaims
}

type AuthMiddleware struct {
    jwtSecret string
}

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
    return &AuthMiddleware{jwtSecret: jwtSecret}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "message": "authorization header is required",
            })
            return
        }

        // check Bearer prefix
        if !strings.HasPrefix(authHeader, "Bearer ") {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "message": "bearer token is required",
            })
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")

        // parse and validate JWT
        claims := &JWTClaims{}
        token, err := jwt.ParseWithClaims(
            tokenString,
            claims,
            func(t *jwt.Token) (interface{}, error) {
                // validate signing method — prevent algorithm confusion attack
                if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
                    return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
                }
                return []byte(m.jwtSecret), nil
            },
        )

        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "message": "invalid or expired token",
            })
            return
        }

        // inject validated claims into context
        // downstream services read these from X-User-* headers
        c.Set("userID",   claims.User.UserID)
        c.Set("userRole", claims.User.Role.Name)
        c.Set("userEmail", claims.User.Email)
        c.Next()
    }
}

// RequireRole checks role after Authenticate runs
func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole := c.GetString("userRole")
        for _, role := range roles {
            if userRole == role {
                c.Next()
                return
            }
        }
        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
            "success": false,
            "message": "insufficient permissions",
        })
    }
}

// RateLimiter implements a token bucket rate limiter using IP address
type RateLimiter struct {
    max     int           // max requests per window
    window  time.Duration // time window
    clients map[string][]time.Time
    mu      sync.RWMutex
}

func NewRateLimiter(max int, windowSeconds int) *RateLimiter {
    return &RateLimiter{
        max:     max,
        window:  time.Duration(windowSeconds) * time.Second,
        clients: make(map[string][]time.Time),
    }
}

func (r *RateLimiter) getClientKey(c *gin.Context) string {
    // Use X-Forwarded-For if behind a proxy, otherwise use RemoteAddr
    if forwarded := c.GetHeader("X-Forwarded-For"); forwarded != "" {
        return strings.Split(forwarded, ",")[0]
    }
    return c.ClientIP()
}

func (r *RateLimiter) Allow() gin.HandlerFunc {
    return func(c *gin.Context) {
        key := r.getClientKey(c)

        r.mu.Lock()
        now := time.Now()
        windowStart := now.Add(-r.window)

        // Filter out old requests outside the window
        requests := r.clients[key]
        validRequests := make([]time.Time, 0)
        for _, t := range requests {
            if t.After(windowStart) {
                validRequests = append(validRequests, t)
            }
        }

        // Check if limit exceeded
        if len(validRequests) >= r.max {
            r.clients[key] = validRequests
            r.mu.Unlock()
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
                "success": false,
                "message": "rate limit exceeded, please try again later",
            })
            return
        }

        // Add current request
        r.clients[key] = append(validRequests, now)
        r.mu.Unlock()

        c.Next()
    }
}

// RequestIDMiddleware generates and sets a unique request ID for each request
func RequestIDMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        c.Set("requestID", requestID)
        c.Header("X-Request-ID", requestID)
        c.Next()
    }
}

func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")

        // only allow configured origins
        for _, allowed := range allowedOrigins {
            if origin == allowed {
                c.Header("Access-Control-Allow-Origin", origin)
                break
            }
        }

        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID")
        c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }
        c.Next()
    }
}