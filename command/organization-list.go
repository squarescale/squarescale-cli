package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// OrganizationListCommand is a cli.Command implementation for listing all Squarescale organizations.
type OrganizationListCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *OrganizationListCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	return c.runWithSpinner("list organizations", endpoint.String(), func(client *squarescale.Client) (string, error) {
		organizations, err := client.ListOrganizations()

		if err != nil {
			return "", err
		}

		var msg string

		for _, o := range organizations {
			msg += fmt.Sprintf("[%s]\n", o.Name)

			msg += fmt.Sprintf("  Projects:\n")
			msg += fmt.Sprintf("    Name\tStatus\tSize\n")
			msg += fmt.Sprintf("    ----\t------\t----\n")
			for _, p := range o.Projects {
				msg += fmt.Sprintf("    %s\t%s\t%d/%d\n", p.Name, p.InfraStatus, p.NomadNodesReady, p.ClusterSize)
			}

			msg += fmt.Sprintf("\n")

			msg += fmt.Sprintf("  Collaborators:\n")
			msg += fmt.Sprintf("    Name\tEmail\n")
			msg += fmt.Sprintf("    ----\t-----\n")
			for _, c := range o.Collaborators {
				msg += fmt.Sprintf("    %s\t%s\n", c.Name, c.Email)
			}

			msg += fmt.Sprintf("\n")
		}

		if len(organizations) == 0 {
			msg = "No organizations found"
		}

		return c.FormatTable(msg, true), nil
	})
}

// Synopsis is part of cli.Command implementation.
func (c *OrganizationListCommand) Synopsis() string {
	return "List all organizations"
}

// Help is part of cli.Command implementation.
func (c *OrganizationListCommand) Help() string {
	helpText := `
usage: sqsc organization list [options]

  List al organizations.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
