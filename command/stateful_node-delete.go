package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// StatefulNodeDeleteCommand is a cli.Command implementation for deleting a stateful-node.
type StatefulNodeDeleteCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *StatefulNodeDeleteCommand) Run(args []string) int {
	// Parse flags
	c.flagSet = newFlagSet(c, c.Ui)
	alwaysYes := yesFlag(c.flagSet)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := projectUUIDFlag(c.flagSet)
	projectName := projectNameFlag(c.flagSet)
	nowait := nowaitFlag(c.flagSet)

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	statefulNodeName, err := statefulNodeNameArg(c.flagSet, 0)
	if err != nil {
		return c.errorWithUsage(err)
	}

	if c.flagSet.NArg() > 4 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()[1:]))
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	c.Ui.Info("Are you sure you want to delete " + statefulNodeName + "?")
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

	res := c.runWithSpinner("deleting stateful-node", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		fmt.Printf("Delete on project `%s` the stateful node `%s`\n", projectToShow, statefulNodeName)
		err = client.DeleteStatefulNode(UUID, statefulNodeName)
		return "", err
	})
	if res != 0 {
		return res
	}

	if !*nowait {
		c.runWithSpinner("wait for stateful-node delete", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
func (c *StatefulNodeDeleteCommand) Synopsis() string {
	return "Delete stateful-node from project."
}

// Help is part of cli.Command implementation.
func (c *StatefulNodeDeleteCommand) Help() string {
	helpText := `
usage: sqsc stateful-node delete [options] <stateful_node_name>

  Delete stateful node from project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
