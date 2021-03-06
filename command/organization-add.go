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
func (c *OrganizationAddCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	name := c.flagSet.String("name", "", "name")
	email := c.flagSet.String("email", "", "contact email")

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if *name == "" {
		return c.errorWithUsage(fmt.Errorf("Name must not be empty, use -name option"))
	}

	if *email == "" {
		return c.errorWithUsage(fmt.Errorf("email must not be empty, use -email option"))
	}

	res := c.runWithSpinner("add organization", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
func (c *OrganizationAddCommand) Synopsis() string {
	return "Add organization"
}

// Help is part of cli.Command implementation.
func (c *OrganizationAddCommand) Help() string {
	helpText := `
usage: sqsc organization add [options]

  Add organization.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
