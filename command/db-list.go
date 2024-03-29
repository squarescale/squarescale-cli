package command

import (
	"errors"
	"flag"
	"fmt"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/squarescale/squarescale-cli/squarescale"
	"github.com/squarescale/squarescale-cli/ui"
)

// DBListCommand is a cli.Command implementation to list the database engines and instances available for use.
type DBListCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *DBListCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	provider := providerFlag(cmd.flagSet)
	region := regionFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if cmd.flagSet.NArg() > 0 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	if *provider == "" {
		return cmd.errorWithUsage(errors.New("Cloud provider is mandatory"))
	}

	if *region == "" {
		return cmd.errorWithUsage(errors.New("Cloud provider region is mandatory"))
	}

	return cmd.runWithSpinner("list available database engines and sizes", endpoint.String(), func(client *squarescale.Client) (string, error) {
		engines, err := client.GetAvailableDBEngines(*provider, *region)
		if err != nil {
			return "", err
		}

		sizes, err := client.GetAvailableDBSizes(*provider, *region)
		if err != nil {
			return "", err
		}

		var out string = "\n\tAvailable engines\n\n"
		out += fmtDbEngineListOutput(engines)
		out += "\n\n\tAvailable sizes\n\n"
		out += fmtDbSizeListOutput(sizes)
		return out, nil
	})
}

// Synopsis is part of cli.Command implementation.
func (cmd *DBListCommand) Synopsis() string {
	return "List available database engines and sizes"
}

// Help is part of cli.Command implementation.
func (cmd *DBListCommand) Help() string {
	helpText := `
usage: sqsc db list

  List database engines and sizes available for use
  on the SquareScale platform.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}

func fmtDbEngineListOutput(engines []squarescale.DataseEngine) string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	table.SetHeader([]string{"Name", "Label", "Version"})
	data := make([][]string, len(engines), len(engines))
	for i, engine := range engines {
		data[i] = []string{engine.Name, engine.Label, engine.Version}
	}

	table.AppendBulk(data)

	ui.FormatTable(table)

	table.Render()

	return string(regexp.MustCompile(`[\n\x09][\n\x09]*$`).ReplaceAll([]byte(tableString.String()), []byte("")))
}

func fmtDbSizeListOutput(sizes []squarescale.DataseSize) string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Name", "Description"})
	data := make([][]string, len(sizes), len(sizes))

	for i, size := range sizes {
		data[i] = []string{size.Name, size.Description}
	}

	table.AppendBulk(data)

	ui.FormatTable(table)

	table.Render()

	return string(regexp.MustCompile(`[\n\x09][\n\x09]*$`).ReplaceAll([]byte(tableString.String()), []byte("")))
}
