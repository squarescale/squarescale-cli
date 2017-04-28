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

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
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
	var last string
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

func getLogs(client *squarescale.Client, project, container string) (string, string, error) {
	return getLogsAfter(client, project, container, "")
}

func getLogsAfter(client *squarescale.Client, project string, container string, last string) (string, string, error) {
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
