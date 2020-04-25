package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// StatefullNodeAddCommand is a cli.Command implementation for creating a Squarescale project.
type StatefullNodeAddCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *StatefullNodeAddCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	project := projectFlag(c.flagSet)

	nodeType := c.flagSet.String("node-type", "t2.micro", "Statefull node type")
	zone := c.flagSet.String("zone", "eu-west-1a", "Statefull node zone")
	nowait := nowaitFlag(c.flagSet)

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	statefullNodeName, err := statefullNodeNameArg(c.flagSet, 0)
	if err != nil {
		return c.errorWithUsage(err)
	}

	if c.flagSet.NArg() > 1 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if err := validateProjectName(*project); err != nil {
		return c.errorWithUsage(err)
	}

	res := c.runWithSpinner("add statefull_node", endpoint.String(), func(client *squarescale.Client) (string, error) {
		msg := fmt.Sprintf("Successfully added statefull_node '%s' to project '%s'", statefullNodeName, *project)
		_, err := client.AddStatefullNode(*project, statefullNodeName, *nodeType, *zone)
		return msg, err
	})
	if res != 0 {
		return res
	}

	if !*nowait {
		c.runWithSpinner("wait for statefull_node add", endpoint.String(), func(client *squarescale.Client) (string, error) {
			statefullNode, err := client.WaitStatefullNode(*project, statefullNodeName, 5)
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
func (c *StatefullNodeAddCommand) Synopsis() string {
	return "Add a statefull_node from the project."
}

// Help is part of cli.Command implementation.
func (c *StatefullNodeAddCommand) Help() string {
	helpText := `
usage: sqsc statefull_node add [options] <statefull_node_name>

  Add a statefull_node on the project.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
