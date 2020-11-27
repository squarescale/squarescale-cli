package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ContainerListCommand is a cli.Command implementation for listing all containers of project.
type ContainerListCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *ContainerListCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := c.flagSet.String("project-uuid", "", "set the uuid of the project")
	projectName := c.flagSet.String("project-name", "", "set the name of the project")
	containerArg := containerFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	return c.runWithSpinner("list containers", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		var msg string = "Name\tSize\tPort\n"
		for _, c := range containers {
			if *containerArg != "" && *containerArg != c.Name {
				continue
			}
			msg += fmt.Sprintf("%s\t%d/%d\t%d\n", c.Name, c.Running, c.Size, c.WebPort)
		}

		if len(containers) == 0 {
			msg = "No containers found"
		}

		return c.FormatTable(msg, true), nil
	})
}

// Synopsis is part of cli.Command implementation.
func (c *ContainerListCommand) Synopsis() string {
	return "List containers of project"
}

// Help is part of cli.Command implementation.
func (c *ContainerListCommand) Help() string {
	helpText := `
usage: sqsc container list [options]

  List containers of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
