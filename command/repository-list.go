package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
	"github.com/squarescale/squarescale-cli/tokenstore"
)

// RepositoryListCommand is a cli.Command implementation for listing Squarescale project repositories.
type RepositoryListCommand struct {
	Meta
}

// Run is part of cli.Command implementation.
func (r *RepositoryListCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("repository list", flag.ContinueOnError)
	cmdFlags.Usage = func() { r.Ui.Output(r.Help()) }
	endpoint := EndpointFlag(cmdFlags)
	project := ProjectFlag(cmdFlags)
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	if *project == "" {
		r.Ui.Error("Project name must be specified.\n")
		r.Ui.Output(r.Help())
		return 1
	}

	token, err := tokenstore.GetToken(*endpoint)
	if err != nil {
		r.Error(err)
		return 1
	}

	s := startSpinner("list repositories")
	repositories, err := squarescale.ListRepositories(*endpoint, token, *project)
	if err != nil {
		s.Stop()
		r.Error(err)
		return 1
	}

	var msg string
	if len(repositories) > 0 {
		msg = strings.Join(repositories, "\n")
	} else {
		msg = fmt.Sprintf("No repositories attached to project '%s'", *project)
	}

	s.Stop()
	r.Ui.Info(msg)
	return 0
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
