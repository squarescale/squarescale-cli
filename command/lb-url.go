package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// LBURLCommand gets the URL of the load balancer associated to a projects and prints it on the standard output.
type LBURLCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *LBURLCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	project := projectFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if err := validateProjectName(*project); err != nil {
		return c.errorWithUsage(err)
	}

	return c.runWithSpinner("project url", *endpoint, func(client *squarescale.Client) (string, error) {
		url, err := client.ProjectURL(*project)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("Project '%s' is available at %s", *project, url), nil
	})
}

// Synopsis is part of cli.Command implementation.
func (c *LBURLCommand) Synopsis() string {
	return "Display project's public URL if available"
}

// Help is part of cli.Command implementation.
func (c *LBURLCommand) Help() string {
	helpText := `
usage: sqsc lb url [options]

  Display load balancer public URL if available in the specified Squarescale project.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
