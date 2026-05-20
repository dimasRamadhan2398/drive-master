package cli

import (
	"log"

	"payment-service/models"

	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Run database migrations for payment-service. Use --reset to reset the database first.`,
	Run:   runMigrate,
}

var (
	migrateReset bool
)

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().BoolVar(&migrateReset, "reset", false, "Drop all tables and recreate them")
}

func runMigrate(cmd *cobra.Command, args []string) {
	// Load config
	LoadConfig()

	// Connect to database
	db, err := gorm.Open(postgres.Open(getDSN()), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Reset database if flag is set
	if migrateReset {
		log.Println("WARNING: Dropping all tables...")
		if err := db.Migrator().DropTable(
			&models.Payment{},
			&models.Transaction{},
			&models.PaymentMethod{},
		); err != nil {
			log.Fatalf("Failed to drop tables: %v", err)
		}
		log.Println("All tables dropped successfully")
	}

	// Run migrations
	log.Println("Running database migrations...")
	if err := db.AutoMigrate(
		&models.Payment{},
		&models.Transaction{},
		&models.PaymentMethod{},
	); err != nil {
		log.Fatalf("Failed to migrate tables: %v", err)
	}

	log.Println("Database migrations completed successfully")
}