package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"api-gateway/pkg/config"
	middleware "api-gateway/pkg/middlewares"
	"api-gateway/routes"

	_ "api-gateway/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Get config path from env or use default path
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		// Default to pkg/config/config.yaml relative to working directory
		configPath = "pkg/config/config.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config from %s: %v", configPath, err)
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

	// ── SWAGGER docs — aggregates all services ──────────────────
	// gateway exposes each service's swagger at these paths
	router.GET("/swagger/user-service/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("http://user-service:8001/swagger/doc.json")))
	router.GET("/swagger/booking-service/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("http://booking-service:8003/swagger/doc.json")))
	router.GET("/swagger/core-service/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("http://core-service:8002/swagger/doc.json")))

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("API Gateway listening on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
