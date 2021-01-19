package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// BatchExecuteCommand is a cli.Command implementation for top level `sqsc batch` command.
type BatchExecuteCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *BatchExecuteCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := c.flagSet.String("project-uuid", "", "set the uuid of the project")
	projectName := c.flagSet.String("project-name", "", "set the name of the project")

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	batchName, err := batchNameArg(c.flagSet, 0)
	if err != nil {
		return c.errorWithUsage(err)
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if c.flagSet.NArg() > 4 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()[1:]))
	}

	return c.runWithSpinner("executing batch", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		err = client.ExecuteBatch(UUID, batchName)
		return "", err
	})

}

// Synopsis is part of cli.Command implementation.
func (c *BatchExecuteCommand) Synopsis() string {
	return "Execute a batch"
}

// Help is part of cli.Command implementation.
func (c *BatchExecuteCommand) Help() string {
	helpText := `
usage: sqsc batch exec [options] <batch_name>

  Schedule a batch.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
