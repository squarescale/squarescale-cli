package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ExternalNodeAddCommand is a cli.Command implementation for creating a external-node.
type ExternalNodeAddCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *ExternalNodeAddCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)

	publicIP := externalNodePublicIP(cmd.flagSet)
	nowait := nowaitFlag(cmd.flagSet)

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

	if *publicIP == "" {
		return cmd.errorWithUsage(errors.New("Public IP is mandatory"))
	}

	if cmd.flagSet.NArg() > 1 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}
	var UUID string

	res := cmd.runWithSpinner("add external node", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		msg := fmt.Sprintf("Successfully added external node '%s' to project '%s'", externalNodeName, projectToShow)
		_, err = client.AddExternalNode(UUID, externalNodeName, *publicIP)
		return msg, err
	})
	if res != 0 {
		return res
	}

	if !*nowait {
		cmd.runWithSpinner("wait for external node add", endpoint.String(), func(client *squarescale.Client) (string, error) {
			// can also be externalNode, err := client.WaitExternalNode(UUID, externalNodeName, 5, []string{"provisionned", "inconsistent"})
			externalNode, err := client.WaitExternalNode(UUID, externalNodeName, 5, []string{})
			if err != nil {
				return "", err
			} else {
				return externalNode.Name, nil
			}
		})
	}

	return 0
}

// Synopsis is part of cli.Command implementation.
func (cmd *ExternalNodeAddCommand) Synopsis() string {
	return "Add external node node to project."
}

// Help is part of cli.Command implementation.
func (cmd *ExternalNodeAddCommand) Help() string {
	helpText := `
usage: sqsc external-node add [options] <external_node_name>

  Add external node node to project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
