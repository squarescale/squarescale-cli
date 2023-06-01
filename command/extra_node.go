package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

// ExtraNodeCommand is a cli.Command implementation for top level `sqsc extra node` command.
type ExtraNodeCommand struct {
}

// Run is part of cli.Command implementation.
func (c *ExtraNodeCommand) Run(args []string) int {
	return cli.RunResultHelp
}

// Synopsis is part of cli.Command implementation.
func (c *ExtraNodeCommand) Synopsis() string {
	return "Commands related to a extra-node"
}

// Help is part of cli.Command implementation.
func (c *ExtraNodeCommand) Help() string {
	helpText := `
usage: sqsc extra node <subcommand>

  Run a project extra node related command.

  List of supported subcommands is available below.

`
	return strings.TrimSpace(helpText)
}
