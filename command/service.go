package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

// ServiceCommand is a cli.Command implementation for top level `sqsc service` command.
type ServiceCommand struct {
}

// Run is part of cli.Command implementation.
func (c *ServiceCommand) Run(args []string) int {
	return cli.RunResultHelp
}

// Synopsis is part of cli.Command implementation.
func (c *ServiceCommand) Synopsis() string {
	return "Commands to interact with services in a project"
}

// Help is part of cli.Command implementation.
func (c *ServiceCommand) Help() string {
	helpText := `
usage: sqsc service <subcommand>

  Run a project service related command.

  List of supported subcommands is available below.

`
	return strings.TrimSpace(helpText)
}
