package command

import (
	"flag"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ProjectListCommand is a cli.Command implementation for listing all Squarescale projects.
type ProjectListCommand struct {
	Meta
}

// Run is part of cli.Command implementation.
func (c *ProjectListCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("project list", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	endpoint := EndpointFlag(cmdFlags)
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	var msg string
	err := c.runWithSpinner("list projects", *endpoint, func(token string) error {
		projects, e := squarescale.ListProjects(*endpoint, token)
		if e != nil {
			return e
		}

		msg = strings.Join(projects, "\n")
		if len(projects) == 0 {
			msg = "No projects found"
		}

		return nil
	})

	if err != nil {
		return c.error(err)
	}

	return c.info(msg)
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

Options:

  -endpoint="http://www.staging.sqsc.squarely.io" Squarescale endpoint
`
	return strings.TrimSpace(helpText)
}
