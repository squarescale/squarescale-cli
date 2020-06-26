package command

import (
	"errors"
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
	projectUUID := c.flagSet.String("project-uuid", "", "uuid of the targeted project")
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if *projectUUID == "" {
		return c.errorWithUsage(errors.New("Project uuid is mandatory"))
	}

	return c.runWithSpinner("project url", endpoint.String(), func(client *squarescale.Client) (string, error) {
		url, err := client.ProjectURL(*projectUUID)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("Project '%s' is available at %s", *projectUUID, url), nil
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
