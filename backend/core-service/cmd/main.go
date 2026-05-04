package main

import (
	"log"
	"os"

	"core-service/models"
	"core-service/repositories"
	"core-service/services"
	"core-service/controllers"
	"core-service/routes"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	dsn := getenv("POSTGRES_DSN", "host=localhost user=app password=app dbname=app port=5432 sslmode=disable")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Migrate region models
	if err := db.AutoMigrate(&models.Province{}, &models.Regency{}, &models.District{}, &models.ProcessedEvent{}); err != nil {
		log.Fatalf("failed to migrate tables: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: getenv("REDIS_ADDR", "localhost:6379"),
	})

	// Initialize repositories
	repoRegistry := repositories.NewRepositoryRegistry(db)
	repoRegistry.SetCacheClient(redisClient)

	// Initialize services
	svcRegistry := services.NewServiceRegistry(repoRegistry)

	// Initialize controllers
	controllerRegistry := controllers.NewControllerRegistry(svcRegistry)

	// Setup HTTP router
	router := gin.Default()

	// Register routes
	apiGroup := router.Group("/api/v1")
	routeRegistry := routes.NewRouteRegistry(controllerRegistry, apiGroup)
	routeRegistry.Serve()

	port := getenv("PORT", "8082")
	log.Printf("core-service listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start core-service: %v", err)
	}
}