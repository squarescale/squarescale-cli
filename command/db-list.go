package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// DBListCommand is a cli.Command implementation to list the database engines and instances available for use.
type DBListCommand struct {
	Meta
}

// Run is part of cli.Command implementation.
func (c *DBListCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("db list", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	endpoint := endpointFlag(cmdFlags)
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	var msg string
	err := c.runWithSpinner("list available database engines and instances", *endpoint, func(token string) error {
		engines, err := squarescale.GetAvailableDBEngines(*endpoint, token)
		if err != nil {
			return err
		}

		msg += fmtDBListOutput("Available engines", engines)

		instances, err := squarescale.GetAvailableDBInstances(*endpoint, token)
		if err != nil {
			return err
		}

		msg += "\n\n"
		msg += fmtDBListOutput("Available instances", instances)
		return nil
	})

	if err != nil {
		return c.error(err)
	}

	return c.info(msg)
}

// Synopsis is part of cli.Command implementation.
func (c *DBListCommand) Synopsis() string {
	return "Lists all available database engines and instances"
}

// Help is part of cli.Command implementation.
func (c *DBListCommand) Help() string {
	helpText := `
usage: sqsc db list

  Lists all database engines and instances available for use
  in the Squarescale services.

Options:

  -endpoint="http://www.staging.sqsc.squarely.io" Squarescale endpoint
`
	return strings.TrimSpace(helpText)
}

func fmtDBListOutput(title string, lines []string) string {
	var res string
	for _, line := range lines {
		res += fmt.Sprintf("  %s\n", line)
	}

	res = strings.Trim(res, "\n")
	return fmt.Sprintf("%s\n\n%s", title, res)
}
