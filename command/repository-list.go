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
}

// Run is part of cli.Command implementation.
func (r *RepositoryListCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("repository list", flag.ContinueOnError)
	cmdFlags.Usage = func() { r.Ui.Output(r.Help()) }
	endpoint := endpointFlag(cmdFlags)
	project := projectFlag(cmdFlags)
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	err := validateProjectName(*project)
	if err != nil {
		return r.errorWithUsage(err, r.Help())
	}

	var msg string
	err = r.runWithSpinner("list repositories", *endpoint, func(token string) error {
		repositories, e := squarescale.ListRepositories(*endpoint, token, *project)
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
		return r.error(err)
	}

	return r.info(msg)
}

// Synopsis is part of cli.Command implementation.
func (r *RepositoryListCommand) Synopsis() string {
	return "Lists all Git repositories attached to a Squarescale project"
}

// Help is part of cli.Command implementation.
func (r *RepositoryListCommand) Help() string {
	helpText := `
usage: sqsc repository list [options]

  Lists all Git repositories attached to the specified Squarescale project.

Options:

  -endpoint="http://www.staging.sqsc.squarely.io" Squarescale endpoint
  -project=""                                     Squarescale project name
`
	return strings.TrimSpace(helpText)
}
