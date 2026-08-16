package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aep-e2e-validator",
	Short: "An AEP validator that ensures compatibility of AEP HTTP APIs end-to-end.",
	Long:  `aep-e2e-validator covers the gap of validation of runtime functionality defined in aep.dev.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
