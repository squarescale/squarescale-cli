package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ServiceListCommand is a cli.Command implementation for listing all services aka Docker containers of project.
type ServiceListCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *ServiceListCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)
	containerArg := serviceFlag(cmd.flagSet)
	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if cmd.flagSet.NArg() > 0 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	return cmd.runWithSpinner("listing services", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		containers, err := client.GetServices(UUID)
		if err != nil {
			return "", err
		}

		var msg string = "Name\tSize\tPort\tScheduling groups\n"
		for _, c := range containers {
			if *containerArg != "" && *containerArg != c.Name {
				continue
			}

			var schedulingGroups []string
			for _, schedulingGroup := range c.SchedulingGroups {
				schedulingGroups = append(schedulingGroups, schedulingGroup.Name)
			}

			msg += fmt.Sprintf("%s\t%d/%d\t%d\t%s\n", c.Name, c.Running, c.Size, c.WebPort, strings.Join(schedulingGroups[:], "\n\t\t\t"))
		}

		if len(containers) == 0 {
			msg = "No service found"
		}

		return cmd.FormatTable(msg, true), nil
	})
}

// Synopsis is part of cli.Command implementation.
func (cmd *ServiceListCommand) Synopsis() string {
	return "List services aka Docker containers of project"
}

// Help is part of cli.Command implementation.
func (cmd *ServiceListCommand) Help() string {
	helpText := `
usage: sqsc service list [options]

  List services aka Docker containers of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
