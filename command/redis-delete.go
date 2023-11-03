package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// RedisDeleteCommand is a cli.Command implementation for creating a SquareScale project.
type RedisDeleteCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *RedisDeleteCommand) Run(args []string) int {
	// Parse flags
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	alwaysYes := yesFlag(cmd.flagSet)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	redisName, err := redisNameArg(cmd.flagSet, 0)
	if err != nil {
		return cmd.errorWithUsage(err)
	}

	if cmd.flagSet.NArg() > 1 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()[1:]))
	}

	cmd.Ui.Info("Are you sure you want to delete " + redisName + "?")
	if *alwaysYes {
		cmd.Ui.Info("(approved from command line)")
	} else {
		res, err := cmd.Ui.Ask("y/N")
		if err != nil {
			return cmd.error(err)
		} else if res != "Y" && res != "y" {
			return cmd.cancelled()
		}
	}

	res := cmd.runWithSpinner("deleting redis", endpoint.String(), func(client *squarescale.Client) (string, error) {
		var UUID string
		var err error
		var projectToShow string
		if *projectUUID == "" {
			UUID, err = client.ProjectByName(*projectName)
			if err != nil {
				return "", err
			}
		} else {
			UUID = *projectUUID
		}

		fmt.Printf("Delete redis `%s` on project `%s`\n", redisName, projectToShow)
		err = client.DeleteRedis(UUID, redisName) //insérer la fonction dans redis
		return "", err
	})
	if res != 0 {
		return res
	}

	return 0
}

// Synopsis is part of cli.Command implementation.
func (cmd *RedisDeleteCommand) Synopsis() string {
	return "Delete a redis on the project."
}

// Help is part of cli.Command implementation.
func (cmd *RedisDeleteCommand) Help() string {
	helpText := `
usage: sqsc redis delete [options] <redis_name>

  Delete the redis from the project.

`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
