package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// BatchDeleteCommand is a cli.Command implementation for creating a SquareScale project.
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
	batchName := batchNameFlag(c.flagSet)
	projectUUID := projectUUIDFlag(c.flagSet)
	projectName := projectNameFlag(c.flagSet)

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if *batchName == "" {
		return c.errorWithUsage(fmt.Errorf(("Batch name is mandatory. Please, chose a batch name.")))
	}

	if c.flagSet.NArg() > 1 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()[1:]))
	}

	c.Ui.Info("Are you sure you want to delete " + *batchName + "?")
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
		var UUID string
		var err error
		var projectToShow string
		if *projectUUID == "" {
			projectToShow = *projectName
			UUID, err = client.ProjectByName(*projectName)
			if err != nil {
				return "", err
			}
		} else {
			projectToShow = *projectUUID
			UUID = *projectUUID
		}

		fmt.Printf("Delete on project `%s` the batch `%s`\n", projectToShow, *batchName)
		err = client.DeleteBatch(UUID, *batchName) //insérer la fonction dans batches
		return "", err
	})
	if res != 0 {
		return res
	}

	return 0
}

// Synopsis is part of cli.Command implementation.
func (c *BatchDeleteCommand) Synopsis() string {
	return "Delete batch from project."
}

// Help is part of cli.Command implementation.
func (c *BatchDeleteCommand) Help() string {
	helpText := `
usage: sqsc batch delete [options] <batch_name>

  Delete batch from project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
