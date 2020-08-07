package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

// DBCommand is a cli.Command implementation for top level `sqsc db` command.
type DBCommand struct {
}

// Run is part of cli.Command implementation.
func (c *DBCommand) Run(args []string) int {
	return cli.RunResultHelp
}

// Synopsis is part of cli.Command implementation.
func (c *DBCommand) Synopsis() string {
	return "Commands to interact with databases in a project"
}

// Help is part of cli.Command implementation.
func (c *DBCommand) Help() string {
	helpText := `
usage: sqsc db <subcommand>

  Run a project database related command.
  List of supported subcommands is available below.
`
	return strings.TrimSpace(helpText)
}
