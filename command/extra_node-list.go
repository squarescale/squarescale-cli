package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ExtraNodeListCommand is a cli.Command implementation for listing extra-nodes.
type ExtraNodeListCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *ExtraNodeListCommand) Run(args []string) int {
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

	if cmd.flagSet.NArg() > 0 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	return cmd.runWithSpinner("list extra-node", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		extraNodes, err := client.GetExtraNodes(UUID)
		if err != nil {
			return "", err
		}

		//msg doit retourner la list des extra-nodes valides : verifier si on recoit une liste des extra-nodes ou des ID
		//msg := fmt.Sprintf("list of availables extra-nodes: %v", extra-node)

		var msg string
		msg = "Name\tNode type\tZone\t\tStatus\n"
		for _, n := range extraNodes {
			msg += fmt.Sprintf("%s\t%s\t%s\t%s\n", n.Name, n.NodeType, n.Zone, n.Status)
		}

		if len(extraNodes) == 0 {
			msg = "No extra-node found"
		}

		return msg, nil
	})
}

// Synopsis is part of cli.Command implementation.
func (cmd *ExtraNodeListCommand) Synopsis() string {
	return "List extra-nodes of project"
}

// Help is part of cli.Command implementation.
func (cmd *ExtraNodeListCommand) Help() string {
	helpText := `
usage: sqsc extra-node list [options]

  List extra-nodes of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
