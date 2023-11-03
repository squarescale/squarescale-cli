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

func (cmd *NetworkPolicyDeployCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if cmd.flagSet.NArg() > 1 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}
	version := networkPolicyVersionArg(cmd.flagSet)
	if version == "" {
		return cmd.errorWithUsage(fmt.Errorf("version is mandator"))
	}

	return cmd.runWithSpinner("deploy network policy version", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
func (cmd *NetworkPolicyDeployCommand) Synopsis() string {
	return "Deploy specific network policy version"
}

// Help is part of cli.Command implementation.
func (cmd *NetworkPolicyDeployCommand) Help() string {
	helpText := `
usage: sqsc network-policy deploy [options] VERSION

  Deploy network policy version
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
