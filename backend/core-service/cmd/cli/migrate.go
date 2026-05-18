package cli

import (
	"core-service/database"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
	Long:  `Database migration commands for core-service: up, down, fresh, seed, reset, status`,
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run database migrations",
	Long:  `Run all database migrations to create/update tables`,
	Run:   runMigrateUp,
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback database migrations",
	Long:  `Drop all tables to rollback migrations (WARNING: destructive)`,
	Run:   runMigrateDown,
}

var migrateFreshCmd = &cobra.Command{
	Use:   "fresh",
	Short: "Fresh migration - drop and recreate all tables",
	Long:  `Drop all tables and run migrations from scratch (WARNING: destructive)`,
	Run:   runMigrateFresh,
}

var migrateSeedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed database with sample data",
	Long:  `Run database seeders to populate sample data`,
	Run:   runMigrateSeed,
}

var migrateResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset database - fresh migration + seed",
	Long:  `Drop all tables, run migrations, and seed data (WARNING: destructive)`,
	Run:   runMigrateReset,
}

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	Long:  `Show current database tables and row counts`,
	Run:   runMigrateStatus,
}

var (
	migrateTarget   string // Target table(s) for migration: all, regions, cars, packages
	migrateForce    bool   // Force operation without confirmation
	migrateSeedForce bool  // Force re-seed (delete existing data)
	migrateTables   []string // Specific tables for seeding
)

func init() {
	// Add subcommands
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateFreshCmd)
	migrateCmd.AddCommand(migrateSeedCmd)
	migrateCmd.AddCommand(migrateResetCmd)
	migrateCmd.AddCommand(migrateStatusCmd)
	rootCmd.AddCommand(migrateCmd)

	// Flags for migrate up
	migrateUpCmd.Flags().StringVarP(&migrateTarget, "target", "t", "all", "Target migration (all/regions/cars/packages)")
	migrateUpCmd.Flags().BoolVar(&migrateForce, "force", false, "Skip confirmation for potentially destructive operations")

	// Flags for migrate down
	migrateDownCmd.Flags().BoolVarP(&migrateForce, "yes", "y", false, "Skip confirmation prompt")

	// Flags for migrate fresh
	migrateFreshCmd.Flags().BoolVarP(&migrateForce, "yes", "y", false, "Skip confirmation prompt")

	// Flags for migrate seed
	migrateSeedCmd.Flags().StringSliceVarP(&migrateTables, "tables", "t", []string{"all"}, "Tables to seed (all/provinces/regencies/districts/cars/packages)")
	migrateSeedCmd.Flags().BoolVar(&migrateSeedForce, "force", false, "Delete existing data before seeding")

	// Flags for migrate reset
	migrateResetCmd.Flags().BoolVarP(&migrateForce, "yes", "y", false, "Skip confirmation prompt")
	migrateResetCmd.Flags().BoolVar(&migrateSeedForce, "force", false, "Force re-seed data")
}

func confirmAction(message string) bool {
	if migrateForce {
		return true
	}
	fmt.Printf("\n⚠️  %s\n", message)
	fmt.Print("Are you sure? [y/N]: ")
	var input string
	fmt.Scanln(&input)
	return input == "y" || input == "Y"
}

func runMigrateUp(cmd *cobra.Command, args []string) {
	LoadConfig()

	db, err := gorm.Open(postgres.Open(getDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer closeDB(db)

	log.Printf("Starting migration with target: %s", migrateTarget)

	if err := database.RunMigration(db, migrateTarget); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Printf("✅ Migration completed successfully for target: %s", migrateTarget)
}

func runMigrateDown(cmd *cobra.Command, args []string) {
	if !confirmAction("This will DROP all tables in the database! All data will be lost.") {
		fmt.Println("❌ Cancelled.")
		return
	}

	LoadConfig()

	db, err := gorm.Open(postgres.Open(getDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer closeDB(db)

	log.Println("📦 Dropping all tables...")

	// Tables in reverse dependency order (child tables first)
	tables := []string{
		"package_benefits",
		"packages",
		"cars",
		"districts",
		"regencies",
		"provinces",
	}

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)).Error; err != nil {
			log.Fatalf("Failed to drop table %s: %v", table, err)
		}
		log.Printf("  ✅ Dropped: %s", table)
	}

	log.Println("✅ All tables dropped successfully")
}

func runMigrateFresh(cmd *cobra.Command, args []string) {
	if !confirmAction("This will DROP and RECREATE all tables! All data will be lost.") {
		fmt.Println("❌ Cancelled.")
		return
	}

	LoadConfig()

	db, err := gorm.Open(postgres.Open(getDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer closeDB(db)

	log.Println("🔄 Starting fresh migration...")

	// Drop all tables first
	tables := []string{
		"package_benefits",
		"packages",
		"cars",
		"districts",
		"regencies",
		"provinces",
	}

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)).Error; err != nil {
			log.Printf("  ⚠️  Warning dropping %s: %v", table, err)
		}
	}

	log.Println("  ✅ All existing tables dropped")

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("✅ Fresh migration completed successfully")
}

func runMigrateSeed(cmd *cobra.Command, args []string) {
	LoadConfig()

	db, err := gorm.Open(postgres.Open(getDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer closeDB(db)

	// Handle force seeding - clear tables first if force flag is set
	if migrateSeedForce {
		log.Println("🔥 Force mode enabled - clearing existing data...")
		clearTablesForSeed(db)
	}

	// Determine which tables to seed
	seedAll := slices.Contains(migrateTables, "all")

	if seedAll {
		log.Println("🌱 Seeding all tables...")
		if err := database.RunSeeders(db); err != nil {
			log.Fatalf("Seeding failed: %v", err)
		}
	} else {
		for _, table := range migrateTables {
			log.Printf("  🌱 Seeding: %s", table)
			if err := database.RunSeederByName(db, table); err != nil {
				log.Fatalf("Failed to seed %s: %v", table, err)
			}
		}
	}

	log.Println("✅ Seeding completed successfully")
}

func runMigrateReset(cmd *cobra.Command, args []string) {
	if !confirmAction("This will DROP all tables, recreate them, and seed data! ALL DATA WILL BE LOST.") {
		fmt.Println("❌ Cancelled.")
		return
	}

	LoadConfig()

	db, err := gorm.Open(postgres.Open(getDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer closeDB(db)

	log.Println("🔄 Starting database reset...")

	// Step 1: Drop all tables
	log.Println("  Step 1/3: Dropping all tables...")
	tables := []string{
		"package_benefits",
		"packages",
		"cars",
		"districts",
		"regencies",
		"provinces",
	}

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)).Error; err != nil {
			log.Printf("  ⚠️  Warning dropping %s: %v", table, err)
		}
	}
	log.Println("  ✅ All tables dropped")

	// Step 2: Run migrations
	log.Println("  Step 2/3: Running migrations...")
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("  ✅ Migrations completed")

	// Step 3: Run seeders
	log.Println("  Step 3/3: Seeding data...")
	if err := database.RunSeeders(db); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}
	log.Println("  ✅ Seeding completed")

	log.Println("✅ Database reset completed successfully")
}

func runMigrateStatus(cmd *cobra.Command, args []string) {
	LoadConfig()

	db, err := gorm.Open(postgres.Open(getDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer closeDB(db)

	tables := []struct {
		name string
		check string
	}{
		{"provinces", "provinces"},
		{"regencies", "regencies"},
		{"districts", "districts"},
		{"cars", "cars"},
		{"packages", "packages"},
		{"package_benefits", "package_benefits"},
	}

	fmt.Println()
	fmt.Println("📊 Database Migration Status")
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Printf("Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println()

	for _, t := range tables {
		var count int64
		result := db.Raw(fmt.Sprintf("SELECT COUNT(*) FROM %s", t.name)).Scan(&count)

		status := "❌"
		statusColor := "\033[31m"

		if result.Error == nil {
			status = "✅"
			statusColor = "\033[32m"
		} else if strings.Contains(result.Error.Error(), "record not found") {
			status = "⚠️"
			statusColor = "\033[33m"
		}

		fmt.Printf("  %s%-20s\033[0m %s Rows: %d\n", statusColor, t.name, status, count)
	}

	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println()
}

func closeDB(db *gorm.DB) {
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
	}
}

func clearTablesForSeed(db *gorm.DB) {
	tables := []string{
		"package_benefits",
		"packages",
		"cars",
		"districts",
		"regencies",
		"provinces",
	}

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			// Fallback to DELETE if TRUNCATE fails
			log.Printf("  ⚠️  Truncate failed for %s, using DELETE", table)
			db.Exec(fmt.Sprintf("DELETE FROM %s", table))
		}
		log.Printf("  🗑️  Cleared: %s", table)
	}
}