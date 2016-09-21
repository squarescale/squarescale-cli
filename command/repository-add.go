package command

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
	"github.com/squarescale/squarescale-cli/tokenstore"
)

// RepositoryAddCommand is a cli.Command implementation for adding a repository to a Squarescale project.
type RepositoryAddCommand struct {
	Meta
}

// Run is part of cli.Command implementation.
func (r *RepositoryAddCommand) Run(args []string) int {
	// Parse flags
	cmdFlags := flag.NewFlagSet("repository:add", flag.ContinueOnError)
	cmdFlags.Usage = func() { r.Ui.Output(r.Help()) }
	endpoint := EndpointFlag(cmdFlags)
	project := ProjectFlag(cmdFlags)
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	// Validate project name
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

	// Retrieve git remote url
	gitRemote, err := findGitRemote()
	if err != nil {
		r.Error(err)
		return 1
	}

	// send to squarescale
	err = squarescale.AddRepository(*endpoint, token, *project, gitRemote)
	if err != nil {
		r.Error(err)
		return 1
	}

	r.Ui.Info(fmt.Sprintf("Successfully attached repository '%s' to project '%s'", gitRemote, *project))
	return 0
}

// Synopsis is part of cli.Command implementation.
func (r *RepositoryAddCommand) Synopsis() string {
	return "Attach the current Git repository to a Squarescale project"
}

// Help is part of cli.Command implementation.
func (r *RepositoryAddCommand) Help() string {
	helpText := `
usage: sqsc repository:add [options]

  Adds current Git repository to the specified Squarescale project. The
  repository must contain a Dockerfile to be attached to a project.

Options:

  -endpoint="http://www.staging.sqsc.squarely.io" Squarescale endpoint
  -project=""                                     Squarescale project name
`
	return strings.TrimSpace(helpText)
}

func findGitRemote() (string, error) {
	remote, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		return "", err
	}

	return strings.Replace(string(remote), "\n", "", -1), nil
}
