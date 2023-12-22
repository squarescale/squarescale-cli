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

func (cmd *NetworkPolicyDeleteCommand) Run(args []string) int {
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

	return cmd.runWithSpinner("delete network policy version", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
func (cmd *NetworkPolicyDeleteCommand) Synopsis() string {
	return "Delete a network policy"
}

// Help is part of cli.Command implementation.
func (cmd *NetworkPolicyDeleteCommand) Help() string {
	helpText := `
usage: sqsc network-policy delete VERSION

  Delete network policy
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
