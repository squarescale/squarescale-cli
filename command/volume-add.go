package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// VolumeAddCommand is a cli.Command implementation for creating a Squarescale project.
type VolumeAddCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *VolumeAddCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	project := projectFlag(c.flagSet)

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

	res := c.runWithSpinner("add volume", endpoint.String(), func(client *squarescale.Client) (string, error) {
		msg := fmt.Sprintf("Successfully added volume '%s' to project '%s'", *name, *project)
		err := client.AddVolume(*project, *name, *size, *volumeType, *zone)
		return msg, err
	})
	if res != 0 {
		return res
	}

	if !*nowait {
		c.runWithSpinner("wait for volume add", endpoint.String(), func(client *squarescale.Client) (string, error) {
			volume, err := client.WaitVolume(*project, *name)
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
	return "Add a volume from the project."
}

// Help is part of cli.Command implementation.
func (c *VolumeAddCommand) Help() string {
	helpText := `
usage: sqsc volume add [options] <volume_name>

  Add a volume on the project.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
