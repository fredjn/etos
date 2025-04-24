# ETOS CLI (Go Implementation)

This is a Go implementation of the ETOS CLI, providing a command-line interface for interacting with the Eiffel Test Orchestration System.

## Building

To build the CLI:

```bash
go build -o etosctl
```

## Development

The CLI is built using the [Cobra](https://github.com/spf13/cobra) library, which provides a robust framework for building command-line applications in Go.

### Project Structure

- `cmd/` - Contains the root command and command definitions
- `pkg/` - Contains shared packages and utilities
  - `command/` - Base command structure and utilities
  - `engine/` - Customization engine for loading commands
  - `models/` - Data models and types

### Adding New Commands

To add a new command:

1. Create a new file in the `cmd` directory
2. Define your command using the `command.Command` base type
3. Register your command in the root command's initialization

Example:

```go
package cmd

import (
	"github.com/eiffel-community/etos/cli-go/pkg/command"
)

var myCommand = command.NewCommand(
	"mycommand",
	"Short description",
	"Long description",
)

func init() {
	rootCmd.AddCommand(myCommand)
}
```

## License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details. 