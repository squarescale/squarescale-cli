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
	var wantedProjectName string
	args = cmdFlags.Args()
	switch len(args) {
	case 0:
	case 1:
		wantedProjectName = args[0]
	default:
		return c.errorWithUsage(errors.New("Too many command line arguments"), c.Help())
	}

	var definitiveName string
	err := c.runWithSpinner("create project", *endpoint, func(token string) error {
		if wantedProjectName != "" {
			valid, same, fmtName, err := squarescale.CheckProjectName(*endpoint, token, wantedProjectName)
			if err != nil {
				return fmt.Errorf("Cannot validate project name '%s'", wantedProjectName)
			}

			if !valid {
				return fmt.Errorf(
					"Project name '%s' is invalid (already taken or not well formed), please choose another one",
					wantedProjectName)
			}

			if !same {
				c.pauseSpinner()
				c.Ui.Warn(fmt.Sprintf("Project will be created as '%s', is this ok?", fmtName))
				_, err := c.Ui.Ask("Enter to accept, Ctrl-c to cancel:")
				if err != nil {
					return err
				}

				c.startSpinner()
			}

			definitiveName = fmtName

		} else {
			generatedName, err := squarescale.FindProjectName(*endpoint, token)
			if err != nil {
				return err
			}

			definitiveName = generatedName
		}

		return squarescale.CreateProject(*endpoint, token, definitiveName)
	})

	if err != nil {
		return c.error(err)
	}

	return c.info(fmt.Sprintf("Created project '%s'", definitiveName))
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
