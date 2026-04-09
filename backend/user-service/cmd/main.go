package main

import (
	"log"
	"os"
	"user-service/clients"
	"user-service/handlers"
	"user-service/models"
	"user-service/repositories"
	"user-service/services"

	"github.com/gin-gonic/gin"
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

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("failed to migrate user table: %v", err)
	}

	kafkaBroker := getenv("KAFKA_BROKER", "localhost:9092")
	kafkaTopic := getenv("KAFKA_TOPIC_USER_CREATED", "user.created")

	userRepo := repositories.NewUserRepository(db)
	producer := clients.NewKafkaProducer(kafkaBroker, kafkaTopic)
	userService := services.NewUserService(userRepo, producer)
	httpHandler := handlers.NewHTTPHandler(userService)

	router := gin.Default()
	httpHandler.RegisterRoutes(router)

	port := getenv("PORT", "8081")
	log.Printf("user-service listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start user-service: %v", err)
	}
}
