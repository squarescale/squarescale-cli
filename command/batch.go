package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

// BatchCommand is a cli.Command implementation for top level `sqsc batch` command.
type BatchCommand struct {
}

// Run is part of cli.Command implementation.
func (c *BatchCommand) Run(args []string) int {
	return cli.RunResultHelp
}

// Synopsis is part of cli.Command implementation.
func (c *BatchCommand) Synopsis() string {
	return "Commands to interact with batch jobs in a project"
}

// Help is part of cli.Command implementation.
func (c *BatchCommand) Help() string {
	helpText := `
usage: sqsc batch <subcommand>

  Run a project batch related command.

  List of supported subcommands is available below.

`
	return strings.TrimSpace(helpText)
}
