package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// SchedulingGroupAddCommand is a cli.Command implementation for creating a scheduling group.
type SchedulingGroupAddCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *SchedulingGroupAddCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	schedulingGroupName, err := schedulingGroupNameArg(cmd.flagSet, 0)
	if err != nil {
		return cmd.errorWithUsage(err)
	}

	if cmd.flagSet.NArg() > 1 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}
	var UUID string

	res := cmd.runWithSpinner("add scheduling group", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		msg := fmt.Sprintf("Successfully added scheduling group '%s' to project '%s'", schedulingGroupName, projectToShow)
		_, err = client.AddSchedulingGroup(UUID, schedulingGroupName)
		return msg, err
	})
	if res != 0 {
		return res
	}

	return 0
}

// Synopsis is part of cli.Command implementation.
func (cmd *SchedulingGroupAddCommand) Synopsis() string {
	return "Add scheduling group to project."
}

// Help is part of cli.Command implementation.
func (cmd *SchedulingGroupAddCommand) Help() string {
	helpText := `
usage: sqsc scheduling-group add [options] <scheduling_group_name>

  Add scheduling group to project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
