package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ContainerDeleteCommand is a cli.Command implementation for creating a Squarescale project.
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
	projectArg := projectFlag(c.flagSet)
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

	if err := validateProjectName(*projectArg); err != nil {
		return c.errorWithUsage(err)
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
		projectName := *projectArg
		container, err := client.GetContainerInfo(projectName, containerName)
		if err != nil {
			return "", err
		}
		err = client.DeleteContainer(container)
		return "", err
	})
}

// Synopsis is part of cli.Command implementation.
func (c *ContainerDeleteCommand) Synopsis() string {
	return "Delete a container from the project."
}

// Help is part of cli.Command implementation.
func (c *ContainerDeleteCommand) Help() string {
	helpText := `
usage: sqsc container delete [options] <container_name>

  Delete the container from the project.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
