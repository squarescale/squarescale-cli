package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// OrganizationDeleteCommand is a cli.Command implementation for remove a Squarescale organization.
type OrganizationDeleteCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *OrganizationDeleteCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	alwaysYes := yesFlag(c.flagSet)
	endpoint := endpointFlag(c.flagSet)

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	organizationName, err := organizationNameArg(c.flagSet, 0)

	if err != nil {
		return c.errorWithUsage(err)
	}

	if c.flagSet.NArg() > 1 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()[1:]))
	}

	c.Ui.Info("Are you sure you want to delete " + organizationName + "?")
	if *alwaysYes {
		c.Ui.Info("(approved from command line)")
	} else {
		res, err := c.Ui.Ask("y/N")

		if err != nil {
			return c.error(err)
		} else if res != "Y" && res != "y" {
			return c.cancelled()
		}
	}

	res := c.runWithSpinner("deleting organization", endpoint.String(), func(client *squarescale.Client) (string, error) {
		err := client.DeleteOrganization(organizationName)
		return "", err
	})

	if res != 0 {
		return res
	}

	return 0
}

// Synopsis is part of cli.Command implementation.
func (c *OrganizationDeleteCommand) Synopsis() string {
	return "Remove a organization"
}

// Help is part of cli.Command implementation.
func (c *OrganizationDeleteCommand) Help() string {
	helpText := `
usage: sqsc organization remove [options] <organization_name>

  Remove a organization.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}