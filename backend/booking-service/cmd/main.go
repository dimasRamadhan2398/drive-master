package main

import (
	"context"
	"log"
	"os"

	"core-service/handlers"
	"core-service/models"
	"core-service/repositories"
	"core-service/services"

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

	if err := db.AutoMigrate(&models.ProcessedEvent{}); err != nil {
		log.Fatalf("failed to migrate processed_events table: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: getenv("REDIS_ADDR", "localhost:6379"),
	})

	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	eventRepo := repositories.NewEventRepository(db)
	cacheRepo := repositories.NewCacheRepository(redisClient)
	eventService := services.NewEventService(eventRepo, cacheRepo)

	consumer := handlers.NewKafkaConsumer(
		getenv("KAFKA_BROKER", "localhost:9092"),
		getenv("KAFKA_TOPIC_USER_CREATED", "user.created"),
		getenv("KAFKA_GROUP_ID", "core-service-group"),
		eventService,
	)
	go consumer.Consume(ctx)

	router := gin.Default()
	httpHandler := handlers.NewHTTPHandler()
	httpHandler.RegisterRoutes(router)

	port := getenv("PORT", "8082")
	log.Printf("core-service listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start core-service: %v", err)
	}
}
