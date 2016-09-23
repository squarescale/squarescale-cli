package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

// ProjectCommand is a cli.Command implementation for top level `sqsc project` command.
type ProjectCommand struct {
}

// Run is part of cli.Command implementation.
func (c *ProjectCommand) Run(args []string) int {
	return cli.RunResultHelp
}

// Synopsis is part of cli.Command implementation.
func (c *ProjectCommand) Synopsis() string {
	return "Commands related to a Squarescale project"
}

// Help is part of cli.Command implementation.
func (c *ProjectCommand) Help() string {
	helpText := `
usage: sqsc project <subcommand>

  Run a Squarescale project related command. List of supported subcommands
  is available below.
`
	return strings.TrimSpace(helpText)
}
