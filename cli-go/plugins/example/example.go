package main

import (
	"github.com/fredjn/etos/cli-go/pkg/command"
)

// Registry is the plugin's command registry
type Registry struct{}

// RegisterCommands returns the commands provided by this plugin
func (r *Registry) RegisterCommands() map[string]*command.Command {
	commands := make(map[string]*command.Command)

	// Create an example command
	exampleCmd := command.NewCommand(
		"example",
		"Example command",
		"An example command that demonstrates plugin functionality",
	)

	// Add the command to the registry
	commands["example"] = exampleCmd

	return commands
}

// Export the Registry symbol
var Registry = &Registry{} 