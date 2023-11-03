package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// RedisListCommand is a cli.Command implementation for listing all redis.
type RedisListCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *RedisListCommand) Run(args []string) int {
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

	return cmd.runWithSpinner("list redis", endpoint.String(), func(client *squarescale.Client) (string, error) {
		var UUID string
		var err error
		var redisList []squarescale.RedisDbConfig
		if *projectUUID == "" {
			UUID, err = client.ProjectByName(*projectName)
			if err != nil {
				return "", err
			}
		} else {
			UUID = *projectUUID
		}

		redisList, err = client.GetRedis(UUID)
		if err != nil {
			return "", err
		}

		var msg string = "Name\n"
		for _, r := range redisList {
			msg += fmt.Sprintf("%s\n", r.Name)
		}

		if len(redisList) == 0 {
			msg = "No redis found"
		}

		return msg, nil

	})

}

// Synopsis is part of cli.Command implementation.
func (cmd *RedisListCommand) Synopsis() string {
	return "List Redis instances for a SquareScale projects"
}

// Help is part of cli.Command implementation.
func (cmd *RedisListCommand) Help() string {
	helpText := `
usage: sqsc redis list [options]

  List all Redis instances for a given SquareScale project.

`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
