package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

// StatefullNodeCommand is a cli.Command implementation for top level `sqsc statefull_node` command.
type StatefullNodeCommand struct {
}

// Run is part of cli.Command implementation.
func (c *StatefullNodeCommand) Run(args []string) int {
	return cli.RunResultHelp
}

// Synopsis is part of cli.Command implementation.
func (c *StatefullNodeCommand) Synopsis() string {
	return "Commands related to a statefull_node"
}

// Help is part of cli.Command implementation.
func (c *StatefullNodeCommand) Help() string {
	helpText := `
usage: sqsc statefull_node <subcommand>

  Run a project statefull_node related command.
  List of supported subcommands is available below.
`
	return strings.TrimSpace(helpText)
}
