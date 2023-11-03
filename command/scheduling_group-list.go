package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// SchedulingGroupListCommand is a cli.Command implementation for listing scheduling groups.
type SchedulingGroupListCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *SchedulingGroupListCommand) Run(args []string) int {
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

	if cmd.flagSet.NArg() > 0 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	return cmd.runWithSpinner("list scheduling groups", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		schedulingGroups, err := client.GetSchedulingGroups(UUID)
		if err != nil {
			return "", err
		}

		var msg string
		for _, sg := range schedulingGroups {
			msg += fmt.Sprintf("[%s]\n", sg.Name)
			msg += fmt.Sprintf("  Nodes:\n")
			if client.GetSchedulingGroupNodes(sg, "\n\t") != "" {
				msg += fmt.Sprintf("\t%s\n\n", client.GetSchedulingGroupNodes(sg, "\n\t"))
			}
			msg += fmt.Sprintf("  Services:\n")
			if client.GetSchedulingGroupServices(sg, "\n\t") != "" {
				msg += fmt.Sprintf("\t%s\n\n", client.GetSchedulingGroupServices(sg, "\n\t"))
			}
		}

		if len(schedulingGroups) == 0 {
			msg = "No scheduling groups found"
		}

		return msg, nil
	})
}

// Synopsis is part of cli.Command implementation.
func (cmd *SchedulingGroupListCommand) Synopsis() string {
	return "List scheduling groups of project"
}

// Help is part of cli.Command implementation.
func (cmd *SchedulingGroupListCommand) Help() string {
	helpText := `
usage: sqsc scheduling-group list [options]

  List scheduling groups of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
