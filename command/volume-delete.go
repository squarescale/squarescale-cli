package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// VolumeDeleteCommand is a cli.Command implementation for creating a SquareScale project.
type VolumeDeleteCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *VolumeDeleteCommand) Run(args []string) int {
	// Parse flags
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	alwaysYes := yesFlag(cmd.flagSet)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)
	nowait := nowaitFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	volumeName, err := volumeNameArg(cmd.flagSet, 0)
	if err != nil {
		return cmd.errorWithUsage(err)
	}

	if cmd.flagSet.NArg() > 4 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()[1:]))
	}

	cmd.Ui.Info("Are you sure you want to delete " + volumeName + "?")
	if *alwaysYes {
		cmd.Ui.Info("(approved from command line)")
	} else {
		res, err := cmd.Ui.Ask("y/N")
		if err != nil {
			return cmd.error(err)
		} else if res != "Y" && res != "y" {
			return cmd.cancelled()
		}
	}

	var UUID string

	res := cmd.runWithSpinner("deleting volume", endpoint.String(), func(client *squarescale.Client) (string, error) {
		var err error
		if *projectUUID == "" {
			UUID, err = client.ProjectByName(*projectName)
			if err != nil {
				return "", err
			}
		} else {
			UUID = *projectUUID
		}

		err = client.DeleteVolume(UUID, volumeName)
		return "", err
	})
	if res != 0 {
		return res
	}

	if !*nowait {
		cmd.runWithSpinner("wait for volume delete", endpoint.String(), func(client *squarescale.Client) (string, error) {
			_, err := client.WaitProject(UUID, 5)
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
func (cmd *VolumeDeleteCommand) Synopsis() string {
	return "Delete volume from project."
}

// Help is part of cli.Command implementation.
func (cmd *VolumeDeleteCommand) Help() string {
	helpText := `
usage: sqsc volume delete [options] <volume_name>

  Delete volume from project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
