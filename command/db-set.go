package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// DBSetCommand is a cli.Command implementation for scaling the db of a project.
type DBSetCommand struct {
	Meta
}

// Run is part of cli.Command implementation.
func (c *DBSetCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("db scale", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	endpoint := endpointFlag(cmdFlags)
	projectNameArg := projectFlag(cmdFlags)
	dbEngineArg := dbEngineFlag(cmdFlags)
	dbInstanceArg := dbEngineInstance(cmdFlags)
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	err := validateDBSetCommandArgs(*projectNameArg, *dbEngineArg, *dbInstanceArg)
	if err != nil {
		return c.errorWithUsage(err, c.Help())
	}

	c.Ui.Warn(fmt.Sprintf("Changing cluster settings for project '%s' will cause a downtime. Is this ok?", *projectNameArg))
	_, err = c.Ui.Ask("Enter to accept, Ctrl-c to cancel:")
	if err != nil {
		return c.error(err)
	}

	var msg string
	err = c.runWithSpinner("scale project database", *endpoint, func(token string) error {
		engine, instance, e := squarescale.GetDBConfig(*endpoint, token, *projectNameArg)
		if e != nil {
			return e
		}

		if *dbEngineArg == engine && *dbInstanceArg == instance {
			msg = fmt.Sprintf("Database for project '%s' is already configured with these parameters", *projectNameArg)
			return nil
		}

		if *dbEngineArg != "" {
			engine = *dbEngineArg
		}

		if *dbInstanceArg != "" {
			instance = "db.t2." + *dbInstanceArg
		}

		msg = fmt.Sprintf("Successfully set database (engine = '%s', instance = '%s') for project '%s'", engine, instance, *projectNameArg)
		return squarescale.ConfigDB(*endpoint, token, *projectNameArg, engine, instance)
	})

	if err != nil {
		return c.error(err)
	}

	return c.info(msg)
}

// Synopsis is part of cli.Command implementation.
func (c *DBSetCommand) Synopsis() string {
	return "Set and scale up/down the database engine attached to a Squarescale project"
}

// Help is part of cli.Command implementation.
func (c *DBSetCommand) Help() string {
	helpText := `
usage: sqsc db set [options]

  Set and scale up/down the database engine attached to a Squarescale project.

Options:

  -endpoint="http://www.staging.sqsc.squarely.io" Squarescale endpoint
  -project=""                                     Squarescale project name
  -engine=<mysql|postgres|aurora|mariadb>         Database engine
  -instance=<micro|small|medium>                  Database instance size
`
	return strings.TrimSpace(helpText)
}

func validateDBSetCommandArgs(project, dbEngine, dbInstance string) error {
	if err := validateProjectName(project); err != nil {
		return err
	}

	if err := validateDBEngine(dbEngine); err != nil {
		return err
	}

	if err := validateDBInstance(dbInstance); err != nil {
		return err
	}

	if dbEngine == "" && dbInstance == "" {
		return errors.New("Instance and engine cannot be both empty.")
	}

	return nil
}
