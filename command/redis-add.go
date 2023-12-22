package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// RedisAddCommand is a struct to define a Redis
type RedisAddCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *RedisAddCommand) Run(args []string) int {
	// Parse flags
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)

	wantedRedisName := cmd.flagSet.String("name", "", "Redis name")

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if *wantedRedisName == "" && cmd.flagSet.Arg(0) != "" {
		name := cmd.flagSet.Arg(0)
		wantedRedisName = &name
	}

	if cmd.flagSet.NArg() > 1 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()[1:]))
	}

	//check rules
	if *wantedRedisName == "" {
		return cmd.errorWithUsage(fmt.Errorf(("Redis name is mandatory. Please, choose a redis name.")))
	}

	res := cmd.runWithSpinner("add redis", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		//add function
		err = client.AddRedis(UUID, *wantedRedisName)

		return fmt.Sprintf("Successfully added redis '%s'", *wantedRedisName), err
	})

	if res != 0 {
		return res
	}

	return res

}

// Synopsis is part of cli.Command implementation.
func (cmd *RedisAddCommand) Synopsis() string {
	return "Add a new Redis in a SquareScale project"
}

// Help is part of cli.Command implementation.
func (cmd *RedisAddCommand) Help() string {
	helpText := `
usage: sqsc redis add [options] <redis_name>

  Add a new redis using the provided redis name (name is mandatory).

`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
