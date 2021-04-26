package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

// SchedulingGroupCommand is a cli.Command implementation for top level `sqsc scheduling-group` command.
type SchedulingGroupCommand struct {
}

// Run is part of cli.Command implementation.
func (c *SchedulingGroupCommand) Run(args []string) int {
	return cli.RunResultHelp
}

// Synopsis is part of cli.Command implementation.
func (c *SchedulingGroupCommand) Synopsis() string {
	return "Commands related to a scheduling groups"
}

// Help is part of cli.Command implementation.
func (c *SchedulingGroupCommand) Help() string {
	helpText := `
usage: sqsc scheduling-group <subcommand>

  Run a project scheduling group related command.

  List of supported subcommands is available below.

`
	return strings.TrimSpace(helpText)
}
