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
func (c *RedisListCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := projectUUIDFlag(c.flagSet)
	projectName := projectNameFlag(c.flagSet)

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	return c.runWithSpinner("list redis", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
func (c *RedisListCommand) Synopsis() string {
	return "List Redis instances for a SquareScale projects"
}

// Help is part of cli.Command implementation.
func (c *RedisListCommand) Help() string {
	helpText := `
usage: sqsc redis list [options]

  List all Redis instances for a given SquareScale project.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
