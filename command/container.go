package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

// ContainerCommand is a cli.Command implementation for top level `sqsc container` command.
type ContainerCommand struct {
}

// Run is part of cli.Command implementation.
func (c *ContainerCommand) Run(args []string) int {
	return cli.RunResultHelp
}

// Synopsis is part of cli.Command implementation.
func (c *ContainerCommand) Synopsis() string {
	return "Commands to interact with containers in a project"
}

// Help is part of cli.Command implementation.
func (c *ContainerCommand) Help() string {
	helpText := `
usage: sqsc container <subcommand>

  Run a project container related command.
  List of supported subcommands is available below.
`
	return strings.TrimSpace(helpText)
}
