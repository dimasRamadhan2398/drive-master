package cli

import (
	"log"

	"payment-service/pkg/config"
	"payment-service/pkg/logger"

	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed database with initial data",
	Long:  `Seed database with initial payment methods and other required data.`,
	Run:   runSeed,
}

func init() {
	rootCmd.AddCommand(seedCmd)
}

func runSeed(cmd *cobra.Command, args []string) {
	// Load config
	configPath := getEnv("CONFIG_PATH", "pkg/config/config.yaml")
	loadedConfig, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	config.Set(loadedConfig)

	// Initialize logger
	if err := logger.Init(&loadedConfig.Log); err != nil {
		panic(err)
	}
	defer logger.Sync()

	// Connect to database
	db, err := gorm.Open(postgres.Open(getDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Starting database seeding...")

	// Seed payment methods
	seedPaymentMethods(db)

	log.Println("Database seeding completed successfully")
}

func seedPaymentMethods(db *gorm.DB) {
	log.Println("Seeding payment methods...")

	paymentMethods := []struct {
		Code        string
		Name        string
		Description string
		IsActive    bool
	}{
		{"credit_card", "Credit Card", "Pay with credit card (Visa, Mastercard, JCB)", true},
		{"debit_card", "Debit Card", "Pay with debit card", true},
		{"bank_transfer", "Bank Transfer", "Transfer via bank ATM or mobile banking", true},
		{"ewallet", "E-Wallet", "Pay with e-wallet (GoPay, OVO, Dana, ShopeePay)", true},
		{"qris", "QRIS", "Scan QR code to pay", true},
		{"cod", "Cash on Delivery", "Pay when you receive the service", true},
		{"virtual_account", "Virtual Account", "Pay via virtual account number", true},
	}

	for _, pm := range paymentMethods {
		// Check if payment method already exists
		var existing struct{}
		db.Raw("SELECT 1 FROM payment_methods WHERE code = ?", pm.Code).Scan(&existing)

		if existing == struct{}{} {
			log.Printf("Payment method '%s' already exists, skipping", pm.Code)
			continue
		}

		result := db.Exec(`
			INSERT INTO payment_methods (code, name, description, is_active, created_at, updated_at)
			VALUES (?, ?, ?, ?, NOW(), NOW())
		`, pm.Code, pm.Name, pm.Description, pm.IsActive)

		if result.Error != nil {
			log.Printf("Failed to seed payment method '%s': %v", pm.Code, result.Error)
		} else {
			log.Printf("Seeded payment method: %s", pm.Name)
		}
	}
}