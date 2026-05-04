package cli

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"booking-service/controllers"
	"booking-service/database/seeders"
	"booking-service/models"
	"booking-service/repositories"
	"booking-service/routes"
	"booking-service/services"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/redis/go-redis/v9"
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
	serveMigrate bool
	serveSeed    bool
)

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVarP(&servePort, "port", "p", "8082", "Port to listen on")
	serveCmd.Flags().StringVar(&serveHost, "host", "0.0.0.0", "Host to bind to")
	serveCmd.Flags().BoolVar(&serveMigrate, "migrate", true, "Run database migrations on startup")
	serveCmd.Flags().BoolVar(&serveSeed, "seed", false, "Run database seeders on startup")
}

func runServe(cmd *cobra.Command, args []string) {
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

	// Run seeders if enabled
	if serveSeed {
		runSeeders(db)
	}

	// Initialize repositories
	bookingRepo := repositories.NewBookingRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)
	entitlementRepo := repositories.NewEntitlementRepository(db)
	certificationRepo := repositories.NewCertificationRepository(db)
	userCertRepo := repositories.NewUserCertificationRepository(db)

	// Initialize services
	bookingService := services.NewBookingService(bookingRepo, entitlementRepo)
	sessionService := services.NewSessionService(sessionRepo)
	entitlementService := services.NewEntitlementService(entitlementRepo)
	certificationService := services.NewCertificationService(certificationRepo, userCertRepo)

	// Create service registry
	serviceRegistry := &serviceRegistryImpl{
		bookingService:        bookingService,
		sessionService:        sessionService,
		entitlementService:    entitlementService,
		certificationService:  certificationService,
	}

	// Initialize controller registry
	controllerRegistry := controllers.NewControllerRegistry(serviceRegistry)

	// Setup Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "booking-service"})
	})

	// Initialize route registry and register all routes
	routeRegistry := routes.NewRouteRegistry(controllerRegistry, router)
	routeRegistry.Serve()

	addr := fmt.Sprintf("%s:%s", serveHost, servePort)
	log.Printf("Booking Service listening on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func runMigrations(db *gorm.DB) {
	log.Println("Running database migrations...")

	if err := db.AutoMigrate(
		&models.Booking{},
		&models.Session{},
		&models.UserEntitlement{},
		&models.Certification{},
		&models.UserCertification{},
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

// serviceRegistryImpl implements services.IServiceRegistry
type serviceRegistryImpl struct {
	bookingService        services.IBookingService
	sessionService        services.ISessionService
	entitlementService    services.IEntitlementService
	certificationService  services.ICertificationService
}

func (s *serviceRegistryImpl) GetBookingService() services.IBookingService {
	return s.bookingService
}

func (s *serviceRegistryImpl) GetSessionService() services.ISessionService {
	return s.sessionService
}

func (s *serviceRegistryImpl) GetEntitlementService() services.IEntitlementService {
	return s.entitlementService
}

func (s *serviceRegistryImpl) GetCertificationService() services.ICertificationService {
	return s.certificationService
}
