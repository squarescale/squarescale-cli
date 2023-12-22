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
func (cmd *ProjectProvisionCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if cmd.flagSet.NArg() > 0 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	return cmd.runWithSpinner("provision project", endpoint.String(), func(client *squarescale.Client) (string, error) {
		var UUID string
		var err error
		if *projectUUID == "" {
			UUID, err = client.ProjectByName(*projectName)
			if err != nil {
				return "", err
			}
		} else {
			UUID = *projectUUID
		}

		err = client.ProvisionProject(UUID)
		return "", err
	})
}

// Synopsis is part of cli.Command implementation.
func (cmd *ProjectProvisionCommand) Synopsis() string {
	return "Provision infrastructure of project"
}

// Help is part of cli.Command implementation.
func (cmd *ProjectProvisionCommand) Help() string {
	helpText := `
usage: sqsc project provision [options]

  Provision the infrastructure of the project.

`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
