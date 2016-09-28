package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
	"github.com/squarescale/squarescale-cli/tokenstore"
)

// ProjectCreateCommand is a cli.Command implementation for creating a Squarescale project.
type ProjectCreateCommand struct {
	Meta
}

// Run is part of cli.Command implementation.
func (c *ProjectCreateCommand) Run(args []string) int {
	// Parse flags
	cmdFlags := flag.NewFlagSet("project create", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	endpoint := EndpointFlag(cmdFlags)
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	// Check for a project name
	var projectName string
	args = cmdFlags.Args()
	switch len(args) {
	case 0:
	case 1:
		projectName = args[0]
	default:
		c.Ui.Error("Too many command line arguments\n")
		c.Ui.Output(c.Help())
		return 1
	}

	token, err := tokenstore.GetToken(*endpoint)
	if err != nil {
		return c.error(err)
	}

	s := startSpinner("create project")
	projectName, err = squarescale.CreateProject(*endpoint, token, projectName)
	if err != nil {
		s.Stop()
		return c.error(err)
	}

	s.Stop()
	return c.info(fmt.Sprintf("Created project '%s'", projectName))
}

// Synopsis is part of cli.Command implementation.
func (c *ProjectCreateCommand) Synopsis() string {
	return "Create Project in Squarescale"
}

// Help is part of cli.Command implementation.
func (c *ProjectCreateCommand) Help() string {
	helpText := `
usage: sqsc project create [options] <project_name>

  Creates a new project using the provided project name. If no project name is
  provided, then a new name is automatically generated.

Options:

  -endpoint="http://www.staging.sqsc.squarely.io" Squarescale endpoint
`
	return strings.TrimSpace(helpText)
}
