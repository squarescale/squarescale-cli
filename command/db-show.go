package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// DBShowCommand is a cli.Command implementation for scaling the db of a project.
type DBShowCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *DBShowCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectNameArg := projectFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	err := validateProjectName(*projectNameArg)
	if err != nil {
		return c.errorWithUsage(err)
	}

	return c.runWithSpinner("fetch database configuration", *endpoint, func(client *squarescale.Client) (string, error) {
		enabled, engine, instance, e := client.GetDBConfig(*projectNameArg)
		if e != nil {
			return "", e
		}

		return c.FormatTable(fmt.Sprintf("DB enabled:\t%v\nDB Engine:\t%s\nDB Instance Class:\t%s", enabled, engine, instance), false), nil
	})
}

// Synopsis is part of cli.Command implementation.
func (c *DBShowCommand) Synopsis() string {
	return "Show the database engine options attached to a Squarescale project"
}

// Help is part of cli.Command implementation.
func (c *DBShowCommand) Help() string {
	helpText := `
usage: sqsc db show [options]

  Show options for the database engine attached to a Squarescale project.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
