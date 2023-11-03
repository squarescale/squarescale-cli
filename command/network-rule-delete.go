package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// NetworkRuleDeleteCommand is a cli.Command implementation for top level `sqsc network-rule` command.
type NetworkRuleDeleteCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *NetworkRuleDeleteCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)
	ruleName := networkRuleNameFlag(cmd.flagSet)
	serviceName := networkServiceNameFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if cmd.flagSet.NArg() > 0 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if *serviceName == "" {
		return cmd.errorWithUsage(fmt.Errorf(("Service name is mandatory.")))
	}

	if *ruleName == "" {
		return cmd.errorWithUsage(fmt.Errorf(("Name is mandatory.")))
	}

	return cmd.runWithSpinner("delete network rule", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		err = client.DeleteNetworkRule(UUID, *serviceName, *ruleName)
		return "", err
	})

}

// Synopsis is part of cli.Command implementation.
func (cmd *NetworkRuleDeleteCommand) Synopsis() string {
	return "Delete a network rule"
}

// Help is part of cli.Command implementation.
func (cmd *NetworkRuleDeleteCommand) Help() string {
	helpText := `
usage: sqsc network-rule delete [options]

  Delete a network rule.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
