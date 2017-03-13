package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ProjectListCommand is a cli.Command implementation for listing all Squarescale projects.
type ProjectListCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *ProjectListCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	return c.runWithSpinner("list projects", *endpoint, func(client *squarescale.Client) (string, error) {
		projects, err := client.ListProjects()
		if err != nil {
			return "", err
		}

		var msg string = "Name\tStatus\tSize\n"
		for _, p := range projects {
			msg += fmt.Sprintf("%s\t%s\t%d/%d\n", p.Name, p.InfraStatus, p.NomadNodesReady, p.ClusterSize)
		}

		if len(projects) == 0 {
			msg = "No projects found"
		}

		return FormatTable(msg, true), nil
	})
}

// Synopsis is part of cli.Command implementation.
func (c *ProjectListCommand) Synopsis() string {
	return "Lists the Squarescale projects"
}

// Help is part of cli.Command implementation.
func (c *ProjectListCommand) Help() string {
	helpText := `
usage: sqsc project list [options]

  Lists all Squarescale projects attached to the authenticated account.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
