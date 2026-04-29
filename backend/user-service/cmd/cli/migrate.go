package cli

import (
	"log"

	"user-service/models"

	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	// Load config first
	LoadConfig()

	db, err := gorm.Open(postgres.Open(getDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if migrateDryRun {
		log.Println("Dry run mode - no changes will be made")
		log.Println("Would run migrations for the following models:")
		modelsList := []string{
			"User", "Role", "MemberProfile", "InstructorProfile",
			"WorkExperience", "InstructorArea",
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

	err = db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.MemberProfile{},
		&models.InstructorProfile{},
		&models.WorkExperience{},
		&models.InstructorArea{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate tables: %v", err)
	}

	log.Println("Database migrations completed successfully")
}