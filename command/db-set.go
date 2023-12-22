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
func (cmd *DBSetCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	nowait := nowaitFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)
	dbEngine := dbEngineFlag(cmd.flagSet)
	dbSize := dbSizeFlag(cmd.flagSet)
	dbVersion := dbVersionFlag(cmd.flagSet)
	dbBackupEnabled := dbBackupFlag(cmd.flagSet)
	dbDisabled := cmd.flagSet.Bool("disable", false, "Disable database")
	alwaysYes := yesFlag(cmd.flagSet)
	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if cmd.flagSet.NArg() > 0 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if *dbDisabled && (*dbEngine != "" || *dbSize != "" || *dbVersion != "") {
		return cmd.errorWithUsage(errors.New("Cannot specify engine or size or version when disabling database."))
	}

	if !*dbDisabled && (*dbEngine == "" && *dbSize == "" && *dbVersion == "" && !*dbBackupEnabled) {
		return cmd.errorWithUsage(errors.New("Size, engine and version are mandatory."))
	}

	var projectToShow string
	if *projectUUID == "" {
		projectToShow = *projectName
	} else {
		projectToShow = *projectUUID
	}

	cmd.Ui.Warn(fmt.Sprintf("Changing cluster settings for project '%s' will cause a downtime.", projectToShow))
	ok, err := AskYesNo(cmd.Ui, alwaysYes, "Is this ok?", false)
	if err != nil {
		return cmd.error(err)
	} else if !ok {
		return cmd.cancelled()
	}

	var UUID string

	res := cmd.runWithSpinner("scale project database", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		engine := checkFlag(db.Engine, dbEngine)
		size := checkFlag(db.Size, dbSize)
		version := checkFlag(db.Version, dbVersion)

		if *dbEngine == db.Engine && *dbSize == db.Size && *dbVersion == db.Version && !*dbDisabled == db.Enabled {
			return fmt.Sprintf("Database for project '%s' is already configured with these parameters", projectToShow), nil
		}

		payload := squarescale.JSONObject{
			"database": map[string]interface{}{"enabled": !*dbDisabled, "engine": engine, "size": size, "version": version, "backup_enabled": dbBackupEnabled},
		}

		_, err = client.ConfigDB(UUID, &payload)

		var msg string
		if *dbDisabled {
			msg = fmt.Sprintf("Successfully disabled database for project '%s'", projectToShow)
		} else {
			msg = fmt.Sprintf(
				"Successfully set database for project '%s': %s",
				projectToShow, db.String())
		}

		return msg, err
	})
	if res != 0 {
		return res
	}

	if !*nowait {
		res = cmd.runWithSpinner("wait for database change", endpoint.String(), func(client *squarescale.Client) (string, error) {
			projectStatus, err := client.WaitProject(UUID, 5)
			if err != nil {
				return projectStatus, err
			} else {
				return projectStatus, nil
			}
		})
	}

	return res
}

// Synopsis is part of cli.Command implementation.
func (cmd *DBSetCommand) Synopsis() string {
	return "Set and scale up/down the database engine attached to project"
}

// Help is part of cli.Command implementation.
func (cmd *DBSetCommand) Help() string {
	helpText := `
usage: sqsc db set [options]

  Set and scale up/down the database engine attached to project.
  The list of database engines and instances available for use can be accessed
  using the 'db list' command.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}

// validateDBSetCommandArgs ensures that the following predicate is satisfied:
// - 'disabled' is true and (dbEngine and dbInstance are both empty)
// - 'disabled' is false and (dbEngine or dbInstance is not empty)
func validateDBSetCommandArgs(projectUUID string, db squarescale.DbConfig) error {

	return nil
}

func checkFlag(dbFlag string, flag *string) string {
	var result string
	if flag == nil {
		result = dbFlag
	} else {
		result = *flag
	}
	return result
}
