package command

import (
	"errors"
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
	projectUUID := c.flagSet.String("project-uuid", "", "uuid of the targeted project")
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if *projectUUID == "" {
		return c.errorWithUsage(errors.New("Project uuid is mandatory"))
	}

	return c.runWithSpinner("fetch database configuration", endpoint.String(), func(client *squarescale.Client) (string, error) {
		db, e := client.GetDBConfig(*projectUUID)
		if e != nil {
			return "", e
		}

		return c.FormatTable(fmt.Sprintf("DB enabled:\t%v\nDB Engine:\t%s\nDB Size:\t%s\nVersion:\t%s", db.Enabled, db.Engine, db.Size, db.Version), false), nil
	})
}

// Synopsis is part of cli.Command implementation.
func (c *DBShowCommand) Synopsis() string {
	return "Show database engine attached to project"
}

// Help is part of cli.Command implementation.
func (c *DBShowCommand) Help() string {
	helpText := `
usage: sqsc db show [options]

  Show database engine attached to project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
