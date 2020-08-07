package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

// LBCommand is a cli.Command implementation for top level `sqsc lb` command.
type LBCommand struct {
}

// Run is part of cli.Command implementation.
func (c *LBCommand) Run(args []string) int {
	return cli.RunResultHelp
}

// Synopsis is part of cli.Command implementation.
func (c *LBCommand) Synopsis() string {
	return "Commands to interact with load balancer in a project"
}

// Help is part of cli.Command implementation.
func (c *LBCommand) Help() string {
	helpText := `
usage: sqsc lb <subcommand>

  Run a project load balancer related command.
  List of supported subcommands is available below.
`
	return strings.TrimSpace(helpText)
}
