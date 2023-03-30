package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ExternalNodeListCommand is a cli.Command implementation for listing external nodes.
type ExternalNodeListCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (e *ExternalNodeListCommand) Run(args []string) int {
	e.flagSet = newFlagSet(e, e.Ui)
	endpoint := endpointFlag(e.flagSet)
	projectUUID := projectUUIDFlag(e.flagSet)
	projectName := projectNameFlag(e.flagSet)

	if err := e.flagSet.Parse(args); err != nil {
		return 1
	}

	if e.flagSet.NArg() > 0 {
		return e.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", e.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return e.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	return e.runWithSpinner("list external nodes", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		var msg string = "Name\t\t\t\tPublicIP\tStatus\n"
		for _, en := range externalNodes {
			msg += fmt.Sprintf("%s\t%s\t%s\n", en.Name, en.PublicIP, en.Status)
		}

		if len(externalNodes) == 0 {
			msg = "No external nodes"
		}

		return msg, nil
	})
}

// Synopsis is part of cli.Command implementation.
func (e *ExternalNodeListCommand) Synopsis() string {
	return "List external nodes of project"
}

// Help is part of cli.Command implementation.
func (e *ExternalNodeListCommand) Help() string {
	helpText := `
usage: sqsc external-node list [options]

  List external nodes of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(e.flagSet))
}
