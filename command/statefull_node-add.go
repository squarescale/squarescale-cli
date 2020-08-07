package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// StatefullNodeAddCommand is a cli.Command implementation for creating a Squarescale project.
type StatefulNodeAddCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *StatefulNodeAddCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := c.flagSet.String("project-uuid", "", "set the uuid of the project")

	nodeType := c.flagSet.String("node-type", "dev", "Statefull node type")
	zone := c.flagSet.String("zone", "eu-west-1a", "Statefull node zone")
	nowait := nowaitFlag(c.flagSet)

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" {
		return c.errorWithUsage(errors.New("Project uuid is mandatory"))
	}

	statefullNodeName, err := statefullNodeNameArg(c.flagSet, 0)
	if err != nil {
		return c.errorWithUsage(err)
	}

	if c.flagSet.NArg() > 1 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	res := c.runWithSpinner("add statefull_node", endpoint.String(), func(client *squarescale.Client) (string, error) {
		msg := fmt.Sprintf("Successfully added statefull_node '%s' to project '%s'", statefullNodeName, *projectUUID)
		_, err := client.AddStatefullNode(*projectUUID, statefullNodeName, *nodeType, *zone)
		return msg, err
	})
	if res != 0 {
		return res
	}

	if !*nowait {
		c.runWithSpinner("wait for statefull_node add", endpoint.String(), func(client *squarescale.Client) (string, error) {
			statefullNode, err := client.WaitStatefullNode(*projectUUID, statefullNodeName, 5)
			if err != nil {
				return "", err
			} else {
				return statefullNode.Name, nil
			}
		})
	}

	return 0
}

// Synopsis is part of cli.Command implementation.
func (c *StatefulNodeAddCommand) Synopsis() string {
	return "Add stateful node to project."
}

// Help is part of cli.Command implementation.
func (c *StatefulNodeAddCommand) Help() string {
	helpText := `
usage: sqsc stateful node add [options] <stateful_node_name>

  Add stateful node to project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
