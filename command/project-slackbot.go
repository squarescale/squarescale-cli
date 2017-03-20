package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ProjectSlackbotCommand is a cli.Command implementation for updating slack
// integration.
type ProjectSlackbotCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *ProjectSlackbotCommand) Run(args []string) int {
	// Parse flags
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	projectName, err := projectNameArg(c.flagSet, 0)
	if err != nil {
		return c.errorWithUsage(err)
	}

	if c.flagSet.NArg() > 2 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()[2:]))
	}

	url := c.flagSet.Arg(1)
	if url == "" {
		return c.runWithSpinner("retrieve slack settings", *endpoint, func(client *squarescale.Client) (string, error) {
			return client.ProjectSlackURL(projectName)
		})
	} else {
		return c.runWithSpinner("change slack settings", *endpoint, func(client *squarescale.Client) (string, error) {
			return url, client.SetProjectSlackURL(projectName, url)
		})
	}
}

// Synopsis is part of cli.Command implementation.
func (c *ProjectSlackbotCommand) Synopsis() string {
	return "Show and update slack integration"
}

// Help is part of cli.Command implementation.
func (c *ProjectSlackbotCommand) Help() string {
	helpText := `
usage: sqsc project slackbot [options] <project_name> [<url>]

  Set the slack URL integration, if an URL is provided, or show the current URL.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
