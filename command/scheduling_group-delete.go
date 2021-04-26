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
func (c *SchedulingGroupDeleteCommand) Run(args []string) int {
	// Parse flags
	c.flagSet = newFlagSet(c, c.Ui)
	alwaysYes := yesFlag(c.flagSet)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := c.flagSet.String("project-uuid", "", "set the uuid of the project")
	projectName := c.flagSet.String("project-name", "", "set the name of the project")

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	scheduligGroupName, err := schedulingGroupNameArg(c.flagSet, 0)
	if err != nil {
		return c.errorWithUsage(err)
	}

	if c.flagSet.NArg() > 4 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()[1:]))
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	c.Ui.Info("Are you sure you want to delete " + scheduligGroupName + "?")
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

	var UUID string

	res := c.runWithSpinner("deleting scheduling group", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
func (c *SchedulingGroupDeleteCommand) Synopsis() string {
	return "Delete scheduling group from project."
}

// Help is part of cli.Command implementation.
func (c *SchedulingGroupDeleteCommand) Help() string {
	helpText := `
usage: sqsc scheduling-group delete [options] <scheduling_group_name>

  Delete scheduling group from project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
