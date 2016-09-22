package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

// RepositoryCommand is a cli.Command implementation for top level `sqsc repository` commands.
type RepositoryCommand struct {
}

// Run is part of cli.Command implementation.
func (r *RepositoryCommand) Run(args []string) int {
	return cli.RunResultHelp
}

// Synopsis is part of cli.Command implementation.
func (r *RepositoryCommand) Synopsis() string {
	return "Commands related to the repository of a Squarescale project"
}

// Help is part of cli.Command implementation.
func (r *RepositoryCommand) Help() string {
	helpText := `
usage: sqsc repository <subcommand>

  Run a Squarescale project repository related command. List of supported
  subcommands is available below.
`
	return strings.TrimSpace(helpText)
}
