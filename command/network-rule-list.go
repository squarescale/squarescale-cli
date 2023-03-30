package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ContainerListCommand is a cli.Command implementation for listing all network rules of service container of project.
type NetworkRuleListCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *NetworkRuleListCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := projectUUIDFlag(c.flagSet)
	projectName := projectNameFlag(c.flagSet)
	serviceName := networkServiceNameFlag(c.flagSet)

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if *serviceName == "" {
		return c.errorWithUsage(fmt.Errorf(("Service name is mandatory.")))
	}

	return c.runWithSpinner("list service container network rules", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		rules, err := client.ListNetworkRules(UUID, *serviceName)
		if err != nil {
			return "", err
		}

		var msg string = "Name\tInt. Proto/Port\tExt. Proto/Port\tDomain\tPath. Prefix\n"
		for _, c := range rules {
			msg += fmt.Sprintf("%s\t%s/%d\t%s/%d\t%s\t%s\n", c.Name, c.InternalProtocol, c.InternalPort, c.ExternalProtocol, c.ExternalPort, c.DomainExpression, c.PathPrefix)
		}

		if len(rules) == 0 {
			msg = "No network rules found"
		}

		return c.FormatTable(msg, true), nil
	})
}

// Synopsis is part of cli.Command implementation.
func (c *NetworkRuleListCommand) Synopsis() string {
	return "List service container network rules of project"
}

// Help is part of cli.Command implementation.
func (c *NetworkRuleListCommand) Help() string {
	helpText := `
usage: sqsc network-rule list [options]

  List service container network rules of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
