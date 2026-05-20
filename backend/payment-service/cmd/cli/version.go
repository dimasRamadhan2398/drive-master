package cli

import (
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of payment-service",
	Long:  `All software has versions. This is payment-service's version information.`,
	Run: func(cmd *cobra.Command, args []string) {
		println(GetVersion())
	},
}