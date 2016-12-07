package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

type UrlCommand struct {
	Meta
	flagSet *flag.FlagSet
}

func (c *UrlCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	project := projectFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	err := validateProjectName(*project)
	if err != nil {
		return c.errorWithUsage(err)
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

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
