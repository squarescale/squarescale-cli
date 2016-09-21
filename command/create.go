package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
	"github.com/squarescale/squarescale-cli/tokenstore"
)

// CreateCommand is a cli.Command implementation for creating a Squarescale project.
type CreateCommand struct {
	Meta
}

// Run is part of cli.Command implementation.
func (c *CreateCommand) Run(args []string) int {
	// Parse flags
	cmdFlags := flag.NewFlagSet("create", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	endpoint := EndpointFlag(cmdFlags)
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	// Check for a payload
	var projectName string
	args = cmdFlags.Args()
	switch len(args) {
	case 0:
	case 1:
		projectName = args[0]
	default:
		c.Ui.Error("Too many command line arguments.\n")
		c.Ui.Output(c.Help())
		return 1
	}

	token, err := tokenstore.GetToken(*endpoint)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error: %v\n", err))
		return 1
	}

	projectName, err = squarescale.CreateProject(*endpoint, token, projectName)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error: %v\n", err))
		return 1
	}

	c.Ui.Info(fmt.Sprintf("Created project %s\n", projectName))
	return 0
}

// Synopsis is part of cli.Command implementation.
func (c *CreateCommand) Synopsis() string {
	return "Create Project in Squarescale"
}

// Help is part of cli.Command implementation.
func (c *CreateCommand) Help() string {
	helpText := `
Usage: sqsc create [options] [project_name]

	Creates a new project using the provided project name. If no project name is
	provided, then a new name is automatically generated.

Options:

	-endpoint="http://www.staging.sqsc.squarely.io" Squarescale endpoint
`
	return strings.TrimSpace(helpText)
}
