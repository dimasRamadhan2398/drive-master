package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware handles Cross-Origin Resource Sharing
// pkg/middlewares/cors.go
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")

        // Allowlisted origins for direct access
        allowedOrigins := map[string]bool{
            "http://localhost:3000": true,
            "https://yourdomain.com": true,
        }

        // Check if request is forwarded by Kong (has Via header) or is from allowed origin
        via := c.Request.Header.Get("Via")
        if via != "" || allowedOrigins[origin] {
            // Request is from Kong proxy or allowed origin - set the actual origin
            if origin != "" {
                c.Header("Access-Control-Allow-Origin", origin)
            } else {
                c.Header("Access-Control-Allow-Origin", "*")
            }
        } else if origin == "" {
            // Same-origin or non-browser request — allow
            c.Header("Access-Control-Allow-Origin", "*")
        }
        // If origin is set but not in allowlist and not via Kong, headers are intentionally omitted

        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Api-Key, X-Request-At, X-Service-Name")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Max-Age", "86400") // cache preflight for 24h

        // ✅ Critical: terminate OPTIONS here — don't pass to route handler
        if c.Request.Method == http.MethodOptions {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }

        c.Next()
    }
}