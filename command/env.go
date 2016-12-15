package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

// EnvCommand is a cli.Command implementation for top level `sqsc db` command.
type EnvCommand struct {
}

// Run is part of cli.Command implementation.
func (c *EnvCommand) Run(args []string) int {
	return cli.RunResultHelp
}

// Synopsis is part of cli.Command implementation.
func (c *EnvCommand) Synopsis() string {
	return "Commands to set or list the project environment variables"
}

// Help is part of cli.Command implementation.
func (c *EnvCommand) Help() string {
	helpText := `
usage: sqsc env <subcommand>

  Run a Squarescale environment variables related command.
  List of supported subcommands is available below.
`
	return strings.TrimSpace(helpText)
}
