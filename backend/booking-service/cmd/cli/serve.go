package cli

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"booking-service/clients/core"
	"booking-service/clients/user"
	"booking-service/controllers"
	"booking-service/docs"
	"booking-service/models"
	"booking-service/pkg/config"
	"booking-service/pkg/kafka"
	"booking-service/pkg/middlewares"
	"booking-service/pkg/scheduler"
	"booking-service/repositories"
	"booking-service/routes"
	"booking-service/services"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Long:  `Start the booking-service HTTP server with all routes and middleware configured.`,
	Run:   runServe,
}

var (
	servePort    string
	serveHost    string
	serveSwagger  bool
	serveMigrate bool
)

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVarP(&servePort, "port", "p", "8003", "Port to listen on")
	serveCmd.Flags().StringVar(&serveHost, "host", "0.0.0.0", "Host to bind to")
	serveCmd.Flags().BoolVar(&serveSwagger, "swagger", true, "Enable Swagger documentation")
	serveCmd.Flags().BoolVar(&serveMigrate, "migrate", true, "Run database migrations on startup")
}

func runServe(cmd *cobra.Command, args []string) {
	// Load config using environment variable or default path
	configPath := getEnv("CONFIG_PATH", "pkg/config/config.yaml")
	loadedConfig, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config from %s: %v", configPath, err)
	}
	config.Set(loadedConfig)

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err)
	}
	time.Local = loc

	log.Printf("Starting Booking Service on %s:%s", serveHost, servePort)

	db, err := gorm.Open(postgres.Open(getDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: getEnv("REDIS_ADDR", "localhost:6379"),
	})
	defer redisClient.Close()

	// Run migrations if enabled
	if serveMigrate {
		runMigrations(db)
	}

	// Initialize repositories
	sessionRepo := repositories.NewSessionRepository(db)
	entitlementRepo := repositories.NewEntitlementRepository(db)
	enrollmentRepo := repositories.NewEnrollmentRepository(db)
	scheduleRepo := repositories.NewScheduleRepository(db)
	paymentRepo := repositories.NewPaymentRepository(db)

	// Initialize user-service client and availability service
	userClient := user.NewUserClient(getEnv("USER_SERVICE_URL", "http://localhost:8001"), loadedConfig.JWT.Secret)
	coreClient := core.NewCoreClient(getEnv("CORE_SERVICE_URL", "http://localhost:8002"))
	availabilityService := services.NewAvailabilityService(userClient)

	// Initialize user anonymization service for handling user.deleted events
	userAnonymizationService := services.NewUserAnonymizationService(enrollmentRepo, sessionRepo, entitlementRepo)

	// Initialize Kafka consumer for handling user.deleted events
	_ = initKafkaConsumer(userAnonymizationService)

	// Initialize services (after Kafka is initialized so we can pass the eventPublisher)
	sessionService := services.NewSessionService(sessionRepo, scheduleRepo, entitlementRepo)
	entitlementService := services.NewEntitlementService(entitlementRepo)
	enrollmentService := services.NewEnrollmentService(enrollmentRepo, entitlementRepo)

	scheduleService := services.NewScheduleService(scheduleRepo, enrollmentRepo, sessionRepo, entitlementRepo, sessionService, availabilityService, userClient, coreClient)
	paymentService := services.NewPaymentService(paymentRepo, enrollmentRepo)
	revenueService := services.NewRevenueService(coreClient)

	// Create service registry
	serviceRegistry := &serviceRegistryImpl{
		sessionService:     sessionService,
		entitlementService: entitlementService,
		enrollmentService:  enrollmentService,
		scheduleService:    scheduleService,
		paymentService:     paymentService,
		revenueService:     revenueService,
	}

	// Initialize schedule generator for automatic schedule slot generation
	_ = initScheduleGenerator(scheduleRepo, userClient)

	// Initialize controller registry
	controllerRegistry := controllers.NewControllerRegistry(serviceRegistry)

	// Initialize auth middleware
	authMiddleware := middlewares.NewAuthMiddleware(loadedConfig.JWT.Secret)

	// Setup Gin router
	router := gin.New()

	// Add logger and recovery middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Explicitly disable trailing slash redirect to prevent CORS issues on redirects
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "booking-service"})
	})

	// Swagger documentation endpoint
	docs.SwaggerInfo.Title = "Booking Service API"
	docs.SwaggerInfo.Description = "API for managing bookings, sessions, entitlements, and certifications"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", servePort)
	docs.SwaggerInfo.BasePath = "/api/v1"
	target := `"securityDefinitions"`
    security := `"security":[{"BearerAuth":[], "XApiKey":[],"XRequestAt":[],"XServiceName":[]}],`
    docs.SwaggerInfo.SwaggerTemplate = strings.Replace(
        docs.SwaggerInfo.SwaggerTemplate,
        target,
        security+target,
        1,
    )

	// Swagger documentation
	if serveSwagger {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Setup routes
	group := router.Group("/api/v1")

	// Initialize route registry and register all routes
	routeRegistry := routes.NewRouteRegistry(controllerRegistry, group, authMiddleware)
	routeRegistry.Serve()

	addr := fmt.Sprintf("%s:%s", serveHost, servePort)
	log.Printf("Booking Service listening on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func runMigrations(db *gorm.DB) {
	log.Println("Running database migrations...")

	// Fix column type mismatch before AutoMigrate
	// The user_entitlements.enrollment_id column might be bigint from old schema
	// but the model expects uuid
	fixEnrollmentIDColumnType(db)

	if err := db.AutoMigrate(
		&models.Enrollment{},
		&models.UserEntitlement{},
		&models.DrivingSession{},
		&models.Schedule{},
		&models.Payment{},
	); err != nil {
		log.Fatalf("Failed to migrate tables: %v", err)
	}

	log.Println("Database migrations completed successfully")
}

// fixEnrollmentIDColumnType checks and alters the enrollment_id column type in user_entitlements table
// This handles schema drift where the enrollment tables might have bigint IDs but should be uuid
func fixEnrollmentIDColumnType(db *gorm.DB) {
	// Check if enrollments.id is bigint (old schema) and needs to be converted
	var enrollmentIDType string
	err := db.Raw(`
		SELECT data_type
		FROM information_schema.columns
		WHERE table_name = 'enrollments'
		AND column_name = 'id'
	`).Scan(&enrollmentIDType).Error

	if err != nil {
		log.Printf("Could not check enrollments.id column type: %v", err)
		return
	}

	// If enrollments.id is still bigint, we have a schema mismatch - drop and recreate both tables
	if enrollmentIDType == "bigint" {
		log.Println("Found schema mismatch: enrollments.id is bigint but should be uuid. Dropping tables for recreation...")
		
		// Drop dependent tables first (due to foreign keys)
		tables := []string{"user_entitlements", "driving_sessions", "payments", "enrollments"}
		for _, table := range tables {
			err = db.Exec(`DROP TABLE IF EXISTS "` + table + `" CASCADE`).Error
			if err != nil {
				log.Printf("Warning: could not drop %s table: %v", table, err)
			} else {
				log.Printf("Successfully dropped %s table", table)
			}
		}
	}
}

// serviceRegistryImpl implements services.IServiceRegistry
type serviceRegistryImpl struct {
	sessionService     services.ISessionService
	entitlementService services.IEntitlementService
	enrollmentService  services.IEnrollmentService
	scheduleService    services.IScheduleService
	paymentService     services.IPaymentService
	revenueService     services.IRevenueService
}

func (s *serviceRegistryImpl) GetSessionService() services.ISessionService {
	return s.sessionService
}

func (s *serviceRegistryImpl) GetEntitlementService() services.IEntitlementService {
	return s.entitlementService
}

func (s *serviceRegistryImpl) GetEnrollmentService() services.IEnrollmentService {
	return s.enrollmentService
}

func (s *serviceRegistryImpl) GetScheduleService() services.IScheduleService {
	return s.scheduleService
}

func (s *serviceRegistryImpl) GetPaymentService() services.IPaymentService {
	return s.paymentService
}

func (s *serviceRegistryImpl) GetRevenueService() services.IRevenueService {
	return s.revenueService
}

// initKafkaConsumer initializes the Kafka consumer for handling user.deleted events
// Returns the event publisher for use in other services
func initKafkaConsumer(anonymizationService services.IUserAnonymizationService) kafka.IEventPublisher {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Get Kafka configuration from environment
	kafkaEnabled := getEnv("KAFKA_ENABLED", "false") == "true"
	if !kafkaEnabled {
		logger.Info("Kafka consumer is disabled, skipping initialization")
		return nil
	}

	kafkaBrokers := strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ",")
	kafkaTopic := getEnv("KAFKA_TOPIC", "user-events")
	kafkaGroupID := getEnv("KAFKA_GROUP_ID", "booking-service")

	// Create Kafka consumer
	consumer, err := kafka.NewConsumer(kafka.ConsumerConfig{
		Brokers:  kafkaBrokers,
		GroupID:  kafkaGroupID,
		Topics:   []string{kafkaTopic},
		Enabled:  true,
		Version:  "3.6.0",
		Assignor: "roundrobin",
	}, logger)
	if err != nil {
		logger.Error("Failed to create Kafka consumer", zap.Error(err))
		return nil
	}

	// Create event publisher with the consumer
	eventPublisher := kafka.NewEventPublisher(kafka.EventPublisherConfig{
		Consumer: consumer,
		Topic:    kafkaTopic,
		Enabled:  true,
	})

	// Create and register the user deleted handler
	userDeletedHandler := kafka.NewUserDeletedHandler(func(ctx context.Context, userID string) error {
		return anonymizationService.AnonymizeUserData(ctx, userID)
	})
	eventPublisher.RegisterHandler(userDeletedHandler)

	// Start the consumer in a goroutine
	ctx := context.Background()
	if err := eventPublisher.StartConsumer(ctx); err != nil {
		logger.Error("Failed to start Kafka consumer", zap.Error(err))
		return nil
	}

	logger.Info("Kafka consumer initialized and started for user.deleted events")
	return eventPublisher
}

// initScheduleGenerator initializes the schedule generator for automatic schedule slot generation
func initScheduleGenerator(scheduleRepo repositories.IScheduleRepository, userClient user.IUserClient) *scheduler.ScheduleGenerator {
	scheduleGenerator := scheduler.NewScheduleGenerator(scheduleRepo, userClient)

	// Get cron expression from environment or use default
	// Default: "5 0 * * *" runs at 00:05 AM every day
	cronExpr := getEnv("SCHEDULE_GENERATOR_CRON", "5 0 * * *")

	// Get number of days to generate ahead (default: 7)
	generationDays := 7
	if days := getEnv("SCHEDULE_GENERATION_DAYS", ""); days != "" {
		if parsed, err := parseInt(days); err == nil && parsed > 0 {
			generationDays = parsed
		}
	}
	scheduleGenerator.SetGenerationDays(generationDays)

	// Check if scheduler is enabled
	enabled := getEnv("SCHEDULE_GENERATOR_ENABLED", "true") == "true"
	if !enabled {
		log.Println("Schedule generator is disabled, skipping initialization")
		return scheduleGenerator
	}

	// Start the scheduler
	if err := scheduleGenerator.Start(cronExpr); err != nil {
		log.Printf("Failed to start schedule generator: %v", err)
		return scheduleGenerator
	}

	log.Printf("Schedule generator started with cron: %s (generating %d days ahead)", cronExpr, generationDays)

	// Run once immediately on startup to populate initial schedules
	go func() {
		// Wait for user-service to be available (max 30 seconds)
		userServiceURL := getEnv("USER_SERVICE_URL", "http://localhost:8001")
		maxRetries := 10
		retryDelay := 3 * time.Second

		for i := 0; i < maxRetries; i++ {
			resp, err := http.Get(userServiceURL + "/health")
			if err == nil {
				resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					break
				}
			}
			if i < maxRetries-1 {
				log.Printf("Waiting for user-service to be available (attempt %d/%d)...", i+1, maxRetries)
				time.Sleep(retryDelay)
			}
		}

		ctx := context.Background()
		if err := scheduleGenerator.RunOnce(ctx); err != nil {
			log.Printf("Initial schedule generation failed: %v", err)
		} else {
			log.Println("Initial schedule generation completed")
		}
	}()

	return scheduleGenerator
}

// parseInt parses a string to int
func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}
