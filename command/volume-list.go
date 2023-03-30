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
	projectUUID :=projectUUIDFlag(b.flagSet)
	projectName := projectNameFlag(b.flagSet)

	if err := b.flagSet.Parse(args); err != nil {
		return 1
	}

	if b.flagSet.NArg() > 0 {
		return b.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", b.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return b.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	return b.runWithSpinner("list volume", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		volumes, err := client.GetVolumes(UUID)
		if err != nil {
			return "", err
		}

		//msg doit retourner la list des volumes valides : verifier si on recoit une liste des volumes ou des ID
		//msg := fmt.Sprintf("list of availables volumes: %v", volume)

		var msg string = "Name\tSize\tType\tZone\t\tNode\tStatus\n"
		for _, v := range volumes {
			msg += fmt.Sprintf("%s\t%d\t%s\t%s\t%s\t%s\n", v.Name, v.Size, v.Type, v.Zone, v.StatefulNodeName, v.Status)
		}

		if len(volumes) == 0 {
			msg = "No volumes found"
		}

		return msg, nil

	})

}

// Synopsis is part of cli.Command implementation.
func (b *VolumeListCommand) Synopsis() string {
	return "List volumes of project"
}

// Help is part of cli.Command implementation.
func (b *VolumeListCommand) Help() string {
	helpText := `
usage: sqsc volume list [options]

  List volumes of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(b.flagSet))
}
