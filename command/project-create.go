package command

import (
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
	wantedProjectName := projectNameFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	return c.runWithSpinner("create project", *endpoint, func(client *squarescale.Client) (string, error) {
		var definitiveName string

		if *wantedProjectName != "" {
			valid, same, fmtName, err := client.CheckProjectName(*wantedProjectName)
			if err != nil {
				return "", fmt.Errorf("Cannot validate project name '%s'", *wantedProjectName)
			}

			if !valid {
				return "", fmt.Errorf(
					"Project name '%s' is invalid (already taken or not well formed), please choose another one",
					*wantedProjectName)
			}

			if !same {
				if err := c.askConfirmName(fmtName); err != nil {
					return "", err
				}
			}

			definitiveName = fmtName

		} else {
			generatedName, err := client.FindProjectName()
			if err != nil {
				return "", err
			}

			if err := c.askConfirmName(generatedName); err != nil {
				return "", err
			}

			definitiveName = generatedName
		}

		return fmt.Sprintf("Created project '%s'", definitiveName), client.CreateProject(definitiveName)
	})
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

func (c *ProjectCreateCommand) askConfirmName(name string) error {
	c.pauseSpinner()
	c.Ui.Warn(fmt.Sprintf("Project will be created as '%s', is this ok?", name))
	_, err := c.Ui.Ask("Enter to accept, Ctrl-c to cancel:")
	if err != nil {
		return err
	}

	c.startSpinner()
	return nil
}
