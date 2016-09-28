package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
	"github.com/squarescale/squarescale-cli/tokenstore"
)

// StatusCommand is a cli.Command implementation for knowing if user is authorized.
type StatusCommand struct {
	Meta
}

// Run is part of cli.Command implementation.
func (c *StatusCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("status", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	endpoint := EndpointFlag(cmdFlags)
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	token, err := tokenstore.GetToken(*endpoint)
	if err != nil {
		c.Error(err)
		return 1
	}

	s := startSpinner("check authorization")
	err = squarescale.ValidateToken(*endpoint, token)
	if err != nil {
		s.Stop()
		c.Error(fmt.Errorf("Invalid token. Use the 'login' command first (%v).", err))
		return 1
	}

	s.Stop()
	c.Ui.Info("Current token is correctly authorized on Squarescale services.")
	return 0
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

Options:

  -endpoint="http://www.staging.sqsc.squarely.io" Squarescale endpoint
`
	return strings.TrimSpace(helpText)
}
