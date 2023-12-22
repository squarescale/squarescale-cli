package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// OrganizationListCommand is a cli.Command implementation for listing all SquareScale organizations.
type OrganizationListCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *OrganizationListCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if cmd.flagSet.NArg() > 0 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	return cmd.runWithSpinner("list organizations", endpoint.String(), func(client *squarescale.Client) (string, error) {
		organizations, err := client.ListOrganizations()

		if err != nil {
			return "", err
		}

		var msg string

		for _, o := range organizations {
			msg += fmt.Sprintf("[%s]\n", o.Name)
			msg += fmt.Sprintf("  Root user email: %s\n\n", o.RootUser.Email)
			msg += fmt.Sprintf("  Projects:\n")
			msg += fmt.Sprintf("    UUID\tName\tStatus\tSize\n")
			msg += fmt.Sprintf("    ----\t----\t------\t----\n")
			for _, p := range o.Projects {
				msg += fmt.Sprintf("    %s\t%s\t%s\t%d/%d\n", p.UUID, p.Name, p.InfraStatus, p.NomadNodesReady, p.ClusterSize)
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

		return cmd.FormatTable(msg, true), nil
	})
}

// Synopsis is part of cli.Command implementation.
func (cmd *OrganizationListCommand) Synopsis() string {
	return "List organizations"
}

// Help is part of cli.Command implementation.
func (cmd *OrganizationListCommand) Help() string {
	helpText := `
usage: sqsc organization list [options]

  List organizations.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
