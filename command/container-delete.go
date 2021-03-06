package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ContainerDeleteCommand is a cli.Command implementation for creating a SquareScale project.
type ContainerDeleteCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *ContainerDeleteCommand) Run(args []string) int {
	// Parse flags
	c.flagSet = newFlagSet(c, c.Ui)
	alwaysYes := yesFlag(c.flagSet)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := c.flagSet.String("project-uuid", "", "set the uuid of the project")
	projectName := c.flagSet.String("project-name", "", "set the name of the project")

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

	c.Ui.Info("Are you sure you want to delete " + containerName + "?")
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

	return c.runWithSpinner("deleting container", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		container, err := client.GetServicesInfo(UUID, containerName)
		if err != nil {
			return "", err
		}
		err = client.DeleteService(container)
		return "", err
	})
}

// Synopsis is part of cli.Command implementation.
func (c *ContainerDeleteCommand) Synopsis() string {
	return "Delete container from project."
}

// Help is part of cli.Command implementation.
func (c *ContainerDeleteCommand) Help() string {
	helpText := `
usage: sqsc container delete [options] <container_name>

  Delete container from project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
