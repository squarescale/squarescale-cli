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

func (b *NetworkPolicyAddCommand) Run(args []string) int {
	b.flagSet = newFlagSet(b, b.Ui)
	endpoint := endpointFlag(b.flagSet)
	projectUUID := projectUUIDFlag(b.flagSet)
	projectName := projectNameFlag(b.flagSet)
	policyName := networkPolicyNameFlag(b.flagSet)

	if err := b.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return b.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if *policyName == "" {
		return b.errorWithUsage(errors.New("Policy name is mandatory"))
	}

	if b.flagSet.NArg() > 2 {
		return b.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", b.flagSet.Args()))
	}

	policy := squarescale.NetworkPolicy{}
	policy.Name = *policyName

	policyFile, err := networkPolicyFileArg(b.flagSet)
	if err != nil {
		return b.errorWithUsage(err)
	}

	if policyFile != "" {
		err := policy.ImportPoliciesFromFile(policyFile)
		if err != nil {
			b.error(err)
		}
	}

	return b.runWithSpinner("add network policy", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
func (b *NetworkPolicyAddCommand) Synopsis() string {
	return "Create a new network policy"
}

// Help is part of cli.Command implementation.
func (b *NetworkPolicyAddCommand) Help() string {
	helpText := `
usage: sqsc network-policy add -network-policy-name NAME FILE-PATH

  Add a new network policy to a project
`
	return strings.TrimSpace(helpText + optionsFromFlags(b.flagSet))
}
