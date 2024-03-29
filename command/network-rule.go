package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

// NetworkRuleCommand is a cli.Command implementation for top level `sqsc servic` command.
type NetworkRuleCommand struct {
}

// Run is part of cli.Command implementation.
func (cmd *NetworkRuleCommand) Run(args []string) int {
	return cli.RunResultHelp
}

// Synopsis is part of cli.Command implementation.
func (cmd *NetworkRuleCommand) Synopsis() string {
	return "Commands to interact with network rules in a project"
}

// Help is part of cli.Command implementation.
func (cmd *NetworkRuleCommand) Help() string {
	helpText := `
usage: sqsc network-rule <subcommand>

  Run a network rule related command.
  List of supported subcommands is available below.
`
	return strings.TrimSpace(helpText)
}
