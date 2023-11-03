package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// NetworkPolicyDeleteCommand is a cli.Command implementation for deleting network policies
type NetworkPolicyDeleteCommand struct {
	Meta
	flagSet *flag.FlagSet
}

func (c *NetworkPolicyDeleteCommand) Run(args []string) int {
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

	if c.flagSet.NArg() > 1 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}
	version := networkPolicyVersionArg(c.flagSet)

	return c.runWithSpinner("delete network policy version", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		err = client.DeleteNetworkPolicy(UUID, version)
		if err != nil {
			return "", err
		}

		return "deleted " + version, nil
	})
}

// Synopsis is part of cli.Command implementation.
func (c *NetworkPolicyDeleteCommand) Synopsis() string {
	return "Delete a network policy"
}

// Help is part of cli.Command implementation.
func (c *NetworkPolicyDeleteCommand) Help() string {
	helpText := `
usage: sqsc network-policy delete VERSION

  Delete network policy
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
