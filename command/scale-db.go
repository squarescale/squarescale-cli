package command

import (
	"flag"
	"strings"
)

// ScaleDBCommand is a cli.Command implementation for scaling the db of a project.
type ScaleDBCommand struct {
	Meta
}

// Run is part of cli.Command implementation.
func (c *ScaleDBCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("project scale-db", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	//endpoint := EndpointFlag(cmdFlags)
	project := ProjectFlag(cmdFlags)
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	err := validateArgs(*project)
	if err != nil {
		return c.errorWithUsage(err, c.Help())
	}

	return 0
}

// Synopsis is part of cli.Command implementation.
func (c *ScaleDBCommand) Synopsis() string {
	return "Scale up/down the database of a Squarescale project"
}

// Help is part of cli.Command implementation.
func (c *ScaleDBCommand) Help() string {
	helpText := `
usage: sqsc scale-db [options] <micro|small|medium>

  Scale the database of the specified Squarescale project.

Options:

  -endpoint="http://www.staging.sqsc.squarely.io" Squarescale endpoint
  -project=""                                     Squarescale project name
`
	return strings.TrimSpace(helpText)
}
