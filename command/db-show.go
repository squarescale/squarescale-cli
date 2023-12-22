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
func (cmd *DBShowCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)
	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if cmd.flagSet.NArg() > 0 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	return cmd.runWithSpinner("fetch database configuration", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		return cmd.FormatTable(
			fmt.Sprintf(
				"DB Enabled:\t%v\nDB Engine:\t%s\nDB Size:\t%s\nDB Version:\t%s\nDB Backup Enabled:\t%v\nDB Backup Retention:\t%v",
				db.Enabled, db.Engine, db.Size, db.Version, db.BackupEnabled, db.BackupRetention,
			), false),
			nil
	})
}

// Synopsis is part of cli.Command implementation.
func (cmd *DBShowCommand) Synopsis() string {
	return "Show database engine attached to project"
}

// Help is part of cli.Command implementation.
func (cmd *DBShowCommand) Help() string {
	helpText := `
usage: sqsc db show [options]

  Show database engine attached to project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
