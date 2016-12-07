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
	cmdFlags := flag.NewFlagSet("db set", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	endpoint := endpointFlag(cmdFlags)
	projectNameArg := projectFlag(cmdFlags)
	dbEngineArg := dbEngineFlag(cmdFlags)
	dbInstanceArg := dbEngineInstanceFlag(cmdFlags)
	dbDisabledArg := dbDisabledFlag(cmdFlags)
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	err := validateDBSetCommandArgs(*projectNameArg, *dbEngineArg, *dbInstanceArg, *dbDisabledArg)
	if err != nil {
		return c.errorWithUsage(err)
	}

	c.Ui.Warn(fmt.Sprintf("Changing cluster settings for project '%s' will cause a downtime. Is this ok?", *projectNameArg))
	_, err = c.Ui.Ask("Enter to accept, Ctrl-c to cancel:")
	if err != nil {
		return c.error(err)
	}

	var msg string
	err = c.runWithSpinner("scale project database", *endpoint, func(token string) error {
		enabled, engine, instance, e := squarescale.GetDBConfig(*endpoint, token, *projectNameArg)
		if e != nil {
			return e
		}

		if *dbDisabledArg && !enabled {
			msg = fmt.Sprintf("Database for project '%s' is already disabled", *projectNameArg)
			return nil
		}

		if *dbEngineArg == engine && *dbInstanceArg == instance {
			msg = fmt.Sprintf("Database for project '%s' is already configured with these parameters", *projectNameArg)
			return nil
		}

		if *dbEngineArg != "" {
			engine = *dbEngineArg
		}

		if *dbInstanceArg != "" {
			instance = *dbInstanceArg
		}

		if *dbDisabledArg {
			msg = fmt.Sprintf("Successfully disabled database for project '%s'", *projectNameArg)
		} else {
			msg = fmt.Sprintf(
				"Successfully set database (engine = '%s', instance = '%s') for project '%s'",
				engine, instance, *projectNameArg)
		}

		enabled = !(*dbDisabledArg)
		return squarescale.ConfigDB(*endpoint, token, *projectNameArg, enabled, engine, instance)
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
  The list of database engines and instances available for use can be accessed
  using the 'db list' command.

Options:

  -endpoint="http://www.staging.sqsc.squarely.io" Squarescale endpoint
  -project=""                                     Squarescale project name
  -engine=""                                      Database engine
  -instance=""                                    Database instance size
  -disabled                                       Disable database
`
	return strings.TrimSpace(helpText)
}

// validateDBSetCommandArgs ensures that the following predicate is satisfied:
// - 'disabled' is true and (dbEngine and dbInstance are both empty)
// - 'disabled' is false and (dbEngine or dbInstance is not empty)
func validateDBSetCommandArgs(project, dbEngine, dbInstance string, disabled bool) error {
	if err := validateProjectName(project); err != nil {
		return err
	}

	if disabled && (dbEngine != "" || dbInstance != "") {
		return errors.New("Cannot specify engine or instance when disabling database.")
	}

	if !disabled && (dbEngine == "" && dbInstance == "") {
		return errors.New("Instance and engine cannot be both empty.")
	}

	return nil
}
