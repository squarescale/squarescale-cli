package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// StatefulNodeListCommand is a cli.Command implementation for listing stateful-nodes.
type StatefulNodeListCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (b *StatefulNodeListCommand) Run(args []string) int {
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

	if b.flagSet.NArg() > 0 {
		return b.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", b.flagSet.Args()))
	}

	return b.runWithSpinner("list stateful-node", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		statefulNodes, err := client.GetStatefulNodes(UUID)
		if err != nil {
			return "", err
		}

		//msg doit retourner la list des stateful-nodes valides : verifier si on recoit une liste des stateful-nodes ou des ID
		//msg := fmt.Sprintf("list of availables stateful-nodes: %v", stateful-node)

		var msg string
		msg = "Name\tNode type\tZone\t\tStatus\n"
		for _, n := range statefulNodes {
			msg += fmt.Sprintf("%s\t%s\t%s\t%s\n", n.Name, n.NodeType, n.Zone, n.Status)
		}

		if len(statefulNodes) == 0 {
			msg = "No stateful-node found"
		}

		return msg, nil
	})
}

// Synopsis is part of cli.Command implementation.
func (b *StatefulNodeListCommand) Synopsis() string {
	return "List stateful nodes of project"
}

// Help is part of cli.Command implementation.
func (b *StatefulNodeListCommand) Help() string {
	helpText := `
usage: sqsc stateful node list [options]

  List stateful nodes of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(b.flagSet))
}
