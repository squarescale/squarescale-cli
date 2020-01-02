package command

import (
	"flag"
	"fmt"
	"os"
	"strings"

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

	var apiKey string
	var err error
	apiKey = os.Getenv("SQSC_API_KEY")
	if apiKey == "" {
		// Retrieve credentials from user input
		apiKey, err = c.askForCredentials()
		if err != nil {
			return c.error(err)
		}
	}

	err = tokenstore.SaveToken(endpoint.String(), apiKey)
	if err != nil {
		return c.error(err)
	} else {
		_, err = c.ensureLogin(endpoint.String())
		if err != nil {
			return c.error(err)
		} else {
			c.Ui.Info(fmt.Sprintf("Successfully authenticated !"))
			return 0
		}
	}
}

// Synopsis is part of cli.Command implementation.
func (c *LoginCommand) Synopsis() string {
	return "Login to Squarescale"
}

// Help is part of cli.Command implementation.
func (c *LoginCommand) Help() string {
	helpText := `
usage: sqsc login [options]

  Logs the user in Squarescale services.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}

func (c *LoginCommand) askForCredentials() (string, error) {
	apiKey, err := c.Ui.AskSecret("API key (typing will be hidden):")
	if err != nil {
		return "", err
	}
	return apiKey, err
}
