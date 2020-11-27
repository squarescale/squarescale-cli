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
func (r *RedisListCommand) Run(args []string) int {
	r.flagSet = newFlagSet(r, r.Ui)
	endpoint := endpointFlag(r.flagSet)
	projectUUID := r.flagSet.String("project-uuid", "", "set the uuid of the project")
	projectName := r.flagSet.String("project-name", "", "set the name of the project")

	if err := r.flagSet.Parse(args); err != nil {
		return 1
	}

	if r.flagSet.NArg() > 0 {
		return r.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", r.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return r.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	return r.runWithSpinner("list redis", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
func (r *RedisListCommand) Synopsis() string {
	return "List Redis instances for a Squarescale projects"
}

// Help is part of cli.Command implementation.
func (r *RedisListCommand) Help() string {
	helpText := `
usage: sqsc redis list [options]

  List all Redis instances for a given Squarescale project.

`
	return strings.TrimSpace(helpText + optionsFromFlags(r.flagSet))
}
