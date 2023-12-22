package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

// ExternalNodeCommand is a cli.Command implementation for top level `sqsc external-node` command.
type ExternalNodeCommand struct {
}

// Run is part of cli.Command implementation.
func (cmd *ExternalNodeCommand) Run(args []string) int {
	return cli.RunResultHelp
}

// Synopsis is part of cli.Command implementation.
func (cmd *ExternalNodeCommand) Synopsis() string {
	return "Commands related to a external nodes"
}

// Help is part of cli.Command implementation.
func (cmd *ExternalNodeCommand) Help() string {
	helpText := `
usage: sqsc external-node <subcommand>

  Run a project external node related command.

  List of supported subcommands is available below.

`
	return strings.TrimSpace(helpText)
}
