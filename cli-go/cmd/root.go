package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// RootCmd is the root command for the CLI
var RootCmd = &cobra.Command{
	Use:   "etosctl",
	Short: "ETOS CLI - A command line interface for ETOS",
	Long: `ETOS CLI is a command line interface for ETOS (Eiffel Test Orchestration System).
It provides commands for managing and interacting with ETOS services.`,
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	RootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug output")
} 