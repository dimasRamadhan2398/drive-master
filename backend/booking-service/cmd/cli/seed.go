package cli

import (
	"log"

	"booking-service/database/seeders"

	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with initial data",
	Long:  `Run database seeders to populate the database with initial data.`,
	Run:   runSeed,
}

var (
	seedDryRun bool
	seedList   []string
)

func init() {
	rootCmd.AddCommand(seedCmd)

	seedCmd.Flags().BoolVar(&seedDryRun, "dry-run", false, "Show what would be seeded without running")
	seedCmd.Flags().StringSliceVarP(&seedList, "only", "", nil, "Seed only specific seeders (e.g., --only=entitlements)")
}

func runSeed(cmd *cobra.Command, args []string) {
	db := connectDB()

	if seedDryRun {
		log.Println("Dry run mode - no changes will be made")
		log.Println("Available seeders:")
		seederList := []string{
			"entitlements", "certifications",
		}
		for _, s := range seederList {
			log.Printf("  - %s", s)
		}
		return
	}

	log.Println("Starting database seeding...")

	seederRunner := seeders.NewSeederRunner(db)

	if len(seedList) > 0 {
		if err := seederRunner.RunSpecific(seedList); err != nil {
			log.Fatalf("Failed to run seeders: %v", err)
		}
	} else {
		if err := seederRunner.RunAll(); err != nil {
			log.Fatalf("Failed to run seeders: %v", err)
		}
	}

	log.Println("Database seeding completed successfully")
}
