package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

// VolumeCommand is a cli.Command implementation for top level `sqsc volume` command.
type VolumeCommand struct {
}

// Run is part of cli.Command implementation.
func (c *VolumeCommand) Run(args []string) int {
	return cli.RunResultHelp
}

// Synopsis is part of cli.Command implementation.
func (c *VolumeCommand) Synopsis() string {
	return "Commands related to a volume"
}

// Help is part of cli.Command implementation.
func (c *VolumeCommand) Help() string {
	helpText := `
usage: sqsc volume <subcommand>

  Run a project volume related command.
  List of supported subcommands is available below.
`
	return strings.TrimSpace(helpText)
}
