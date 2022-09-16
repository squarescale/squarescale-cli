package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ContainerShowCommand is a cli.Command implementation for listing all projects.
type ContainerShowCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *ContainerShowCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := c.flagSet.String("project-uuid", "", "set the uuid of the project")
	projectName := c.flagSet.String("project-name", "", "set the name of the project")
	containerArg := filterNameFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	return c.runWithSpinner("showing Docker container", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		var msg string
		for _, co := range containers {
			if *containerArg != "" && *containerArg != co.Name {
				continue
			}
			if msg != "" {
				msg += "\n-----------\n\n"
			}
			tbl := ""
			tbl += fmt.Sprintf("Name:\t%s\n", co.Name)
			tbl += fmt.Sprintf("Size:\t%d/%d\n", co.Running, co.Size)
			tbl += fmt.Sprintf("Run Command:\t%s\n", co.RunCommand)
			tbl += fmt.Sprintf("Web Port:\t%d\n", co.WebPort)
			tbl += fmt.Sprintf("Memory limit:\t%d MB\n", co.Limits.Memory)
			tbl += fmt.Sprintf("CPU limit:\t%d MHz\n", co.Limits.CPU)
			tbl += fmt.Sprintf("Network limit:\t%d Mbps\n", co.Limits.Net)
			msg += c.FormatTable(tbl, false)
			msg += "\n"

			if len(co.RefreshCallbacks) > 0 {
				msg += fmt.Sprintf("Refresh callbacks:\n")
				for _, url := range co.RefreshCallbacks {
					msg += fmt.Sprintf("  - %s\n", url)
				}
			}
		}

		if len(containers) == 0 {
			msg = "No containers found"
		}

		return msg, nil
	})
}

// Synopsis is part of cli.Command implementation.
func (c *ContainerShowCommand) Synopsis() string {
	return "Show Docker container of project"
}

// Help is part of cli.Command implementation.
func (c *ContainerShowCommand) Help() string {
	helpText := `
usage: sqsc container show [options]

  Show Docker container of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
