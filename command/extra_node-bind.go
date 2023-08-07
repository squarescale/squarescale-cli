package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ExtraNodeBindCommand is a cli.Command implementation for binding a extra-node.
type ExtraNodeBindCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *ExtraNodeBindCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := projectUUIDFlag(c.flagSet)
	projectName := projectNameFlag(c.flagSet)
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

	extraNodeName, err := extraNodeNameArg(c.flagSet, 0)
	if err != nil {
		return c.errorWithUsage(err)
	}

	if c.flagSet.NArg() > 1 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	res := c.runWithSpinner("bind extra-node", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		msg := fmt.Sprintf("Successfully binded volume '%s' to extra-node '%s'", *volumeName, extraNodeName)
		err = client.BindVolumeOnExtraNode(UUID, extraNodeName, *volumeName)
		return msg, err
	})
	if res != 0 {
		return res
	}

	return 0
}

// Synopsis is part of cli.Command implementation.
func (c *ExtraNodeBindCommand) Synopsis() string {
	return "Bind extra-node to project."
}

// Help is part of cli.Command implementation.
func (c *ExtraNodeBindCommand) Help() string {
	helpText := `
usage: sqsc extra-node bind [options] <extra_node_name>

  Bind extra-node to project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
