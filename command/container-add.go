package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/squarescale/squarescale-cli/squarescale"
)

// ContainerAddCommand is a cli.Command implementation for adding a docker container to a Squarescale project.
type ContainerAddCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *ContainerAddCommand) Run(args []string) int {
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

	volumesSplited := strings.Split(*volumes, ",")
	volumesLength := len(volumesSplited)
	if *volumes == "" {
		volumesLength = 0
	}
	volumesToBind := make([]squarescale.VolumeToBind, volumesLength)
	if volumesLength >= 1 {
		for index, volume := range volumesSplited {
			volumeSplited := strings.Split(volume, ":")
			volumeName := volumeSplited[0]
			mountPoint := volumeSplited[1]
			var readOnly bool
			if len(volumeSplited) > 2 {
				switch strings.ToLower(volumeSplited[2]) {
				case "ro":
					readOnly = true
				case "rw":
					readOnly = false
				default:
					log.Fatal("Read only parameter must be RO or RW")
				}
			} else {
				readOnly = false
			}
			volumesToBind[index] = squarescale.VolumeToBind{volumeName, mountPoint, readOnly}
		}
	}

	return c.runWithSpinner("add docker image", endpoint.String(), func(client *squarescale.Client) (string, error) {
		msg := fmt.Sprintf("Successfully added docker image '%s' to project '%s' (%v instance(s))", *image, *project, *instances)
		return msg, client.AddImage(*project, *image, *username, *password, *instances, *serviceName, volumesToBind)
	})
}

// Synopsis is part of cli.Command implementation.
func (c *ContainerAddCommand) Synopsis() string {
	return "Add a docker container to a Squarescale project"
}

// Help is part of cli.Command implementation.
func (c *ContainerAddCommand) Help() string {
	helpText := `
usage: sqsc container add [options]

  Adds docker container to the specified Squarescale project.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}

func (c *ContainerAddCommand) validateArgs(project, image string, instances int) error {
	if image == "" {
		return errors.New("Docker image name cannot be empty")
	}

	if instances <= 0 {
		return errors.New("Invalid value provided for instances count: it cannot be 0 or negative")
	}

	return validateProjectName(project)
}
