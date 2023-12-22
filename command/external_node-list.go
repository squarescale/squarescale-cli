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

// ExternalNodeListCommand is a cli.Command implementation for listing external nodes.
type ExternalNodeListCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *ExternalNodeListCommand) Run(args []string) int {
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

	return cmd.runWithSpinner("list external nodes", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		externalNodes, err := client.GetExternalNodes(UUID)
		if err != nil {
			return "", err
		}

		tableString := &strings.Builder{}
		table := tablewriter.NewWriter(tableString)
		// reset by ui/table.go FormatTable function: table.SetAutoFormatHeaders(false)
		// seems like this should be taken into account earlier than in the ui/table.go FormatTable function to have effect on fields
		table.SetAutoWrapText(false)
		table.SetHeader([]string{"Name", "PublicIP", "Status"})

		for _, en := range externalNodes {
			table.Append([]string{
				en.Name,
				en.PublicIP,
				en.Status,
			})
		}

		if len(externalNodes) == 0 {
			return "No external nodes", nil
		}

		ui.FormatTable(table)

		table.Render()
		// Remove trailing \n and HT
		return string(regexp.MustCompile(`[\n\x09][\n\x09]*$`).ReplaceAll([]byte(tableString.String()), []byte(""))), nil
	})
}

// Synopsis is part of cli.Command implementation.
func (cmd *ExternalNodeListCommand) Synopsis() string {
	return "List external nodes of project"
}

// Help is part of cli.Command implementation.
func (cmd *ExternalNodeListCommand) Help() string {
	helpText := `
usage: sqsc external-node list [options]

  List external nodes of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
