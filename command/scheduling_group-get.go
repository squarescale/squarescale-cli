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

func (c *SchedulingGroupGetCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := projectUUIDFlag(c.flagSet)
	projectName := projectNameFlag(c.flagSet)

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	schedulingGroupName, err := schedulingGroupNameArg(c.flagSet, 0)
	if err != nil {
		return c.errorWithUsage(err)
	}

	if c.flagSet.NArg() > 1 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	return c.runWithSpinner("show scheduling group", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
func (c *SchedulingGroupGetCommand) Synopsis() string {
	return "Get a scheduling group of project"
}

// Help is part of cli.Command implementation.
func (c *SchedulingGroupGetCommand) Help() string {
	helpText := `
usage: sqsc scheduling-group get [options] <scheduling_group_name>

  Get a scheduling group of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
