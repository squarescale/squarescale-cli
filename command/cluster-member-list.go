package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ClusterMemberListCommand is a cli.Command implementation for listing external nodes.
type ClusterMemberListCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *ClusterMemberListCommand) Run(args []string) int {
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

	return cmd.runWithSpinner("list cluster nodes", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		clusterMembers, err := client.GetClusterMembers(UUID)
		if err != nil {
			return "", err
		}

		var msg string = "Name\t\t\t\tPublicIP\tPrivateIP\tStatus\tScheduling group\n"
		for _, cm := range clusterMembers {
			if cm.PublicIP == "" {
				msg += fmt.Sprintf("%s\t\t\t\t%s\t%s\t%s\n", cm.Name, cm.PrivateIP, cm.Status, cm.SchedulingGroup.Name)
			} else {
				msg += fmt.Sprintf("%s\t%s\t%s\t%s\t%s\n", cm.Name, cm.PublicIP, cm.PrivateIP, cm.Status, cm.SchedulingGroup.Name)
			}
		}

		if len(clusterMembers) == 0 {
			msg = "No cluster members"
		}

		return msg, nil
	})
}

// Synopsis is part of cli.Command implementation.
func (cmd *ClusterMemberListCommand) Synopsis() string {
	return "List cluster nodes of project"
}

// Help is part of cli.Command implementation.
func (cmd *ClusterMemberListCommand) Help() string {
	helpText := `
usage: sqsc cluster node list [options]

  List cluster nodes of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
