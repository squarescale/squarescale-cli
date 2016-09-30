package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
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
		return c.errorWithUsage(errors.New("Too many command line arguments"), c.Help())
	}

	var msg string
	err := c.runWithSpinner("create project", *endpoint, func(token string) error {

		pName, e := squarescale.CreateProject(*endpoint, token, projectName)
		if e != nil {
			return e
		}

		msg = fmt.Sprintf("Created project '%s'", pName)
		return nil
	})

	if err != nil {
		return c.error(err)
	}

	return c.info(msg)
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
