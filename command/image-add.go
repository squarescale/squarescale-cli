package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ImageAddCommand is a cli.Command implementation for adding a Docker image to project.
type ImageAddCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *ImageAddCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := c.flagSet.String("project-uuid", "", "set the uuid of the project")
	projectName := c.flagSet.String("project-name", "", "set the name of the project")
	runCommand := c.flagSet.String("run-command", "", "command / arguments that are used for execution")
	entrypoint := c.flagSet.String("entrypoint", "", "This is the script / program that will be executed")
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

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

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
		return msg, client.AddImage(UUID, *image, *username, *password, *entrypoint, *runCommand, *instances, *serviceName, volumesToBind)
	})
}

// Synopsis is part of cli.Command implementation.
func (c *ImageAddCommand) Synopsis() string {
	return "Add Docker image to project"
}

// Help is part of cli.Command implementation.
func (c *ImageAddCommand) Help() string {
	helpText := `
usage: sqsc image add [options]

  Add Docker image to project.
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
