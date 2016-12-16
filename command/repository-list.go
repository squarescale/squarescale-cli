package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// RepositoryListCommand is a cli.Command implementation for listing Squarescale project repositories.
type RepositoryListCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *RepositoryListCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	project := projectFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	err := validateProjectName(*project)
	if err != nil {
		return c.errorWithUsage(err)
	}

	var msg string
	err = c.runWithSpinner("list repositories", *endpoint, func(client *squarescale.Client) error {
		repositories, e := client.ListRepositories(*project)
		if e != nil {
			return e
		}

		msg = strings.Join(repositories, "\n")
		if len(repositories) == 0 {
			msg = fmt.Sprintf("No repositories attached to project '%s'", *project)
		}

		return nil
	})

	if err != nil {
		return c.error(err)
	}

	return c.info(msg)
}

// Synopsis is part of cli.Command implementation.
func (c *RepositoryListCommand) Synopsis() string {
	return "Lists all Git repositories attached to a Squarescale project"
}

// Help is part of cli.Command implementation.
func (c *RepositoryListCommand) Help() string {
	helpText := `
usage: sqsc repository list [options]

  Lists all Git repositories attached to the specified Squarescale project.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
