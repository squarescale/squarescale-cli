package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ExternalNodeAddCommand is a cli.Command implementation for creating a external-node.
type ExternalNodeAddCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *ExternalNodeAddCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := c.flagSet.String("project-uuid", "", "set the uuid of the project")
	projectName := c.flagSet.String("project-name", "", "set the name of the project")

	publicIP := c.flagSet.String("public-ip", "", "External node public IP")
	nowait := nowaitFlag(c.flagSet)

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	externalNodeName, err := externalNodeNameArg(c.flagSet, 0)
	if err != nil {
		return c.errorWithUsage(err)
	}

	if *publicIP == "" {
		return c.errorWithUsage(errors.New("Public IP is mandatory"))
	}

	if c.flagSet.NArg() > 1 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}
	var UUID string

	res := c.runWithSpinner("add external node", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		msg := fmt.Sprintf("Successfully added external node '%s' to project '%s'", externalNodeName, projectToShow)
		_, err = client.AddExternalNode(UUID, externalNodeName, *publicIP)
		return msg, err
	})
	if res != 0 {
		return res
	}

	if !*nowait {
		c.runWithSpinner("wait for external node add", endpoint.String(), func(client *squarescale.Client) (string, error) {
			externalNode, err := client.WaitExternalNode(UUID, externalNodeName, 5)
			if err != nil {
				return "", err
			} else {
				return externalNode.Name, nil
			}
		})
	}

	return 0
}

// Synopsis is part of cli.Command implementation.
func (c *ExternalNodeAddCommand) Synopsis() string {
	return "Add external node node to project."
}

// Help is part of cli.Command implementation.
func (c *ExternalNodeAddCommand) Help() string {
	helpText := `
usage: sqsc external-node add [options] <external_node_name>

  Add external node node to project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}