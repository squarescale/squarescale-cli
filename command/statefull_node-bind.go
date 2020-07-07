package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// StatefullNodeBindCommand is a cli.Command implementation for creating a Squarescale project.
type StatefullNodeBindCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *StatefullNodeBindCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := c.flagSet.String("project-uuid", "", "set the uuid of the project")
	volumeName := c.flagSet.String("volume-name", "", "Volume name to bind")

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" {
		return c.errorWithUsage(errors.New("Project uuid is mandatory"))
	}

	if *volumeName == "" {
		return c.errorWithUsage(errors.New("Volume name cannot be empty"))
	}

	statefullNodeName, err := statefullNodeNameArg(c.flagSet, 0)
	if err != nil {
		return c.errorWithUsage(err)
	}

	if c.flagSet.NArg() > 1 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	res := c.runWithSpinner("bind statefull node", endpoint.String(), func(client *squarescale.Client) (string, error) {
		msg := fmt.Sprintf("Successfully binded volume '%s' to statefull_node '%s'", *volumeName, statefullNodeName)
		err := client.BindVolumeOnStatefullNode(*projectUUID, statefullNodeName, *volumeName)
		return msg, err
	})
	if res != 0 {
		return res
	}

	return 0
}

// Synopsis is part of cli.Command implementation.
func (c *StatefullNodeBindCommand) Synopsis() string {
	return "Bind a statefull_node from the project."
}

// Help is part of cli.Command implementation.
func (c *StatefullNodeBindCommand) Help() string {
	helpText := `
usage: sqsc statefull_node bind [options] <statefull_node_name>

  Bind a statefull_node on the project.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
