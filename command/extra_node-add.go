package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ExtraNodeAddCommand is a cli.Command implementation for creating a extra-node.
type ExtraNodeAddCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *ExtraNodeAddCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)

	nodeType := cmd.flagSet.String("node-type", "dev", "Extra-node type")
	zone := cmd.flagSet.String("zone", "eu-west-1a", "Extra-node zone")
	nowait := nowaitFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	extraNodeName, err := extraNodeNameArg(cmd.flagSet, 0)
	if err != nil {
		return cmd.errorWithUsage(err)
	}

	if cmd.flagSet.NArg() > 1 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}
	var UUID string

	res := cmd.runWithSpinner("add extra-node", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		msg := fmt.Sprintf("Successfully added extra-node '%s' to project '%s'", extraNodeName, projectToShow)
		_, err = client.AddExtraNode(UUID, extraNodeName, *nodeType, *zone)
		return msg, err
	})
	if res != 0 {
		return res
	}

	if !*nowait {
		cmd.runWithSpinner("wait for extra-node add", endpoint.String(), func(client *squarescale.Client) (string, error) {
			extraNode, err := client.WaitExtraNode(UUID, extraNodeName, 5)
			if err != nil {
				return "", err
			} else {
				return extraNode.Name, nil
			}
		})
	}

	return 0
}

// Synopsis is part of cli.Command implementation.
func (cmd *ExtraNodeAddCommand) Synopsis() string {
	return "Add extra-node to project."
}

// Help is part of cli.Command implementation.
func (cmd *ExtraNodeAddCommand) Help() string {
	helpText := `
usage: sqsc extra-node add [options] <extra_node_name>

  Add extra-node to project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
