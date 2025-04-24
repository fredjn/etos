package engine

import (
	"os"
	"path/filepath"
	"plugin"

	"github.com/fredjn/etos/cli-go/pkg/command"
	"github.com/fredjn/etos/cli-go/pkg/logging"
)

// CommandRegistry is the interface that plugins must implement to register commands
type CommandRegistry interface {
	RegisterCommands() map[string]*command.Command
}

// CustomizationEngine handles loading and managing custom commands
type CustomizationEngine struct {
	commands map[string]*command.Command
}

// NewCustomizationEngine creates a new customization engine
func NewCustomizationEngine() *CustomizationEngine {
	return &CustomizationEngine{
		commands: make(map[string]*command.Command),
	}
}

// Start initializes the customization engine and loads commands
func (e *CustomizationEngine) Start() {
	e.discover()
}

// GetCommands returns the map of loaded commands
func (e *CustomizationEngine) GetCommands() map[string]*command.Command {
	return e.commands
}

// discover loads plugins and custom commands
func (e *CustomizationEngine) discover() {
	// Look for plugins in the plugins directory
	pluginDir := "plugins"
	if _, err := os.Stat(pluginDir); os.IsNotExist(err) {
		// Create plugins directory if it doesn't exist
		if err := os.Mkdir(pluginDir, 0755); err != nil {
			logging.DefaultLogger.Error("Failed to create plugins directory: %v", err)
			return
		}
	}

	// Walk through the plugins directory
	err := filepath.Walk(pluginDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-plugin files
		if info.IsDir() || filepath.Ext(path) != ".so" {
			return nil
		}

		logging.DefaultLogger.Debug("Loading plugin: %s", path)

		// Load the plugin
		p, err := plugin.Open(path)
		if err != nil {
			logging.DefaultLogger.Error("Failed to load plugin %s: %v", path, err)
			return nil
		}

		// Look for the registry symbol
		sym, err := p.Lookup("Registry")
		if err != nil {
			logging.DefaultLogger.Error("Plugin %s does not export Registry symbol: %v", path, err)
			return nil
		}

		// Cast to CommandRegistry
		registry, ok := sym.(CommandRegistry)
		if !ok {
			logging.DefaultLogger.Error("Plugin %s Registry does not implement CommandRegistry interface", path)
			return nil
		}

		// Register commands from the plugin
		commands := registry.RegisterCommands()
		for name, cmd := range commands {
			if _, exists := e.commands[name]; exists {
				logging.DefaultLogger.Warning("Command %s already exists, skipping", name)
				continue
			}
			e.commands[name] = cmd
			logging.DefaultLogger.Debug("Registered command %s from plugin %s", name, path)
		}

		return nil
	})

	if err != nil {
		logging.DefaultLogger.Error("Error walking plugins directory: %v", err)
	}
} 