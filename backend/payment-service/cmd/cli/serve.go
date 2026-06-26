package cli

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"payment-service/controllers"
	"payment-service/models"
	"payment-service/pkg/config"
	pkgKafka "payment-service/pkg/kafka"
	"payment-service/pkg/logger"
	"payment-service/pkg/middlewares"
	"payment-service/pkg/redis"
	"payment-service/repositories"
	"payment-service/routes"
	"payment-service/services"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Long:  `Start the payment-service HTTP server with all routes and middleware configured.`,
	Run:   runServe,
}

var (
	servePort     string
	serveHost     string
	serveSwagger  bool
	serveMigrate  bool
	serveSeed     bool
)

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVarP(&servePort, "port", "p", "8004", "Port to listen on")
	serveCmd.Flags().StringVar(&serveHost, "host", "0.0.0.0", "Host to bind to")
	serveCmd.Flags().BoolVar(&serveSwagger, "swagger", true, "Enable Swagger documentation")
	serveCmd.Flags().BoolVar(&serveMigrate, "migrate", true, "Run database migrations on startup")
	serveCmd.Flags().BoolVar(&serveSeed, "seed", false, "Run database seeders on startup")
}

func runServe(cmd *cobra.Command, args []string) {
	// Load config using environment variable or default path
	configPath := getEnv("CONFIG_PATH", "pkg/config/config.yaml")
	loadedConfig, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config from %s: %v", configPath, err)
	}
	config.Set(loadedConfig)

	// Initialize logger
	loggerConfig := &logger.Config{
		Level:  loadedConfig.Log.Level,
		Format: loadedConfig.Log.Format,
	}
	if err := logger.Init(loggerConfig); err != nil {
		panic(err)
	}
	defer logger.Sync()

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err)
	}
	time.Local = loc

	logger.Info("Starting Payment Service",
		logger.LogField("port", loadedConfig.Server.Port),
		logger.LogField("host", loadedConfig.Database.Host),
	)

	// Set Gin mode
	if loadedConfig.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else if loadedConfig.Server.Mode == "debug" {
		gin.SetMode(gin.DebugMode)
	}

	db, err := gorm.Open(postgres.Open(getDSN()), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		logger.Fatal("Failed to connect to database", logger.LogField("error", err))
	}

	// Initialize Redis
	redisClient, err := redis.NewRedisConnection(&loadedConfig.Redis)
	if err != nil {
		logger.Fatal("Failed to connect to Redis", logger.LogField("error", err))
	}
	defer redisClient.Close()

	// Initialize Kafka producer
	kafkaProducer, err := pkgKafka.NewProducer(pkgKafka.Config{
		Brokers:     loadedConfig.Kafka.Brokers,
		Topic:       loadedConfig.Kafka.Topic,
		ServiceName: loadedConfig.Kafka.ServiceName,
		Enabled:     loadedConfig.Kafka.Enabled,
		UseAsync:    false,
	}, logger.GetLogger())
	if err != nil {
		logger.Warn("Failed to initialize Kafka producer", logger.LogField("error", err))
	}
	defer kafkaProducer.Close()

	// Initialize Kafka consumer
	kafkaConsumer, err := pkgKafka.NewConsumer(pkgKafka.ConsumerConfig{
		Brokers:  loadedConfig.Kafka.Brokers,
		GroupID:  "payment-service-consumer",
		Topics:   []string{loadedConfig.Kafka.Topic},
		Enabled:  loadedConfig.Kafka.Enabled,
		Version:  "3.6.0",
		Assignor: "roundrobin",
	}, logger.GetLogger())
	if err != nil {
		logger.Warn("Failed to initialize Kafka consumer", logger.LogField("error", err))
	}

	// Initialize event publisher
	eventPublisher := pkgKafka.NewEventPublisher(pkgKafka.EventPublisherConfig{
		Producer: kafkaProducer,
		Consumer: kafkaConsumer,
		Topic:    loadedConfig.Kafka.Topic,
		Enabled:  loadedConfig.Kafka.Enabled,
	})

	// Run migrations if enabled
	if serveMigrate {
		runMigrations(db)
	}

	// Run seeders if enabled
	if serveSeed {
		runSeeders(db)
	}

	// Initialize repositories
	repoRegistry := repositories.NewRepositoryRegistry(db)

	// Initialize auth middleware with JWT secret
	authMiddleware := middlewares.NewAuthMiddleware(loadedConfig.JWT.Secret)

	router := gin.Default()

	// Rate limiter configuration
	maxRequests := float64(loadedConfig.App.RateLimiterMax)
	expirationTTL := time.Duration(loadedConfig.App.RateLimiterTime) * time.Second

	// CORS and rate limiter middleware - MUST be added BEFORE routes
	router.Use(middlewares.CORSMiddleware())
	router.Use(middlewares.RateLimiter(maxRequests, expirationTTL))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "payment-service",
		})
	})

	// Setup API routes (basic for now)
	apiGroup := router.Group("/api/v1")
	apiGroup.Use(authMiddleware.Authenticate())

	// Initialize service registry
	serviceRegistry := services.NewServiceRegistry(repoRegistry, &loadedConfig.Midtrans)

	// Create public API router group for unauthenticated endpoints (like webhook callbacks)
	publicApi := router.Group("/api/v1")

	// Initialize and register transaction routes
	transactionController := controllers.NewTransactionController(
		repoRegistry.GetTransaction(),
		repoRegistry.GetPayment(),
		repoRegistry.GetPaymentMethod(),
		serviceRegistry.GetMidtrans(),
		eventPublisher,
	)
	transactionRoute := routes.NewTransactionRoute(transactionController, apiGroup, publicApi, authMiddleware)
	transactionRoute.Run()

	addr := fmt.Sprintf("%s:%d", serveHost, loadedConfig.Server.Port)
	log.Printf("Payment Service listening on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	_ = redisClient    // unused for now
}

func runMigrations(db *gorm.DB) {
	log.Println("Running database migrations...")

	if err := db.AutoMigrate(
		&models.Payment{},
		&models.Transaction{},
		&models.PaymentMethod{},
		&models.Refund{},
	); err != nil {
		log.Fatalf("Failed to migrate tables: %v", err)
	}

	log.Println("Database migrations completed successfully")
}

func runSeeders(db *gorm.DB) {
	log.Println("Starting database seeding...")
	// Add seeders here if needed
	log.Println("Database seeding completed successfully")
}