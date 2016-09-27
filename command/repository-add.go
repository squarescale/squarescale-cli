package command

import (
	"errors"
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
	cmdFlags := flag.NewFlagSet("repository add", flag.ContinueOnError)
	cmdFlags.Usage = func() { r.Ui.Output(r.Help()) }
	endpoint := EndpointFlag(cmdFlags)
	project := ProjectFlag(cmdFlags)
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	err := validateArgs(*project)
	if err != nil {
		r.ErrorWithUsage(err, r.Help())
		return 1
	}

	token, err := tokenstore.GetToken(*endpoint)
	if err != nil {
		r.Error(err)
		return 1
	}

	gitRemote, err := findGitRemote()
	if err != nil {
		r.Error(err)
		return 1
	}

	s := startSpinner(fmt.Sprintf("add repository '%s' to project '%s'", gitRemote, *project))
	err = squarescale.AddRepository(*endpoint, token, *project, gitRemote)
	if err != nil {
		s.Stop()
		r.Error(err)
		return 1
	}

	s.Stop()
	r.Ui.Info(fmt.Sprintf("Successfully added repository '%s' to project '%s'", gitRemote, *project))
	return 0
}

// Synopsis is part of cli.Command implementation.
func (r *RepositoryAddCommand) Synopsis() string {
	return "Attach the current Git repository to a Squarescale project"
}

// Help is part of cli.Command implementation.
func (r *RepositoryAddCommand) Help() string {
	helpText := `
usage: sqsc repository add [options]

  Adds current Git repository to the specified Squarescale project. The
  repository must contain a Dockerfile to be attached to a project.

Options:

  -endpoint="http://www.staging.sqsc.squarely.io" Squarescale endpoint
  -project=""                                     Squarescale project name
`
	return strings.TrimSpace(helpText)
}

func findGitRemote() (string, error) {
	output, err := exec.Command("git", "remote", "get-url", "origin").CombinedOutput()
	formattedOutput := strings.Replace(string(output), "\n", "", -1)

	if err != nil {
		formattedOutput = strings.Replace(formattedOutput, "fatal: ", "", -1)
		return "", errors.New(formattedOutput)
	}

	return formattedOutput, nil
}
