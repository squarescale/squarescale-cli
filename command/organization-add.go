package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// OrganizationAddCommand is a cli.Command implementation for creating a SquareScale organization.
type OrganizationAddCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *OrganizationAddCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	name := organizationNameFlag(cmd.flagSet)
	email := organizationEmailFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if cmd.flagSet.NArg() > 0 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	if *name == "" {
		return cmd.errorWithUsage(fmt.Errorf("Name must not be empty, use -name option"))
	}

	if *email == "" {
		return cmd.errorWithUsage(fmt.Errorf("email must not be empty, use -email option"))
	}

	res := cmd.runWithSpinner("add organization", endpoint.String(), func(client *squarescale.Client) (string, error) {
		msg := fmt.Sprintf("Successfully added organization '%s'", *name)
		err := client.AddOrganization(*name, *email)
		return msg, err
	})

	if res != 0 {
		return res
	}

	return 0
}

// Synopsis is part of cli.Command implementation.
func (cmd *OrganizationAddCommand) Synopsis() string {
	return "Add organization"
}

// Help is part of cli.Command implementation.
func (cmd *OrganizationAddCommand) Help() string {
	helpText := `
usage: sqsc organization add [options]

  Add organization.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
