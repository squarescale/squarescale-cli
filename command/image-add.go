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
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if err := c.validateArgs(*project, *image); err != nil {
		return c.errorWithUsage(err)
	}

	return c.runWithSpinner("add docker image", *endpoint, func(client *squarescale.Client) (string, error) {
		msg := fmt.Sprintf("Successfully added docker image '%s' to project '%s'", *image, *project)
		return msg, client.AddImage(*project, *image)
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

func (c *ImageAddCommand) validateArgs(project, image string) error {
	if image == "" {
		return errors.New("Docker image name cannot be empty")
	}

	return validateProjectName(project)
}
