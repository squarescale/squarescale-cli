package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// DBScaleCommand is a cli.Command implementation for scaling the db of a project.
type DBScaleCommand struct {
	Meta
}

// Run is part of cli.Command implementation.
func (c *DBScaleCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("db scale", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	endpoint := EndpointFlag(cmdFlags)
	project := ProjectFlag(cmdFlags)
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	err := validateArgs(*project)
	if err != nil {
		return c.errorWithUsage(err, c.Help())
	}

	// check for a database instance size
	var dbInstance string
	args = cmdFlags.Args()
	switch len(args) {
	case 0:
	case 1:
		dbInstance = args[0]
	default:
		return c.errorWithUsage(errors.New("Too many command line arguments"), c.Help())
	}

	err = validateDbInstance(dbInstance)
	if err != nil {
		return c.errorWithUsage(err, c.Help())
	}

	c.Ui.Warn(fmt.Sprintf("Changing cluster settings for project '%s' will cause a downtime. Is this ok?", *project))
	_, err = c.Ui.Ask("Enter to accept, Ctrl-c to cancel:")
	if err != nil {
		return c.error(err)
	}

	err = c.runWithSpinner("scale project database", *endpoint, func(token string) error {
		return squarescale.ScaleDB(*endpoint, token, *project, dbInstance)
	})

	if err != nil {
		return c.error(err)
	}

	return c.info(fmt.Sprintf("Successfully scaled database to '%s' instance for project '%s'", dbInstance, *project))
}

// Synopsis is part of cli.Command implementation.
func (c *DBScaleCommand) Synopsis() string {
	return "Scale up/down the database of a Squarescale project"
}

// Help is part of cli.Command implementation.
func (c *DBScaleCommand) Help() string {
	helpText := `
usage: sqsc scale-db [options] <micro|small|medium>

  Scale the database of the specified Squarescale project.

Options:

  -endpoint="http://www.staging.sqsc.squarely.io" Squarescale endpoint
  -project=""                                     Squarescale project name
`
	return strings.TrimSpace(helpText)
}

func validateDbInstance(instance string) error {
	switch instance {
	case "micro", "small", "medium":
		return nil
	default:
		return errors.New("Invalid value for database instance")
	}
}
