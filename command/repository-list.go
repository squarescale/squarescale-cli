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

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if err := validateProjectName(*project); err != nil {
		return c.errorWithUsage(err)
	}

	return c.runWithSpinner("list repositories", *endpoint, func(client *squarescale.Client) (string, error) {
		repositories, err := client.ListRepositories(*project)
		if err != nil {
			return "", err
		}

		msg := strings.Join(repositories, "\n")
		if len(repositories) == 0 {
			msg = fmt.Sprintf("No repositories attached to project '%s'", *project)
		}

		return msg, nil
	})
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
