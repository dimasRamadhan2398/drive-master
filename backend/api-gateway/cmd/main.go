package main

import (
	"fmt"
	"log"
	"net/http"

	"api-gateway/pkg/config"
	middleware "api-gateway/pkg/middlewares"
	"api-gateway/routes"

	"github.com/gin-gonic/gin"
)

func main() {
    cfg, err := config.Load("config.yaml")
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    if cfg.Server.Mode == "release" {
        gin.SetMode(gin.ReleaseMode)
    }

    router := gin.Default()

    // global middleware — order matters
    router.Use(middleware.RequestIDMiddleware())
    router.Use(middleware.CORSMiddleware(cfg.CORS.AllowedOrigins))
    router.Use(middleware.NewRateLimiter(cfg.RateLimiter.Max, cfg.RateLimiter.Time).Allow())

    // health check — no auth
    router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status":  "healthy",
            "service": "api-gateway",
        })
    })

    // register all proxy routes
    routes.Register(router, cfg)

    addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
    log.Printf("API Gateway listening on %s", addr)
    if err := router.Run(addr); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}