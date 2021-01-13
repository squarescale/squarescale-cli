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
	projectUUID := c.flagSet.String("project-uuid", "", "uuid of the targeted project")
	projectName := c.flagSet.String("project-name", "", "set the name of the project")
	dbEngine := c.flagSet.String("engine", "", "Database engine")
	dbSize := c.flagSet.String("size", "", "Database size")
	dbVersion := c.flagSet.String("engine-version", "", "Database version")
	dbDisabled := c.flagSet.Bool("disable", false, "Disable database")
	alwaysYes := yesFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if *dbDisabled && (*dbEngine != "" || *dbSize != "" || *dbVersion != "") {
		return c.errorWithUsage(errors.New("Cannot specify engine or size or version when disabling database."))
	}

	if !*dbDisabled && (*dbEngine == "" && *dbSize == "" && *dbVersion == "") {
		return c.errorWithUsage(errors.New("Size, engine and version are mandatory."))
	}

	var projectToShow string
	if *projectUUID == "" {
		projectToShow = *projectName
	} else {
		projectToShow = *projectUUID
	}

	c.Ui.Warn(fmt.Sprintf("Changing cluster settings for project '%s' will cause a downtime.", projectToShow))
	ok, err := AskYesNo(c.Ui, alwaysYes, "Is this ok?", false)
	if err != nil {
		return c.error(err)
	} else if !ok {
		return c.cancelled()
	}

	var taskId int

	res := c.runWithSpinner("scale project database", endpoint.String(), func(client *squarescale.Client) (string, error) {
		var UUID string
		var err error
		if *projectUUID == "" {
			UUID, err = client.ProjectByName(*projectName)
			if err != nil {
				return "", err
			}
		} else {
			UUID = *projectUUID
		}

		db, e := client.GetDBConfig(UUID)
		if e != nil {
			return "", e
		}

		if *dbDisabled && !db.Enabled {
			return fmt.Sprintf("Database for project '%s' is already disabled", projectToShow), nil
		}

		if *dbEngine == db.Engine && *dbSize == db.Size && *dbVersion == db.Version && !*dbDisabled == db.Enabled {
			return fmt.Sprintf("Database for project '%s' is already configured with these parameters", projectToShow), nil
		}

		payload := squarescale.JSONObject{
			"database": map[string]interface{}{"enabled": !*dbDisabled, "engine": *dbEngine, "size": *dbSize, "version": *dbVersion},
		}

		taskId, err = client.ConfigDB(UUID, &payload)

		var msg string
		if *dbDisabled {
			msg = fmt.Sprintf("[#%d] Successfully disabled database for project '%s'", taskId, projectToShow)
		} else {
			msg = fmt.Sprintf(
				"[#%d] Successfully set database for project '%s': %s",
				taskId, projectToShow, db.String())
		}

		return msg, err
	})
	if res != 0 {
		return res
	}

	return res
}

// Synopsis is part of cli.Command implementation.
func (c *DBSetCommand) Synopsis() string {
	return "Set and scale up/down the database engine attached to project"
}

// Help is part of cli.Command implementation.
func (c *DBSetCommand) Help() string {
	helpText := `
usage: sqsc db set [options]

  Set and scale up/down the database engine attached to project.
  The list of database engines and instances available for use can be accessed
  using the 'db list' command.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}

// validateDBSetCommandArgs ensures that the following predicate is satisfied:
// - 'disabled' is true and (dbEngine and dbInstance are both empty)
// - 'disabled' is false and (dbEngine or dbInstance is not empty)
func validateDBSetCommandArgs(projectUUID string, db squarescale.DbConfig) error {

	return nil
}
