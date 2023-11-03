package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

type ClusterMemberSetCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *ClusterMemberSetCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)
	nameArg := clusterNameFlag(cmd.flagSet)
	schedulingGroupArg := clusterSchedulingGroupsFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if cmd.flagSet.NArg() > 0 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	if *nameArg == "" {
		return cmd.errorWithUsage(errors.New("Node name is mandatory."))
	}

	return cmd.runWithSpinner("configure node settings", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		clusterMember, err := client.GetClusterMemberInfo(UUID, *nameArg)
		if err != nil {
			return "", err
		}

		schedulingGroup, err := client.GetSchedulingGroupInfo(UUID, *schedulingGroupArg)
		if err != nil {
			return "", err
		}

		newClusterMember := squarescale.ClusterMember{
			SchedulingGroup: schedulingGroup,
		}

		msg := fmt.Sprintf("Successfully configured cluster node '%s'", clusterMember.Name)
		err = client.ConfigClusterMember(UUID, clusterMember, newClusterMember)

		return msg, err
	})
}

// Synopsis is part of cli.Command implementation.
func (cmd *ClusterMemberSetCommand) Synopsis() string {
	return "Set node settings attached to a project cluster."
}

// Help is part of cli.Command implementation.
func (cmd *ClusterMemberSetCommand) Help() string {
	helpText := `
usage: sqsc cluster node set [options]

  Set node settings attached to a project cluster.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
