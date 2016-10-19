package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

type LogsCommand struct {
	Meta
}

func ContainerFlag(f *flag.FlagSet) *string {
	return f.String("container", "", "Optional container short url to filter logs, ex: <githubuser>/<repository>")
}

func (c *LogsCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("logs", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	endpoint := EndpointFlag(cmdFlags)
	project := ProjectFlag(cmdFlags)
	container := ContainerFlag(cmdFlags)
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	err := validateArgs(*project)
	if err != nil {
		return c.errorWithUsage(err, c.Help())
	}

	var msg string
	var waitText string
	if *container == "" {
		waitText = fmt.Sprintf("last logs for project '%s'", *project)
	} else {
		waitText = fmt.Sprintf("last logs for container '%s'", *container)
	}
	err = c.runWithSpinner(waitText, *endpoint, func(token string) error {
		logs, e := squarescale.ProjectLogs(*endpoint, token, *project, *container)
		if e != nil {
			return e
		}

		msg = strings.Join(logs, "\n")
		if len(logs) == 0 {
			msg = fmt.Sprintf("No available logs")
		}

		return nil
	})

	if err != nil {
		return c.error(err)
	}

	return c.info(msg)
}

func (c *LogsCommand) Synopsis() string {
	return "Display last logs of a Squarescale project or one specific container"
}

func (c *LogsCommand) Help() string {
	helpText := `
usage: sqsc project logs [options]

  Display last logs of the specified Squarescale project.

Options:

  -endpoint="http://www.staging.sqsc.squarely.io" Squarescale endpoint
  -project=""                                     Squarescale project name
  -container=""                                   Optional short url to filter
                                                  logs for a specific container.
                                                  Should be like:
                                                  <github user>/<repo name>
`
	return strings.TrimSpace(helpText)
}
