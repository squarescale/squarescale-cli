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

func (b *SchedulingGroupGetCommand) Run(args []string) int {
	b.flagSet = newFlagSet(b, b.Ui)
	endpoint := endpointFlag(b.flagSet)
	projectUUID := projectUUIDFlag(b.flagSet)
	projectName := projectNameFlag(b.flagSet)

	if err := b.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return b.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	schedulingGroupName, err := schedulingGroupNameArg(b.flagSet, 0)
	if err != nil {
		return b.errorWithUsage(err)
	}

	if b.flagSet.NArg() > 1 {
		return b.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", b.flagSet.Args()))
	}

	return b.runWithSpinner("show scheduling group", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
func (b *SchedulingGroupGetCommand) Synopsis() string {
	return "Get a scheduling group of project"
}

// Help is part of cli.Command implementation.
func (b *SchedulingGroupGetCommand) Help() string {
	helpText := `
usage: sqsc scheduling-group get [options] <scheduling_group_name>

  Get a scheduling group of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(b.flagSet))
}
