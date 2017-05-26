package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ContainerShowCommand is a cli.Command implementation for listing all Squarescale projects.
type ContainerShowCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *ContainerShowCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectArg := projectFlag(c.flagSet)
	containerArg := filterNameFlag(c.flagSet)
	typeArg := filterTypeFlag(c.flagSet)
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

		var msg string = "Type\tName\tStatus\tSize\tWeb\tPort\n"
		for _, c := range containers {
			if *containerArg != "" && *containerArg != c.ShortName {
				continue
			}
			if *typeArg != "" && *typeArg != c.Type {
				continue
			}
			st, _ := c.Status()
			msg += fmt.Sprintf("Type:     %s\n", c.Type)
			msg += fmt.Sprintf("Name:     %s\n", c.ShortName)
			msg += fmt.Sprintf("Status:   %s\n", st)
			msg += fmt.Sprintf("Size:     %d/%d\n", c.Running, c.Size)
			msg += fmt.Sprintf("Web:      %v\n", c.Web)
			msg += fmt.Sprintf("Web Port: %d\n", c.WebPort)
			if len(c.RefreshCallbacks) > 0 {
				msg += fmt.Sprintf("Refresh callbacks:\n")
				for _, url := range c.RefreshCallbacks {
					msg += fmt.Sprintf("  - %s\n", url)
				}
			}
			if len(c.BuildCallbacks) > 0 {
				msg += fmt.Sprintf("Rebuild callbacks:\n")
				for _, url := range c.BuildCallbacks {
					msg += fmt.Sprintf("  - %s\n", url)
				}
			}
		}

		if len(containers) == 0 {
			msg = "No containers found"
		}

		return c.FormatTable(msg, false), nil
	})
}

// Synopsis is part of cli.Command implementation.
func (c *ContainerShowCommand) Synopsis() string {
	return "Lists the containers of a Squarescale projects"
}

// Help is part of cli.Command implementation.
func (c *ContainerShowCommand) Help() string {
	helpText := `
usage: sqsc container show [options]

  List all containers of a given Squarescale project.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
