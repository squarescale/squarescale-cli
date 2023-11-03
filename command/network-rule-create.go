package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// NetworkRuleCreateCommand is a cli.Command implementation for top level `sqsc network-rule` command.
type NetworkRuleCreateCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *NetworkRuleCreateCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)
	ruleName := networkRuleNameFlag(cmd.flagSet)
	serviceName := networkServiceNameFlag(cmd.flagSet)
	externalProtocol := networkExternalProtocolFlag(cmd.flagSet)
	internalProtocol := networkInternalProtocolFlag(cmd.flagSet)
	internalPort := networkInternalPortFlag(cmd.flagSet)
	domainExpression := networkDomainFlag(cmd.flagSet)
	pathPrefix := networkPathFlag(cmd.flagSet)
	externalPort := 0

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

	if *internalProtocol == "" && !(*internalProtocol == "ssh" || *internalProtocol == "http" || *internalProtocol == "https") {
		return cmd.errorWithUsage(fmt.Errorf(("Internal protocol is mandatory and must be ssh, http or https")))
	}

	if *internalPort == 0 || *internalPort < 1 || *internalPort > 65535 {
		return cmd.errorWithUsage(fmt.Errorf(("Internal port is mandatory and must be contains between 0 and 65535")))
	}

	if *externalProtocol == "" {
		return cmd.errorWithUsage(fmt.Errorf(("External protocol is mandatory.")))
	} else if *externalProtocol == "ssh" {
		externalPort = 22
	} else if *externalProtocol == "http" {
		externalPort = 80
	} else if *externalProtocol == "https" {
		externalPort = 443
	} else {
		return cmd.errorWithUsage(fmt.Errorf(("External protocol is ssh, http or https.")))
	}

	return cmd.runWithSpinner("create network rule", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		newRule := squarescale.NetworkRule{
			Name:             *ruleName,
			InternalPort:     *internalPort,
			InternalProtocol: *internalProtocol,
			ExternalPort:     externalPort,
			ExternalProtocol: *externalProtocol,
			DomainExpression: *domainExpression,
			PathPrefix:       *pathPrefix,
		}
		err = client.CreateNetworkRule(UUID, *serviceName, newRule)
		return "", err
	})

}

// Synopsis is part of cli.Command implementation.
func (cmd *NetworkRuleCreateCommand) Synopsis() string {
	return "Create a network rule"
}

// Help is part of cli.Command implementation.
func (cmd *NetworkRuleCreateCommand) Help() string {
	helpText := `
usage: sqsc network-rule create [options]

  Create a network rule.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
