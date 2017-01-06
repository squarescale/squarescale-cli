package command

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// LogsCommand outputs the logs of a container for a given project.
type LogsCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *LogsCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	project := projectFlag(c.flagSet)
	container := containerFlag(c.flagSet)
	var follow bool
	c.flagSet.BoolVar(&follow, "f", false, "Wait for next logs")
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	err := validateProjectName(*project)
	if err != nil {
		return c.errorWithUsage(err)
	}

	var waitText string
	if *container == "" {
		waitText = fmt.Sprintf("last logs for project '%s'", *project)
	} else {
		waitText = fmt.Sprintf("last logs for container '%s'", *container)
	}

	var sqscClient *squarescale.Client
	var msg string
	var last int
	retCode := c.runWithSpinner(waitText, *endpoint, func(client *squarescale.Client) (string, error) {
		sqscClient = client
		msg, last, err = getLogs(client, *project, *container)
		return msg, err
	})

	if retCode != 0 {
		return retCode
	}

	if !follow {
		return 0
	}

	for {
		time.Sleep(time.Second)
		msg, last, err = getLogsAfter(sqscClient, *project, *container, last)
		if err != nil {
			return c.error(err)
		}

		if msg != "" {
			c.info(msg)
		}
	}
}

func getLogs(client *squarescale.Client, project, container string) (string, int, error) {
	return getLogsAfter(client, project, container, -1)
}

func getLogsAfter(client *squarescale.Client, project string, container string, last int) (string, int, error) {
	logs, lastOffset, e := client.ProjectLogs(project, container, last)
	if e != nil {
		return "", -1, e
	}

	msg := "\033[0m" + strings.Join(logs, "\n")
	if len(logs) == 0 {
		if last < 0 {
			msg = fmt.Sprintf("No available logs")
		} else {
			lastOffset = last
			msg = ""
		}
	}

	return msg, lastOffset, nil
}

// Synopsis is part of cli.Command implementation.
func (c *LogsCommand) Synopsis() string {
	return "Display last logs of a Squarescale project or one specific container"
}

// Help is part of cli.Command implementation.
func (c *LogsCommand) Help() string {
	helpText := `
usage: sqsc logs [options]

  Display last logs of the specified Squarescale project.
  Container short url to filter logs follows the pattern: <user>/<repository>

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
