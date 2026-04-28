package main

import (
	"api-gateway/services"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	router := gin.Default()

	routerService := services.NewRouterService(
		getenv("USER_SERVICE_URL", "http://localhost:8081"),
		getenv("CORE_SERVICE_URL", "http://localhost:8082"),
	)
	routerService.RegisterRoutes(router)

	port := getenv("PORT", "8080")
	log.Printf("api-gateway listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start api-gateway: %v", err)
	}
}
