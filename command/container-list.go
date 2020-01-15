package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ContainerListCommand is a cli.Command implementation for listing all Squarescale projects.
type ContainerListCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *ContainerListCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectArg := projectFlag(c.flagSet)
	containerArg := containerFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if err := validateProjectName(*projectArg); err != nil {
		return c.errorWithUsage(err)
	}

	return c.runWithSpinner("list containers", endpoint.String(), func(client *squarescale.Client) (string, error) {
		containers, err := client.GetContainers(*projectArg)
		if err != nil {
			return "", err
		}

		var msg string = "Name\tSize\tWeb\tPort\n"
		for _, c := range containers {
			if *containerArg != "" && *containerArg != c.ShortName {
				continue
			}
			msg += fmt.Sprintf("%s\t%d/%d\t%v\t%d\n", c.ShortName, c.Running, c.Size, c.Web, c.WebPort)
		}

		if len(containers) == 0 {
			msg = "No containers found"
		}

		return c.FormatTable(msg, true), nil
	})
}

// Synopsis is part of cli.Command implementation.
func (c *ContainerListCommand) Synopsis() string {
	return "Lists the containers of a Squarescale projects"
}

// Help is part of cli.Command implementation.
func (c *ContainerListCommand) Help() string {
	helpText := `
usage: sqsc container list [options]

  List all containers of a given Squarescale project.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
