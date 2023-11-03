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
func (cmd *StatusCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if cmd.flagSet.NArg() > 0 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	fmt.Printf("SquareScale endpoint: %s\n", endpoint.String())
	return cmd.runWithSpinner("check authorization", endpoint.String(), func(client *squarescale.Client) (string, error) {
		return "You're currently logged in", nil // do nothing as we are already authenticated here.
	})
}

// Synopsis is part of cli.Command implementation.
func (cmd *StatusCommand) Synopsis() string {
	return "authorization check to SquareScale platform"
}

// Help is part of cli.Command implementation.
func (cmd *StatusCommand) Help() string {
	helpText := `
usage: sqsc status [options]

  Asks the SquareScale platform to check whether the current user is
  authenticated. This command checks the validity of the credentials
  stored in the $HOME/.netrc file.

`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
