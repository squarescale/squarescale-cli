package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ServiceScheduleCommand is a cli.Command implementation for top level `sqsc servic` command.
type ServiceScheduleCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *ServiceScheduleCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := c.flagSet.String("project-uuid", "", "set the uuid of the project")

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	serviceName, err := serviceNameArg(c.flagSet, 0)
	if err != nil {
		return c.errorWithUsage(err)
	}

	if c.flagSet.NArg() > 4 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()[1:]))
	}

	return c.runWithSpinner("scheduling service", endpoint.String(), func(client *squarescale.Client) (string, error) {
		uuid := *projectUUID
		err := client.ScheduleService(uuid, serviceName)
		return "", err
	})

}

// Synopsis is part of cli.Command implementation.
func (c *ServiceScheduleCommand) Synopsis() string {
	return "Schedule a service"
}

// Help is part of cli.Command implementation.
func (c *ServiceScheduleCommand) Help() string {
	helpText := `
usage: sqsc service schedule [options] <service_name>

  Schedule a service.
`
	return strings.TrimSpace(helpText)
}
