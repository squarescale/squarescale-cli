package command

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/squarescale/squarescale-cli/squarescale"
)

type LogsCommand struct {
	Meta
	flagSet *flag.FlagSet
}

func ContainerFlag(f *flag.FlagSet) *string {
	return f.String("container", "", "Optional container short url to filter logs, ex: <githubuser>/<repository>")
}

func FollowFlag(f *flag.FlagSet) *string {
	return f.String("f", "", "Follow logs")
}

func (c *LogsCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	project := projectFlag(c.flagSet)
	container := ContainerFlag(c.flagSet)
	var follow bool
	c.flagSet.BoolVar(&follow, "f", false, "Wait for next logs")
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	err := validateProjectName(*project)
	if err != nil {
		return c.errorWithUsage(err)
	}

	var msg string
	var waitText string
	if *container == "" {
		waitText = fmt.Sprintf("last logs for project '%s'", *project)
	} else {
		waitText = fmt.Sprintf("last logs for container '%s'", *container)
	}

	var sqscClient *squarescale.Client
	var last string
	err = c.runWithSpinner(waitText, *endpoint, func(client *squarescale.Client) error {
		var e error
		sqscClient = client
		msg, last, e = getLogs(client, *project, *container)
		return e
	})

	if err != nil {
		return c.error(err)
	}

	c.info(msg)
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

func getLogs(client *squarescale.Client, project, container string) (string, string, error) {
	return getLogsAfter(client, project, container, "")
}

func getLogsAfter(client *squarescale.Client, project, container, last string) (string, string, error) {
	logs, lastTimestamp, e := client.ProjectLogs(project, container, last)
	if e != nil {
		return "", "", e
	}

	msg := "\033[0m" + strings.Join(logs, "\n")
	if len(logs) == 0 {
		if last == "" {
			msg = fmt.Sprintf("No available logs")
		} else {
			lastTimestamp = last
			msg = ""
		}
	}

	return msg, lastTimestamp, nil
}

func (c *LogsCommand) Synopsis() string {
	return "Display last logs of a Squarescale project or one specific container"
}

func (c *LogsCommand) Help() string {
	helpText := `
usage: sqsc logs [options]

  Display last logs of the specified Squarescale project.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
