package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// BatchListCommand is a cli.Command implementation for listing all batches.
type BatchListCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (b *BatchListCommand) Run(args []string) int {
	b.flagSet = newFlagSet(b, b.Ui)
	endpoint := endpointFlag(b.flagSet)
	projectUUID := b.flagSet.String("project-uuid", "", "set the uuid of the project")

	if err := b.flagSet.Parse(args); err != nil {
		return 1
	}

	if b.flagSet.NArg() > 0 {
		return b.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", b.flagSet.Args()))
	}

	if *projectUUID == "" {
		return b.errorWithUsage(errors.New("Project uuid is mandatory"))
	}

	return b.runWithSpinner("list batch", endpoint.String(), func(client *squarescale.Client) (string, error) {
		batches, err := client.GetBatches(*projectUUID)
		if err != nil {
			return "", err
		}

		var msg string = "Name\t\tDocker image\tPeriodic\n"
		for _, b := range batches {
			msg += fmt.Sprintf("%s\t%s\t%t\n", b.BatchCommon.Name, b.DockerImage.Name, b.BatchCommon.Periodic)
		}

		if len(batches) == 0 {
			msg = "No batch found"
		}

		return msg, nil

	})

}

// Synopsis is part of cli.Command implementation.
func (b *BatchListCommand) Synopsis() string {
	return "Lists the batches of a Squarescale projects"
}

// Help is part of cli.Command implementation.
func (b *BatchListCommand) Help() string {
	helpText := `
usage: sqsc batch list [options]

  List all batches of a given Squarescale project.

`
	return strings.TrimSpace(helpText + optionsFromFlags(b.flagSet))
}
