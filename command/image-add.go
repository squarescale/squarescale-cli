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
	serviceName := c.flagSet.String("servicename", "", "service name")
	image := c.flagSet.String("name", "", "Docker image name")
	username := c.flagSet.String("username", "", "Username")
	password := c.flagSet.String("password", "", "Password")
	volumes := c.flagSet.String("volumes", "", "Volumes")
	instances := repoOrImageInstancesFlag(c.flagSet)

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if err := c.validateArgs(*project, *image, *instances); err != nil {
		return c.errorWithUsage(err)
	}

	volumesToBind := parseVolumesToBind(*volumes)

	return c.runWithSpinner("add docker image", endpoint.String(), func(client *squarescale.Client) (string, error) {
		msg := fmt.Sprintf("Successfully added docker image '%s' to project '%s' (%v instance(s))", *image, *project, *instances)
		return msg, client.AddImage(*project, *image, *username, *password, *instances, *serviceName, volumesToBind)
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

func (c *ImageAddCommand) validateArgs(project, image string, instances int) error {
	if image == "" {
		return errors.New("Docker image name cannot be empty")
	}

	if instances <= 0 {
		return errors.New("Invalid value provided for instances count: it cannot be 0 or negative")
	}

	return validateProjectName(project)
}
