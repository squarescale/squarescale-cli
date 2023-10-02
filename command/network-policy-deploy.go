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

func (b *NetworkPolicyDeployCommand) Run(args []string) int {
	b.flagSet = newFlagSet(b, b.Ui)
	endpoint := endpointFlag(b.flagSet)
	projectUUID := projectUUIDFlag(b.flagSet)
	projectName := projectNameFlag(b.flagSet)

	if err := b.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return b.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if b.flagSet.NArg() > 1 {
		return b.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", b.flagSet.Args()))
	}
	version := networkPolicyVersionArg(b.flagSet)
	if version == "" {
		return b.errorWithUsage(fmt.Errorf("version is mandator"))
	}

	return b.runWithSpinner("deploy network policy version", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
func (b *NetworkPolicyDeployCommand) Synopsis() string {
	return "Deploy specific network policy version"
}

// Help is part of cli.Command implementation.
func (b *NetworkPolicyDeployCommand) Help() string {
	helpText := `
usage: sqsc network-policy deploy [options] VERSION

  Deploy network policy version
`
	return strings.TrimSpace(helpText + optionsFromFlags(b.flagSet))
}
