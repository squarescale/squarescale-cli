package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ProjectRemoveCommand is a cli.Command implementation for creating a Squarescale project.
type ProjectRemoveCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *ProjectRemoveCommand) Run(args []string) int {
	// Parse flags
	c.flagSet = newFlagSet(c, c.Ui)
	alwaysYes := yesFlag(c.flagSet)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := c.flagSet.String("project-uuid", "", "set the uuid of the project")
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" {
		return c.errorWithUsage(errors.New("Project uuid is mandatory"))
	}

	if c.flagSet.NArg() > 1 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()[1:]))
	}

	c.Ui.Info("Destroy infrastructure and configuration for project " + *projectUUID + "?")
	ok, err := AskYesNo(c.Ui, alwaysYes, "Proceed ?", true)
	if err != nil {
		return c.error(err)
	} else if !ok {
		return c.error(CancelledError)
	}

	res := c.runWithSpinner("unprovision project", endpoint.String(), func(client *squarescale.Client) (string, error) {
		err := client.ProjectUnprovision(*projectUUID)
		if err != nil {
			return "", err
		}

		var infraStatus string = "provisionning"

	loop:
		for {
			status, err := client.ProjectStatus(*projectUUID)
			if err != nil {
				return "", err
			}
			if status.InfraStatus == "error" {
				return "", fmt.Errorf("Unknown infrastructure error.")
			} else if status.InfraStatus == "no_infra" {
				break loop
			} else if infraStatus != status.InfraStatus {
				infraStatus = status.InfraStatus
				c.Ui.Info("Infrastructure status: " + infraStatus)
			}

			time.Sleep(time.Second)
		}

		return "", nil
	})
	if res != 0 {
		return res
	}

	return c.runWithSpinner("delete project", endpoint.String(), func(client *squarescale.Client) (string, error) {
		err := client.ProjectDelete(*projectUUID)
		if err != nil {
			return "", err
		}
		return "", nil
	})
}

// Synopsis is part of cli.Command implementation.
func (c *ProjectRemoveCommand) Synopsis() string {
	return "Remove Project in Squarescale"
}

// Help is part of cli.Command implementation.
func (c *ProjectRemoveCommand) Help() string {
	helpText := `
usage: sqsc project remove [options] <project_name>

  Destroy infrastructure and remove project completely.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
