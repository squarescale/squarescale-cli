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
func (cmd *ServiceScheduleCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	serviceName, err := serviceNameArg(cmd.flagSet, 0)
	if err != nil {
		return cmd.errorWithUsage(err)
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if cmd.flagSet.NArg() > 4 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()[1:]))
	}

	return cmd.runWithSpinner("scheduling service", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
func (cmd *ServiceScheduleCommand) Synopsis() string {
	return "Schedule a service aka Docker container"
}

// Help is part of cli.Command implementation.
func (cmd *ServiceScheduleCommand) Help() string {
	helpText := `
usage: sqsc service schedule [options] <service_name>

  Schedule a service aka Docker container.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
