package command

import (
	"errors"
	"flag"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// RepositoryAddCommand is a cli.Command implementation for adding a repository to a Squarescale project.
type RepositoryAddCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *RepositoryAddCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	project := projectFlag(c.flagSet)
	buildService := buildServiceFlag(c.flagSet)
	url := repoUrlFlag(c.flagSet)
	instances := repoOrImageInstancesFlag(c.flagSet)
	branch := c.flagSet.String("branch", "master", "Git branch")

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if *buildService != "travis" && *buildService != "internal" {
		return c.errorWithUsage(fmt.Errorf("Unknown build service: %v. Correct values are travis or internal", *buildService))
	}

	if *instances <= 0 {
		err := errors.New("Invalid value provided for instances count: it cannot be 0 or negative")
		return c.errorWithUsage(err)
	}

	err := validateProjectName(*project)
	if err != nil {
		return c.errorWithUsage(err)
	}

	gitRemote := *url
	if gitRemote == "" {
		gitRemote, err = findGitRemote()
		if err != nil {
			return c.error(err)
		}
	}

	label := fmt.Sprintf("add repository '%s' to project '%s'", gitRemote, *project)
	return c.runWithSpinner(label, endpoint.String(), func(client *squarescale.Client) (string, error) {
		msg := fmt.Sprintf("Successfully added repository '%s' to project '%s' (%v instance(s))", gitRemote, *project, *instances)
		res := client.AddRepository(*project, gitRemote, *branch, *buildService, *instances)
		time.Sleep(30 * time.Second)
		return msg, res
	})
}

// Synopsis is part of cli.Command implementation.
func (c *RepositoryAddCommand) Synopsis() string {
	return "Attach the current Git repository to a Squarescale project"
}

// Help is part of cli.Command implementation.
func (c *RepositoryAddCommand) Help() string {
	helpText := `
usage: sqsc repository add [options]

  Adds current Git repository to the specified Squarescale project. The
  repository must contain a Dockerfile to be attached to a project.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
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

func buildServiceFlag(f *flag.FlagSet) *string {
	return f.String("build-service", "travis", "Set the service used to build project repositories.")
}
