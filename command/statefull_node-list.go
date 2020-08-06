package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// StatefullNodeListCommand is a cli.Command implementation for listing all statefull_nodes.
type StatefulNodeListCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (b *StatefulNodeListCommand) Run(args []string) int {
	b.flagSet = newFlagSet(b, b.Ui)
	endpoint := endpointFlag(b.flagSet)
	projectUUID := b.flagSet.String("project-uuid", "", "set the uuid of the project")

	if err := b.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" {
		return b.errorWithUsage(errors.New("Project uuid is mandatory"))
	}

	if b.flagSet.NArg() > 0 {
		return b.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", b.flagSet.Args()))
	}

	return b.runWithSpinner("list statefull_node", endpoint.String(), func(client *squarescale.Client) (string, error) {
		statefullNodes, err := client.GetStatefullNodes(*projectUUID)
		if err != nil {
			return "", err
		}

		//msg doit retourner la list des statefull_nodes valides : vérifier si on reçoit une liste des statefull_nodes ou des ID
		//msg := fmt.Sprintf("list of availables statefull_nodes: %v", statefull_node)

		var msg string
		msg = "Name\tNode type\tZone\t\tStatus\n"
		for _, n := range statefullNodes {
			msg += fmt.Sprintf("%s\t%s\t%s\t%s\n", n.Name, n.NodeType, n.Zone, n.Status)
		}

		if len(statefullNodes) == 0 {
			msg = "No statefull nodes found"
		}

		return msg, nil
	})
}

// Synopsis is part of cli.Command implementation.
func (b *StatefulNodeListCommand) Synopsis() string {
	return "Lists the statefull nodes of a Squarescale projects"
}

// Help is part of cli.Command implementation.
func (b *StatefulNodeListCommand) Help() string {
	helpText := `
usage: sqsc stateful node list [options]

  List all stateful nodes of a given Squarescale project.

`
	return strings.TrimSpace(helpText + optionsFromFlags(b.flagSet))
}
