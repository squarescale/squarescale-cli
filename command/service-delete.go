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
func (c *ServiceDeleteCommand) Run(args []string) int {
	// Parse flags
	c.flagSet = newFlagSet(c, c.Ui)
	alwaysYes := yesFlag(c.flagSet)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := projectUUIDFlag(c.flagSet)
	projectName := projectNameFlag(c.flagSet)

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	containerName, err := containerNameArg(c.flagSet, 0)
	if err != nil {
		return c.errorWithUsage(err)
	}

	if c.flagSet.NArg() > 4 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()[1:]))
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	c.Ui.Info("Are you sure you want to delete service " + containerName + "?")
	if *alwaysYes {
		c.Ui.Info("(approved from command line)")
	} else {
		res, err := c.Ui.Ask("y/N")
		if err != nil {
			return c.error(err)
		} else if res != "Y" && res != "y" {
			return c.cancelled()
		}
	}

	return c.runWithSpinner("deleting service", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
func (c *ServiceDeleteCommand) Synopsis() string {
	return "Delete service aka Docker container from project."
}

// Help is part of cli.Command implementation.
func (c *ServiceDeleteCommand) Help() string {
	helpText := `
usage: sqsc service delete [options] <service_name>

  Delete service aka Docker container from project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
