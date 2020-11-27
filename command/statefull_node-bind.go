package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// StatefullNodeBindCommand is a cli.Command implementation for creating a Squarescale project.
type StatefulNodeBindCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *StatefulNodeBindCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := c.flagSet.String("project-uuid", "", "set the uuid of the project")
	projectName := c.flagSet.String("project-name", "", "set the name of the project")
	volumeName := c.flagSet.String("volume-name", "", "Volume name to bind")

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
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

		msg := fmt.Sprintf("Successfully binded volume '%s' to statefull_node '%s'", *volumeName, statefullNodeName)
		err = client.BindVolumeOnStatefullNode(UUID, statefullNodeName, *volumeName)
		return msg, err
	})
	if res != 0 {
		return res
	}

	return 0
}

// Synopsis is part of cli.Command implementation.
func (c *StatefulNodeBindCommand) Synopsis() string {
	return "Bind stateful node to project."
}

// Help is part of cli.Command implementation.
func (c *StatefulNodeBindCommand) Help() string {
	helpText := `
usage: sqsc stateful node bind [options] <stateful_node_name>

  Bind stateful node to project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
