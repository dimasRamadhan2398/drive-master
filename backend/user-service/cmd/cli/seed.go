package cli

import (
	"log"

	"user-service/database/seeders"

	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with initial data",
	Long:  `Run database seeders to populate the database with initial data like roles and sample users.`,
	Run:   runSeed,
}

var (
	seedDryRun bool
	seedList   []string
)

func init() {
	rootCmd.AddCommand(seedCmd)

	seedCmd.Flags().BoolVar(&seedDryRun, "dry-run", false, "Show what would be seeded without running")
	seedCmd.Flags().StringSliceVarP(&seedList, "only", "", nil, "Seed only specific seeders (e.g., --only=roles,users)")
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
			"roles", "users", "work_experiences", "instructor_areas",
		}
		for _, s := range seederList {
			log.Printf("  - %s", s)
		}
		return
	}

	log.Println("Starting database seeding...")

	seederRunner := seeders.NewSeederRunner(db)

	// Filter seeders if specified
	if len(seedList) > 0 {
		// TODO: Implement selective seeder running
		log.Printf("Running only: %v", seedList)
	}

	if err := seederRunner.RunAll(); err != nil {
		log.Fatalf("Failed to run seeders: %v", err)
	}

	log.Println("Database seeding completed successfully")
}