package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

// StatefulNodeCommand is a cli.Command implementation for top level `sqsc stateful-node` command.
type StatefulNodeCommand struct {
}

// Run is part of cli.Command implementation.
func (c *StatefulNodeCommand) Run(args []string) int {
	return cli.RunResultHelp
}

// Synopsis is part of cli.Command implementation.
func (c *StatefulNodeCommand) Synopsis() string {
	return "Commands related to a stateful-node"
}

// Help is part of cli.Command implementation.
func (c *StatefulNodeCommand) Help() string {
	helpText := `
usage: sqsc stateful-node <subcommand>

  Run a project stateful node related command.

  List of supported subcommands is available below.

`
	return strings.TrimSpace(helpText)
}
