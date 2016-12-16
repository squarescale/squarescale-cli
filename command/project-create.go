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
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *ProjectCreateCommand) Run(args []string) int {
	// Parse flags
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	// Check for a project name
	var wantedProjectName string
	args = c.flagSet.Args()
	switch len(args) {
	case 0:
	case 1:
		wantedProjectName = args[0]
	default:
		return c.errorWithUsage(errors.New("Too many command line arguments"))
	}

	var definitiveName string
	err := c.runWithSpinner("create project", *endpoint, func(client *squarescale.Client) error {
		if wantedProjectName != "" {
			valid, same, fmtName, err := client.CheckProjectName(wantedProjectName)
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
			generatedName, err := client.FindProjectName()
			if err != nil {
				return err
			}

			definitiveName = generatedName
		}

		return client.CreateProject(definitiveName)
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

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
