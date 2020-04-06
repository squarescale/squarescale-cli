package command

import (
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
	projectArg := projectFlag(b.flagSet)

	if err := b.flagSet.Parse(args); err != nil {
		return 1
	}

	if b.flagSet.NArg() > 0 {
		return b.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", b.flagSet.Args()))
	}

	return b.runWithSpinner("list batch", endpoint.String(), func(client *squarescale.Client) (string, error) {
		batch, err := client.GetBatches(*projectArg)
		if err != nil {
			return "", err
		}

		//msg doit retourner la list des batches valides : vérifier si on reçoit une liste des batches ou des ID
		msg := fmt.Sprintf("list of availables batches: %v", batch)

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
