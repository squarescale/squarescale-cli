package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// VolumeAddCommand is a cli.Command implementation for creating a SquareScale project.
type VolumeAddCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *VolumeAddCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := c.flagSet.String("project-uuid", "", "set the uuid of the project")
	projectName := c.flagSet.String("project-name", "", "set the name of the project")

	name := c.flagSet.String("name", "", "Volume name")
	size := c.flagSet.Int("size", 1, "Volume size (in Gb)")
	volumeType := c.flagSet.String("type", "gp2", "Volume type")
	zone := c.flagSet.String("zone", "eu-west-1a", "Volume zone")
	nowait := nowaitFlag(c.flagSet)

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	var UUID string

	res := c.runWithSpinner("add volume", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		msg := fmt.Sprintf("Successfully added volume '%s' to project '%s'", *name, projectToShow)
		err = client.AddVolume(UUID, *name, *size, *volumeType, *zone)
		return msg, err
	})
	if res != 0 {
		return res
	}

	if !*nowait {
		c.runWithSpinner("wait for volume add", endpoint.String(), func(client *squarescale.Client) (string, error) {
			volume, err := client.WaitVolume(UUID, *name, 5)
			if err != nil {
				return "", err
			} else {
				return volume.Name, nil
			}
		})
	}

	return 0
}

// Synopsis is part of cli.Command implementation.
func (c *VolumeAddCommand) Synopsis() string {
	return "Add volume to project."
}

// Help is part of cli.Command implementation.
func (c *VolumeAddCommand) Help() string {
	helpText := `
usage: sqsc volume add [options] <volume_name>

  Add volume to project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
