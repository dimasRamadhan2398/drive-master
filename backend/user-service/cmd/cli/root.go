package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Variables set during build
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "user-service",
	Short: "User Service CLI",
	Long:  `User Service - A microservice for managing users, authentication, members, and instructors.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Global pre-run logic can go here
	},
	SilenceUsage: true,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add version command
	rootCmd.AddCommand(versionCmd)

	// Persistent flags for all commands
	rootCmd.PersistentFlags().Bool("verbose", false, "Enable verbose output")
	rootCmd.PersistentFlags().String("config", "", "Config file path (default is ./config.yaml)")
}

// GetVersion returns version info
func GetVersion() string {
	return fmt.Sprintf("Version: %s\nBuild Time: %s\nGit Commit: %s", Version, BuildTime, GitCommit)
}

// GetRootCmd returns the root command for testing
func GetRootCmd() *cobra.Command {
	return rootCmd
}