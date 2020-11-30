package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

// ImageCommand is a cli.Command implementation for top level `sqsc image` command.
type ImageCommand struct {
}

// Run is part of cli.Command implementation.
func (c *ImageCommand) Run(args []string) int {
	return cli.RunResultHelp
}

// Synopsis is part of cli.Command implementation.
func (c *ImageCommand) Synopsis() string {
	return "Commands to interact with docker images in a project"
}

// Help is part of cli.Command implementation.
func (c *ImageCommand) Help() string {
	helpText := `
usage: sqsc image <subcommand>

  Run a project docker image related command.

  List of supported subcommands is available below.

`
	return strings.TrimSpace(helpText)
}
