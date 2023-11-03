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

func (cmd *NetworkPolicyAddCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)
	policyName := networkPolicyNameFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if *policyName == "" {
		return cmd.errorWithUsage(errors.New("Policy name is mandatory"))
	}

	if cmd.flagSet.NArg() > 2 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	policy := squarescale.NetworkPolicy{}
	policy.Name = *policyName

	policyFile, err := networkPolicyFileArg(cmd.flagSet)
	if err != nil {
		return cmd.errorWithUsage(err)
	}

	if policyFile != "" {
		err := policy.ImportPoliciesFromFile(policyFile)
		if err != nil {
			return cmd.error(err)
		}
	}

	return cmd.runWithSpinner("add network policy", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
func (cmd *NetworkPolicyAddCommand) Synopsis() string {
	return "Create a new network policy"
}

// Help is part of cli.Command implementation.
func (cmd *NetworkPolicyAddCommand) Help() string {
	helpText := `
usage: sqsc network-policy add -network-policy-name NAME FILE-PATH

  Add a new network policy to a project
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
