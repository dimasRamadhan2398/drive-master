package cli

import (
	"fmt"
	"log"
	"os"

	"user-service/pkg/config"
)

// getEnv gets an environment variable or returns a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// setEnvIfNotSet sets environment variable only if not already set
func setEnvIfNotSet(key, value string) {
	if os.Getenv(key) == "" {
		os.Setenv(key, value)
	}
}

// joinStrings joins a string slice with a separator
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

// LoadConfig loads the configuration from file or environment
func LoadConfig() {
	configPath := getEnv("CONFIG_PATH", "config.yaml")
	if cfg, err := config.Load(configPath); err != nil {
		// Config file not required, env vars will be used
		log.Printf("Config file not found at %s, using environment variables", configPath)
	} else {
		config.Set(cfg)
	}
}

// getDSN returns the database connection string from config or environment
func getDSN() string {
	cfg := config.Get()
	if cfg == nil {
		return fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			getEnv("POSTGRES_HOST", "localhost"),
			getEnv("POSTGRES_USER", "postgres"),
			getEnv("POSTGRES_PASSWORD", "postgres"),
			getEnv("POSTGRES_DB", "postgres"),
			getEnv("POSTGRES_PORT", "5432"),
			getEnv("POSTGRES_SSLMODE", "disable"),
		)
	}

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		getEnv("POSTGRES_HOST", cfg.Database.Host),
		getEnv("POSTGRES_USER", cfg.Database.Username),
		getEnv("POSTGRES_PASSWORD", cfg.Database.Password),
		getEnv("POSTGRES_DB", cfg.Database.Name),
		cfg.Database.Port,
		getEnv("POSTGRES_SSLMODE", cfg.Database.SSLMode),
	)
}