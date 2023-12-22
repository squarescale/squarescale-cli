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
func (cmd *BatchDeleteCommand) Run(args []string) int {
	// Parse flags
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	alwaysYes := yesFlag(cmd.flagSet)
	endpoint := endpointFlag(cmd.flagSet)
	batchName := batchNameFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if *batchName == "" {
		return cmd.errorWithUsage(fmt.Errorf(("Batch name is mandatory. Please, chose a batch name.")))
	}

	if cmd.flagSet.NArg() > 1 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()[1:]))
	}

	cmd.Ui.Info("Are you sure you want to delete " + *batchName + "?")
	if *alwaysYes {
		cmd.Ui.Info("(approved from command line)")
	} else {
		res, err := cmd.Ui.Ask("y/N")
		if err != nil {
			return cmd.error(err)
		} else if res != "Y" && res != "y" {
			return cmd.cancelled()
		}
	}

	res := cmd.runWithSpinner("deleting batch", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
func (cmd *BatchDeleteCommand) Synopsis() string {
	return "Delete batch from project."
}

// Help is part of cli.Command implementation.
func (cmd *BatchDeleteCommand) Help() string {
	helpText := `
usage: sqsc batch delete [options] <batch_name>

  Delete batch from project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
