package cli

import (
	"fmt"
	"log"
	"net/http"

	"user-service/controllers"
	"user-service/database/seeders"
	docs "user-service/docs"
	"user-service/models"
	"user-service/pkg/config"
	pkgKafka "user-service/pkg/kafka"
	"user-service/pkg/logger"
	"user-service/pkg/middlewares"
	"user-service/pkg/redis"
	"user-service/repositories"
	"user-service/routes"
	"user-service/services"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Long:  `Start the user-service HTTP server with all routes and middleware configured.`,
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

	serveCmd.Flags().StringVarP(&servePort, "port", "p", "8081", "Port to listen on")
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
	if err := logger.Init(&loadedConfig.Log); err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info("Starting User Service", logger.LogField("port", loadedConfig.Server.Port), logger.LogField("host", loadedConfig.Database.Host))

	// Set Gin mode
	if loadedConfig.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else if loadedConfig.Server.Mode == "debug" {
		gin.SetMode(gin.DebugMode)
	}

	db, err := gorm.Open(postgres.Open(getDSN()), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect to database: %v", logger.LogField("error", err))
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
		GroupID:  "user-service-consumer",
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

	// TODO: Pass eventPublisher to services that need it
	_ = eventPublisher
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

	// Initialize services
	serviceRegistry := services.NewServiceRegistry(repoRegistry)

	// Initialize controller
	userController := controllers.NewUserController(
		serviceRegistry.GetUserService(),
		serviceRegistry.GetAuthService(),
		serviceRegistry.GetMemberService(),
		serviceRegistry.GetInstructorService(),
		serviceRegistry.GetRoleService(),
		serviceRegistry.GetEmailService(),
		serviceRegistry.GetMediaService(),
	)

	router := gin.Default()

	// Initialize Swagger docs
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "User Service API"
	docs.SwaggerInfo.Description = "API documentation for User Service"
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", servePort)
	docs.SwaggerInfo.BasePath = "/"

	// Setup routes
	routes.SetupRoutes(router, userController)

	router.Use(middlewares.CORSMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "auth-service",
		})
	})

	// Swagger documentation
	if serveSwagger {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	addr := fmt.Sprintf("%s:%s", serveHost, servePort)
	log.Printf("User Service listening on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func runMigrations(db *gorm.DB) {
	log.Println("Running database migrations...")

	if err := db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.MemberProfile{},
		&models.InstructorProfile{},
		&models.WorkExperience{},
		&models.InstructorArea{},
	); err != nil {
		log.Fatalf("Failed to migrate tables: %v", err)
	}

	log.Println("Database migrations completed successfully")
}

func runSeeders(db *gorm.DB) {
	log.Println("Starting database seeding...")

	seederRunner := seeders.NewSeederRunner(db)
	if err := seederRunner.RunAll(); err != nil {
		log.Fatalf("Failed to run seeders: %v", err)
	}
}