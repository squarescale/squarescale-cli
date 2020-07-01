package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// BatchDeleteCommand is a cli.Command implementation for creating a Squarescale project.
type BatchDeleteCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *BatchDeleteCommand) Run(args []string) int {
	// Parse flags
	c.flagSet = newFlagSet(c, c.Ui)
	alwaysYes := yesFlag(c.flagSet)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := c.flagSet.String("project-uuid", "", "set the uuid of the project")

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" {
		return c.errorWithUsage(errors.New("Project uuid is mandatory"))
	}

	batchName, err := batchNameArg(c.flagSet, 0)
	if err != nil {
		return c.errorWithUsage(err)
	}

	if c.flagSet.NArg() > 1 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()[1:]))
	}

	if err := validateProjectName(*projectUUID); err != nil {
		return c.errorWithUsage(err)
	}

	c.Ui.Info("Are you sure you want to delete " + batchName + "?")
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

	res := c.runWithSpinner("deleting batch", endpoint.String(), func(client *squarescale.Client) (string, error) {
		fmt.Printf("Delete on project `%s` the batch `%s`\n", *projectUUID, batchName)
		err := client.DeleteBatch(*projectUUID, batchName) //insérer la fonction dans batches
		return "", err
	})
	if res != 0 {
		return res
	}

	return 0
}

// Synopsis is part of cli.Command implementation.
func (c *BatchDeleteCommand) Synopsis() string {
	return "Delete a batch from the project."
}

// Help is part of cli.Command implementation.
func (c *BatchDeleteCommand) Help() string {
	helpText := `
usage: sqsc batch delete [options] <batch_name>

  Delete the batch from the project.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
