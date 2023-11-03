package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ExternalNodeGetCommand is a cli.Command implementation for get a external node.
type ExternalNodeGetCommand struct {
	Meta
	flagSet *flag.FlagSet
}

func (c *ExternalNodeGetCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := projectUUIDFlag(c.flagSet)
	projectName := projectNameFlag(c.flagSet)

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

	if c.flagSet.NArg() > 1 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	return c.runWithSpinner("show external nodes", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		externalNode, err := client.GetExternalNodeInfo(UUID, externalNodeName)
		if err != nil {
			return "", err
		}

		var msg string
		msg = ""
		msg += fmt.Sprintf("Name:\t\t%s\n", externalNode.Name)
		msg += fmt.Sprintf("Public IP:\t%s\n", externalNode.PublicIP)
		msg += fmt.Sprintf("Status:\t\t%s\n", externalNode.Status)

		return msg, nil
	})
}

// Synopsis is part of cli.Command implementation.
func (c *ExternalNodeGetCommand) Synopsis() string {
	return "Get a external node of project"
}

// Help is part of cli.Command implementation.
func (c *ExternalNodeGetCommand) Help() string {
	helpText := `
usage: sqsc external-node get [options] <external_node_name>

  Get a external node of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
