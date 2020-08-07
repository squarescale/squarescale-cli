package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ContainerAddCommand is a cli.Command implementation for adding a docker container to project.
type ContainerAddCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *ContainerAddCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := c.flagSet.String("project-uuid", "", "set the uuid of the project")
	serviceName := c.flagSet.String("servicename", "", "service name")
	image := c.flagSet.String("name", "", "Docker image name")
	runCommand := c.flagSet.String("run-command", "", "command / arguments that are used for execution")
	entrypoint := c.flagSet.String("entrypoint", "", "This is the script / program that will be executed")
	username := c.flagSet.String("username", "", "Username")
	password := c.flagSet.String("password", "", "Password")
	volumes := c.flagSet.String("volumes", "", "Volumes. format: ${VOL2_NAME}:${VOL2_MOUNT_POINT}:ro,${VOL2_NAME}:${VOL2_MOUNT_POINT}:rw")
	instances := repoOrImageInstancesFlag(c.flagSet)

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if *projectUUID == "" {
		return c.errorWithUsage(errors.New("Project uuid is mandatory"))
	}

	if err := c.validateArgs(*image, *instances); err != nil {
		return c.errorWithUsage(err)
	}

	volumesToBind := parseVolumesToBind(*volumes)

	return c.runWithSpinner("add docker image", endpoint.String(), func(client *squarescale.Client) (string, error) {
		msg := fmt.Sprintf("Successfully added docker image '%s' to project '%s' (%v instance(s))", *image, *projectUUID, *instances)
		return msg, client.AddImage(*projectUUID, *image, *username, *password, *entrypoint, *runCommand, *instances, *serviceName, volumesToBind)
	})
}

// Synopsis is part of cli.Command implementation.
func (c *ContainerAddCommand) Synopsis() string {
	return "Add docker container to project"
}

// Help is part of cli.Command implementation.
func (c *ContainerAddCommand) Help() string {
	helpText := `
usage: sqsc container add [options]

  Add docker container to project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}

func (c *ContainerAddCommand) validateArgs(image string, instances int) error {
	if image == "" {
		return errors.New("Docker image name cannot be empty")
	}

	if instances <= 0 {
		return errors.New("Invalid value provided for instances count: it cannot be 0 or negative")
	}
	return nil
}
