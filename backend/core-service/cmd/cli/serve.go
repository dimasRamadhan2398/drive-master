package cli

import (
	"context"
	"core-service/controllers"
	"core-service/database"
	"core-service/handlers"
	"core-service/pkg/config"
	"core-service/pkg/kafka"
	"core-service/pkg/logger"
	"core-service/pkg/middlewares"
	"core-service/repositories"
	"core-service/routes"
	"core-service/services"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"core-service/docs"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Long:  `Start the core-service HTTP server with all routes and middleware configured.`,
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

	serveCmd.Flags().StringVarP(&servePort, "port", "p", "8001", "Port to listen on")
	serveCmd.Flags().StringVar(&serveHost, "host", "0.0.0.0", "Host to bind to")
	serveCmd.Flags().BoolVar(&serveSwagger, "swagger", true, "Enable Swagger documentation")
	serveCmd.Flags().BoolVar(&serveMigrate, "migrate", true, "Run database migrations on startup")
	serveCmd.Flags().BoolVar(&serveSeed, "seed", true, "Run database seeders on startup")
}

func runServe(cmd *cobra.Command, args []string) {
	// Load config using environment variable or default path
	configPath := getEnv("CONFIG_PATH", "pkg/config/config.yaml")
	loadedConfig, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config from %s: %v", configPath, err)
	}
	config.Set(loadedConfig)

	// Initialize logger with LogConfig
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

	logger.Info("Starting Core Service",
		logger.NewLogField("port", loadedConfig.Server.Port),
		logger.NewLogField("host", loadedConfig.Database.Host))

	// Set Gin mode
	if loadedConfig.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else if loadedConfig.Server.Mode == "debug" {
		gin.SetMode(gin.DebugMode)
	}

	db, err := gorm.Open(postgres.Open(getDSN()), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect to database",
			logger.NewLogField("error", err))
	}

	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", loadedConfig.Redis.Host, loadedConfig.Redis.Port),
		Password: loadedConfig.Redis.Password,
		DB:       loadedConfig.Redis.DB,
	})
	defer redisClient.Close()

	// Initialize Kafka producer
	var kafkaProducer *kafka.KafkaProducer
	if loadedConfig.Kafka.Enabled {
		kafkaProducer = kafka.NewKafkaProducer(
			loadedConfig.Kafka.Brokers[0],
			loadedConfig.Kafka.Topic,
		)
		defer kafkaProducer.Close()
	}

	// Initialize event publisher
	var eventPublisher *kafka.EventPublisher
	if kafkaProducer != nil {
		eventPublisher = kafka.NewEventPublisher(kafkaProducer)
	}

	// Run migrations if enabled
	if serveMigrate {
		if err := database.Migrate(db); err != nil {
			logger.Error("Failed to run migrations",
				logger.NewLogField("error", err))
		} else {
			logger.Info("Database migrations completed successfully")
		}
	}

	// Run seeders if enabled
	if serveSeed {
		if err := database.RunSeeders(db); err != nil {
			logger.Error("Failed to run seeders",
				logger.NewLogField("error", err))
		} else {
			logger.Info("Database seeders completed successfully")
		}
	}

	// Initialize repositories
	repoRegistry := repositories.NewRepositoryRegistry(db)
	repoRegistry.SetCacheClient(redisClient)

	// Initialize services
	serviceRegistry := services.NewServiceRegistry(repoRegistry, eventPublisher)

	// Initialize and start Kafka consumer if enabled
	var kafkaConsumer *handlers.KafkaConsumer
	kafkaCtx, kafkaCancel := context.WithCancel(context.Background())
	if loadedConfig.Kafka.Enabled && len(loadedConfig.Kafka.Topics) > 0 {
		kafkaConsumer = handlers.NewKafkaConsumer(
			loadedConfig.Kafka.Brokers,
			loadedConfig.Kafka.Topics,
			loadedConfig.Kafka.ConsumerGroup,
			serviceRegistry.GetEventService(),
		)
		go kafkaConsumer.Consume(kafkaCtx)
		logger.Info("Kafka consumer started",
			logger.NewLogField("topics", loadedConfig.Kafka.Topics),
			logger.NewLogField("group", loadedConfig.Kafka.ConsumerGroup))
	}

	// Initialize controller
	controllerRegistry := controllers.NewControllerRegistry(serviceRegistry)

	// Use gin.New() instead of gin.Default() to have full control
	router := gin.New()

	// Add logger and recovery middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Explicitly disable trailing slash redirect to prevent CORS issues on redirects
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	// // CORS middleware - MUST be added BEFORE routes
	// router.Use(func(c *gin.Context) {
	// 	c.Header("Access-Control-Allow-Origin", "*")
	// 	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	// 	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
	// 	if c.Request.Method == "OPTIONS" {
	// 		c.AbortWithStatus(http.StatusNoContent)
	// 		return
	// 	}
	// 	c.Next()
	// })

	maxRequests := float64(loadedConfig.App.RateLimiterMax)
	expirationTTL := time.Duration(loadedConfig.App.RateLimiterTime) * time.Second
	router.Use(middlewares.RateLimiter(maxRequests, expirationTTL))

	// Initialize Swagger docs
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "Core Service API"
	docs.SwaggerInfo.Description = "API documentation for Core Service"
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", loadedConfig.Server.Port)
	docs.SwaggerInfo.BasePath = "/api/v1"

	// Swagger documentation
	if serveSwagger {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "core-service",
		})
	})

	// Setup routes
	group := router.Group("/api/v1")
	route := routes.NewRouteRegistry(controllerRegistry, group)
	route.Serve()

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	addr := fmt.Sprintf("%s:%d", serveHost, loadedConfig.Server.Port)
	go func() {
		logger.Info("Server started", logger.NewLogField("address", addr))
		if err := router.Run(addr); err != nil {
			logger.Error("Server error", logger.NewLogField("error", err))
		}
	}()

	// Wait for shutdown signal
	<-quit
	logger.Info("Shutting down server...")

	// Stop Kafka consumer first
	if kafkaConsumer != nil {
		kafkaCancel()
		kafkaConsumer.Stop()
		logger.Info("Kafka consumer stopped")
	}

	logger.Info("Server shutdown complete")
}
