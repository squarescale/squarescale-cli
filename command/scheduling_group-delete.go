package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// StatefulNodeDeleteCommand is a cli.Command implementation for deleting a scheduling group.
type SchedulingGroupDeleteCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *SchedulingGroupDeleteCommand) Run(args []string) int {
	// Parse flags
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	alwaysYes := yesFlag(cmd.flagSet)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	scheduligGroupName, err := schedulingGroupNameArg(cmd.flagSet, 0)
	if err != nil {
		return cmd.errorWithUsage(err)
	}

	if cmd.flagSet.NArg() > 1 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()[1:]))
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	cmd.Ui.Info("Are you sure you want to delete " + scheduligGroupName + "?")
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

	res := cmd.runWithSpinner("deleting scheduling group", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		fmt.Printf("Delete on project `%s` the scheduling group `%s`\n", projectToShow, scheduligGroupName)
		err = client.DeleteSchedulingGroup(UUID, scheduligGroupName)
		return "", err
	})
	if res != 0 {
		return res
	}

	return 0
}

// Synopsis is part of cli.Command implementation.
func (cmd *SchedulingGroupDeleteCommand) Synopsis() string {
	return "Delete scheduling group from project."
}

// Help is part of cli.Command implementation.
func (cmd *SchedulingGroupDeleteCommand) Help() string {
	helpText := `
usage: sqsc scheduling-group delete [options] <scheduling_group_name>

  Delete scheduling group from project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
