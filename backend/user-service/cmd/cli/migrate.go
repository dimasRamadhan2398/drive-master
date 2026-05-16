package cli

import (
	"log"
	"slices"

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
	migrateDown  bool
	migrateSteps int
	migrateDryRun bool
	migrateReset bool
)

func init() {
	rootCmd.AddCommand(migrateCmd)

	migrateCmd.Flags().BoolVar(&migrateReset, "reset", false, "Drop and recreate all tables")
	migrateCmd.Flags().BoolVar(&migrateDown, "down", false, "Roll back migrations")
	migrateCmd.Flags().IntVarP(&migrateSteps, "steps", "n", 0, "Number of migrations to run (positive for up, negative for down)")
	migrateCmd.Flags().BoolVar(&migrateDryRun, "dry-run", false, "Show what would be migrated without running")
}

func runMigrate(cmd *cobra.Command, args []string) {
	// Load config first
	LoadConfig()

	db, err := gorm.Open(postgres.Open(getDSN()), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
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

	if migrateReset {
		log.Println("Resetting database - dropping all tables...")
		if err := dropAllTables(db); err != nil {
			log.Fatalf("Failed to drop tables: %v", err)
		}
		log.Println("All tables dropped successfully")
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

// dropAllTables drops all tables in the database using GORM
func dropAllTables(db *gorm.DB) error {
	// Get all table names from the database
	var tables []string

	// Use the underlying SQL DB to get tables
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	rows, err := sqlDB.Query(`
		SELECT tablename FROM pg_tables
		WHERE schemaname = 'public'
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return err
		}
		tables = append(tables, tableName)
	}

	// Drop tables in reverse dependency order to avoid FK issues
	// Order matters for foreign key dependencies
	order := []string{
		// "users",                    // depends on nothing
		// "instructor_areas",         // depends on regions
		// "work_experiences",          // depends on users
		// "member_profiles",           // depends on users
		// "instructor_profiles",       // depends on users
		"roles",                    // depends on nothing
	}

	for _, table := range order {
		if err := db.Exec("DROP TABLE IF EXISTS " + table + " CASCADE").Error; err != nil {
			log.Printf("Warning: Failed to drop table %s: %v", table, err)
		} else {
			log.Printf("Dropped table: %s", table)
		}
	}

	// Also drop any tables not in our list
	for _, table := range tables {
		if !slices.Contains(order, table) {
			if err := db.Exec("DROP TABLE IF EXISTS " + table + " CASCADE").Error; err != nil {
				log.Printf("Warning: Failed to drop table %s: %v", table, err)
			} else {
				log.Printf("Dropped table: %s", table)
			}
		}
	}

	return nil
}