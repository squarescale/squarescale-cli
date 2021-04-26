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
func (b *SchedulingGroupListCommand) Run(args []string) int {
	b.flagSet = newFlagSet(b, b.Ui)
	endpoint := endpointFlag(b.flagSet)
	projectUUID := b.flagSet.String("project-uuid", "", "set the uuid of the project")
	projectName := b.flagSet.String("project-name", "", "set the name of the project")

	if err := b.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return b.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if b.flagSet.NArg() > 0 {
		return b.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", b.flagSet.Args()))
	}

	return b.runWithSpinner("list scheduling groups", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		schedulingGroup, err := client.GetSchedulingGroups(UUID)
		if err != nil {
			return "", err
		}

		var msg string
		msg = "Name\n"
		for _, n := range schedulingGroup {
			msg += fmt.Sprintf("%s\n", n.Name)
		}

		if len(schedulingGroup) == 0 {
			msg = "No scheduling groups found"
		}

		return msg, nil
	})
}

// Synopsis is part of cli.Command implementation.
func (b *SchedulingGroupListCommand) Synopsis() string {
	return "List scheduling groups of project"
}

// Help is part of cli.Command implementation.
func (b *SchedulingGroupListCommand) Help() string {
	helpText := `
usage: sqsc scheduling-group list [options]

  List scheduling groups of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(b.flagSet))
}
