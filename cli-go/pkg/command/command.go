package command

import (
	"github.com/spf13/cobra"
)

// Command represents a base command that other commands can inherit from
type Command struct {
	*cobra.Command
	parent *Command
}

// NewCommand creates a new command with the given name and description
func NewCommand(name, short, long string) *Command {
	cmd := &Command{
		Command: &cobra.Command{
			Use:   name,
			Short: short,
			Long:  long,
		},
	}
	return cmd
}

// AddCommand adds a subcommand to this command
func (c *Command) AddCommand(subcommand *Command) {
	subcommand.parent = c
	c.Command.AddCommand(subcommand.Command)
}

// GetParent returns the parent command
func (c *Command) GetParent() *Command {
	return c.parent
} 