package cli

import (
	"log"

	"booking-service/models"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Run database migrations to create or update database schema.`,
	Run:   runMigrate,
}

var (
	migrateDown   bool
	migrateSteps  int
	migrateDryRun bool
)

func init() {
	rootCmd.AddCommand(migrateCmd)

	migrateCmd.Flags().BoolVar(&migrateDown, "down", false, "Roll back migrations")
	migrateCmd.Flags().IntVarP(&migrateSteps, "steps", "n", 0, "Number of migrations to run (positive for up, negative for down)")
	migrateCmd.Flags().BoolVar(&migrateDryRun, "dry-run", false, "Show what would be migrated without running")
}

func runMigrate(cmd *cobra.Command, args []string) {
	db := connectDB()

	if migrateDryRun {
		log.Println("Dry run mode - no changes will be made")
		log.Println("Would run migrations for the following models:")
		modelsList := []string{
			"Booking", "Session", "UserEntitlement",
			"Certification", "UserCertification",
		}
		for _, m := range modelsList {
			log.Printf("  - %s", m)
		}
		return
	}

	if migrateDown {
		log.Println("Rollback is not yet implemented. Use --steps for incremental migrations.")
		return
	}

	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&models.Booking{},
		&models.Session{},
		&models.UserEntitlement{},
		&models.Certification{},
		&models.UserCertification{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate tables: %v", err)
	}

	log.Println("Database migrations completed successfully")
}
