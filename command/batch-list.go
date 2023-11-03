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
func (cmd *BatchListCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if cmd.flagSet.NArg() > 0 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	return cmd.runWithSpinner("list batch", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		batches, err := client.GetBatches(UUID)
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
func (cmd *BatchListCommand) Synopsis() string {
	return "List batches of project"
}

// Help is part of cli.Command implementation.
func (cmd *BatchListCommand) Help() string {
	helpText := `
usage: sqsc batch list [options]

  List batches of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
