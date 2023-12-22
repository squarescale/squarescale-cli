package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

// ClusterCommand is a cli.Command implementation for top level `sqsc cluster` command.
type ClusterCommand struct {
}

// Run is part of cli.Command implementation.
func (cmd *ClusterCommand) Run(args []string) int {
	return cli.RunResultHelp
}

// Synopsis is part of cli.Command implementation.
func (cmd *ClusterCommand) Synopsis() string {
	return "Commands to interact with cluster instances in a project"
}

// Help is part of cli.Command implementation.
func (cmd *ClusterCommand) Help() string {
	helpText := `
usage: sqsc cluster <subcommand>

  Run a project cluster related command.

  List of supported subcommands is available below.

`
	return strings.TrimSpace(helpText)
}
