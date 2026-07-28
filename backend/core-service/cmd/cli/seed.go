package cli

import (
	"core-service/database"
	"log"

	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with initial data",
	Long:  `Run database seeders to populate the database with initial data of cars, packages and regions.`,
	Run:   runSeed,
}

var (
	seedDryRun bool
	seedList   []string
)

func init() {
	rootCmd.AddCommand(seedCmd)

	seedCmd.Flags().BoolVar(&seedDryRun, "dry-run", false, "Show what would be seeded without running")
	seedCmd.Flags().StringSliceVarP(&seedList, "only", "", nil, "Seed only specific seeders (e.g., --only=packages,cars)")
}

func runSeed(cmd *cobra.Command, args []string) {
	// Load config first
	LoadConfig()

	db, err := gorm.Open(postgres.Open(getDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if seedDryRun {
		log.Println("Dry run mode - no changes will be made")
		log.Println("Available seeders:")
		seederList := []string{
			"provinces", "regencies", "districts",
			"cars", "packages",
		}
		for _, s := range seederList {
			log.Printf("  - %s", s)
		}
		return
	}

	log.Println("Starting database seeding...")

	// Run all seeders by default, or specific ones if specified
	if len(seedList) > 0 {
		for _, name := range seedList {
			log.Printf("Running seeder: %s", name)
			if err := database.RunSeederByName(db, name); err != nil {
				log.Fatalf("Failed to run seeder '%s': %v", name, err)
			}
		}
	} else {
		if err := database.RunSeeders(db); err != nil {
			log.Fatalf("Failed to run seeders: %v", err)
		}
	}

	log.Println("Database seeding completed successfully")
}