package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `All software has versions. This is booking-service's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(GetVersion())
	},
}

func init() {
	// This command is added in root.go
}
