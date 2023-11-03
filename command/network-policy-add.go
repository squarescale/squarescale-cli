package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// NetworkPolicyAddCommand is a cli.Command implementation for adding network policies
type NetworkPolicyAddCommand struct {
	Meta
	flagSet *flag.FlagSet
}

func (c *NetworkPolicyAddCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := projectUUIDFlag(c.flagSet)
	projectName := projectNameFlag(c.flagSet)
	policyName := networkPolicyNameFlag(c.flagSet)

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if *policyName == "" {
		return c.errorWithUsage(errors.New("Policy name is mandatory"))
	}

	if c.flagSet.NArg() > 2 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	policy := squarescale.NetworkPolicy{}
	policy.Name = *policyName

	policyFile, err := networkPolicyFileArg(c.flagSet)
	if err != nil {
		return c.errorWithUsage(err)
	}

	if policyFile != "" {
		err := policy.ImportPoliciesFromFile(policyFile)
		if err != nil {
			return c.error(err)
		}
	}

	return c.runWithSpinner("add network policy", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		policyVersion, err := client.AddNetworkPolicy(UUID, policy.Name, policy.Policies)
		if err != nil {
			return "", err
		}

		return "Network policy version: " + policyVersion, nil
	})
}

// Synopsis is part of cli.Command implementation.
func (c *NetworkPolicyAddCommand) Synopsis() string {
	return "Create a new network policy"
}

// Help is part of cli.Command implementation.
func (c *NetworkPolicyAddCommand) Help() string {
	helpText := `
usage: sqsc network-policy add -network-policy-name NAME FILE-PATH

  Add a new network policy to a project
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
