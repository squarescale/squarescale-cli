package command

import (
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
func (c *NetworkRuleDeleteCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := c.flagSet.String("project-uuid", "", "set the uuid of the project")
	ruleName := c.flagSet.String("name", "", "name of the rule")
	serviceName := c.flagSet.String("service-name", "", "name of the service the rule will be attached")

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if *projectUUID == "" {
		return c.errorWithUsage(fmt.Errorf(("Project uuid is mandatory.")))
	}

	if *serviceName == "" {
		return c.errorWithUsage(fmt.Errorf(("Service name is mandatory.")))
	}

	if *ruleName == "" {
		return c.errorWithUsage(fmt.Errorf(("Name is mandatory.")))
	}

	return c.runWithSpinner("delete network rule", endpoint.String(), func(client *squarescale.Client) (string, error) {

		err := client.DeleteNetworkRule(*projectUUID, *serviceName, *ruleName)
		return "", err
	})

}

// Synopsis is part of cli.Command implementation.
func (c *NetworkRuleDeleteCommand) Synopsis() string {
	return "Delete a network rule"
}

// Help is part of cli.Command implementation.
func (c *NetworkRuleDeleteCommand) Help() string {
	helpText := `
usage: sqsc network-rule delete [options]

  Delete a network rule.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
