package command

import (
	"errors"
	"flag"
	"fmt"
	"os/exec"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// RepositoryAddCommand is a cli.Command implementation for adding a repository to a Squarescale project.
type RepositoryAddCommand struct {
	Meta
}

// Run is part of cli.Command implementation.
func (r *RepositoryAddCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("repository add", flag.ContinueOnError)
	cmdFlags.Usage = func() { r.Ui.Output(r.Help()) }
	endpoint := endpointFlag(cmdFlags)
	project := projectFlag(cmdFlags)
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	err := validateProjectName(*project)
	if err != nil {
		return r.errorWithUsage(err)
	}

	gitRemote, err := findGitRemote()
	if err != nil {
		return r.error(err)
	}

	err = r.runWithSpinner(fmt.Sprintf("add repository '%s' to project '%s'", gitRemote, *project), *endpoint, func(token string) error {
		return squarescale.AddRepository(*endpoint, token, *project, gitRemote)
	})

	if err != nil {
		return r.error(err)
	}

	return r.info(fmt.Sprintf("Successfully added repository '%s' to project '%s'", gitRemote, *project))
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
