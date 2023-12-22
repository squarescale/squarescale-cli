package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// OrganizationDeleteCommand is a cli.Command implementation for remove a SquareScale organization.
type OrganizationDeleteCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *OrganizationDeleteCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	alwaysYes := yesFlag(cmd.flagSet)
	endpoint := endpointFlag(cmd.flagSet)
	name := organizationNameFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if cmd.flagSet.NArg() > 0 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	if *name == "" {
		return cmd.errorWithUsage(fmt.Errorf("Name must not be empty, use -name option"))
	}

	cmd.Ui.Info("Are you sure you want to delete " + *name + "?")
	if *alwaysYes {
		cmd.Ui.Info("(approved from command line)")
	} else {
		res, err := cmd.Ui.Ask("y/N")

		if err != nil {
			return cmd.error(err)
		} else if res != "Y" && res != "y" {
			return cmd.cancelled()
		}
	}

	res := cmd.runWithSpinner("deleting organization", endpoint.String(), func(client *squarescale.Client) (string, error) {
		err := client.DeleteOrganization(*name)
		return "", err
	})

	if res != 0 {
		return res
	}

	return 0
}

// Synopsis is part of cli.Command implementation.
func (cmd *OrganizationDeleteCommand) Synopsis() string {
	return "Remove organization"
}

// Help is part of cli.Command implementation.
func (cmd *OrganizationDeleteCommand) Help() string {
	helpText := `
usage: sqsc organization delete [options]

  Delete organization.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
