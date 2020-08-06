package command

import (
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
func (c *NetworkRuleCreateCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := c.flagSet.String("project-uuid", "", "set the uuid of the project")
	ruleName := c.flagSet.String("name", "", "name of the rule")
	serviceName := c.flagSet.String("service-name", "", "name of the service the rule will be attached")
	externalProtocol := c.flagSet.String("external-protocol", "", "name of the exposed protocol")
	internalProtocol := c.flagSet.String("internal-protocol", "", "name of the internal protocol")
	internalPort := c.flagSet.Int("internal-port", 0, "value of the internal port")
	domainExpression := c.flagSet.String("domain-expression", "", "custom domain the service is accessed")
	externalPort := 0

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" {
		return c.errorWithUsage(fmt.Errorf(("Project uuid is mandatory.")))
	}

	if *ruleName == "" {
		return c.errorWithUsage(fmt.Errorf(("Name is mandatory.")))
	}

	if *internalProtocol == "" && !(*internalProtocol == "ssh" || *internalProtocol == "http" || *internalProtocol == "https") {
		return c.errorWithUsage(fmt.Errorf(("Internal protocol is mandatory and must be ssh, http or https")))
	}

	if *internalPort == 0 || *internalPort < 1 || *internalPort > 65535 {
		return c.errorWithUsage(fmt.Errorf(("Internal port is mandatory and must be contains between 0 and 65535")))
	}

	if *externalProtocol == "" {
		return c.errorWithUsage(fmt.Errorf(("External protocol is mandatory.")))
	} else if *externalProtocol == "ssh" {
		externalPort = 22
	} else if *externalProtocol == "http" {
		externalPort = 80
	} else if *externalProtocol == "https" {
		externalPort = 443
	} else {
		return c.errorWithUsage(fmt.Errorf(("External protocol is ssh, http or https.")))
	}

	return c.runWithSpinner("create network rule", endpoint.String(), func(client *squarescale.Client) (string, error) {

		newRule := squarescale.NetworkRule{
			Name:             *ruleName,
			InternalPort:     *internalPort,
			InternalProtocol: *internalProtocol,
			ExternalPort:     externalPort,
			ExternalProtocol: *externalProtocol,
			DomainExpression: *domainExpression,
		}
		err := client.CreateNetworkRule(*projectUUID, *serviceName, newRule)
		return "", err
	})

}

// Synopsis is part of cli.Command implementation.
func (c *NetworkRuleCreateCommand) Synopsis() string {
	return "Create a network rule"
}

// Help is part of cli.Command implementation.
func (c *NetworkRuleCreateCommand) Help() string {
	helpText := `
usage: sqsc network-rule create [options]

  Create a network rule.
`
	return strings.TrimSpace(helpText)
}
