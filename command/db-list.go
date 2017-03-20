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
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *DBListCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	return c.runWithSpinner("list available database engines and instances", *endpoint, func(client *squarescale.Client) (string, error) {
		engines, err := client.GetAvailableDBEngines()
		if err != nil {
			return "", err
		}

		msg := fmtDBListOutput("Available engines", engines)

		instances, err := client.GetAvailableDBInstances()
		if err != nil {
			return "", err
		}

		msg += "\n\n"
		msg += fmtDBListOutput("Available instances", instances)
		return msg, nil
	})
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

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}

func fmtDBListOutput(title string, lines []string) string {
	var res string
	for _, line := range lines {
		res += fmt.Sprintf("  %s\n", line)
	}

	res = strings.Trim(res, "\n")
	return fmt.Sprintf("%s\n\n%s", title, res)
}
