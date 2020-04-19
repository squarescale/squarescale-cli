package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// VolumeDeleteCommand is a cli.Command implementation for creating a Squarescale project.
type VolumeDeleteCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *VolumeDeleteCommand) Run(args []string) int {
	// Parse flags
	c.flagSet = newFlagSet(c, c.Ui)
	alwaysYes := yesFlag(c.flagSet)
	endpoint := endpointFlag(c.flagSet)
	projectArg := projectFlag(c.flagSet)
	nowait := nowaitFlag(c.flagSet)

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	volumeName, err := volumeNameArg(c.flagSet, 0)
	if err != nil {
		return c.errorWithUsage(err)
	}

	if c.flagSet.NArg() > 4 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()[1:]))
	}

	if err := validateProjectName(*projectArg); err != nil {
		return c.errorWithUsage(err)
	}

	c.Ui.Info("Are you sure you want to delete " + volumeName + "?")
	if *alwaysYes {
		c.Ui.Info("(approved from command line)")
	} else {
		res, err := c.Ui.Ask("y/N")
		if err != nil {
			return c.error(err)
		} else if res != "Y" && res != "y" {
			return c.cancelled()
		}
	}

	res := c.runWithSpinner("deleting volume", endpoint.String(), func(client *squarescale.Client) (string, error) {
		projectName := *projectArg
		err := client.DeleteVolume(projectName, volumeName)
		return "", err
	})
	if res != 0 {
		return res
	}

	if !*nowait {
		c.runWithSpinner("wait for volume delete", endpoint.String(), func(client *squarescale.Client) (string, error) {
			_, err := client.WaitProject(*projectArg)
			if err != nil {
				return "", err
			} else {
				return "", nil
			}
		})
	}

	return 0
}

// Synopsis is part of cli.Command implementation.
func (c *VolumeDeleteCommand) Synopsis() string {
	return "Delete a volume from the project."
}

// Help is part of cli.Command implementation.
func (c *VolumeDeleteCommand) Help() string {
	helpText := `
usage: sqsc volume delete [options] <volume_name>

  Delete the volume from the project.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
