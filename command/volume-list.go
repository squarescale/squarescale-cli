package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// VolumeListCommand is a cli.Command implementation for listing all volumes.
type VolumeListCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (b *VolumeListCommand) Run(args []string) int {
	b.flagSet = newFlagSet(b, b.Ui)
	endpoint := endpointFlag(b.flagSet)
	projectUUID := b.flagSet.String("project-uuid", "", "set the uuid of the project")

	if err := b.flagSet.Parse(args); err != nil {
		return 1
	}

	if b.flagSet.NArg() > 0 {
		return b.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", b.flagSet.Args()))
	}

	if *projectUUID == "" {
		return b.errorWithUsage(errors.New("Project uuid is mandatory"))
	}

	return b.runWithSpinner("list volume", endpoint.String(), func(client *squarescale.Client) (string, error) {
		volumes, err := client.GetVolumes(*projectUUID)
		if err != nil {
			return "", err
		}

		//msg doit retourner la list des volumes valides : vérifier si on reçoit une liste des volumes ou des ID
		//msg := fmt.Sprintf("list of availables volumes: %v", volume)

		var msg string = "Name\tSize\tType\tZone\t\tNode\tStatus\n"
		for _, v := range volumes {
			msg += fmt.Sprintf("%s\t%d\t%s\t%s\t%s\t%s\n", v.Name, v.Size, v.Type, v.Zone, v.StatefullNodeName, v.Status)
		}

		if len(volumes) == 0 {
			msg = "No volumes found"
		}

		return msg, nil

	})

}

// Synopsis is part of cli.Command implementation.
func (b *VolumeListCommand) Synopsis() string {
	return "Lists the volumes of a Squarescale projects"
}

// Help is part of cli.Command implementation.
func (b *VolumeListCommand) Help() string {
	helpText := `
usage: sqsc volume list [options]

  List all volumes of a given Squarescale project.

`
	return strings.TrimSpace(helpText + optionsFromFlags(b.flagSet))
}
