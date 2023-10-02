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

func (b *NetworkPolicyDeleteCommand) Run(args []string) int {
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

	return b.runWithSpinner("delete network policy version", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
func (b *NetworkPolicyDeleteCommand) Synopsis() string {
	return "Delete a network policy"
}

// Help is part of cli.Command implementation.
func (b *NetworkPolicyDeleteCommand) Help() string {
	helpText := `
usage: sqsc network-policy delete VERSION

  Delete network policy
`
	return strings.TrimSpace(helpText + optionsFromFlags(b.flagSet))
}
