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
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *DBSetCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	nowait := nowaitFlag(c.flagSet)
	projectNameArg := projectFlag(c.flagSet)
	dbEngineArg := dbEngineFlag(c.flagSet)
	dbInstanceArg := dbEngineInstanceFlag(c.flagSet)
	alwaysYes := yesFlag(c.flagSet)
	dbDisabledArg := disabledFlag(c.flagSet, "Disable database")
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	err := validateDBSetCommandArgs(*projectNameArg, *dbEngineArg, *dbInstanceArg, *dbDisabledArg)
	if err != nil {
		return c.errorWithUsage(err)
	}

	c.Ui.Warn(fmt.Sprintf("Changing cluster settings for project '%s' will cause a downtime.", *projectNameArg))
	ok, err := AskYesNo(c.Ui, alwaysYes, "Is this ok?", false)
	if err != nil {
		return c.error(err)
	} else if !ok {
		return c.cancelled()
	}

	var taskId int

	res := c.runWithSpinner("scale project database", endpoint.String(), func(client *squarescale.Client) (string, error) {
		var err error
		enabled, engine, instance, e := client.GetDBConfig(*projectNameArg)
		if e != nil {
			return "", e
		}

		if *dbDisabledArg && !enabled {
			*nowait = true
			return fmt.Sprintf("Database for project '%s' is already disabled", *projectNameArg), nil
		}

		if *dbEngineArg == engine && *dbInstanceArg == instance {
			*nowait = true
			return fmt.Sprintf("Database for project '%s' is already configured with these parameters", *projectNameArg), nil
		}

		if *dbEngineArg != "" {
			engine = *dbEngineArg
		}

		if *dbInstanceArg != "" {
			instance = *dbInstanceArg
		}

		enabled = !(*dbDisabledArg)
		taskId, err = client.ConfigDB(*projectNameArg, enabled, engine, instance)

		var msg string
		if *dbDisabledArg {
			msg = fmt.Sprintf("[#%d] Successfully disabled database for project '%s'", taskId, *projectNameArg)
		} else {
			msg = fmt.Sprintf(
				"[#%d] Successfully set database (engine = '%s', instance = '%s') for project '%s'",
				taskId, engine, instance, *projectNameArg)
		}

		return msg, err
	})
	if res != 0 {
		return res
	}

	if !*nowait {
		res = c.runWithSpinner("wait for database change", endpoint.String(), func(client *squarescale.Client) (string, error) {
			task, err := client.WaitTask(taskId)
			if err != nil {
				return "", err
			} else {
				return task.Status, nil
			}
		})
	}

	return res
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

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
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
