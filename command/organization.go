package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

// OrganizationCommand is a cli.Command implementation for top level `sqsc organization` command.
type OrganizationCommand struct {
}

// Run is part of cli.Command implementation.
func (c *OrganizationCommand) Run(args []string) int {
	return cli.RunResultHelp
}

// Synopsis is part of cli.Command implementation.
func (c *OrganizationCommand) Synopsis() string {
	return "Commands related to a organization"
}

// Help is part of cli.Command implementation.
func (c *OrganizationCommand) Help() string {
	helpText := `
usage: sqsc organization <subcommand>

  Run a project volume related command.
  List of supported subcommands is available below.
`
	return strings.TrimSpace(helpText)
}
