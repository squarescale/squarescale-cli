package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// StatefullNodeDeleteCommand is a cli.Command implementation for creating a Squarescale project.
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
	projectUUID := c.flagSet.String("project-uuid", "", "set the uuid of the project")
	nowait := nowaitFlag(c.flagSet)

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	statefullNodeName, err := statefullNodeNameArg(c.flagSet, 0)
	if err != nil {
		return c.errorWithUsage(err)
	}

	if c.flagSet.NArg() > 4 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()[1:]))
	}

	if *projectUUID == "" {
		return c.errorWithUsage(errors.New("Project uuid is mandatory"))
	}

	c.Ui.Info("Are you sure you want to delete " + statefullNodeName + "?")
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

	res := c.runWithSpinner("deleting statefull node", endpoint.String(), func(client *squarescale.Client) (string, error) {
		uuid := *projectUUID
		fmt.Printf("Delete on project `%s` the statefull node `%s`\n", uuid, statefullNodeName)
		err := client.DeleteStatefullNode(uuid, statefullNodeName)
		return "", err
	})
	if res != 0 {
		return res
	}

	if !*nowait {
		c.runWithSpinner("wait for statefull node delete", endpoint.String(), func(client *squarescale.Client) (string, error) {
			_, err := client.WaitProject(*projectUUID, 5)
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
	return "Delete stateful node from project."
}

// Help is part of cli.Command implementation.
func (c *StatefulNodeDeleteCommand) Help() string {
	helpText := `
usage: sqsc stateful node delete [options] <stateful_node_name>

  Delete stateful node from project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
