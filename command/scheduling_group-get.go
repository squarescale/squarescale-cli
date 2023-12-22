package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// SchedulingGroupGetCommand is a cli.Command implementation for get a scheduling group.
type SchedulingGroupGetCommand struct {
	Meta
	flagSet *flag.FlagSet
}

func (cmd *SchedulingGroupGetCommand) Run(args []string) int {
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

	return cmd.runWithSpinner("show scheduling group", endpoint.String(), func(client *squarescale.Client) (string, error) {
		var UUID string
		var err error
		if *projectUUID == "" {
			UUID, err = client.ProjectByName(*projectName)
			if err != nil {
				return "", err
			}
		} else {
			UUID = *projectUUID
		}

		schedulingGroup, err := client.GetSchedulingGroupInfo(UUID, schedulingGroupName)
		if err != nil {
			return "", err
		}

		var msg string
		msg = ""
		msg += fmt.Sprintf("Name:\t\t%s\n", schedulingGroup.Name)
		msg += fmt.Sprintf("Nodes:\t\t%s\n", client.GetSchedulingGroupNodes(schedulingGroup, "\n\t\t"))
		msg += fmt.Sprintf("Services:\t%s\n", client.GetSchedulingGroupServices(schedulingGroup, "\n\t\t"))

		return msg, nil
	})
}

// Synopsis is part of cli.Command implementation.
func (cmd *SchedulingGroupGetCommand) Synopsis() string {
	return "Get a scheduling group of project"
}

// Help is part of cli.Command implementation.
func (cmd *SchedulingGroupGetCommand) Help() string {
	helpText := `
usage: sqsc scheduling-group get [options] <scheduling_group_name>

  Get a scheduling group of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
