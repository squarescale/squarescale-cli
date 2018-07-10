package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ImageAddCommand is a cli.Command implementation for adding a docker image to a Squarescale project.
type ImageAddCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *ImageAddCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	project := projectFlag(c.flagSet)
	image := c.flagSet.String("name", "", "Docker image name")
	instances := containerInstancesFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if err := c.validateArgs(*project, *image, &instances); err != nil {
		return c.errorWithUsage(err)
	}

	return c.runWithSpinner("add docker image", endpoint.String(), func(client *squarescale.Client) (string, error) {
		msg := fmt.Sprintf("Successfully added docker image '%s' to project '%s'", *image, *project)
		return msg, client.AddImage(*project, *image, instances)
	})
}

// Synopsis is part of cli.Command implementation.
func (c *ImageAddCommand) Synopsis() string {
	return "Add a docker image to a Squarescale project"
}

// Help is part of cli.Command implementation.
func (c *ImageAddCommand) Help() string {
	helpText := `
usage: sqsc image add [options]

  Adds docker image to the specified Squarescale project.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}

func (c *ImageAddCommand) validateArgs(project, image string, instancesP **int) error {
	if image == "" {
		return errors.New("Docker image name cannot be empty")
	}

	if **instancesP <= 0 {
		if **instancesP == -1 { // default value, ignored
			c.Ui.Warn("Number of instances cannot be 0 or negative. Ignored and using default value")
			*instancesP = nil
		} else {
			return errors.New("Number of instances cannot be 0 or negative")
		}
	}

	return validateProjectName(project)
}
