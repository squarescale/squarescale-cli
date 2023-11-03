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
func (cmd *VolumeAddCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)

	name := volumeNameFlag(cmd.flagSet)
	size := volumeSizeFlag(cmd.flagSet)
	volumeType := volumeTypeFlag(cmd.flagSet)
	zone := volumeZoneFlag(cmd.flagSet)
	nowait := nowaitFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if cmd.flagSet.NArg() > 0 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	var UUID string

	res := cmd.runWithSpinner("add volume", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
		cmd.runWithSpinner("wait for volume add", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
func (cmd *VolumeAddCommand) Synopsis() string {
	return "Add volume to project."
}

// Help is part of cli.Command implementation.
func (cmd *VolumeAddCommand) Help() string {
	helpText := `
usage: sqsc volume add [options] <volume_name>

  Add volume to project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
