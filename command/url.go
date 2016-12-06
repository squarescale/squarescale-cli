package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

type UrlCommand struct {
	Meta
}

func (c *UrlCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("url", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	endpoint := endpointFlag(cmdFlags)
	project := projectFlag(cmdFlags)
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	err := validateProjectName(*project)
	if err != nil {
		return c.errorWithUsage(err, c.Help())
	}

	var msg string
	err = c.runWithSpinner("project url", *endpoint, func(token string) error {
		url, e := squarescale.ProjectUrl(*endpoint, token, *project)
		if e != nil {
			return e
		}

		msg = fmt.Sprintf("Project '%s' is available at %s", *project, url)

		return nil
	})

	if err != nil {
		return c.error(err)
	}

	return c.info(msg)
}

func (c *UrlCommand) Synopsis() string {
	return "Display project's public URL if available"
}

func (c *UrlCommand) Help() string {
	helpText := `
usage: sqsc project url [options]

  Display load balancer public URL if available in the specified Squarescale project.

Options:

  -endpoint="http://www.staging.sqsc.squarely.io" Squarescale endpoint
  -project=""                                     Squarescale project name
`
	return strings.TrimSpace(helpText)
}
