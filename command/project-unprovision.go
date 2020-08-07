package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ProjectProvisionCommand is a cli.Command implementation for provision one project.
type ProjectUNProvisionCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *ProjectUNProvisionCommand) Run(args []string) int {
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

	return c.runWithSpinner("unprovision project", endpoint.String(), func(client *squarescale.Client) (string, error) {
		err := client.UNProvisionProject(*projectUUID)
		return "", err
	})
}

// Synopsis is part of cli.Command implementation.
func (c *ProjectUNProvisionCommand) Synopsis() string {
	return "Unprovision infrastructure of project"
}

// Help is part of cli.Command implementation.
func (c *ProjectUNProvisionCommand) Help() string {
	helpText := `
usage: sqsc project unprovision [options]

  Unprovision infrasctructure of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
