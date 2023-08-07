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
func (b *ExtraNodeListCommand) Run(args []string) int {
	b.flagSet = newFlagSet(b, b.Ui)
	endpoint := endpointFlag(b.flagSet)
	projectUUID := projectUUIDFlag(b.flagSet)
	projectName := projectNameFlag(b.flagSet)

	if err := b.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return b.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if b.flagSet.NArg() > 0 {
		return b.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", b.flagSet.Args()))
	}

	return b.runWithSpinner("list extra-node", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
func (b *ExtraNodeListCommand) Synopsis() string {
	return "List extra-nodes of project"
}

// Help is part of cli.Command implementation.
func (b *ExtraNodeListCommand) Help() string {
	helpText := `
usage: sqsc extra-node list [options]

  List extra-nodes of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(b.flagSet))
}
