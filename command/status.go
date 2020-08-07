package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// StatusCommand is a cli.Command implementation for knowing if user is authorized.
type StatusCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *StatusCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	fmt.Printf("Squarescale endpoint: %s\n", endpoint.String())
	return c.runWithSpinner("check authorization", endpoint.String(), func(client *squarescale.Client) (string, error) {
		return "You're currently logged in", nil // do nothing as we are already authenticated here.
	})
}

// Synopsis is part of cli.Command implementation.
func (c *StatusCommand) Synopsis() string {
	return "authorization check to Squarescale platform"
}

// Help is part of cli.Command implementation.
func (c *StatusCommand) Help() string {
	helpText := `
usage: sqsc status [options]

  Asks the Squarescale platform to check whether the current user is
  authenticated. This command checks the validity of the credentials
  stored in the $HOME/.netrc file.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
