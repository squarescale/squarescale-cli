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
func (r *RedisAddCommand) Run(args []string) int {
	// Parse flags
	r.flagSet = newFlagSet(r, r.Ui)
	endpoint := endpointFlag(r.flagSet)
	projectUUID := r.flagSet.String("project-uuid", "", "set the uuid of the project")
	projectName := r.flagSet.String("project-name", "", "set the name of the project")

	wantedRedisName := r.flagSet.String("name", "", "Redis name")

	if err := r.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return r.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if *wantedRedisName == "" && r.flagSet.Arg(0) != "" {
		name := r.flagSet.Arg(0)
		wantedRedisName = &name
	}

	if r.flagSet.NArg() > 1 {
		return r.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", r.flagSet.Args()[1:]))
	}

	//check rules
	if *wantedRedisName == "" {
		return r.errorWithUsage(fmt.Errorf(("Redis name is mandatory. Please, choose a redis name.")))
	}

	res := r.runWithSpinner("add redis", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
func (r *RedisAddCommand) Synopsis() string {
	return "Add a new Redis in a SquareScale project"
}

// Help is part of cli.Command implementation.
func (r *RedisAddCommand) Help() string {
	helpText := `
usage: sqsc redis add [options] <redis_name>

  Add a new redis using the provided redis name (name is mandatory).

`
	return strings.TrimSpace(helpText + optionsFromFlags(r.flagSet))
}
