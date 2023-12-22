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
func (cmd *VolumeListCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if cmd.flagSet.NArg() > 0 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	return cmd.runWithSpinner("list volume", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
			msg += fmt.Sprintf("%s\t%d\t%s\t%s\t%s\t%s\n", v.Name, v.Size, v.Type, v.Zone, v.ExtraNodeName, v.Status)
		}

		if len(volumes) == 0 {
			msg = "No volumes found"
		}

		return msg, nil

	})

}

// Synopsis is part of cli.Command implementation.
func (cmd *VolumeListCommand) Synopsis() string {
	return "List volumes of project"
}

// Help is part of cli.Command implementation.
func (cmd *VolumeListCommand) Help() string {
	helpText := `
usage: sqsc volume list [options]

  List volumes of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
