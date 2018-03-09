package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/kballard/go-shellquote"
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

		var msg string
		for _, co := range containers {
			if msg != "" {
				msg += "\n\n-----------\n\n\n"
			}
			if *containerArg != "" && *containerArg != co.ShortName {
				continue
			}
			if *typeArg != "" && *typeArg != co.Type {
				continue
			}
			st, _ := co.Status()
			if co.NomadStatus.Service.PlacementFailure {
				st += ", service error"
			}
			tbl := ""
			tbl += fmt.Sprintf("Type:\t%s\n", co.Type)
			tbl += fmt.Sprintf("Name:\t%s\n", co.ShortName)
			tbl += fmt.Sprintf("Status:\t%s\n", st)
			tbl += fmt.Sprintf("Size:\t%d/%d\n", co.Running, co.Size)
			tbl += fmt.Sprintf("Run Command:\t%s\n", shellquote.Join(co.RunCommand...))
			tbl += fmt.Sprintf("Web:\t%v\n", co.Web)
			tbl += fmt.Sprintf("Web Port:\t%d\n", co.WebPort)
			tbl += fmt.Sprintf("Memory limit:\t%d MB\n", co.Limits.Memory)
			tbl += fmt.Sprintf("CPU limit:\t%d MHz\n", co.Limits.CPU)
			tbl += fmt.Sprintf("Network limit:\t%d Mbps\n", co.Limits.Net)
			msg += c.FormatTable(tbl, false)
			msg += "\n\n"

			if len(co.RefreshCallbacks) > 0 {
				msg += fmt.Sprintf("Refresh callbacks:\n")
				for _, url := range co.RefreshCallbacks {
					msg += fmt.Sprintf("  - %s\n", url)
				}
			}
			if len(co.BuildCallbacks) > 0 {
				msg += fmt.Sprintf("Rebuild callbacks:\n")
				for _, url := range co.BuildCallbacks {
					msg += fmt.Sprintf("  - %s\n", url)
				}
			}

			placementErrors := false

			if co.NomadStatus.Service.PlacementFailure {
				placementErrors = true
				msg += fmt.Sprintf("Service Errors:\n")
				for _, e := range co.NomadStatus.Service.Errors("service") {
					msg += fmt.Sprintf("  - %s\n", e)
				}
			}

			if placementErrors {
				msg += fmt.Sprintf("\nThese errors means that your cluster doesn't have enough resources to fit your containers. To solve these errors, you can try to reduce the resource required by your service or increase the cluster size.\n\n")
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
