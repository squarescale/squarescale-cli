package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// NetworkPolicyDeployCommand is a cli.Command implementation for deploying network policies
type NetworkPolicyDeployCommand struct {
	Meta
	flagSet *flag.FlagSet
}

func (c *NetworkPolicyDeployCommand) Run(args []string) int {
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
	if version == "" {
		return c.errorWithUsage(fmt.Errorf("version is mandator"))
	}

	return c.runWithSpinner("deploy network policy version", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		err = client.DeployNetworkPolicy(UUID, version)
		if err != nil {
			return "", err
		}

		return "deploying " + version, nil
	})
}

// Synopsis is part of cli.Command implementation.
func (c *NetworkPolicyDeployCommand) Synopsis() string {
	return "Deploy specific network policy version"
}

// Help is part of cli.Command implementation.
func (c *NetworkPolicyDeployCommand) Help() string {
	helpText := `
usage: sqsc network-policy deploy [options] VERSION

  Deploy network policy version
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
