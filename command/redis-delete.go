package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// RedisDeleteCommand is a cli.Command implementation for creating a Squarescale project.
type RedisDeleteCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *RedisDeleteCommand) Run(args []string) int {
	// Parse flags
	c.flagSet = newFlagSet(c, c.Ui)
	alwaysYes := yesFlag(c.flagSet)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := c.flagSet.String("project-uuid", "", "set the uuid of the project")
	projectName := c.flagSet.String("project-name", "", "set the name of the project")

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	redisName, err := redisNameArg(c.flagSet, 0)
	if err != nil {
		return c.errorWithUsage(err)
	}

	if c.flagSet.NArg() > 1 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()[1:]))
	}

	c.Ui.Info("Are you sure you want to delete " + redisName + "?")
	if *alwaysYes {
		c.Ui.Info("(approved from command line)")
	} else {
		res, err := c.Ui.Ask("y/N")
		if err != nil {
			return c.error(err)
		} else if res != "Y" && res != "y" {
			return c.cancelled()
		}
	}

	res := c.runWithSpinner("deleting redis", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
func (c *RedisDeleteCommand) Synopsis() string {
	return "Delete a redis on the project."
}

// Help is part of cli.Command implementation.
func (c *RedisDeleteCommand) Help() string {
	helpText := `
usage: sqsc redis delete [options] <redis_name>

  Delete the redis from the project.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
