package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ServiceAddCommand is a cli.Command implementation for adding a service aka Docker container to project.
type ServiceAddCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *ServiceAddCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := c.flagSet.String("project-uuid", "", "set the uuid of the project")
	projectName := c.flagSet.String("project-name", "", "set the name of the project")
	serviceName := c.flagSet.String("servicename", "", "service name")
	image := c.flagSet.String("image", "", "Docker image name")
	runCommand := c.flagSet.String("run-command", "", "command / arguments that are used for execution")
	entrypoint := c.flagSet.String("entrypoint", "", "This is the script / program that will be executed")
	username := c.flagSet.String("username", "", "Username")
	password := c.flagSet.String("password", "", "Password")
	volumes := c.flagSet.String("volumes", "", "Volumes. format: ${VOL2_NAME}:${VOL2_MOUNT_POINT}:ro,${VOL2_NAME}:${VOL2_MOUNT_POINT}:rw")
	dockerCapabilities := dockerCapabilitiesFlag(c.flagSet)
	instances := repoOrImageInstancesFlag(c.flagSet)

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if err := c.validateArgs(*image, *instances); err != nil {
		return c.errorWithUsage(err)
	}

	dockerCapabilitiesArray := getDockerCapabilitiesArray(*dockerCapabilities)

	volumesToBind := parseVolumesToBind(*volumes)

	return c.runWithSpinner("adding Docker image", endpoint.String(), func(client *squarescale.Client) (string, error) {
		var UUID string
		var err error
		var projectToShow string
		if *projectUUID == "" {
			projectToShow = *projectName
			UUID, err = client.ProjectByName(*projectName)
			if err != nil {
				return "", err
			}
		} else {
			projectToShow = *projectUUID
			UUID = *projectUUID
		}

		msg := fmt.Sprintf("Successfully added Docker image '%s' to project '%s' (%v instance(s))", *image, projectToShow, *instances)
		return msg, client.AddImage(UUID, *image, *username, *password, *entrypoint, *runCommand, *instances, *serviceName, volumesToBind, dockerCapabilitiesArray)
	})
}

// Synopsis is part of cli.Command implementation.
func (c *ServiceAddCommand) Synopsis() string {
	return "Add service aka Docker container to project"
}

// Help is part of cli.Command implementation.
func (c *ServiceAddCommand) Help() string {
	helpText := `
usage: sqsc service add [options]

  Add service aka Docker container to project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}

func (c *ServiceAddCommand) validateArgs(image string, instances int) error {
	if image == "" {
		return errors.New("Docker image name cannot be empty")
	}

	if instances <= 0 {
		return errors.New("Invalid value provided for instances count: it cannot be 0 or negative")
	}
	return nil
}
