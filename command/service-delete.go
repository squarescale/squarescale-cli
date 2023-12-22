package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ServiceDeleteCommand is a cli.Command implementation for deleting a service aka Docker container of a project.
type ServiceDeleteCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *ServiceDeleteCommand) Run(args []string) int {
	// Parse flags
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	alwaysYes := yesFlag(cmd.flagSet)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	containerName, err := containerNameArg(cmd.flagSet, 0)
	if err != nil {
		return cmd.errorWithUsage(err)
	}

	if cmd.flagSet.NArg() > 4 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()[1:]))
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	cmd.Ui.Info("Are you sure you want to delete service " + containerName + "?")
	if *alwaysYes {
		cmd.Ui.Info("(approved from command line)")
	} else {
		res, err := cmd.Ui.Ask("y/N")
		if err != nil {
			return cmd.error(err)
		} else if res != "Y" && res != "y" {
			return cmd.cancelled()
		}
	}

	return cmd.runWithSpinner("deleting service", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		container, err := client.GetServiceInfo(UUID, containerName)
		if err != nil {
			return "", err
		}
		err = client.DeleteService(container)
		return "", err
	})
}

// Synopsis is part of cli.Command implementation.
func (cmd *ServiceDeleteCommand) Synopsis() string {
	return "Delete service aka Docker container from project."
}

// Help is part of cli.Command implementation.
func (cmd *ServiceDeleteCommand) Help() string {
	helpText := `
usage: sqsc service delete [options] <service_name>

  Delete service aka Docker container from project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
