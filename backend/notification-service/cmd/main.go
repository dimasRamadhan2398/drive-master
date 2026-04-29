package main

import (
	"log"
	"os"
	"notification-service/models"
	"notification-service/pkg/kafka"
	"notification-service/pkg/handlers"
	"notification-service/repositories"
	"notification-service/services"

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

	// Auto migrate models
	if err := db.AutoMigrate(
		&models.Notification{},
		&models.NotificationTemplate{},
		&models.NewsletterSubscription{},
		&models.NotificationPreference{},
	); err != nil {
		log.Fatalf("failed to migrate tables: %v", err)
	}

	// Initialize repository
	notificationRepo := repositories.NewNotificationRepository(db)

	// Initialize services
	emailService := services.NewEmailService(
		getenv("MAILTRAP_API_TOKEN", ""),
		getenv("MAILTRAP_FROM_EMAIL", "noreply@example.com"),
		getenv("MAILTRAP_FROM_NAME", "Drive Master"),
	)
	whatsappService := services.NewWhatsAppService(
		getenv("WHATSAPP_API_URL", ""),
		getenv("WHATSAPP_API_TOKEN", ""),
		getenv("WHATSAPP_PHONE_NUMBER", ""),
	)
	notificationService := services.NewNotificationService(notificationRepo, emailService, whatsappService)

	// Initialize HTTP handler
	httpHandler := handlers.NewHTTPHandler(notificationService)

	// Setup router
	router := gin.Default()
	httpHandler.RegisterRoutes(router)

	// Start Kafka consumer in goroutine
	go func() {
		kafkaBrokers := getenv("KAFKA_BROKERS", "localhost:9092")
		consumer := kafka.NewConsumer(
			kafkaBrokers,
			"notification-service",
			[]string{"session-upcoming", "booking-created", "promotional"},
		)
		consumer.Start(notificationService)
	}()

	port := getenv("PORT", "8083")
	log.Printf("notification-service listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start notification-service: %v", err)
	}
}