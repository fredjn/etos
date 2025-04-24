package main

import (
	"os"

	"github.com/fredjn/etos/cli-go/cmd"
	"github.com/fredjn/etos/cli-go/pkg/config"
	"github.com/fredjn/etos/cli-go/pkg/engine"
	"github.com/fredjn/etos/cli-go/pkg/logging"
)

func main() {
	// Initialize configuration
	cfg := config.NewConfig()
	if err := cfg.Load(); err != nil {
		logging.DefaultLogger.Error("Failed to load configuration: %v", err)
		os.Exit(1)
	}

	// Set up logging level
	switch cfg.GetString("log.level") {
	case "debug":
		logging.DefaultLogger.SetLevel(logging.LogLevelDebug)
	case "info":
		logging.DefaultLogger.SetLevel(logging.LogLevelInfo)
	case "warning":
		logging.DefaultLogger.SetLevel(logging.LogLevelWarning)
	case "error":
		logging.DefaultLogger.SetLevel(logging.LogLevelError)
	default:
		logging.DefaultLogger.SetLevel(logging.LogLevelInfo)
	}

	// Initialize the customization engine
	customEngine := engine.NewCustomizationEngine()
	customEngine.Start()

	// Get the loaded commands
	commands := customEngine.GetCommands()

	// Add the commands to the root command
	for name, command := range commands {
		cmd.RootCmd.AddCommand(command.Command)
		logging.DefaultLogger.Info("Added command: %s", name)
	}

	// Execute the root command
	cmd.Execute()
} 