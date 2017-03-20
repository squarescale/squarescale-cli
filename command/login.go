package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/github"
	"github.com/squarescale/squarescale-cli/squarescale"
	"github.com/squarescale/squarescale-cli/tokenstore"
)

// LoginCommand is a cli.Command implementation for authenticating the user into the Squarescale platform.
type LoginCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *LoginCommand) Run(args []string) int {
	// Parse flags
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	// Retrieve credentials from user input
	login, pw, err := c.askForCredentials()
	if err != nil {
		return c.error(err)
	}

	oneTimePW, ghToken, ghTokenURL, err := github.GeneratePersonalToken(login, pw, c.Ui)
	if err != nil {
		return c.error(err)
	}

	c.Ui.Info("Forward GitHub authorization to Squarescale")

	res := 0
	sqscToken, err := squarescale.ObtainTokenFromGitHub(*endpoint, ghToken)
	if err != nil {
		res = c.error(err)

	} else {
		err = tokenstore.SaveToken(*endpoint, sqscToken)
		if err != nil {
			res = c.error(err)

		} else {
			c.Ui.Info(fmt.Sprintf("Successfully authenticated as user %s", login))
		}
	}

	err = github.RevokePersonalToken(ghTokenURL, login, pw, oneTimePW)
	if err == nil {
		c.Ui.Info("Revoke temporary GitHub token")
	}

	return res
}

// Synopsis is part of cli.Command implementation.
func (c *LoginCommand) Synopsis() string {
	return "Login to Squarescale"
}

// Help is part of cli.Command implementation.
func (c *LoginCommand) Help() string {
	helpText := `
usage: sqsc login [options]

  Logs the user in Squarescale services. Uses GitHub credentials to register.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}

func (c *LoginCommand) askForCredentials() (string, string, error) {
	login, err := c.Ui.Ask("GitHub login:")
	if err != nil {
		return "", "", err
	}

	pw, err := c.Ui.AskSecret("GitHub password (typing will be hidden):")
	if err != nil {
		return "", "", err
	}

	return login, pw, err
}
