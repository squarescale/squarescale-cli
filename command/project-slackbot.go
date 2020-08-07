package command

import (
	"errors"
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
	projectUUID := c.flagSet.String("project-uuid", "", "set the uuid of the project")

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" {
		return c.errorWithUsage(errors.New("Project uuid is mandatory"))
	}

	if c.flagSet.NArg() > 2 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()[2:]))
	}

	url := c.flagSet.Arg(1)
	if url == "" {
		return c.runWithSpinner("retrieve slack settings", endpoint.String(), func(client *squarescale.Client) (string, error) {
			return client.ProjectSlackURL(*projectUUID)
		})
	} else {
		return c.runWithSpinner(fmt.Sprintf("change slack integration URL to %s", url), endpoint.String(), func(client *squarescale.Client) (string, error) {
			return url, client.SetProjectSlackURL(*projectUUID, url)
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
