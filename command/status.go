package command

import (
	"flag"
	"strings"
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

	err := c.runWithSpinner("check authorization", *endpoint, func(token string) error {
		return nil // do nothing as we are already authenticated here.
	})

	if err != nil {
		return c.error(err)
	}

	return c.info("You're currently logged in")
}

// Synopsis is part of cli.Command implementation.
func (c *StatusCommand) Synopsis() string {
	return "Check Squarescale authorization"
}

// Help is part of cli.Command implementation.
func (c *StatusCommand) Help() string {
	helpText := `
usage: sqsc status [options]

  Asks the Squarescale services to check whether the current user is
  authenticated. This command checks the validity of the credentials
  stored in the $HOME/.netrc file.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
