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

func (cmd *ExternalNodeGetCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	externalNodeName, err := externalNodeNameArg(cmd.flagSet, 0)
	if err != nil {
		return cmd.errorWithUsage(err)
	}

	if cmd.flagSet.NArg() > 1 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	return cmd.runWithSpinner("show external nodes", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
func (cmd *ExternalNodeGetCommand) Synopsis() string {
	return "Get a external node of project"
}

// Help is part of cli.Command implementation.
func (cmd *ExternalNodeGetCommand) Help() string {
	helpText := `
usage: sqsc external-node get [options] <external_node_name>

  Get a external node of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
