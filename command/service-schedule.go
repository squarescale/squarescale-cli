package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ServiceScheduleCommand is a cli.Command implementation for top level `sqsc service` command.
type ServiceScheduleCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *ServiceScheduleCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := projectUUIDFlag(c.flagSet)
	projectName := projectNameFlag(c.flagSet)

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	serviceName, err := serviceNameArg(c.flagSet, 0)
	if err != nil {
		return c.errorWithUsage(err)
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if c.flagSet.NArg() > 4 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()[1:]))
	}

	return c.runWithSpinner("scheduling service", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		err = client.ScheduleService(UUID, serviceName)
		return "", err
	})

}

// Synopsis is part of cli.Command implementation.
func (c *ServiceScheduleCommand) Synopsis() string {
	return "Schedule a service aka Docker container"
}

// Help is part of cli.Command implementation.
func (c *ServiceScheduleCommand) Help() string {
	helpText := `
usage: sqsc service schedule [options] <service_name>

  Schedule a service aka Docker container.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
