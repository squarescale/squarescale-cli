package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ExternalNodeGetCommand is a cli.Command implementation for get a external node.
type ExternalNodeGetCommand struct {
	Meta
	flagSet *flag.FlagSet
}

func (b *ExternalNodeGetCommand) Run(args []string) int {
	b.flagSet = newFlagSet(b, b.Ui)
	endpoint := endpointFlag(b.flagSet)
	projectUUID := b.flagSet.String("project-uuid", "", "set the uuid of the project")
	projectName := b.flagSet.String("project-name", "", "set the name of the project")

	if err := b.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return b.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	externalNodeName, err := externalNodeNameArg(b.flagSet, 0)
	if err != nil {
		return b.errorWithUsage(err)
	}

	if b.flagSet.NArg() > 1 {
		return b.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", b.flagSet.Args()))
	}

	return b.runWithSpinner("show external nodes", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		externalNode, err := client.GetExternalNodeInfo(UUID, externalNodeName)
		if err != nil {
			return "", err
		}

		var msg string
		msg = ""
		msg += fmt.Sprintf("Name:\t\t%s\n", externalNode.Name)
		msg += fmt.Sprintf("Public IP:\t%s\n", externalNode.PublicIP)
		msg += fmt.Sprintf("Status:\t\t%s\n", externalNode.Status)

		return msg, nil
	})
}

// Synopsis is part of cli.Command implementation.
func (b *ExternalNodeGetCommand) Synopsis() string {
	return "Get a external node of project"
}

// Help is part of cli.Command implementation.
func (b *ExternalNodeGetCommand) Help() string {
	helpText := `
usage: sqsc external-node get [options] <external_node_name>

  Get a external node of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(b.flagSet))
}
