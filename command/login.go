package command

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/squarescale/squarescale-cli/tokenstore"
)

// LoginCommand is a cli.Command implementation for authenticating the user into the SquareScale platform.
type LoginCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *LoginCommand) Run(args []string) int {
	// Parse flags
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if cmd.flagSet.NArg() > 0 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	var apiKey string
	var err error
	apiKey = os.Getenv("SQSC_TOKEN")
	if apiKey == "" {
		// Retrieve credentials from previous session (token storage)
		apiKey, err = tokenstore.GetToken(endpoint.String())
		if err != nil || apiKey == "" {
			// Retrieve credentials from user input
			apiKey, err = cmd.askForCredentials()
			if err != nil {
				return cmd.error(err)
			}
		}
	}

	err = tokenstore.SaveToken(endpoint.String(), apiKey)
	if err != nil {
		return cmd.error(err)
	} else {
		_, err = cmd.ensureLogin(endpoint.String())
		if err != nil {
			return cmd.error(err)
		} else {
			cmd.Ui.Info(fmt.Sprintf("Successfully authenticated !"))
			return 0
		}
	}
}

// Synopsis is part of cli.Command implementation.
func (cmd *LoginCommand) Synopsis() string {
	return "login to SquareScale platform"
}

// Help is part of cli.Command implementation.
func (cmd *LoginCommand) Help() string {
	helpText := `
usage: sqsc login [options]

  Logs the user in SquareScale platform.

`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}

func (cmd *LoginCommand) askForCredentials() (string, error) {
	apiKey, err := cmd.Ui.AskSecret("API key (typing will be hidden):")
	if err != nil {
		return "", err
	}
	return apiKey, err
}
