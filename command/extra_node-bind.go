package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ExtraNodeBindCommand is a cli.Command implementation for binding a extra-node.
type ExtraNodeBindCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *ExtraNodeBindCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)
	volumeName := cmd.flagSet.String("volume-name", "", "Volume name to bind")

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if *volumeName == "" {
		return cmd.errorWithUsage(errors.New("Volume name cannot be empty"))
	}

	extraNodeName, err := extraNodeNameArg(cmd.flagSet, 0)
	if err != nil {
		return cmd.errorWithUsage(err)
	}

	if cmd.flagSet.NArg() > 1 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	res := cmd.runWithSpinner("bind extra-node", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		msg := fmt.Sprintf("Successfully binded volume '%s' to extra-node '%s'", *volumeName, extraNodeName)
		err = client.BindVolumeOnExtraNode(UUID, extraNodeName, *volumeName)
		return msg, err
	})
	if res != 0 {
		return res
	}

	return 0
}

// Synopsis is part of cli.Command implementation.
func (cmd *ExtraNodeBindCommand) Synopsis() string {
	return "Bind extra-node to project."
}

// Help is part of cli.Command implementation.
func (cmd *ExtraNodeBindCommand) Help() string {
	helpText := `
usage: sqsc extra-node bind [options] <extra_node_name>

  Bind extra-node to project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
