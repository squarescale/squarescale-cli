package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ProjectProvisionCommand is a cli.Command implementation for provision one project.
type ProjectProvisionCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *ProjectProvisionCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := c.flagSet.String("project-uuid", "", "set the uuid of the project")

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" {
		return c.errorWithUsage(errors.New("Project uuid is mandatory"))
	}

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	return c.runWithSpinner("provision project", endpoint.String(), func(client *squarescale.Client) (string, error) {
		err := client.ProvisionProject(*projectUUID)
		return "", err
	})
}

// Synopsis is part of cli.Command implementation.
func (c *ProjectProvisionCommand) Synopsis() string {
	return "Provision one project"
}

// Help is part of cli.Command implementation.
func (c *ProjectProvisionCommand) Help() string {
	helpText := `
usage: sqsc project provision [options]

  Provision the infrasctructure of the project.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
